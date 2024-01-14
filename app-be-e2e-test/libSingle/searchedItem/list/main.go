package searchedItemList

import (
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"diskon-hunter/price-monitoring-e2e-test/shared/delivery/searchedItem"
	dto "diskon-hunter/price-monitoring-e2e-test/shared/delivery/searchedItem/list"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestDTOV1 = dto.RequestDTOV1

var DefaultRequestObject = GenerateRequestObject(GenerateRequestObjectArgs{})

type GenerateRequestObjectArgs struct {
}

// GenerateRequestObject allow only a few parameter to be customized
// so that you don't have to create the whole
// dto.RequestDTOV1 with all the fields.
func GenerateRequestObject(args GenerateRequestObjectArgs) dto.RequestDTOV1 {
	dto := dto.RequestDTOV1{}
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
func Execute(body dto.RequestDTOV1, jwtToken string) (result ExecuteResult, err error) {
	endpoint := shared.GetBackendUrl() + dto.PathV1
	req, err := shared.CreatePostRequestWithJsonBody(
		shared.CreatePostRequestWithJsonBodyArgs{
			Body:     body,
			Endpoint: endpoint,
			JwtToken: jwtToken,
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

func CheckResultIsCorrect(
	request dto.RequestDTOV1, result ExecuteResult, requesterUserId string,
) (
	searchedItemDAOList []searchedItem.StlUserSearchesItemDetailDAOV1, err error,
) {
	//Should directly get data from database (or other data stores)
	// and check the correctness of the data in that database (or other data stores).
	return
}

func RevertResult(request dto.RequestDTOV1, result ExecuteResult, requesterUserId string) (err error) {
	//Should remove data from database (or other data stores) directly
	// and undo other side effects that can be reverted (e.g. data mutations, etc).
	//Some side effects can't be reverted: sending push notification/email/SMS, etc.
	return
}
