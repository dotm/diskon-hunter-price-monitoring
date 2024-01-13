package delivery

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"diskon-hunter/price-monitoring/shared/dynamodbhelper"
	"diskon-hunter/price-monitoring/shared/jwttoken"
	"diskon-hunter/price-monitoring/shared/lambdahelper"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"

	"github.com/aws/aws-lambda-go/events"

	userEdit "diskon-hunter/price-monitoring/src/user/edit/function"
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

	var res userEdit.CommandV1DataResponse
	if errObj == nil {
		cmd := userEdit.NewCommandV1(
			"1",                       //Version
			jwttoken.GetUserId(token), //RequesterUserId
			reqBody.Password,          //Password
			reqBody.WhatsAppNumber,    //WhatsAppNumber
		)

		res, errObj = userEdit.CommandV1Handler(
			ctx,
			userEdit.CommandV1Dependencies{
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
