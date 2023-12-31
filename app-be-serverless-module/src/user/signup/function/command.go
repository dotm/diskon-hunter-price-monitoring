package userSignUp

import (
	"context"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/emailutil"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/password"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/shared/stringmasker"
	"diskon-hunter/price-monitoring/src/user"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/google/uuid"
)

/*
Commands represent input from client through API requests.
Addition, change, or removal of struct fields might cause version increment
*/
type CommandV1 struct {
	Version  string `json:"version"` //should follow the struct name suffix
	Email    string `json:"email"`
	Password string `json:"password"`
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
	HubUserId string
	Email     string
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

	errObj, err := dynamodbhelper.ValidateUserEmailHasNotBeenRegistered(
		dependencies.DynamoDBClient,
		command.Email,
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	/* Business Logic
	Perform business logic preferably through domain model's methods.
	*/
	hashedPassword, err := password.Hash(command.Password)
	if err != nil {
		err = fmt.Errorf("error hashing password: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	newUser := user.StlUserDetailDAOV1{
		HubUserId:      uuid.NewString(),
		Email:          command.Email,
		HashedPassword: hashedPassword,
	}
	otp, err := emailutil.GenerateOTP()
	if err != nil {
		err = fmt.Errorf("error generating otp: %v", err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}

	/* Persisting Data
	Persist event to event store.
	If write model is used, also persist write model with atomic transaction.
	*/

	errObj, err = emailutil.SendEmail(
		emailutil.CreateClientFromSession(),
		emailutil.SendEmailArgs{
			Sender:        emailutil.DefaultEmailSender,
			RecipientList: []*string{aws.String(newUser.Email)},
			CcList:        []*string{},
			Subject:       "Your OTP for Diskon Hunter",
			HtmlBody:      fmt.Sprintf("<p>OTP anda adalah %s (berlaku %v menit). Mohon tidak membagikannya kepada <strong>siapapun termasuk admin Diskon Hunter</strong>.</p><p>Your OTP is %s (will expire in %v minutes). Please don't share it to <strong>anyone, even to Diskon Hunter's admin</strong>.</p>", otp, emailutil.OtpExpiredInMinutes, otp, emailutil.OtpExpiredInMinutes),
			TextBody:      fmt.Sprintf("OTP anda adalah %s (berlaku %v menit). Mohon tidak membagikannya kepada siapapun termasuk admin Diskon Hunter. \nYour OTP is %s (will expire in %v minutes). Please don't share it to anyone, even to Diskon Hunter's admin.", otp, emailutil.OtpExpiredInMinutes, otp, emailutil.OtpExpiredInMinutes),
		},
	)
	if errObj != nil {
		//error already well described on the calling method
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, errObj
	}

	transactItems, errObj, err := dynamodbhelper.CreateTransactionItemsForUserCreateOtp(newUser, otp)
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

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{
		HubUserId: newUser.HubUserId,
		Email:     newUser.Email,
	}, nil
}
