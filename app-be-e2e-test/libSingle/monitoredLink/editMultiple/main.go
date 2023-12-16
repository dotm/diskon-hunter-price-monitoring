package monitoredLinkAddMultiple

import (
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"diskon-hunter/price-monitoring-e2e-test/shared/constenum"
	"diskon-hunter/price-monitoring-e2e-test/shared/currencyutil"
	dto "diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink/editMultiple"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"
)

type RequestDTOV1 = dto.RequestDTOV1
type MonitoredLinkRequestDTOV1 = dto.MonitoredLinkRequestDTOV1

var DefaultRequestObject = GenerateRequestObject(GenerateRequestObjectArgs{
	MonitoredLinkList: []dto.MonitoredLinkRequestDTOV1{
		{HubMonitoredLinkUrl: "https://mock.com/product/1", AlertPrice: currencyutil.NewFromNumberString("150000", "IDR"), ActiveAlertMethodList: []constenum.AlertMethod{constenum.AlertMethodEmail}},
		{HubMonitoredLinkUrl: "https://mock.com/product/2", AlertPrice: currencyutil.NewFromNumberString("140000", "IDR"), ActiveAlertMethodList: []constenum.AlertMethod{}},
	},
})

type GenerateRequestObjectArgs struct {
	MonitoredLinkList []dto.MonitoredLinkRequestDTOV1
}

// GenerateRequestObject allow only a few parameter to be customized
// so that you don't have to create the whole
// dto.RequestDTOV1 with all the fields.
func GenerateRequestObject(args GenerateRequestObjectArgs) dto.RequestDTOV1 {
	dto := dto.RequestDTOV1{
		MonitoredLinkList: args.MonitoredLinkList,
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

	// GetStlMonitoredLinkDetailDynamoDBTableV1
	cleanedMonitoredLinkUrlList := []string{}
	cleanedMonitoredLinkUrlMapToNewDetail := map[string]dto.MonitoredLinkRequestDTOV1{}
	for _, monitoredLinkData := range request.MonitoredLinkList {
		cleanedMonitoredLinkUrlList = append(cleanedMonitoredLinkUrlList, monitoredLinkData.HubMonitoredLinkUrl)
		cleanedMonitoredLinkUrlMapToNewDetail[monitoredLinkData.HubMonitoredLinkUrl] = monitoredLinkData
	}

	userId := requesterUserId
	// GetStlUserMonitorsLinkDetailDynamoDBTableV1
	userMonitorsLinkDAOList, errObj, err := dynamodbhelper.GetStlUserMonitorsLinkDetailList(
		dynamodbhelper.CreateClientFromSession(), userId, cleanedMonitoredLinkUrlList)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	for _, userMonitorsLinkSavedData := range userMonitorsLinkDAOList {
		requestData := cleanedMonitoredLinkUrlMapToNewDetail[userMonitorsLinkSavedData.HubMonitoredLinkUrl]
		if requestData.AlertPrice != userMonitorsLinkSavedData.AlertPrice {
			return fmt.Errorf(
				"__exception__ error updating AlertPrice for %v: request is %v, stored is %v",
				requestData.HubMonitoredLinkUrl, requestData.AlertPrice, userMonitorsLinkSavedData.AlertPrice,
			)
		}
		if len(requestData.ActiveAlertMethodList) != len(userMonitorsLinkSavedData.ActiveAlertMethodList) {
			//TODO: check not only length but elements
			return fmt.Errorf(
				"__exception__ error updating ActiveAlertMethodList for %v: request is %v, stored is %v",
				requestData.HubMonitoredLinkUrl, requestData.ActiveAlertMethodList, userMonitorsLinkSavedData.ActiveAlertMethodList,
			)
		}
	}
	fmt.Printf("__success__ updating GetStlUserMonitorsLinkDetailDynamoDBTableV1:\n%+v\n\n", userMonitorsLinkDAOList)

	return
}

func RevertResult(request dto.RequestDTOV1, result ExecuteResult, requesterUserId string) (err error) {
	//Should remove data from database (or other data stores) directly
	// and undo other side effects that can be reverted (e.g. data mutations, etc).
	//Some side effects can't be reverted: sending push notification/email/SMS, etc.

	//no need to revert for edit because we'll delete the edited data when we revert the add operation
	return nil
}
