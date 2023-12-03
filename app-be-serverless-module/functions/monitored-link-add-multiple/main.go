package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	monitoredLinkAddMultipleLambda "diskon-hunter/price-monitoring/src/monitoredLink/addMultiple/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(monitoredLinkAddMultipleLambda.LambdaHandlerV1)
}
