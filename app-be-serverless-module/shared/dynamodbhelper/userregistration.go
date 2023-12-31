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
	existingEmailUserMappingDAO user.StlUserEmailAuthenticationDAOV1,
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
		return existingEmailUserMappingDAO, createerror.InternalException(err), err
	}
	if existingEmailUserMapping.Item == nil {
		err = fmt.Errorf("email has not been registered")
		return existingEmailUserMappingDAO, createerror.UserEmailNotRegistered(), err
	}

	err = dynamodbattribute.UnmarshalMap(existingEmailUserMapping.Item, &existingEmailUserMappingDAO)
	if err != nil {
		err = fmt.Errorf("error unmarshaling existingEmailUserMappingDAO: %v", err)
		return existingEmailUserMappingDAO, createerror.InternalException(err), err
	}
	return existingEmailUserMappingDAO, nil, nil
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
