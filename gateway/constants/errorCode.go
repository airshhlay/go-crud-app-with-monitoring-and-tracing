package constants

// service code starts with 2<error group><error category><number>
const (
	// 400 errors
	// user parameters

	// ErrorBadRequest service error code
	ErrorBadRequest = 140011
	// ErrorNoCookie service error code
	ErrorNoCookie = 140012
	// ErrorNoUserIDInToken service error code
	ErrorNoUserIDInToken = 140013
	// ErrorInvalidRequest service error code
	ErrorInvalidRequest = 140014

	// 401 errors
	// ErrorUnauthorized service error code
	ErrorUnauthorized = 140111
	// ErrorTokenInvalid service error code
	ErrorTokenInvalid = 140112

	// 500 errors
	// server errors

	// ErrorUserserviceConnection service error code
	ErrorUserserviceConnection = 150011
	// ErrorItemserviceConnection service error code
	ErrorItemserviceConnection = 150012
	// ErrorCreateGRPCChannel service error code
	ErrorCreateGRPCChannel = 150013
	// ErrorNoUserIDReturned service error code
	ErrorNoUserIDReturned = 150014
	// JWT errors

	// ErrorJWTSignatureInvalid service error code
	ErrorJWTSignatureInvalid = 150021
	// ErrorGenerateJWTToken service error code
	ErrorGenerateJWTToken = 150022

	// ErrorUserIDNotInToken
	ErrorGetUserIDFromToken = 150031

	// function errors

	// ErrorParseInt service error code
	ErrorParseInt = 150041
	// ErrorTypeAssertion service error code
	ErrorTypeAssertion = 150051
)
