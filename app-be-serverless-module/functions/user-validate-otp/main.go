package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	userValidateOtpLambda "diskon-hunter/price-monitoring/src/user/validateOtp/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(userValidateOtpLambda.LambdaHandlerV1)
}
