package shared

import (
	"bytes"
	"diskon-hunter/price-monitoring-e2e-test/shared/jwttoken"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CreatePostRequestWithJsonBodyArgs struct {
	Body     interface{}
	Endpoint string
	JwtToken string
}

func CreatePostRequestWithJsonBody(args CreatePostRequestWithJsonBodyArgs) (*http.Request, error) {
	body, err := json.Marshal(args.Body)
	if err != nil {
		fmt.Printf("error creating request body: %v", err.Error())
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, args.Endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("error creating request: %v", err.Error())
		return nil, err
	}
	req.Header.Set("User-Agent", "go-e2e-backend")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwttoken.CleanUpJWT(args.JwtToken)))

	return req, nil
}

func LogResponse(res *http.Response) error {
	fmt.Printf("response Status:\n%v\n\n", res.Status)
	responseHeader, err := (json.Marshal(res.Header))
	if err != nil {
		fmt.Printf("error reading response header: %v", err.Error())
		return err
	}
	fmt.Printf("response Headers:\n%+v\n\n", string(responseHeader))
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v", err.Error())
		return err
	}
	fmt.Printf("response Body:\n%+v\n", string(body))

	return nil
}

func LogResponseAndReturnBody(res *http.Response) (body []byte, err error) {
	fmt.Printf("response Status:\n%v\n\n", res.Status)
	responseHeader, err := (json.Marshal(res.Header))
	if err != nil {
		fmt.Printf("error reading response header: %v", err.Error())
		return body, err
	}
	fmt.Printf("response Headers:\n%+v\n\n", string(responseHeader))
	body, err = io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v", err.Error())
		return body, err
	}
	fmt.Printf("response Body:\n%+v\n", string(body))

	return body, nil
}
