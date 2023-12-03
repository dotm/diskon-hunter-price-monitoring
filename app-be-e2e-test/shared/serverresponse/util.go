package serverresponse

type Obj struct {
	Ok  bool      `json:"ok"`
	Err *ErrorObj `json:"err,omitempty"` //prefixed with * because nullable

	//Data should be specified in single e2e test request
	// Data interface{} `json:"data,omitempty"` //nullable
}

type ErrorObj struct {
	Code        string          `json:"code"`
	Message     *string         `json:"msg,omitempty"` //prefixed with * because nullable
	SendErrorTo map[string]bool `json:"-"`
}

const SendErrorToLog = "to-log"
const SendErrorToDevEmail = "to-dev-email"
