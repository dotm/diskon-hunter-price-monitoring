package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"fmt"

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

// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_BatchGetItem.html
// If you request more than 100 items,
// BatchGetItem returns a ValidationException with the message
// "Too many items requested for the BatchGetItem call."
const MaxBatchGetItem = 100 //items

// https://docs.aws.amazon.com/amazondynamodb/latest/APIReference/API_TransactWriteItems.html
// TransactWriteItems is a synchronous write operation that groups up to 100 action requests.
const MaxTransactWriteItems = 100 //items

func TransactWriteItemsInWaves(
	dynamoDBClient *dynamodb.DynamoDB,
	transactItems []*dynamodb.TransactWriteItem,
) (
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(transactItems) <= MaxTransactWriteItems {
		input := &dynamodb.TransactWriteItemsInput{
			TransactItems: transactItems,
		}
		_, err = dynamoDBClient.TransactWriteItems(input)
		if err != nil {
			err = fmt.Errorf("error TransactWriteItems: %v", err)
			return createerror.InternalException(err), err
		}
	} else {
		//If we use this pattern, make sure the operations is idempotent:
		// no uuid generation, no using time.Now() unless it's not a concern
		// (DynamoDB Put, Update, and Delete operations are already idempotent if the items' attributes are the same).
		//So that in case of failure, we can retry it again without worrying about multiple entry.
		// (multiple entry: adding 10 items -> 5 succeed, 5 failed -> retry -> end up with 15 items instead of 10)
		lastIndexOfLastWave := -1
		lastIndexOfTransactItems := len(transactItems) - 1
		for {
			startIndexOfThisWave := lastIndexOfLastWave + 1
			lastIndexOfThisWave := lastIndexOfLastWave + MaxTransactWriteItems
			isLastWave := false
			if lastIndexOfThisWave >= lastIndexOfTransactItems {
				lastIndexOfThisWave = lastIndexOfTransactItems
				isLastWave = true
			}

			input := &dynamodb.TransactWriteItemsInput{
				TransactItems: transactItems[startIndexOfThisWave : lastIndexOfThisWave+1],
				//+1 because in [start:end], the start is inclusive but the end is exclusive.
			}
			_, err = dynamoDBClient.TransactWriteItems(input)
			if err != nil {
				err = fmt.Errorf("error TransactWriteItems: %v", err)
				return createerror.InternalException(err), err
			}

			if isLastWave {
				break
			} else {
				lastIndexOfLastWave = lastIndexOfThisWave
			}
		}
	}
	return nil, nil
}

func BatchGetItemInWaves(
	dynamoDBClient *dynamodb.DynamoDB,
	tableName string,
	batchGetItemKeys []map[string]*dynamodb.AttributeValue,
	parseResponseToDAO func(response []map[string]*dynamodb.AttributeValue) (
		errObj *serverresponse.ErrorObj,
		err error,
	),
) (
	errObj *serverresponse.ErrorObj,
	err error,
) {
	if len(batchGetItemKeys) <= MaxBatchGetItem {
		batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
			RequestItems: map[string]*dynamodb.KeysAndAttributes{
				tableName: {
					Keys: batchGetItemKeys,
				},
			},
		})
		if err != nil {
			err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
			return createerror.InternalException(err), err
		}

		errObj, err := parseResponseToDAO(batchGetItemOutput.Responses[tableName])
		if errObj != nil || err != nil {
			return errObj, err
		}
	} else {
		//If we use this pattern, make sure the operations is idempotent:
		// no uuid generation, no using time.Now() unless it's not a concern
		// (DynamoDB Put, Update, and Delete operations are already idempotent if the items' attributes are the same).
		//So that in case of failure, we can retry it again without worrying about multiple entry.
		// (multiple entry: adding 10 items -> 5 succeed, 5 failed -> retry -> end up with 15 items instead of 10)
		lastIndexOfLastWave := -1
		lastIndexOfAll := len(batchGetItemKeys) - 1
		for {
			startIndexOfThisWave := lastIndexOfLastWave + 1
			lastIndexOfThisWave := lastIndexOfLastWave + MaxBatchGetItem
			isLastWave := false
			if lastIndexOfThisWave >= lastIndexOfAll {
				lastIndexOfThisWave = lastIndexOfAll
				isLastWave = true
			}

			batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
				RequestItems: map[string]*dynamodb.KeysAndAttributes{
					tableName: {
						Keys: batchGetItemKeys[startIndexOfThisWave : lastIndexOfThisWave+1],
						//+1 because in [start:end], the start is inclusive but the end is exclusive.
					},
				},
			})
			if err != nil {
				err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
				return createerror.InternalException(err), err
			}

			errObj, err := parseResponseToDAO(batchGetItemOutput.Responses[tableName])
			if errObj != nil || err != nil {
				return errObj, err
			}

			if isLastWave {
				break
			} else {
				lastIndexOfLastWave = lastIndexOfThisWave
			}
		}
	}

	return nil, nil
}
