package cronPopulateMonitoredLinkLatestPrice

import (
	"context"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/emailutil"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version string //should follow the struct name suffix
}

func NewCommandV1(
	Version string, //should follow the struct name suffix
) CommandV1 {
	return CommandV1{
		Version: Version,
	}
}

func (x CommandV1) createLoggableString() (string, error) {
	//strip any sensitive information.
	//strip any fields that are too large to be printed (e.g. image blob).
	loggableCommand := x //no sensitive info and no large fields so we'll just use x
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

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/
	errObj, err := emailutil.SendEmail(
		emailutil.CreateClientFromSession(),
		emailutil.SendEmailArgs{
			Sender:        emailutil.DefaultEmailSender,
			RecipientList: []*string{aws.String("diskon.hunter.e2e.test@yopmail.com")},
			CcList:        []*string{},
			Subject:       "kodok",
			HtmlBody:      "test",
			TextBody:      "test",
		},
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{}, nil
}
