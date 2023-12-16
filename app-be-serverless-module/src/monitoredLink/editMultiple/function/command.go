package monitoredLinkEditMultiple

import (
	"context"
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version           string
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
	HubMonitoredLinkUrl   string
	ActiveAlertMethodList []constenum.AlertMethod
	AlertPrice            currencyutil.Currency
}

func NewMonitoredLinkDetailCommandV1(
	HubMonitoredLinkUrl string,
	ActiveAlertMethodList []constenum.AlertMethod,
	AlertPrice currencyutil.Currency,
) MonitoredLinkDetailCommandV1 {
	return MonitoredLinkDetailCommandV1{
		HubMonitoredLinkUrl:   HubMonitoredLinkUrl,
		ActiveAlertMethodList: ActiveAlertMethodList,
		AlertPrice:            AlertPrice,
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

type CommandV1DataResponse = []monitoredLink.StlUserMonitorsLinkDetailDAOV1

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

	monitoredLinkUrlMapToNewDetail := map[string]MonitoredLinkDetailCommandV1{}
	cleanedUrlList := []string{}
	for _, monitoredLinkData := range command.MonitoredLinkList {
		monitoredLinkUrlMapToNewDetail[monitoredLinkData.HubMonitoredLinkUrl] = monitoredLinkData
		//HubMonitoredLinkUrl is assumed to be cleaned because it comes from our database
		cleanedUrlList = append(cleanedUrlList, monitoredLinkData.HubMonitoredLinkUrl)
	}
	existingMonitoredLinkList, errObj, err := dynamodbhelper.GetStlUserMonitorsLinkDetailList(
		dependencies.DynamoDBClient,
		command.RequesterUserId,
		cleanedUrlList,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	transactItems := []*dynamodb.TransactWriteItem{}
	editedMonitoredLinkDAOList := []monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
	for i := 0; i < len(existingMonitoredLinkList); i++ {
		url := existingMonitoredLinkList[i].HubMonitoredLinkUrl

		//validate ActiveAlertMethod is already paid
		for _, activeAlertMethod := range monitoredLinkUrlMapToNewDetail[url].ActiveAlertMethodList {
			activeAlertMethodHasBeenPaid := false
			for _, paidAlertMethod := range existingMonitoredLinkList[i].PaidAlertMethodList {
				if activeAlertMethod == paidAlertMethod {
					activeAlertMethodHasBeenPaid = true
				}
			}
			if !activeAlertMethodHasBeenPaid {
				return emptyResponse, createerror.AlertMethodHasNotBeenPaid(fmt.Errorf("%v is unpaid", activeAlertMethod))
			}
		}

		//use existingDAO[i].field for fields that can't be changed
		//use idMapToNewDetail[id].field for fields that can be changed
		editedMonitoredLinkDAO := monitoredLink.NewStlUserMonitorsLinkDetailDAOV1(
			existingMonitoredLinkList[i].HubUserId,                    //HubUserId
			existingMonitoredLinkList[i].HubMonitoredLinkUrl,          //HubMonitoredLinkUrl
			monitoredLinkUrlMapToNewDetail[url].AlertPrice,            //AlertPrice
			monitoredLinkUrlMapToNewDetail[url].ActiveAlertMethodList, //ActiveAlertMethodList
			existingMonitoredLinkList[i].PaidAlertMethodList,          //PaidAlertMethodList
			existingMonitoredLinkList[i].TimeExpired,                  //TimeExpired
		)
		editedMonitoredLinkDAOList = append(editedMonitoredLinkDAOList, editedMonitoredLinkDAO)
		editedMonitoredLinkDAOItem, err := dynamodbattribute.MarshalMap(editedMonitoredLinkDAO)
		if err != nil {
			err = fmt.Errorf("error marshaling editedMonitoredLinkDAO: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
		transactItems = append(transactItems, &dynamodb.TransactWriteItem{
			Put: &dynamodb.Put{
				Item:      editedMonitoredLinkDAOItem,
				TableName: aws.String(monitoredLink.GetStlUserMonitorsLinkDetailDynamoDBTableV1()),
			},
		})
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	errObj, err = dynamodbhelper.TransactWriteItemsInWaves(dependencies.DynamoDBClient, transactItems)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	return editedMonitoredLinkDAOList, nil
}
