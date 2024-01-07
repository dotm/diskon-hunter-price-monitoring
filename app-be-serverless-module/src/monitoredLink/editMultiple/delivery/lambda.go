package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"diskon-hunter/price-monitoring/shared/constenum"
	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/jwttoken"
	"diskon-hunter/price-monitoring/shared/lambdahelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"

	"github.com/aws/aws-lambda-go/events"

	monitoredLinkEditMultiple "diskon-hunter/price-monitoring/src/monitoredLink/editMultiple/function"
)

func LambdaHandlerV1(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	logger := lazylogger.New(req.Path)

	var reqBody RequestDTOV1
	err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&reqBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	//parse jwt
	authorizationHeaders := req.Headers["authorization"] //header key is auto-lowercased
	token, errObj := jwttoken.ParseFromAuthorizationHeader(authorizationHeaders)

	var res monitoredLinkEditMultiple.CommandV1DataResponse
	if errObj == nil {
		monitoredLinkList := []monitoredLinkEditMultiple.MonitoredLinkDetailCommandV1{}
		for i := 0; i < len(reqBody.MonitoredLinkList); i++ {
			activeAlertMethodList := []constenum.AlertMethod{}
			for j := 0; j < len(reqBody.MonitoredLinkList[i].ActiveAlertMethodList); j++ {
				activeAlertMethod := constenum.NewAlertMethod(reqBody.MonitoredLinkList[i].ActiveAlertMethodList[j])
				if activeAlertMethod != constenum.UnknownAlertMethod {
					activeAlertMethodList = append(activeAlertMethodList, activeAlertMethod)
				}
			}
			monitoredLinkList = append(monitoredLinkList, monitoredLinkEditMultiple.NewMonitoredLinkDetailCommandV1(
				reqBody.MonitoredLinkList[i].HubMonitoredLinkUrl, //HubMonitoredLinkUrl
				activeAlertMethodList,                            //ActiveAlertMethodList
				reqBody.MonitoredLinkList[i].AlertPrice,          //AlertPrice
			))
		}
		cmd := monitoredLinkEditMultiple.NewCommandV1(
			"1",                       //Version
			jwttoken.GetUserId(token), //RequesterUserId
			monitoredLinkList,         //MonitoredLinkList
		)

		res, errObj = monitoredLinkEditMultiple.CommandV1Handler(
			ctx,
			monitoredLinkEditMultiple.CommandV1Dependencies{
				Logger:         logger,
				DynamoDBClient: dynamodbhelper.CreateClientFromSession(),
			},
			cmd,
		)
	}

	var resObj serverresponse.Obj
	if errObj != nil {
		resObj.Ok = false
		resObj.Err = errObj
	} else {
		resObj.Ok = true
		resObj.Data = res
	}

	lambdaResp := lambdahelper.HandleLogAndPanic(logger, errObj)
	if lambdaResp != nil {
		//will only be executed when panic because lambdaResp is only not nil when panic happened.
		return *lambdaResp, nil
	}

	return lambdahelper.WriteResponseFn(resObj, ""), nil

}
