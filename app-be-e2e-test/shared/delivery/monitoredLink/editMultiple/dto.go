package delivery

import (
	"diskon-hunter/price-monitoring-e2e-test/shared/currencyutil"
	"diskon-hunter/price-monitoring-e2e-test/shared/delivery/monitoredLink"
)

const PathV1 = "/v1/monitoredLink.editMultiple"

// Data Transfer Object is used for API contract with clients
// e.g. frontend app, mobile app, external API call.

// RequestDTO is for data coming from clients to the server.
// Keep in mind that there are other mechanism for incoming data transfer (the most common one is JWT claim).
type RequestDTOV1 struct {
	MonitoredLinkList []MonitoredLinkRequestDTOV1
}

type MonitoredLinkRequestDTOV1 struct {
	HubMonitoredLinkUrl   string //assumed to be cleaned because this comes from our database
	ActiveAlertMethodList []string
	AlertPrice            currencyutil.Currency
}

// ResponseDTO is for data going from the server to clients.
// This will be wrapped in the data field of server response object.

type ResponseDTOV1 = []monitoredLink.StlUserMonitorsLinkDetailDAOV1
