package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userMeLambda "diskon-hunter/price-monitoring/src/user/me/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userMeLambda.LambdaHandlerV1)
}
