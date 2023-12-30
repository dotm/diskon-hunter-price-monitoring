package userResetPassword

import (
	"diskon-hunter/price-monitoring-e2e-test/shared"
	dto "diskon-hunter/price-monitoring-e2e-test/shared/delivery/user/resetPassword"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestDTOV1 = dto.RequestDTOV1

var DefaultRequestObject = GenerateRequestObject(GenerateRequestObjectArgs{
	Email:    "diskon.hunter.e2e@yopmail.com",
	Password: "Test1234!",
})

type GenerateRequestObjectArgs struct {
	Email    string
	Password string
}

// GenerateRequestObject allow only a few parameter to be customized
// so that you don't have to create the whole
// dto.RequestDTOV1 with all the fields.
func GenerateRequestObject(args GenerateRequestObjectArgs) dto.RequestDTOV1 {
	dto := dto.RequestDTOV1{
		Email:    args.Email,
		Password: args.Password,
	}
	return dto
}

type ResponseObj struct {
	serverresponse.Obj //based of serverresponse.Obj (will include all of its fields)

	Data dto.ResponseDTOV1 `json:"data,omitempty"` //nullable
}

// used to return data for the next steps (e.g. checking result, reverting result, etc.)
type ExecuteResult struct {
	ResponseOk           bool
	ResponseErrorCode    string
	ResponseErrorMessage string
	ResponseData         dto.ResponseDTOV1
}

// Execute the request
func Execute(body dto.RequestDTOV1) (result ExecuteResult, err error) {
	endpoint := shared.GetBackendUrl() + dto.PathV1
	req, err := shared.CreatePostRequestWithJsonBody(
		shared.CreatePostRequestWithJsonBodyArgs{
			Body:     body,
			Endpoint: endpoint,
			JwtToken: shared.JwtToken,
		},
	)
	if err != nil {
		return result, err //error already logged from shared util function
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error executing request: %v", err.Error())
		return result, err
	}
	defer res.Body.Close()

	var resObj ResponseObj
	resBody, err := shared.LogResponseAndReturnBody(res)
	if err != nil {
		return result, err //error already logged from shared util function
	}
	err = json.Unmarshal(resBody, &resObj)
	if err != nil {
		fmt.Printf("error unmarshall response: %v", err.Error())
		return result, err
	}

	result.ResponseOk = resObj.Ok
	if resObj.Err != nil {
		result.ResponseErrorCode = resObj.Err.Code
		if resObj.Err.Message != nil {
			result.ResponseErrorMessage = *resObj.Err.Message
		}
	}
	result.ResponseData = resObj.Data

	return result, nil
}

func CheckResultIsCorrect(request dto.RequestDTOV1, result ExecuteResult) (err error) {
	//Should directly get data from database (or other data stores)
	// and check the correctness of the data in that database (or other data stores).
	email := result.ResponseData.Email
	userEmailHasOtpDetailList, errObj, err := dynamodbhelper.GetUserEmailHasOtpDetailListByEmailList(
		dynamodbhelper.CreateClientFromSession(), []string{email})
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	if len(userEmailHasOtpDetailList) < 1 {
		fmt.Printf("__exception__ otp for email %s not added to database", email)
		return
	}

	fmt.Printf("__success__ adding to user email has otp detail table:\n%+v\n\n", userEmailHasOtpDetailList[0])
	return nil
}

func RevertResult(request dto.RequestDTOV1, result ExecuteResult) (err error) {
	//Should remove data from database (or other data stores) directly
	// and undo other side effects that can be reverted (e.g. data mutations, etc).
	//Some side effects can't be reverted: sending push notification/email/SMS, etc.
	emailList := []string{result.ResponseData.Email}
	errObj, err := dynamodbhelper.DeleteUserOtpByEmailList(dynamodbhelper.CreateClientFromSession(), emailList)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}

	fmt.Printf("__success__ DeleteUserOtpByEmailList:\n%+v\n\n", emailList)
	return
}
