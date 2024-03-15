package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/shared/sliceutil"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetStlUserMonitorsLinkDetailList(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
	cleanedMonitoredLinkUrlList []string,
) (
	userMonitorsLinkDetailList []monitoredLink.StlUserMonitorsLinkDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(cleanedMonitoredLinkUrlList) == 0 {
		return userMonitorsLinkDetailList, nil, nil
	}
	cleanedMonitoredLinkUrlList = sliceutil.RemoveDuplicateElements(cleanedMonitoredLinkUrlList)

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(cleanedMonitoredLinkUrlList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubMonitoredLinkUrl": {
				S: aws.String(cleanedMonitoredLinkUrlList[i]),
			},
			"HubUserId": {
				S: aws.String(userId),
			},
		})
	}
	tableName := monitoredLink.GetStlUserMonitorsLinkDetailDynamoDBTableV1()
	parseResponseToDAO := func(response []map[string]*dynamodb.AttributeValue) (
		errObj *serverresponse.ErrorObj,
		err error,
	) {
		for i := 0; i < len(response); i++ {
			userMonitorsLinkDetailDAO := monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
			err = dynamodbattribute.UnmarshalMap(
				response[i],
				&userMonitorsLinkDetailDAO,
			)
			if err != nil {
				err = fmt.Errorf("error unmarshaling DAO from %s: %v", tableName, err)
				return createerror.InternalException(err), err
			}
			userMonitorsLinkDetailList = append(userMonitorsLinkDetailList, userMonitorsLinkDetailDAO)
		}
		return nil, nil
	}
	errObj, err = BatchGetItemInWaves(dynamoDBClient, tableName, batchGetItemKeys, parseResponseToDAO)
	if errObj != nil || err != nil {
		return userMonitorsLinkDetailList, errObj, err
	}

	if len(userMonitorsLinkDetailList) < len(cleanedMonitoredLinkUrlList) {
		subsetUrlListThatIsNotInSuperset := []string{}
		supersetUrlMap := map[string]bool{}
		for i := 0; i < len(cleanedMonitoredLinkUrlList); i++ {
			supersetUrlMap[cleanedMonitoredLinkUrlList[i]] = true
		}
		for i := 0; i < len(userMonitorsLinkDetailList); i++ {
			exist, ok := supersetUrlMap[userMonitorsLinkDetailList[i].HubMonitoredLinkUrl]
			if !ok || !exist {
				subsetUrlListThatIsNotInSuperset = append(subsetUrlListThatIsNotInSuperset, userMonitorsLinkDetailList[i].HubMonitoredLinkUrl)
			}
		}
		err := fmt.Errorf("error can't find StlUserMonitorsLinkDetail urls: %s", strings.Join(subsetUrlListThatIsNotInSuperset, ", "))
		return userMonitorsLinkDetailList, createerror.ClientBadRequest(err), err
	}

	return
}

