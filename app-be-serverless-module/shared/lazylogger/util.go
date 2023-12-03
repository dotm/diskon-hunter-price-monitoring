package lazylogger

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// will be runned when this shared package is imported
func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel) //set log level
}

type Instance struct {
	endpointUrl string //http endpoint url
	queue       []string
}

func New(endpointUrl string) *Instance {
	return &Instance{endpointUrl: endpointUrl, queue: []string{}}
}

func formatTime(t time.Time) string {
	return strings.TrimSuffix(t.UTC().String(), " +0000 UTC")
}
func addNowAsFormattedTime() string {
	return formatTime(time.Now())
}

// withTime means whether we need to also log the time or not
func (x *Instance) EnqueueErrorLog(item error, withTime bool) {
	x.EnqueueStringLog(item.Error(), withTime)
}

// withTime means whether we need to also log the time or not
func (x *Instance) EnqueueJsonLog(item interface{}, withTime bool) error {
	byteSlice, err := json.Marshal(item)
	if err != nil {
		err = fmt.Errorf("error marshaling enqueued struct log: %v", err)
		x.EnqueueErrorLog(err, withTime)
		return err
		//LogQueueAsErrorAndDequeueAllItems will be handled by the function that calls New()
	}
	x.EnqueueStringLog(string(byteSlice), withTime)
	return nil
}

// withTime means whether we need to also log the time or not.
func (x *Instance) EnqueuePanicLog(errFromRecover interface{}, stackTrace []byte, withTime bool) {
	stackTraceString := string(stackTrace)
	stackTraceString = strings.ReplaceAll(stackTraceString, "\n", " >> ")
	stackTraceString = strings.ReplaceAll(stackTraceString, "\t", "")
	x.enqueueString(
		"pnc",
		fmt.Sprintf("panic occurred: %v >>> stack trace: %v", errFromRecover, stackTraceString),
		true,
	)
}

// withTime means whether we need to also log the time or not
func (x *Instance) EnqueueCommandLog(item string, withTime bool) {
	x.enqueueString("cmd", item, withTime)
}

// withTime means whether we need to also log the time or not
func (x *Instance) EnqueueQueryLog(item string, withTime bool) {
	x.enqueueString("qry", item, withTime)
}

// withTime means whether we need to also log the time or not
func (x *Instance) EnqueueStringLog(item string, withTime bool) {
	x.enqueueString("log", item, withTime)
}

func (x *Instance) enqueueString(keyName string, item string, withTime bool) {
	if withTime {
		x.queue = append(x.queue, fmt.Sprintf(`{"%v":"%v", "time":"%v"}`, keyName, item, addNowAsFormattedTime()))
	} else {
		x.queue = append(x.queue, fmt.Sprintf(`{"%v":"%v"}`, keyName, item))
	}
}

// don't run in src. for easier debugging purpose only. loggedQueue is what's inside "error" field in the log
func UnescapeAndPrintLoggedQueue(loggedQueue string) {
	unescapedString := strings.ReplaceAll(loggedQueue, `\"`, `"`)
	parsedString := strings.ReplaceAll(unescapedString, `}, {`, "}\n{")
	fmt.Println(parsedString)
}

// should be called inside the body of the same function
// that calls the New() function that instantiate the logger
func (x *Instance) LogQueueAsErrorAndDequeueAllItems() {
	errMsg := ""
	for i, item := range x.queue {
		if i != 0 {
			errMsg += ", " + item
		} else {
			errMsg += item
		}
	}

	log.Error().
		Str("path", x.endpointUrl).
		Err(fmt.Errorf(errMsg)).
		Send()
	x.queue = []string{}
}

// main functions in shared directory are used as example usage.
// to run main function, rename the package as main and then run: go run ./path/to/this/util/*.go
//
//lint:ignore U1000 Ignore unused function temporarily for debugging
func main() {
	logger := New("/v1/test_function")
	UnescapeAndPrintLoggedQueue(`{\"log\":\"1\"}, {\"log\":\"2\", \"time\":\"2022-09-11 10:26:18.468101514\"}, {\"log\":\"3\", \"time\":\"2022-09-11 10:26:18.468107238\"}`)
	submain1(logger)
	submain2(logger)
	submain3(logger)
	err := fmt.Errorf("mock error")
	defer func() {
		if err != nil {
			logger.LogQueueAsErrorAndDequeueAllItems()
		}
	}()
}
func submain1(logger *Instance) {
	logger.EnqueueStringLog("1", true)
}
func submain2(logger *Instance) {
	logger.EnqueueStringLog("2", true)
}
func submain3(logger *Instance) {
	logger.EnqueueStringLog("3", false)
}
