package createerror

import (
	"diskon-hunter/price-monitoring/shared/serverresponse"
)

func Response(code string, err error, sendErrorTo map[string]bool) *serverresponse.ErrorObj {
	res := &serverresponse.ErrorObj{
		Code:        code,
		SendErrorTo: sendErrorTo,
	}
	if err != nil {
		errMsg := err.Error()
		res.Message = &errMsg
	}
	return res
}

// something that is unexpected and originate from our system
func InternalException(err error) *serverresponse.ErrorObj {
	return Response(
		"internal/exception",
		err,
		map[string]bool{
			serverresponse.SendErrorToLog:      true,
			serverresponse.SendErrorToDevEmail: true,
		},
	)
}

// something that is unexpected and originate from systems other than ours (e.g. partner)
func ExternalException(err error) *serverresponse.ErrorObj {
	return Response(
		"external/exception",
		err,
		map[string]bool{
			serverresponse.SendErrorToLog:      true,
			serverresponse.SendErrorToDevEmail: true,
		},
	)
}

// exceptions caused by bad request from clients.
// if possible don't use this code; instead create new and more domain specific error code
func ClientBadRequest(err error) *serverresponse.ErrorObj {
	return Response(
		"client/bad-request",
		err,
		map[string]bool{
			serverresponse.SendErrorToLog: true,
		},
	)
}

func PageNotFound(err error) *serverresponse.ErrorObj {
	return Response(
		"page/not-found",
		err,
		map[string]bool{
			serverresponse.SendErrorToLog: true,
		},
	)
}