func GetUserMonitorsLinkListOfUserId(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
) (
	stlUserMonitorsLinkDetailList []monitoredLink.StlUserMonitorsLinkDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue = nil
	for {
		companyToMonitoredLinkMappingRes, err := dynamoDBClient.Query(&dynamodb.QueryInput{
			TableName:              aws.String(monitoredLink.GetStlUserMonitorsLinkDetailDynamoDBTableV1()),
			KeyConditionExpression: aws.String("HubUserId = :pk"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pk": {S: aws.String(userId)},
			},
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			err = fmt.Errorf("error query companyToMonitoredLinkMappingRes: %v", err)
			return stlUserMonitorsLinkDetailList, createerror.InternalException(err), err
		}
		for i := 0; i < len(companyToMonitoredLinkMappingRes.Items); i++ {
			stlUserMonitorsLinkDetailDAO := monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
			err = dynamodbattribute.UnmarshalMap(companyToMonitoredLinkMappingRes.Items[i], &stlUserMonitorsLinkDetailDAO)
			if err != nil {
				err = fmt.Errorf("error unmarshaling stlUserMonitorsLinkDetailDAO: %v", err)
				return stlUserMonitorsLinkDetailList, createerror.InternalException(err), err
			}
			stlUserMonitorsLinkDetailList = append(stlUserMonitorsLinkDetailList, stlUserMonitorsLinkDetailDAO)
		}

		lastEvaluatedKey = companyToMonitoredLinkMappingRes.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	return stlUserMonitorsLinkDetailList, nil, nil
}

func GetMonitoredLinkListOfUrl(
	dynamoDBClient *dynamodb.DynamoDB,
	urlList []string,
) (
	stlMonitoredLinkDetailList []monitoredLink.StlMonitoredLinkDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(urlList) == 0 {
		return stlMonitoredLinkDetailList, nil, nil
	}
	urlList = sliceutil.RemoveDuplicateElements(urlList)

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(urlList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubMonitoredLinkUrl": {
				S: aws.String(urlList[i]),
			},
		})
	}
	tableName := monitoredLink.GetStlMonitoredLinkDetailDynamoDBTableV1()
	parseResponseToDAO := func(response []map[string]*dynamodb.AttributeValue) (
		errObj *serverresponse.ErrorObj,
		err error,
	) {
		for i := 0; i < len(response); i++ {
			monitoredLinkDAO := monitoredLink.StlMonitoredLinkDetailDAOV1{}
			err = dynamodbattribute.UnmarshalMap(
				response[i],
				&monitoredLinkDAO,
			)
			if err != nil {
				err = fmt.Errorf("error unmarshaling DAO from %s: %v", tableName, err)
				return createerror.InternalException(err), err
			}
			stlMonitoredLinkDetailList = append(stlMonitoredLinkDetailList, monitoredLinkDAO)
		}
		return nil, nil
	}
	errObj, err = BatchGetItemInWaves(dynamoDBClient, tableName, batchGetItemKeys, parseResponseToDAO)
	if errObj != nil || err != nil {
		return stlMonitoredLinkDetailList, errObj, err
	}

	if len(stlMonitoredLinkDetailList) < len(urlList) {
		subsetIdListThatIsNotInSuperset := []string{}
		supersetIdMap := map[string]bool{}
		for i := 0; i < len(urlList); i++ {
			supersetIdMap[urlList[i]] = true
		}
		for i := 0; i < len(stlMonitoredLinkDetailList); i++ {
			exist, ok := supersetIdMap[stlMonitoredLinkDetailList[i].HubMonitoredLinkUrl]
			if !ok || !exist {
				subsetIdListThatIsNotInSuperset = append(subsetIdListThatIsNotInSuperset, stlMonitoredLinkDetailList[i].HubMonitoredLinkUrl)
			}
		}
		err := fmt.Errorf("error can't find monitored link urls: %s", strings.Join(subsetIdListThatIsNotInSuperset, ", "))
		return stlMonitoredLinkDetailList, createerror.ClientBadRequest(err), err
	}

	return stlMonitoredLinkDetailList, nil, nil
}

func GetCombinedUserMonitoredLinkDataListOfUserId(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
	excludeLinksWithoutAlertMethod bool,
) (
	combinedUserMonitoredLinkDataList []monitoredLink.CombinedUserMonitoredLinkDataV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	unfilteredUserMonitorsLinkList, errObj, err := GetUserMonitorsLinkListOfUserId(
		dynamoDBClient, userId,
	)
	if errObj != nil || err != nil {
		return combinedUserMonitoredLinkDataList, errObj, err
	}

	userMonitorsLinkList := []monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
	if excludeLinksWithoutAlertMethod {
		for _, userMonitorsLinkDetail := range unfilteredUserMonitorsLinkList {
			if len(userMonitorsLinkDetail.ActiveAlertMethodList) > 0 {
				userMonitorsLinkList = append(userMonitorsLinkList, userMonitorsLinkDetail)
			}
		}
	} else {
		userMonitorsLinkList = unfilteredUserMonitorsLinkList
	}

	urlList := []string{}
	for _, userMonitorsLinkDetail := range userMonitorsLinkList {
		urlList = append(urlList, userMonitorsLinkDetail.HubMonitoredLinkUrl)
	}
	monitoredLinkList, errObj, err := GetMonitoredLinkListOfUrl(
		dynamoDBClient, urlList,
	)
	if errObj != nil || err != nil {
		return combinedUserMonitoredLinkDataList, errObj, err
	}
	urlToMonitoredLinkDetailMap := map[string]monitoredLink.StlMonitoredLinkDetailDAOV1{}
	for _, monitoredLinkDetail := range monitoredLinkList {
		urlToMonitoredLinkDetailMap[monitoredLinkDetail.HubMonitoredLinkUrl] = monitoredLinkDetail
	}
	for _, userMonitorsLinkDetail := range userMonitorsLinkList {
		monitoredLinkDetail, ok := urlToMonitoredLinkDetailMap[userMonitorsLinkDetail.HubMonitoredLinkUrl]
		if !ok {
			err = fmt.Errorf("url not found in urlToMonitoredLinkDetailMap: %v", userMonitorsLinkDetail.HubMonitoredLinkUrl)
			return combinedUserMonitoredLinkDataList, createerror.InternalException(err), err
		}
		combinedUserMonitoredLinkDataList = append(
			combinedUserMonitoredLinkDataList,
			monitoredLink.CombinedUserMonitoredLinkDataV1{
				StlUserMonitorsLinkDetailDAOV1: userMonitorsLinkDetail,
				LatestPrice:                    monitoredLinkDetail.LatestPrice,
				TimeLatestScrapped:             monitoredLinkDetail.TimeLatestScrapped,
			},
		)
	}

	return combinedUserMonitoredLinkDataList, nil, nil
}
