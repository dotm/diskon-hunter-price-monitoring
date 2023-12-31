package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userEditLambda "diskon-hunter/price-monitoring/src/user/edit/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userEditLambda.LambdaHandlerV1)
}
