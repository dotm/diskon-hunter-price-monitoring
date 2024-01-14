package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/shared/sliceutil"
	"diskon-hunter/price-monitoring/src/searchedItem"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetStlUserSearchesItemDetailList(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
	idList []string,
) (
	userSearchesItemDetailList []searchedItem.StlUserSearchesItemDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(idList) == 0 {
		return userSearchesItemDetailList, nil, nil
	}
	idList = sliceutil.RemoveDuplicateElements(idList)

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(idList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubSearchedItemId": {
				S: aws.String(idList[i]),
			},
			"HubUserId": {
				S: aws.String(userId),
			},
		})
	}
	tableName := searchedItem.GetStlUserSearchesItemDetailDynamoDBTableV1()
	parseResponseToDAO := func(response []map[string]*dynamodb.AttributeValue) (
		errObj *serverresponse.ErrorObj,
		err error,
	) {
		for i := 0; i < len(response); i++ {
			userSearchesItemDetailDAO := searchedItem.StlUserSearchesItemDetailDAOV1{}
			err = dynamodbattribute.UnmarshalMap(
				response[i],
				&userSearchesItemDetailDAO,
			)
			if err != nil {
				err = fmt.Errorf("error unmarshaling DAO from %s: %v", tableName, err)
				return createerror.InternalException(err), err
			}
			userSearchesItemDetailList = append(userSearchesItemDetailList, userSearchesItemDetailDAO)
		}
		return nil, nil
	}
	errObj, err = BatchGetItemInWaves(dynamoDBClient, tableName, batchGetItemKeys, parseResponseToDAO)
	if errObj != nil || err != nil {
		return userSearchesItemDetailList, errObj, err
	}

	if len(userSearchesItemDetailList) < len(idList) {
		subsetIdListThatIsNotInSuperset := []string{}
		supersetIdMap := map[string]bool{}
		for i := 0; i < len(idList); i++ {
			supersetIdMap[idList[i]] = true
		}
		for i := 0; i < len(userSearchesItemDetailList); i++ {
			exist, ok := supersetIdMap[userSearchesItemDetailList[i].HubSearchedItemId]
			if !ok || !exist {
				subsetIdListThatIsNotInSuperset = append(subsetIdListThatIsNotInSuperset, userSearchesItemDetailList[i].HubSearchedItemId)
			}
		}
		err := fmt.Errorf("error can't find StlUserSearchesItemDetail ids: %s", strings.Join(subsetIdListThatIsNotInSuperset, ", "))
		return userSearchesItemDetailList, createerror.ClientBadRequest(err), err
	}

	return
}

func GetUserSearchesItemListOfUserId(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
) (
	stlUserSearchesItemDetailList []searchedItem.StlUserSearchesItemDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	var lastEvaluatedKey map[string]*dynamodb.AttributeValue = nil
	for {
		companyToSearchedItemMappingRes, err := dynamoDBClient.Query(&dynamodb.QueryInput{
			TableName:              aws.String(searchedItem.GetStlUserSearchesItemDetailDynamoDBTableV1()),
			KeyConditionExpression: aws.String("HubUserId = :pk"),
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":pk": {S: aws.String(userId)},
			},
			ExclusiveStartKey: lastEvaluatedKey,
		})
		if err != nil {
			err = fmt.Errorf("error query companyToSearchedItemMappingRes: %v", err)
			return stlUserSearchesItemDetailList, createerror.InternalException(err), err
		}
		for i := 0; i < len(companyToSearchedItemMappingRes.Items); i++ {
			stlUserSearchesItemDetailDAO := searchedItem.StlUserSearchesItemDetailDAOV1{}
			err = dynamodbattribute.UnmarshalMap(companyToSearchedItemMappingRes.Items[i], &stlUserSearchesItemDetailDAO)
			if err != nil {
				err = fmt.Errorf("error unmarshaling stlUserSearchesItemDetailDAO: %v", err)
				return stlUserSearchesItemDetailList, createerror.InternalException(err), err
			}
			stlUserSearchesItemDetailList = append(stlUserSearchesItemDetailList, stlUserSearchesItemDetailDAO)
		}

		lastEvaluatedKey = companyToSearchedItemMappingRes.LastEvaluatedKey
		if lastEvaluatedKey == nil {
			break
		}
	}

	return stlUserSearchesItemDetailList, nil, nil
}
