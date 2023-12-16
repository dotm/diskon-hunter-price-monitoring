package main

import (
	singleRequest "diskon-hunter/price-monitoring-e2e-test/libSingle/monitoredLink/editMultiple"
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"fmt"
)

func main() {
	executeSingleRequest()
}

func executeSingleRequest() {
	result, err := singleRequest.Execute(singleRequest.DefaultRequestObject, shared.JwtToken)
	fmt.Println()
	fmt.Printf("result:\n%+v\n\nerror:\n%v\n\n-----------\n", result, err)
}
