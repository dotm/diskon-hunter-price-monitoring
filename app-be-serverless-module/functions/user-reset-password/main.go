package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userResetPasswordLambda "diskon-hunter/price-monitoring/src/user/resetPassword/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userResetPasswordLambda.LambdaHandlerV1)
}
