package emailutil

import (
	"crypto/rand"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func CreateClientFromSession() *ses.SES {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{},
		SharedConfigState: session.SharedConfigEnable,
	}))
	return ses.New(session)
}

const DefaultEmailSender = "diskon.hunter.official@gmail.com"

type SendEmailArgs struct {
	Sender        string
	RecipientList []*string //slice of aws.String
	CcList        []*string //slice of aws.String
	Subject       string
	HtmlBody      string
	TextBody      string
}

func SendEmail(sesClient *ses.SES, args SendEmailArgs) (errObj *serverresponse.ErrorObj, err error) {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: args.CcList,
			ToAddresses: args.RecipientList,
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(args.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(args.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(args.Subject),
			},
		},
		Source: aws.String(args.Sender),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	_, err = sesClient.SendEmail(input)
	if err != nil {
		err = fmt.Errorf("error sending email %s: %v", args.Subject, err)
		return createerror.InternalException(err), err
	}

	return nil, nil
}

const OtpExpiredInMinutes = 30
const otpChars = "1234567890"

func GenerateOTP() (string, error) {
	length := 6

	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	otpCharsLength := len(otpChars)
	for i := 0; i < length; i++ {
		buffer[i] = otpChars[int(buffer[i])%otpCharsLength]
	}

	return string(buffer), nil
}
