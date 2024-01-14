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
	searchedItemAddMultiple "diskon-hunter/price-monitoring/src/searchedItem/addMultiple/function"
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

	var res searchedItemAddMultiple.CommandV1DataResponse
	if errObj == nil {
		searchedItemList := []searchedItemAddMultiple.SearchedItemDetailCommandV1{}
		for i := 0; i < len(reqBody.SearchedItemList); i++ {
			searchedItemList = append(searchedItemList, searchedItemAddMultiple.NewSearchedItemDetailCommandV1(
				reqBody.SearchedItemList[i].Name,        //Name
				reqBody.SearchedItemList[i].Description, //Description
				reqBody.SearchedItemList[i].AlertPrice,  //AlertPrice
			))
		}
		cmd := searchedItemAddMultiple.NewCommandV1(
			"1",                       //Version
			jwttoken.GetUserId(token), //RequesterUserId
			searchedItemList,          //SearchedItemList
		)

		res, errObj = searchedItemAddMultiple.CommandV1Handler(
			ctx,
			searchedItemAddMultiple.CommandV1Dependencies{
				Logger:         logger,
				DynamoDBClient: dynamodbhelper.CreateClientFromSession(),
			},
			cmd,
		)
	}
	resDTO := ResponseDTOV1{
		SearchedItemIdList: res.SearchedItemIdList,
	}

	var resObj serverresponse.Obj
	if errObj != nil {
		resObj.Ok = false
		resObj.Err = errObj
	} else {
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
