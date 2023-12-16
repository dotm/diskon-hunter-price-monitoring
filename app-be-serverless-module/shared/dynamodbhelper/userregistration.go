package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/user"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
		err = fmt.Errorf("email has not been registered")
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
