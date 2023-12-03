package monitoredLinkAddMultiple

import (
	"context"
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	monitoredLink "diskon-hunter/price-monitoring/src/monitoredLink"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version           string //should follow the struct name suffix
	RequesterUserId   string
	MonitoredLinkList []MonitoredLinkDetailCommandV1
}

func NewCommandV1(
	Version string, //should follow the struct name suffix
	RequesterUserId string,
	MonitoredLinkList []MonitoredLinkDetailCommandV1,
) CommandV1 {
	return CommandV1{
		Version:           Version,
		RequesterUserId:   RequesterUserId,
		MonitoredLinkList: MonitoredLinkList,
	}
}

type MonitoredLinkDetailCommandV1 struct {
	HubMonitoredLinkUrl string
	AlertPrice          currencyutil.Currency
	AlertMethodList     []constenum.AlertMethod
}

func NewMonitoredLinkDetailCommandV1(
	HubMonitoredLinkUrl string,
	AlertPrice currencyutil.Currency,
	AlertMethodList []constenum.AlertMethod,
) MonitoredLinkDetailCommandV1 {
	return MonitoredLinkDetailCommandV1{
		HubMonitoredLinkUrl: HubMonitoredLinkUrl,
		AlertPrice:          AlertPrice,
		AlertMethodList:     AlertMethodList,
	}
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := x //no sensitive info and no large fields so we'll just use x
	byteSlice, err := json.Marshal(loggableCommand)
	if err != nil {
		return "", err
	} else {
		return string(byteSlice), nil
	}
}

type CommandV1Dependencies struct {
	Logger         *lazylogger.Instance
	DynamoDBClient *dynamodb.DynamoDB
}

type CommandV1DataResponse struct {
	HubUserId                    string
	MonitoredLinkRawToCleanedMap map[string]string
}

/*
Addition, change, or removal of validation might cause version increment
*/
func CommandV1Handler(
	ctx context.Context,
	dependencies CommandV1Dependencies,
	command CommandV1,
) (CommandV1DataResponse, *serverresponse.ErrorObj) {
	//don't mutate this. emptyResponse should be used when returning error.
	emptyResponse := CommandV1DataResponse{}
	//log the command
	loggableCommand, err := command.createLoggableString()
	if err != nil {
		err = fmt.Errorf("error creating loggable string: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	dependencies.Logger.EnqueueCommandLog(loggableCommand, true)

	/* Validations
	Validations from auth, write model,
	or domain model's business logic (from projections or from events replay).
	*/

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/

	stlUserMonitorsLinkDAOList := []monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
	timeExpired := time.Now().AddDate(1, 0, 0) //hardcode 1 year subscription
	monitoredLinkRawToCleanedMap := map[string]string{}
	for i := 0; i < len(command.MonitoredLinkList); i++ {
		rawUrl := command.MonitoredLinkList[i].HubMonitoredLinkUrl
		cleanedUrl := rawUrl //remove irrelevant query params ~kodok
		monitoredLinkRawToCleanedMap[rawUrl] = cleanedUrl
		stlUserMonitorsLinkDAOList = append(stlUserMonitorsLinkDAOList, monitoredLink.NewStlUserMonitorsLinkDetailDAOV1(
			command.RequesterUserId,                      //HubUserId
			cleanedUrl,                                   //HubMonitoredLinkUrl
			command.MonitoredLinkList[i].AlertPrice,      //AlertPrice
			command.MonitoredLinkList[i].AlertMethodList, //AlertMethodList
			timeExpired,                                  //TimeExpired
		))
	}

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(stlUserMonitorsLinkDAOList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubMonitoredLinkUrl": {
				S: aws.String(stlUserMonitorsLinkDAOList[i].HubMonitoredLinkUrl),
			},
		})
	}
	tableName := monitoredLink.GetStlMonitoredLinkDetailDynamoDBTableV1()
	batchGetItemOutput, err := dependencies.DynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return emptyResponse, createerror.InternalException(err)
	}
	existingStlMonitoredLinkDetailMap := map[string]monitoredLink.StlMonitoredLinkDetailDAOV1{}
	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		monitoredLinkDAO := monitoredLink.StlMonitoredLinkDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&monitoredLinkDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling monitoredLinkDAO: %v", err)
			return emptyResponse, createerror.InternalException(err)
		}
		existingStlMonitoredLinkDetailMap[monitoredLinkDAO.HubMonitoredLinkUrl] = monitoredLinkDAO
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/
	transactItems := []*dynamodb.TransactWriteItem{}

	//persist StlUserMonitorsLinkDAOList
	for i := 0; i < len(stlUserMonitorsLinkDAOList); i++ {
		stlUserMonitorsLinkDAO := stlUserMonitorsLinkDAOList[i]
		stlUserMonitorsLinkDAOItem, err := dynamodbattribute.MarshalMap(stlUserMonitorsLinkDAO)
		if err != nil {
			err = fmt.Errorf("error marshaling stlUserMonitorsLinkDAO: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}

		transactItems = append(
			transactItems,
			&dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					Item:      stlUserMonitorsLinkDAOItem,
					TableName: aws.String(monitoredLink.GetStlUserMonitorsLinkDetailDynamoDBTableV1()),
				},
			},
		)
	}
	//insert or update StlMonitoredLinkDetailDAOV1
	for i := 0; i < len(stlUserMonitorsLinkDAOList); i++ {
		//if StlMonitoredLink exist in db,
		// just update the timeExpired (don't update LatestPrice or TimeLatestScrapped)
		// else insert cleanedUrl to HubMonitoredLinkUrl plus TimeExpired (LatestPrice and TimeLatestScrapped should be nil)

		monitoredLinkUrl := stlUserMonitorsLinkDAOList[i].HubMonitoredLinkUrl
		existingStlMonitoredLinkDetailDAO, ok := existingStlMonitoredLinkDetailMap[monitoredLinkUrl]
		var stlMonitoredLinkDAO monitoredLink.StlMonitoredLinkDetailDAOV1
		if ok {
			stlMonitoredLinkDAO = monitoredLink.NewStlMonitoredLinkDetailDAOV1(
				monitoredLinkUrl, //HubMonitoredLinkUrl
				existingStlMonitoredLinkDetailDAO.LatestPrice,        //LatestPrice
				existingStlMonitoredLinkDetailDAO.TimeLatestScrapped, //TimeLatestScrapped
				timeExpired, //TimeExpired
			)
		} else {
			stlMonitoredLinkDAO = monitoredLink.NewStlMonitoredLinkDetailDAOV1(
				monitoredLinkUrl, //HubMonitoredLinkUrl
				nil,              //LatestPrice
				nil,              //TimeLatestScrapped
				timeExpired,      //TimeExpired
			)
		}

		stlMonitoredLinkDAOItem, err := dynamodbattribute.MarshalMap(stlMonitoredLinkDAO)
		if err != nil {
			err = fmt.Errorf("error marshaling stlMonitoredLinkDAO: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
		transactItems = append(
			transactItems,
			&dynamodb.TransactWriteItem{
				Put: &dynamodb.Put{
					Item:      stlMonitoredLinkDAOItem,
					TableName: aws.String(monitoredLink.GetStlMonitoredLinkDetailDynamoDBTableV1()),
				},
			},
		)
	}

	errObj, err := dynamodbhelper.TransactWriteItemsInWaves(dependencies.DynamoDBClient, transactItems)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{
		HubUserId:                    command.RequesterUserId,
		MonitoredLinkRawToCleanedMap: monitoredLinkRawToCleanedMap,
	}, nil
}
