package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/user"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUserById(
	dynamoDBClient *dynamodb.DynamoDB,
	userId string,
) (
	userDAO user.StlUserDetailDAOV1,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	userList, errObj, err := GetUserListByFilter(dynamoDBClient, []string{userId})
	if errObj != nil || err != nil {
		return userDAO, errObj, err
	}
	//user not found error is handled in GetUserListByFilter
	return userList[0], nil, nil
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

func CreateTransactionItemsForEditUser(userDAO user.StlUserDetailDAOV1) ([]*dynamodb.TransactWriteItem, *serverresponse.ErrorObj, error) {
	//don't mutate this. emptyTransaction should be used when returning error.
	emptyTransaction := []*dynamodb.TransactWriteItem{}

	userDAOItem, err := dynamodbattribute.MarshalMap(userDAO)
	if err != nil {
		err = fmt.Errorf("error marshaling userDAO: %v", err)
		return emptyTransaction, createerror.InternalException(err), err
	}
	userEmailAuthenticationDAOItem, err := dynamodbattribute.MarshalMap(user.StlUserEmailAuthenticationDAOV1{
		Email:          userDAO.Email,
		HubUserId:      userDAO.HubUserId,
		HashedPassword: userDAO.HashedPassword,
	})
	if err != nil {
		err = fmt.Errorf("error marshaling userEmailAuthenticationDAO: %v", err)
		return emptyTransaction, createerror.InternalException(err), err
	}

	return []*dynamodb.TransactWriteItem{
		{
			Put: &dynamodb.Put{
				Item:      userDAOItem,
				TableName: aws.String(user.GetStlUserDetailDynamoDBTableV1()),
			},
		},
		{
			Put: &dynamodb.Put{
				Item:      userEmailAuthenticationDAOItem,
				TableName: aws.String(user.GetStlUserEmailAuthenticationDynamoDBTableV1()),
			},
		},
	}, nil, nil
}
