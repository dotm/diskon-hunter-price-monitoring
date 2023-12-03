package dynamodbhelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func CreateClientFromSession() *dynamodb.DynamoDB {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return dynamodb.New(session)
}
