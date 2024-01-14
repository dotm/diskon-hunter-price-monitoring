package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	searchedItemAddMultipleLambda "diskon-hunter/price-monitoring/src/searchedItem/addMultiple/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(searchedItemAddMultipleLambda.LambdaHandlerV1)
}
