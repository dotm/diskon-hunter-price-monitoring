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
	userSignIn "diskon-hunter/price-monitoring/src/user/signin/function"
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

	var res userSignIn.CommandV1DataResponse
	if errObj == nil {
		cmd := userSignIn.CommandV1{
			Version:  "1",
			Email:    reqBody.Email,
			Password: reqBody.Password,
		}

		res, errObj = userSignIn.CommandV1Handler(
			ctx,
			userSignIn.CommandV1Dependencies{
				Logger:         logger,
				DynamoDBClient: dynamodbhelper.CreateClientFromSession(),
			},
			cmd,
		)
	}

	var resObj serverresponse.Obj
	var jwtCookieString string
	if errObj != nil {
		resObj.Ok = false
		resObj.Err = errObj
	} else {
		jwtCookie := http.Cookie{
			Name:     "jwt",
			Value:    res.SignedJwtToken,
			Expires:  res.JwtExpiration,
			SameSite: http.SameSiteNoneMode,
			HttpOnly: true, //true will mitigate the risk of client side script accessing the protected cookie
		}
		jwtCookieString = jwtCookie.String()

		resObj.Ok = true
		resObj.Data = ResponseDTOV1{
			HubUserId:       res.HubUserId,
			Email:           res.Email,
			JwtCookieString: jwtCookieString,
		}
	}

	lambdaResp := lambdahelper.HandleLogAndPanic(logger, errObj)
	if lambdaResp != nil {
		//will only be executed when panic because lambdaResp is only not nil when panic happened.
		return *lambdaResp, nil
	}
	return lambdahelper.WriteResponseFn(resObj, jwtCookieString), nil
}
