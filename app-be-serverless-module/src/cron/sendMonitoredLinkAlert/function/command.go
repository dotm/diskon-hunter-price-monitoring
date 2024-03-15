package cronSendMonitoredLinkAlert

import (
	"context"
	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/emailutil"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"diskon-hunter/price-monitoring/src/monitoredLink"
	"diskon-hunter/price-monitoring/src/user"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	handlebars "github.com/aymerick/raymond"
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

	tableName := user.GetStlUserDetailDynamoDBTableV1()
	scanOutput, err := dependencies.DynamoDBClient.Scan(&dynamodb.ScanInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		err = fmt.Errorf("error scanning %v: %v", tableName, err)
		dependencies.Logger.EnqueueErrorLog(err, true)
		return emptyResponse, createerror.InternalException(err)
	}
	stlUserDetailList := []user.StlUserDetailDAOV1{}
	for i := 0; i < len(scanOutput.Items); i++ {
		stlUserDetailDAO := user.StlUserDetailDAOV1{}
		err = dynamodbattribute.UnmarshalMap(scanOutput.Items[i], &stlUserDetailDAO)
		if err != nil {
			err = fmt.Errorf("error unmarshaling stlUserDetailDAO: %v", err)
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, createerror.InternalException(err)
		}
		stlUserDetailList = append(stlUserDetailList, stlUserDetailDAO)
	}

	emailClient := emailutil.CreateClientFromSession()
	type linkData struct {
		HubMonitoredLinkUrl string
		LatestPriceString   string
		TimeLatestScrapped  string
	}
	type templateDataContext struct {
		LinkDataList []linkData
	}
	htmlTemplate := `<div>
		<p>
			Halo Diskon Hunter! Kami mendeteksi link yang kamu monitor ada dibawah harga yang kamu input:
		</p>
		<ol>
			{{#each linkDataList}}
			<li>
				<ul>
					<li>{{HubMonitoredLinkUrl}}</li>
					<li>
						Harga terakhir: {{LatestPriceString}} (pada {{TimeLatestScrapped}})
					</li>
				</ul>
			</li>
			{{/each}}
		</ol>
		<p>
			Jika kamu tidak ingin menerima email ini lagi, silahkan matikan notifikasi email untuk semua link yang kamu monitor dari website kami.
		</p>
		<p>
			Jika kamu perlu bantuan, silahkan Kontak Kami.
		</p>
	</div>`
	textTemplate := `Halo Diskon Hunter! Kami mendeteksi link yang kamu monitor ada dibawah harga yang kamu input:
		
		{{#each linkDataList}}
			- {{HubMonitoredLinkUrl}}
				- Harga terakhir: {{LatestPriceString}} (pada {{TimeLatestScrapped}})
		{{/each}}

		Jika kamu tidak ingin menerima email ini lagi, silahkan matikan notifikasi email untuk semua link yang kamu monitor dari website kami.
		Jika kamu perlu bantuan, silahkan Kontak Kami.
	`
	for _, stlUserDetail := range stlUserDetailList {
		combinedUserMonitoredLinkDataList, errObj, err := dynamodbhelper.GetCombinedUserMonitoredLinkDataListOfUserId(
			dependencies.DynamoDBClient, stlUserDetail.HubUserId, true,
		)
		if errObj != nil {
			//error already well described on the calling method
			dependencies.Logger.EnqueueErrorLog(err, true)
			return emptyResponse, errObj
		}

		linksAlertedByEmail := []monitoredLink.CombinedUserMonitoredLinkDataV1{}
		for _, combinedUserMonitoredLinkData := range combinedUserMonitoredLinkDataList {
			latestPrice := combinedUserMonitoredLinkData.LatestPrice
			if latestPrice == nil {
				continue
			}
			if latestPrice.IsGreaterThan(combinedUserMonitoredLinkData.AlertPrice) {
				continue
			}
			if len(combinedUserMonitoredLinkData.ActiveAlertMethodList) == 0 {
				continue
			}

			for _, alertMethod := range combinedUserMonitoredLinkData.ActiveAlertMethodList {
				if alertMethod == constenum.AlertMethodEmail {
					linksAlertedByEmail = append(linksAlertedByEmail, combinedUserMonitoredLinkData)
				}
				//add other alert methods here
			}
		}

		if len(linksAlertedByEmail) > 0 {
			templateData := templateDataContext{
				LinkDataList: []linkData{},
			}
			for _, alertedLink := range linksAlertedByEmail {
				templateData.LinkDataList = append(templateData.LinkDataList, linkData{
					HubMonitoredLinkUrl: alertedLink.HubMonitoredLinkUrl,
					LatestPriceString:   alertedLink.LatestPrice.ToDisplayString(),
					TimeLatestScrapped:  alertedLink.TimeLatestScrapped.Format("2006-01-02"),
				})
			}

			htmlBody, err := handlebars.Render(htmlTemplate, templateData)
			if err != nil {
				err = fmt.Errorf("error rendering html template: %v", err)
				dependencies.Logger.EnqueueErrorLog(err, true)
				return emptyResponse, createerror.InternalException(err)
			}
			textBody, err := handlebars.Render(textTemplate, templateData)
			if err != nil {
				err = fmt.Errorf("error rendering text template: %v", err)
				dependencies.Logger.EnqueueErrorLog(err, true)
				return emptyResponse, createerror.InternalException(err)
			}

			errObj, err := emailutil.SendEmail(
				emailClient,
				emailutil.SendEmailArgs{
					Sender:        emailutil.DefaultEmailSender,
					RecipientList: []*string{aws.String(stlUserDetail.Email)},
					CcList:        []*string{},
					Subject:       "Diskon Terdeteksi!",
					HtmlBody:      htmlBody,
					TextBody:      textBody,
				},
			)
			if errObj != nil {
				//error already well described on the calling method
				dependencies.Logger.EnqueueErrorLog(err, true)
				return emptyResponse, errObj
			}
		}
	}

	//You can send the event id back to the requester
	//so that they can periodically check the status of the event.
	return CommandV1DataResponse{}, nil
}
