package userSignIn

import (
	"context"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/jwttoken"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/password"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/shared/stringmasker"
	"diskon-hunter/price-monitoring/src/user"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version  string `json:"version"` //should follow the struct name suffix
	Email    string `json:"email"`
	Password string `json:"password"`
	// Email    string `json:"email"`
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := CommandV1{
		Version:  x.Version,
		Email:    stringmasker.Email(x.Email),
		Password: stringmasker.Password(x.Password),
	}
	byteSlice, err := json.Marshal(loggableCommand)
	if err != nil {
		return "", err
	} else {
		return string(byteSlice), nil
	}
}

type CommandV1Dependencies struct {
	Logger         *lazylogger.Instance
	DynamoDBClient *dynamodb.DynamoDB
}

type CommandV1DataResponse struct {
	HubUserId      string
	SignedJwtToken string
	JwtExpiration  time.Time
}

/*
Addition, change, or removal of validation might cause version increment
*/
func CommandV1Handler(
	ctx context.Context,
	dependencies CommandV1Dependencies,
	command CommandV1,
) (CommandV1DataResponse, *serverresponse.ErrorObj) {
	//don't mutate this. emptyResponse should be used when returning error.
	emptyResponse := CommandV1DataResponse{}
	//log the command
	loggableCommand, err := command.createLoggableString()
	if err != nil {
		err = fmt.Errorf("error creating loggable string: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	dependencies.Logger.EnqueueCommandLog(loggableCommand, true)

	/* Validations
	Validations from auth, write model,
	or domain model's business logic (from projections or from events replay).
	*/
	if command.Email == "" {
		return emptyResponse, createerror.UserAuthenticationShouldSpecifyEmail()
	}

	existingEmailUserMappingItem, errObj, err := dynamodbhelper.ValidateUserEmailIsRegistered(dependencies.DynamoDBClient, command.Email)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	//check password is correct
	existingEmailUserMappingDAO := user.StlUserEmailAuthenticationDAOV1{}
	err = dynamodbattribute.UnmarshalMap(existingEmailUserMappingItem, &existingEmailUserMappingDAO)
	if err != nil {
		err = fmt.Errorf("error unmarshaling existingEmailUserMappingDAO: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	existingUserResp, err := dependencies.DynamoDBClient.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"HubUserId": {
				S: aws.String(existingEmailUserMappingDAO.HubUserId),
			},
		},
		TableName: aws.String(user.GetStlUserDetailDynamoDBTableV1()),
	})
	if err != nil {
		err = fmt.Errorf("error get user item: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	if existingUserResp.Item == nil {
		err = fmt.Errorf("user not found")
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.UserNotFound(err)
	}
	existingUserDAO := user.StlUserDetailDAOV1{}
	err = dynamodbattribute.UnmarshalMap(existingUserResp.Item, &existingUserDAO)
	if err != nil {
		err = fmt.Errorf("error unmarshaling user response: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	passwordCorrect := password.MatchPasswordToHash(command.Password, existingUserDAO.HashedPassword)
	if !passwordCorrect {
		return emptyResponse, createerror.UserCredentialIncorrect()
	}

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/

	jwtExpiration := time.Now().Add(time.Hour * 24 * 365) //365 days
	signedJwtToken, err := jwttoken.BuildAndSign(
		jwttoken.BuildCustomClaims(existingUserDAO.HubUserId),
		jwtExpiration,
	)
	if err != nil {
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{
		HubUserId:      existingUserDAO.HubUserId,
		SignedJwtToken: signedJwtToken,
		JwtExpiration:  jwtExpiration,
	}, nil
}
