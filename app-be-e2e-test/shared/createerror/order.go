package createerror

import "diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"

const OrderNotFoundErrorCode = "order/not-found"

func OrderNotFound(err error) *serverresponse.ErrorObj {
	return Response(
		OrderNotFoundErrorCode,
		err,
		map[string]bool{},
	)
}
