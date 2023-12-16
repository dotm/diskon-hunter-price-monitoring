package createerror

import "diskon-hunter/price-monitoring/shared/serverresponse"

const AlertMethodHasNotBeenPaidErrorCode = "alert-method/has-not-been-paid"

func AlertMethodHasNotBeenPaid(err error) *serverresponse.ErrorObj {
	return Response(
		AlertMethodHasNotBeenPaidErrorCode,
		err,
		map[string]bool{},
	)
}
