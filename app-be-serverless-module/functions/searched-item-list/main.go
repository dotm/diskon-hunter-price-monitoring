package main

import (
	"github.com/aws/aws-lambda-go/lambda"

	searchedItemListLambda "diskon-hunter/price-monitoring/src/searchedItem/list/delivery"
)

// All Lambda handler should be put inside main.go
func main() {
	lambda.Start(searchedItemListLambda.LambdaHandlerV1)
}
