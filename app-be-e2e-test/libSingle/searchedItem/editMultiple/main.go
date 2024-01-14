package searchedItemAddMultiple

import (
	"diskon-hunter/price-monitoring-e2e-test/shared"
	dto "diskon-hunter/price-monitoring-e2e-test/shared/delivery/searchedItem/editMultiple"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestDTOV1 = dto.RequestDTOV1
type SearchedItemRequestDTOV1 = dto.SearchedItemRequestDTOV1

type GenerateRequestObjectArgs struct {
	SearchedItemList []dto.SearchedItemRequestDTOV1
}

// GenerateRequestObject allow only a few parameter to be customized
// so that you don't have to create the whole
// dto.RequestDTOV1 with all the fields.
func GenerateRequestObject(args GenerateRequestObjectArgs) dto.RequestDTOV1 {
	dto := dto.RequestDTOV1{
		SearchedItemList: args.SearchedItemList,
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
	err error,
) {
	//Should directly get data from database (or other data stores)
	// and check the correctness of the data in that database (or other data stores).

	// GetStlSearchedItemDetailDynamoDBTableV1
	idList := []string{}
	idMapToNewDetail := map[string]dto.SearchedItemRequestDTOV1{}
	for _, searchedItemData := range request.SearchedItemList {
		idList = append(idList, searchedItemData.HubSearchedItemId)
		idMapToNewDetail[searchedItemData.HubSearchedItemId] = searchedItemData
	}

	userId := requesterUserId
	// GetStlUserSearchesItemDetailDynamoDBTableV1
	userSearchesItemDAOList, errObj, err := dynamodbhelper.GetStlUserSearchesItemDetailList(
		dynamodbhelper.CreateClientFromSession(), userId, idList)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	for _, userSearchesItemSavedData := range userSearchesItemDAOList {
		requestData := idMapToNewDetail[userSearchesItemSavedData.HubSearchedItemId]
		if requestData.AlertPrice != userSearchesItemSavedData.AlertPrice {
			return fmt.Errorf(
				"__exception__ error updating AlertPrice for %v: request is %v, stored is %v",
				requestData.HubSearchedItemId, requestData.AlertPrice, userSearchesItemSavedData.AlertPrice,
			)
		}
	}
	fmt.Printf("__success__ updating GetStlUserSearchesItemDetailDynamoDBTableV1:\n%+v\n\n", userSearchesItemDAOList)

	return
}

func RevertResult(request dto.RequestDTOV1, result ExecuteResult, requesterUserId string) (err error) {
	//Should remove data from database (or other data stores) directly
	// and undo other side effects that can be reverted (e.g. data mutations, etc).
	//Some side effects can't be reverted: sending push notification/email/SMS, etc.

	//no need to revert for edit because we'll delete the edited data when we revert the add operation
	return nil
}
