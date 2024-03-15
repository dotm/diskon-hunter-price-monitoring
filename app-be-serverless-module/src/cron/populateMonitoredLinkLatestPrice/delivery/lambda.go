package delivery

import (
	"context"

	"github.com/aws/aws-lambda-go/events"

	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lambdahelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	cronPopulateMonitoredLinkLatestPrice "diskon-hunter/price-monitoring/src/cron/populateMonitoredLinkLatestPrice/function"
)

/*
DON'T FORGET to create new lambda.go handlers in the app-be-serverless-module/functions directory
to make sure they are deployed by Terraform to AWS Lambda
*/

func LambdaHandlerV1(ctx context.Context, req events.CloudWatchEvent) {
	//logging and panic handling needs to be copied
	//for all delivery methods (HTTP server, Serverless Function, etc.)
	var errObj *serverresponse.ErrorObj
	logger := lazylogger.New("cron/populateLatestPrice")

	cmd := cronPopulateMonitoredLinkLatestPrice.NewCommandV1(
		"1", //Version
	)

	_, errObj = cronPopulateMonitoredLinkLatestPrice.CommandV1Handler(
		ctx,
		cronPopulateMonitoredLinkLatestPrice.CommandV1Dependencies{
			Logger:         logger,
			DynamoDBClient: dynamodbhelper.CreateClientFromSession(),
		},
		cmd,
	)

	lambdaResp := lambdahelper.HandleLogAndPanic(logger, errObj)
	if lambdaResp != nil {
		//will only be executed when panic because lambdaResp is only not nil when panic happened.
		return
	}
	return
}
