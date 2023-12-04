package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	monitoredLinkListLambda "diskon-hunter/price-monitoring/src/monitoredLink/list/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(monitoredLinkListLambda.LambdaHandlerV1)
}
