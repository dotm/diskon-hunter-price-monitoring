package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	cronSendMonitoredLinkAlert "diskon-hunter/price-monitoring/src/cron/sendMonitoredLinkAlert/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(cronSendMonitoredLinkAlert.LambdaHandlerV1)
}
