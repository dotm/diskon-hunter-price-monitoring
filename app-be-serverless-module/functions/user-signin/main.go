package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userSignInLambda "diskon-hunter/price-monitoring/src/user/signin/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userSignInLambda.LambdaHandlerV1)
}
