package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	searchedItemEditMultipleLambda "diskon-hunter/price-monitoring/src/searchedItem/editMultiple/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(searchedItemEditMultipleLambda.LambdaHandlerV1)
}
