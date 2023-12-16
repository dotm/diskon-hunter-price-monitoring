package dynamodbhelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/emailutil"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/user"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func GetUserEmailHasOtpDetailByEmail(
	dynamoDBClient *dynamodb.DynamoDB,
	email string,
) (
	user.StlUserEmailHasOtpDetailDAOV1,
	*serverresponse.ErrorObj,
	error,
) {
	userEmailHasOtpDetailList, errObj, err := GetUserEmailHasOtpDetailListByEmailList(dynamoDBClient, []string{email})
	if errObj != nil || err != nil || len(userEmailHasOtpDetailList) == 0 {
		return user.StlUserEmailHasOtpDetailDAOV1{}, errObj, err
	}
	return userEmailHasOtpDetailList[0], nil, nil
}

func GetUserEmailHasOtpDetailListByEmailList(
	dynamoDBClient *dynamodb.DynamoDB,
	emailList []string,
) (
	[]user.StlUserEmailHasOtpDetailDAOV1,
	*serverresponse.ErrorObj,
	error,
) {
	userEmailHasOtpDetailList := []user.StlUserEmailHasOtpDetailDAOV1{}
	if len(emailList) == 0 {
		return userEmailHasOtpDetailList, nil, nil
	}

	batchGetItemKeys := []map[string]*dynamodb.AttributeValue{}
	for i := 0; i < len(emailList); i++ {
		batchGetItemKeys = append(batchGetItemKeys, map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(emailList[i]),
			},
		})
	}
	tableName := user.GetStlUserEmailHasOtpDetailDynamoDBTableV1()
	batchGetItemOutput, err := dynamoDBClient.BatchGetItem(&dynamodb.BatchGetItemInput{
		RequestItems: map[string]*dynamodb.KeysAndAttributes{
			tableName: {
				Keys: batchGetItemKeys,
			},
		},
	})
	if err != nil {
		err = fmt.Errorf("error batchGetItemOutput from %s: %v", tableName, err)
		return userEmailHasOtpDetailList, createerror.InternalException(err), err
	}
	for i := 0; i < len(batchGetItemOutput.Responses[tableName]); i++ {
		userEmailHasOtpDetailDAO := user.StlUserEmailHasOtpDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(
			batchGetItemOutput.Responses[tableName][i],
			&userEmailHasOtpDetailDAO,
		)
		if err != nil {
			err = fmt.Errorf("error unmarshaling userEmailHasOtpDetailDAO: %v", err)
			return userEmailHasOtpDetailList, createerror.InternalException(err), err
		}
		userEmailHasOtpDetailList = append(userEmailHasOtpDetailList, userEmailHasOtpDetailDAO)
	}
	// if len(userEmailHasOtpDetailList) < len(encryptedEmailList) {
	// 	//for backend code, need to validate all encrypted email is equal in subset and superset
	// }

	return userEmailHasOtpDetailList, nil, nil
}

func CreateTransactionItemsForUserCreateOtp(userDAO user.StlUserDetailDAOV1, otp string) ([]*dynamodb.TransactWriteItem, *serverresponse.ErrorObj, error) {
	//don't mutate this. emptyTransaction should be used when returning error.
	emptyTransaction := []*dynamodb.TransactWriteItem{}

	userEmailHasOtpDetailDAOItem, err := dynamodbattribute.MarshalMap(user.StlUserEmailHasOtpDetailDAOV1{
		Email:          userDAO.Email,
		HubUserId:      userDAO.HubUserId,
		HashedPassword: userDAO.HashedPassword,
		OTP:            otp,
		TimeExpired:    time.Now().Add(emailutil.OtpExpiredInMinutes * time.Minute),
	})
	if err != nil {
		err = fmt.Errorf("error marshaling userEmailHasOtpDetailDAO: %v", err)
		return emptyTransaction, createerror.InternalException(err), err
	}

	return []*dynamodb.TransactWriteItem{
		{
			Put: &dynamodb.Put{
				Item:      userEmailHasOtpDetailDAOItem,
				TableName: aws.String(user.GetStlUserEmailHasOtpDetailDynamoDBTableV1()),
			},
		},
	}, nil, nil
}

func CreateTransactionItemsForUserValidateOTP(userDAO user.StlUserDetailDAOV1) ([]*dynamodb.TransactWriteItem, *serverresponse.ErrorObj, error) {
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
		{
			Delete: &dynamodb.Delete{
				Key: map[string]*dynamodb.AttributeValue{
					"Email": {
						S: aws.String(userDAO.Email),
					},
				},
				TableName: aws.String(user.GetStlUserEmailHasOtpDetailDynamoDBTableV1()),
			},
		},
	}, nil, nil
}
