package main

import (
	singleRequest "diskon-hunter/price-monitoring-e2e-test/libSingle/user/signin"
	"fmt"
)

func main() {
	executeSingleRequest()
}

func executeSingleRequest() {
	result, err := singleRequest.Execute(singleRequest.GenerateRequestObject(singleRequest.GenerateRequestObjectArgs{
		Email:    "diskon.hunter.e2e@yopmail.com",
		Password: "Test123!",
	}))
	fmt.Println()
	fmt.Printf("result:\n%+v\n\nerror:\n%v\n\n-----------\n", result, err)
}
