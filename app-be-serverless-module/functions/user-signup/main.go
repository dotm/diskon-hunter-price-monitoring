package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userSignUpLambda "diskon-hunter/price-monitoring/src/user/signup/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userSignUpLambda.LambdaHandlerV1)
}
