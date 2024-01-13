package userEdit

import (
	"context"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/password"
	"diskon-hunter/price-monitoring/shared/phoneutil"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/shared/stringmasker"
	"diskon-hunter/price-monitoring/src/user"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version         string
	RequesterUserId string
	Password        string
	WhatsAppNumber  string
}

func NewCommandV1(
	Version string, //should follow the struct name suffix
	RequesterUserId string,
	Password string,
	WhatsAppNumber string,
) CommandV1 {
	return CommandV1{
		Version:         Version,
		RequesterUserId: RequesterUserId,
		Password:        Password,
		WhatsAppNumber:  WhatsAppNumber,
	}
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := CommandV1{
		Version:        x.Version,
		Password:       stringmasker.Password(x.Password),
		WhatsAppNumber: stringmasker.Mobile(x.WhatsAppNumber),
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

type CommandV1DataResponse = user.StlUserDetailDAOV1

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

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/

	existingUserDAO, errObj, err := dynamodbhelper.GetUserById(
		dependencies.DynamoDBClient, command.RequesterUserId,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	hashedPassword := existingUserDAO.HashedPassword
	if command.Password != "" {
		hashedPassword, err = password.Hash(command.Password)
		if err != nil {
			err = fmt.Errorf("error hashing password: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
	}

	standardizedWhatsAppNumber := ""
	if command.WhatsAppNumber != "" && command.WhatsAppNumber != "+62" {
		standardizedWhatsAppNumber = phoneutil.StandardizePhoneNumberPrefix(command.WhatsAppNumber)
	}
	newUser := user.StlUserDetailDAOV1{
		HubUserId:      existingUserDAO.HubUserId,
		Email:          existingUserDAO.Email, //email is NOT editable for security and profit concern
		HashedPassword: hashedPassword,
		WhatsAppNumber: standardizedWhatsAppNumber,
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	transactItems, errObj, err := dynamodbhelper.CreateTransactionItemsForEditUser(newUser)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	errObj, err = dynamodbhelper.TransactWriteItemsInWaves(dependencies.DynamoDBClient, transactItems)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	return newUser, nil
}
