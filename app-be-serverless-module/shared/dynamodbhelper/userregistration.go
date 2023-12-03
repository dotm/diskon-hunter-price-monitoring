package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/user"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func ValidateUserEmailIsRegistered(
	dynamoDBClient *dynamodb.DynamoDB,
	email string,
) (
	existingEmailUserMappingItem map[string]*dynamodb.AttributeValue,
	errObj *serverresponse.ErrorObj,
	err error,
) {
	existingEmailUserMapping, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:      aws.String(user.GetStlUserEmailAuthenticationDynamoDBTableV1()),
		ConsistentRead: aws.Bool(true),
	})
	if err != nil {
		err = fmt.Errorf("error get email user mapping item: %v", err)
		return existingEmailUserMappingItem, createerror.InternalException(err), err
	}
	if existingEmailUserMapping.Item == nil {
		err = fmt.Errorf("email already registered")
		return existingEmailUserMappingItem, createerror.UserEmailNotRegistered(), err
	}

	existingEmailUserMappingItem = existingEmailUserMapping.Item
	return existingEmailUserMappingItem, nil, nil
}

func ValidateUserEmailHasNotBeenRegistered(dynamoDBClient *dynamodb.DynamoDB, email string) (*serverresponse.ErrorObj, error) {
	//validate existing email should not be registered twice
	getItemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName:      aws.String(user.GetStlUserEmailAuthenticationDynamoDBTableV1()),
		ConsistentRead: aws.Bool(true),
	}
	existingEmailUserMapping, err := dynamoDBClient.GetItem(getItemInput)
	if err != nil {
		err = fmt.Errorf("error get existingEmailUserMapping: %v", err)
		return createerror.InternalException(err), err
	}
	if existingEmailUserMapping.Item != nil {
		err = fmt.Errorf("email already registered")
		return createerror.UserEmailAlreadyRegistered(), err
	}

	return nil, nil
}

func CreateTransactionItemsForUserRegistration(userDAO user.StlUserDetailDAOV1) ([]*dynamodb.TransactWriteItem, *serverresponse.ErrorObj, error) {
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
