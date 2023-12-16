package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	monitoredLinkEditMultipleLambda "diskon-hunter/price-monitoring/src/monitoredLink/editMultiple/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(monitoredLinkEditMultipleLambda.LambdaHandlerV1)
}
