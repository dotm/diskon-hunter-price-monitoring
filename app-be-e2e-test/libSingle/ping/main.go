package main

import (
	"bytes"
	"diskon-hunter/price-monitoring-e2e-test/shared"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	endpoint := shared.GetBackendUrl() + "/ping"

	body := []byte(``)
	req, err := http.NewRequest(http.MethodGet, endpoint, bytes.NewBuffer(body))
	if err != nil {
		fmt.Printf("error creating request: %v", err.Error())
		return
	}
	req.Header.Set("User-Agent", "go-e2e-backend")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("error executing request: %v", err.Error())
		return
	}
	defer res.Body.Close()

	fmt.Printf("response Status:\n%v\n\n", res.Status)
	responseHeader, err := (json.Marshal(res.Header))
	if err != nil {
		fmt.Printf("error reading response header: %v", err.Error())
		return
	}
	fmt.Printf("response Headers:\n%+v\n\n", string(responseHeader))
	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %v", err.Error())
		return
	}
	fmt.Printf("response Body:\n%+v\n", string(body))
}
