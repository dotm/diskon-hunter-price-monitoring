package main

import (
	singleRequest "diskon-hunter/price-monitoring-e2e-test/libSingle/searchedItem/list"
	signInRequest "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signIn"
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"fmt"
)

func main() {
	// executeSignIn()
	executeSingleRequest()
}

func executeSingleRequest() {
	result, err := singleRequest.Execute(singleRequest.DefaultRequestObject, shared.JwtToken)
	fmt.Println()
	fmt.Printf("result:\n%+v\n\nerror:\n%v\n\n-----------\n", result, err)
}

func executeSignIn() {
	result, err := signInRequest.Execute(signInRequest.DefaultRequestObject)
	fmt.Println()
	fmt.Printf("result:\n%+v\n\nerror:\n%v\n\n-----------\n", result, err)
}
