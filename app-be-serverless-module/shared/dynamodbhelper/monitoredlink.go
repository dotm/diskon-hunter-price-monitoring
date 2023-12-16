package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
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

func GetMonitoredLinkIdListOfUserId(
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
