package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/lambdahelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	userResetPassword "diskon-hunter/price-monitoring/src/user/resetPassword/function"
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

	var reqBody RequestDTOV1

	err := json.NewDecoder(strings.NewReader(req.Body)).Decode(&reqBody)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	var res userResetPassword.CommandV1DataResponse
	if errObj == nil {
		cmd := userResetPassword.CommandV1{
			Version:  "1",
			Email:    reqBody.Email,
			Password: reqBody.Password,
		}

		res, errObj = userResetPassword.CommandV1Handler(
			ctx,
			userResetPassword.CommandV1Dependencies{
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
		resDTO := ResponseDTOV1{
			Email: res.Email,
		}
		resObj.Ok = true
		resObj.Data = resDTO
	}

	lambdaResp := lambdahelper.HandleLogAndPanic(logger, errObj)
	if lambdaResp != nil {
		//will only be executed when panic because lambdaResp is only not nil when panic happened.
		return *lambdaResp, nil
	}
	return lambdahelper.WriteResponseFn(resObj, ""), nil
}
