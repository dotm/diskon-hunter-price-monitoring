package httphelper

import (
	"diskon-hunter/price-monitoring/shared/createerror"
	"diskon-hunter/price-monitoring/shared/lazylogger"
	"diskon-hunter/price-monitoring/shared/serverresponse"
	"fmt"
	"net/http"
	"runtime/debug"
)

func HandleLogAndPanic(w http.ResponseWriter, logger *lazylogger.Instance, errObj *serverresponse.ErrorObj) {
	if errObj != nil && errObj.SendErrorTo[serverresponse.SendErrorToLog] {
		logger.LogQueueAsErrorAndDequeueAllItems()
	}
	if err := recover(); err != nil {
		logger.EnqueuePanicLog(err, debug.Stack(), true)
		logger.LogQueueAsErrorAndDequeueAllItems()

		WriteResponseFn(w, serverresponse.Obj{
			Ok:  false,
			Err: createerror.InternalException(fmt.Errorf("pnc")),
		})
	}
}

func WriteResponseFn(w http.ResponseWriter, resObj serverresponse.Obj) {
	w.Header().Set("Content-Type", "application/json")
	w.Write(resObj.ToByteSlice())
}
