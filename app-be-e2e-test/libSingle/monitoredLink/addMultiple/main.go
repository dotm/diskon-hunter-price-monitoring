package monitoredLinkAddMultiple

import (
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"diskon-hunter/price-monitoring-e2e-test/shared/constenum"
	"diskon-hunter/price-monitoring-e2e-test/shared/currencyutil"
	"diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink"
	dto "diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink/addMultiple"
	"diskon-hunter/price-monitoring-e2e-test/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/maps"
)

type RequestDTOV1 = dto.RequestDTOV1
type MonitoredLinkRequestDTOV1 = dto.MonitoredLinkRequestDTOV1

var DefaultRequestObject = GenerateRequestObject(GenerateRequestObjectArgs{
	MonitoredLinkList: []dto.MonitoredLinkRequestDTOV1{
		{HubMonitoredLinkUrl: "https://mock.com/product/1", AlertPrice: currencyutil.NewFromNumberString("50000", "IDR"), AlertMethodList: []constenum.AlertMethod{constenum.AlertMethodEmail, constenum.AlertMethodPushNotification}},
		{HubMonitoredLinkUrl: "https://mock.com/product/2", AlertPrice: currencyutil.NewFromNumberString("40000", "IDR"), AlertMethodList: []constenum.AlertMethod{}},
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
	monitoredLinkDAOList []monitoredLink.StlMonitoredLinkDetailDAOV1, err error,
) {
	//Should directly get data from database (or other data stores)
	// and check the correctness of the data in that database (or other data stores).

	cleanedMonitoredLinkUrlList := []string{}
	for i := 0; i < len(request.MonitoredLinkList); i++ {
		rawUrl := request.MonitoredLinkList[i].HubMonitoredLinkUrl
		cleanedUrl, ok := result.ResponseData.MonitoredLinkRawToCleanedMap[rawUrl]
		if !ok {
			err = fmt.Errorf("__exception__ error can't find monitoredLink in response: %s", rawUrl)
			fmt.Printf("%+v\n", err)
			return
		}
		cleanedMonitoredLinkUrlList = append(cleanedMonitoredLinkUrlList, cleanedUrl)
	}
	// GetStlMonitoredLinkDetailDynamoDBTableV1
	monitoredLinkDAOList, errObj, err := dynamodbhelper.GetStlMonitoredLinkDetailByUrlList(
		dynamodbhelper.CreateClientFromSession(), cleanedMonitoredLinkUrlList)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	fmt.Printf("__success__ adding to GetStlMonitoredLinkDetailDynamoDBTableV1:\n%+v\n\n", monitoredLinkDAOList)

	userId := requesterUserId
	// GetStlUserMonitorsLinkDetailDynamoDBTableV1
	userMonitorsLinkDAOList, errObj, err := dynamodbhelper.GetStlUserMonitorsLinkDetailList(
		dynamodbhelper.CreateClientFromSession(), userId, cleanedMonitoredLinkUrlList)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	fmt.Printf("__success__ adding to GetStlUserMonitorsLinkDetailDynamoDBTableV1:\n%+v\n\n", userMonitorsLinkDAOList)

	return
}

func RevertResult(request dto.RequestDTOV1, result ExecuteResult, requesterUserId string) (err error) {
	//Should remove data from database (or other data stores) directly
	// and undo other side effects that can be reverted (e.g. data mutations, etc).
	//Some side effects can't be reverted: sending push notification/email/SMS, etc.
	cleanedMonitoredLinkUrlList := maps.Values(result.ResponseData.MonitoredLinkRawToCleanedMap)
	errObj, err := dynamodbhelper.DeleteStlMonitoredLinkDetailByUrlList(
		dynamodbhelper.CreateClientFromSession(),
		cleanedMonitoredLinkUrlList,
	)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	fmt.Printf("__success__ DeleteStlMonitoredLinkDetailByUrlList:\n%+v\n\n", cleanedMonitoredLinkUrlList)

	userId := requesterUserId
	errObj, err = dynamodbhelper.DeleteStlUserMonitorsLinkDetailByUrlList(
		dynamodbhelper.CreateClientFromSession(),
		userId,
		cleanedMonitoredLinkUrlList,
	)
	if errObj != nil || err != nil {
		fmt.Printf("__exception__ errObj: %+v\n", errObj)
		fmt.Printf("__exception__ err: %v\n", err)
		return
	}
	fmt.Printf("__success__ DeleteStlUserMonitorsLinkDetailByUrlList:\n%+v\n%+v\n\n", userId, cleanedMonitoredLinkUrlList)
	return
}
