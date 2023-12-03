package serverresponse

import (
	"encoding/json"
	"fmt"
)

type Obj struct {
	Ok   bool        `json:"ok"`
	Err  *ErrorObj   `json:"err,omitempty"`  //prefixed with * because nullable
	Data interface{} `json:"data,omitempty"` //nullable
}

type ErrorObj struct {
	Code        string          `json:"code"`
	Message     *string         `json:"msg,omitempty"` //prefixed with * because nullable
	SendErrorTo map[string]bool `json:"-"`
}

const SendErrorToLog = "to-log"
const SendErrorToDevEmail = "to-dev-email"

func (x Obj) ToByteSlice() []byte {
	byteSlice, err := json.Marshal(x)
	if err != nil {
		err = fmt.Errorf("error marshaling command response: %v", err)
		errMsg := fmt.Sprintf(
			`{"err": {"code": "internal/exception", "message": "%v"}}`,
			err,
		)
		fmt.Println(errMsg)
		return []byte(errMsg)
	}

	return byteSlice
}
