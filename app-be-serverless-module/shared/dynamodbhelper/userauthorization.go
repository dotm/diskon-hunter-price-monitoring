package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/user"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CheckIfUserIsRegistered(dynamoDBClient *dynamodb.DynamoDB, userId string) (bool, *serverresponse.ErrorObj, error) {
	output, err := dynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"UserId": {
				S: aws.String(userId),
			},
		},
		TableName: aws.String(user.GetStlUserDetailDynamoDBTableV1()),
	})
	if err != nil {
		err = fmt.Errorf("error get item CheckIfUserIsRegistered: %v", err)
		return false, createerror.InternalException(err), err
	}
	return output.Item != nil, nil, nil
}
