package monitoredLinkList

import (
	"context"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/currencyutil"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*
Query represents input from client through api to read, filter, map, reduce, and move data
from a data store to another place (web UI, pdf, xls, data lake, other SaaS, etc.).
*/
type QueryV1 struct {
	Version         string //should follow the struct name suffix
	RequesterUserId string
}

func (x QueryV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableQuery := x //no sensitive info and no large fields so we'll just use x
	byteSlice, err := json.Marshal(loggableQuery)
	if err != nil {
		return "", err
	} else {
		return string(byteSlice), nil
	}
}

type QueryV1Dependencies struct {
	Logger         *lazylogger.Instance
	DynamoDBClient *dynamodb.DynamoDB
}

type QueryV1DataResponse = []QueryV1UserLinkDetail
type QueryV1UserLinkDetail struct {
	monitoredLink.StlUserMonitorsLinkDetailDAOV1
	LatestPrice        *currencyutil.Currency //can be nil
	TimeLatestScrapped *time.Time
}

/*
Addition, change, or removal of validation and presentation might cause version increment
*/

func QueryV1Handler(
	ctx context.Context,
	dependencies QueryV1Dependencies,
	query QueryV1,
) (QueryV1DataResponse, *serverresponse.ErrorObj) {
	//don't mutate this. emptyResponse should be used when returning error.
	emptyResponse := QueryV1DataResponse{}

	//log the query
	loggableQuery, err := query.createLoggableString()
	if err != nil {
		err = fmt.Errorf("error creating loggable string: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	dependencies.Logger.EnqueueQueryLog(loggableQuery, true)

	/* Validations
	Validations from authorization only.
	*/

	/* Retrieving Data
	Retrieving data from a data store.
	*/

	userMonitorsLinkList, errObj, err := dynamodbhelper.GetUserMonitorsLinkListOfUserId(
		dependencies.DynamoDBClient, query.RequesterUserId,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}
	urlList := []string{}
	for _, userMonitorsLinkDetail := range userMonitorsLinkList {
		urlList = append(urlList, userMonitorsLinkDetail.HubMonitoredLinkUrl)
	}
	monitoredLinkList, errObj, err := dynamodbhelper.GetMonitoredLinkListOfUrl(
		dependencies.DynamoDBClient, urlList,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}
	urlToMonitoredLinkDetailMap := map[string]monitoredLink.StlMonitoredLinkDetailDAOV1{}
	for _, monitoredLinkDetail := range monitoredLinkList {
		urlToMonitoredLinkDetailMap[monitoredLinkDetail.HubMonitoredLinkUrl] = monitoredLinkDetail
	}
	userLinkList := []QueryV1UserLinkDetail{}
	for _, userMonitorsLinkDetail := range userMonitorsLinkList {
		monitoredLinkDetail, ok := urlToMonitoredLinkDetailMap[userMonitorsLinkDetail.HubMonitoredLinkUrl]
		if !ok {
			err = fmt.Errorf("url not found in urlToMonitoredLinkDetailMap: %v", userMonitorsLinkDetail.HubMonitoredLinkUrl)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
		userLinkList = append(userLinkList, QueryV1UserLinkDetail{
			StlUserMonitorsLinkDetailDAOV1: userMonitorsLinkDetail,
			LatestPrice:                    monitoredLinkDetail.LatestPrice,
			TimeLatestScrapped:             monitoredLinkDetail.TimeLatestScrapped,
		})
	}

	/* Presenting Data
	Processing the data (filter, map, reduce, etc.)
	so that the data is ready to be consumed by the client.
	*/

	/* Sending Data
	Unless the data is presented to the client (using return from this function),
	you'll need to move it to the intended data store.
	*/
	return userLinkList, nil
}
