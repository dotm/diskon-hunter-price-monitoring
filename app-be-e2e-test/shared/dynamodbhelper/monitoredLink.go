package dynamodbhelper

import (
	"diskon-hunter/price-monitoring-e2e-test/shared/createerror"
	monitoredLink "diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"diskon-hunter/price-monitoring-e2e-test/shared/sliceutil"
	"strings"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetStlMonitoredLinkDetailByUrlList(
	dynamoDBClient *dynamodb.DynamoDB,
	cleanedMonitoredLinkUrlList []string,
) (
	monitoredLinkList []monitoredLink.StlMonitoredLinkDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(cleanedMonitoredLinkUrlList) == 0 {
		return monitoredLinkList, nil, nil
	}
	cleanedMonitoredLinkUrlList = sliceutil.RemoveDuplicateElements(cleanedMonitoredLinkUrlList)

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(cleanedMonitoredLinkUrlList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubMonitoredLinkUrl": {
				S: aws.String(cleanedMonitoredLinkUrlList[i]),
			},
		})
	}
	tableName := monitoredLink.GetStlMonitoredLinkDetailDynamoDBTableV1()
	batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return monitoredLinkList, createerror.InternalException(err), err
	}
	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		monitoredLinkDAO := monitoredLink.StlMonitoredLinkDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&monitoredLinkDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling monitoredLinkDAO: %v", err)
			return monitoredLinkList, createerror.InternalException(err), err
		}
		monitoredLinkList = append(monitoredLinkList, monitoredLinkDAO)
	}
	if len(monitoredLinkList) < len(cleanedMonitoredLinkUrlList) {
		subsetUrlListThatIsNotInSuperset := []string{}
		supersetUrlMap := map[string]bool{}
		for i := 0; i < len(cleanedMonitoredLinkUrlList); i++ {
			supersetUrlMap[cleanedMonitoredLinkUrlList[i]] = true
		}
		for i := 0; i < len(monitoredLinkList); i++ {
			exist, ok := supersetUrlMap[monitoredLinkList[i].HubMonitoredLinkUrl]
			if !ok || !exist {
				subsetUrlListThatIsNotInSuperset = append(subsetUrlListThatIsNotInSuperset, monitoredLinkList[i].HubMonitoredLinkUrl)
			}
		}
		err := fmt.Errorf("error can't find StlMonitoredLinkDetail urls: %s", strings.Join(subsetUrlListThatIsNotInSuperset, ", "))
		return monitoredLinkList, createerror.ClientBadRequest(err), err
	}

	return monitoredLinkList, nil, nil
}

func GetStlUserMonitorsLinkDetailList(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
	cleanedMonitoredLinkUrlList []string,
) (
	userMonitorsLinkDetailList []monitoredLink.StlUserMonitorsLinkDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
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
	batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return userMonitorsLinkDetailList, createerror.InternalException(err), err
	}

	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		userMonitorsLinkDetailDAO := monitoredLink.StlUserMonitorsLinkDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&userMonitorsLinkDetailDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling userMonitorsLinkDetailDAO: %v", err)
			return userMonitorsLinkDetailList, createerror.InternalException(err), err
		}
		userMonitorsLinkDetailList = append(userMonitorsLinkDetailList, userMonitorsLinkDetailDAO)
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

func DeleteStlMonitoredLinkDetailByUrlList(
	dynamoDBClient *dynamodb.DynamoDB,
	cleanedMonitoredLinkUrlList []string,
) (
	*serverresponse.ErrorObj,
	error,
) {
	if len(cleanedMonitoredLinkUrlList) == 0 {
		return nil, nil
	}

	for _, url := range cleanedMonitoredLinkUrlList {
		var tableName string
		var err error

		tableName = monitoredLink.GetStlMonitoredLinkDetailDynamoDBTableV1()
		_, err = dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"HubMonitoredLinkUrl": {
					S: aws.String(url),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			err = fmt.Errorf("error DeleteItem from %s: %v", tableName, err)
			return createerror.InternalException(err), err
		}
	}

	return nil, nil
}

func DeleteStlUserMonitorsLinkDetailByUrlList(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
	cleanedMonitoredLinkUrlList []string,
) (
	*serverresponse.ErrorObj,
	error,
) {
	if len(cleanedMonitoredLinkUrlList) == 0 {
		return nil, nil
	}

	for _, url := range cleanedMonitoredLinkUrlList {
		var tableName string
		var err error

		tableName = monitoredLink.GetStlUserMonitorsLinkDetailDynamoDBTableV1()
		_, err = dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"HubUserId": {
					S: aws.String(userId),
				},
				"HubMonitoredLinkUrl": {
					S: aws.String(url),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			err = fmt.Errorf("error DeleteItem from %s: %v", tableName, err)
			return createerror.InternalException(err), err
		}
	}

	return nil, nil
}
