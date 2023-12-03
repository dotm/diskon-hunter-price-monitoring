package createerror

import "diskon-hunter/price-monitoring/shared/serverresponse"

const UnknownLocationTypeErrorCode = "location/unknown-type"

func UnknownLocationType(err error) *serverresponse.ErrorObj {
	return Response(
		UnknownLocationTypeErrorCode,
		err,
		map[string]bool{},
	)
}
