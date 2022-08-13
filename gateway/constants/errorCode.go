package constants

// service code starts with 2<error group><error category><number>
const (
	// 400 errors
	// user parameters

	// ErrorBadRequest service error code
	ErrorBadRequest = 140011
	// ErrorNoCookie service error code
	ErrorNoCookie = 140012
	// ErrorInvalidRequest service error code
	ErrorInvalidRequest = 140013

	// 401 errors
	// ErrorUnauthorized service error code
	ErrorUnauthorized = 140111
	// ErrorTokenInvalid service error code
	ErrorTokenInvalid = 140112

	// 500 errors
	// server errors

	// ErrorUserserviceConnection service error code
	ErrorUserserviceConnection = 150011

	// JWT errors

	// ErrorJWTSignatureInvalid service error code
	ErrorJWTSignatureInvalid = 150021
	// ErrorGenerateJWTToken service error code
	ErrorGenerateJWTToken = 150022
)
