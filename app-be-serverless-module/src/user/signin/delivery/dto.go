package delivery

const PathV1 = "/v1/user.signin"

// Data Transfer Object is used for API contract with clients
// e.g. frontend app, mobile app, external API call.

// RequestDTO is for data coming from clients to the server.
// Keep in mind that there are other mechanism for incoming data transfer (the most common one is JWT claim).
type RequestDTOV1 struct {
	Email    string
	Password string
}

// ResponseDTO is for data going from the server to clients.
// This will be wrapped in the data field of server response object.
type ResponseDTOV1 struct {
	HubUserId       string
	Email           string
	JwtCookieString string //workaround because flutter web can't read set-cookie http response header
}
