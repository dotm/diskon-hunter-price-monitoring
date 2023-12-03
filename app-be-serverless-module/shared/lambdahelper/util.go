package lambdahelper

import (
	"bytes"
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
)

func HandleLogAndPanic(logger *lazylogger.Instance, errObj *serverresponse.ErrorObj) *events.APIGatewayProxyResponse {
	if errObj != nil && errObj.SendErrorTo[serverresponse.SendErrorToLog] {
		logger.LogQueueAsErrorAndDequeueAllItems()
	}
	if err := recover(); err != nil {
		logger.EnqueuePanicLog(err, debug.Stack(), true)
		logger.LogQueueAsErrorAndDequeueAllItems()

		panicRes := WriteResponseFn(serverresponse.Obj{
			Ok:  false,
			Err: createerror.InternalException(fmt.Errorf("pnc")),
		}, "")
		return &panicRes
	}
	return nil
}

// Disable escaping character such as &, <, and > into \u0026, \u003c, and \u003e when marshaling JSON
// Used for example when sending S3 presigned url to frontend.
func WriteResponseFnWithUnescapedJSON(resObj serverresponse.Obj, cookie string) events.APIGatewayProxyResponse {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(resObj)
	jsonResp := buffer.Bytes()

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}
	}
	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonResp),
	}
	if cookie != "" {
		resp.Headers["Set-Cookie"] = cookie
	}
	return resp
}

func WriteResponseFn(resObj serverresponse.Obj, cookie string) events.APIGatewayProxyResponse {
	jsonResp, err := json.Marshal(resObj)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}
	}
	resp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonResp),
	}
	if cookie != "" {
		resp.Headers["Set-Cookie"] = cookie
	}
	return resp
}
