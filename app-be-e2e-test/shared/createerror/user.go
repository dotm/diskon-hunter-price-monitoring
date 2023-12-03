package createerror

import "diskon-hunter/price-monitoring-e2e-test/shared/serverresponse"

func UserNotFound(err error) *serverresponse.ErrorObj {
	return Response(
		"user/not-found",
		err,
		map[string]bool{
			serverresponse.SendErrorToLog:      true,
			serverresponse.SendErrorToDevEmail: true,
		},
	)
}

const UserCredentialIncorrectErrorCode = "user/credential-incorrect"

func UserCredentialIncorrect() *serverresponse.ErrorObj {
	return Response(
		UserCredentialIncorrectErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserCredentialEmptyErrorCode = "user/credential-empty"

func UserCredentialEmpty() *serverresponse.ErrorObj {
	return Response(
		UserCredentialEmptyErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserCredentialMalformedErrorCode = "user/credential-malformed"

func UserCredentialMalformed() *serverresponse.ErrorObj {
	return Response(
		UserCredentialMalformedErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserInvalidCredentialIssuerErrorCode = "user/invalid-credential-issuer"

func UserInvalidCredentialIssuer() *serverresponse.ErrorObj {
	return Response(
		UserInvalidCredentialIssuerErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserCredentialExpiredErrorCode = "user/credential-expired"

func UserCredentialExpired() *serverresponse.ErrorObj {
	return Response(
		UserCredentialExpiredErrorCode,
		nil,
		map[string]bool{},
	)
}

const UserAuthenticationShouldSpecifyEmailCode = "user/auth-no-email"

func UserAuthenticationShouldSpecifyEmail() *serverresponse.ErrorObj {
	return Response(
		UserAuthenticationShouldSpecifyEmailCode,
		nil,
		map[string]bool{},
	)
}

const UserEmailAlreadyRegisteredCode = "user/email-already-registered"

func UserEmailAlreadyRegistered() *serverresponse.ErrorObj {
	return Response(
		UserEmailAlreadyRegisteredCode,
		nil,
		map[string]bool{},
	)
}

const UserEmailNotRegisteredCode = "user/email-not-registered"

func UserEmailNotRegistered() *serverresponse.ErrorObj {
	return Response(
		UserEmailNotRegisteredCode,
		nil,
		map[string]bool{},
	)
}

const UserUnauthorizedCode = "user/unauthorized"

func UserUnauthorized(err error) *serverresponse.ErrorObj {
	return Response(
		UserUnauthorizedCode,
		err,
		map[string]bool{},
	)
}
