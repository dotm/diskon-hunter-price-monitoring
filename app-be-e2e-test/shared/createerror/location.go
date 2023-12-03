package createerror

import "diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"

const UnknownLocationTypeErrorCode = "location/unknown-type"

func UnknownLocationType(err error) *serverresponse.ErrorObj {
	return Response(
		UnknownLocationTypeErrorCode,
		err,
		map[string]bool{},
	)
}
