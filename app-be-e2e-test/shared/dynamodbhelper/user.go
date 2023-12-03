package dynamodbhelper

import (
	"diskon-hunter/price-monitoring-e2e-test/shared/createerror"
	"diskon-hunter/price-monitoring-e2e-test/shared/delivery/user"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"

	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func DeleteUserListByFilter(
	dynamoDBClient *dynamodb.DynamoDB,
	userIdList []string,
) (
	*serverresponse.ErrorObj,
	error,
) {
	if len(userIdList) == 0 {
		return nil, nil
	}
	userList, errObj, err := GetUserListByFilter(dynamoDBClient, userIdList)
	if err != nil {
		return errObj, err
	}

	for _, userData := range userList {
		var tableName string
		var err error

		tableName = user.GetStlUserDetailDynamoDBTableV1()
		_, err = dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"HubUserId": {
					S: aws.String(userData.HubUserId),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			err = fmt.Errorf("error DeleteItem from %s: %v", tableName, err)
			return createerror.InternalException(err), err
		}

		tableName = user.GetStlUserEmailAuthenticationDynamoDBTableV1()
		_, err = dynamoDBClient.DeleteItem(&dynamodb.DeleteItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"Email": {
					S: aws.String(userData.Email),
				},
			},
			TableName: aws.String(tableName),
		})
		if err != nil {
			err = fmt.Errorf("error DeleteItem from %s: %v", tableName, err)
			return createerror.InternalException(err), err
		}
		//other authentication method should be added below
	}

	return nil, nil
}

func GetUserListByFilter(
	dynamoDBClient *dynamodb.DynamoDB,
	userIdList []string,
) (
	[]user.StlUserDetailDAOV1,
	*serverresponse.ErrorObj,
	error,
) {
	userList := []user.StlUserDetailDAOV1{}
	if len(userIdList) == 0 {
		return userList, nil, nil
	}

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(userIdList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"HubUserId": {
				S: aws.String(userIdList[i]),
			},
		})
	}
	tableName := user.GetStlUserDetailDynamoDBTableV1()
	batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return userList, createerror.InternalException(err), err
	}
	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		userEmailAuthenticationDAO := user.StlUserDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&userEmailAuthenticationDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling userEmailAuthenticationDAO: %v", err)
			return userList, createerror.InternalException(err), err
		}
		userList = append(userList, userEmailAuthenticationDAO)
	}
	if len(userList) < len(userIdList) {
		subsetIdListThatIsNotInSuperset := []string{}
		supersetIdMap := map[string]bool{}
		for i := 0; i < len(userIdList); i++ {
			supersetIdMap[userIdList[i]] = true
		}
		for i := 0; i < len(userList); i++ {
			exist, ok := supersetIdMap[userList[i].HubUserId]
			if !ok || !exist {
				subsetIdListThatIsNotInSuperset = append(subsetIdListThatIsNotInSuperset, userList[i].HubUserId)
			}
		}
		err := fmt.Errorf("error can't find user of ids: %s", strings.Join(subsetIdListThatIsNotInSuperset, ", "))
		return userList, createerror.ClientBadRequest(err), err
	}

	return userList, nil, nil
}

func GetUserEmailAuthenticationListByEmailList(
	dynamoDBClient *dynamodb.DynamoDB,
	emailList []string,
) (
	[]user.StlUserEmailAuthenticationDAOV1,
	*serverresponse.ErrorObj,
	error,
) {
	userEmailAuthenticationList := []user.StlUserEmailAuthenticationDAOV1{}
	if len(emailList) == 0 {
		return userEmailAuthenticationList, nil, nil
	}

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(emailList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(emailList[i]),
			},
		})
	}
	tableName := user.GetStlUserEmailAuthenticationDynamoDBTableV1()
	batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return userEmailAuthenticationList, createerror.InternalException(err), err
	}
	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		userEmailAuthenticationDAO := user.StlUserEmailAuthenticationDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&userEmailAuthenticationDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling userEmailAuthenticationDAO: %v", err)
			return userEmailAuthenticationList, createerror.InternalException(err), err
		}
		userEmailAuthenticationList = append(userEmailAuthenticationList, userEmailAuthenticationDAO)
	}
	// if len(userEmailAuthenticationList) < len(encryptedEmailList) {
	// 	//for backend code, need to validate all encrypted email is equal in subset and superset
	// }

	return userEmailAuthenticationList, nil, nil
}
