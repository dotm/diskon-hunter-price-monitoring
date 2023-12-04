package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

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
