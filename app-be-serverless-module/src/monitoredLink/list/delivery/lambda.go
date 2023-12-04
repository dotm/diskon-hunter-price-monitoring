package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/jwttoken"
	"diskon-hunter/price-monitoring/shared/lambdahelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	monitoredLinkList "diskon-hunter/price-monitoring/src/monitoredLink/list/function"
)

/*
DON'T FORGET to create new lambda.go handlers in the app-be-serverless-module/functions directory
to make sure they are deployed by Terraform to AWS Lambda
*/

func LambdaHandlerV1(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	//logging and panic handling needs to be copied
	//for all delivery methods (HTTP server, Serverless Function, etc.)
	var errObj *serverresponse.ErrorObj
	logger := lazylogger.New(req.Path)

	//parse request body
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

	var res monitoredLinkList.QueryV1DataResponse
	if errObj == nil {
		cmd := monitoredLinkList.QueryV1{
			Version:         "1",
			RequesterUserId: jwttoken.GetUserId(token),
		}

		res, errObj = monitoredLinkList.QueryV1Handler(
			ctx,
			monitoredLinkList.QueryV1Dependencies{
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
