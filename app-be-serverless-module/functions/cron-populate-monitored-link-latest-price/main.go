package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	cronPopulateMonitoredLinkLatestPrice "diskon-hunter/price-monitoring/src/cron/populateMonitoredLinkLatestPrice/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(cronPopulateMonitoredLinkLatestPrice.LambdaHandlerV1)
}
