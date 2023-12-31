package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	LatestVersion    = "v1_20231231_1357" //change this to force frontend to reload
	v1_20231231_1357 = "v1_20231231_1357"
)

type Response struct {
	Version string `json:"version"`
}

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

func HandleRequest(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp := Response{
		Version: LatestVersion,
	}

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		return serverError(err)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(jsonResp),
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}
