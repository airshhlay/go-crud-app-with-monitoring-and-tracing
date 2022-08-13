package constants

const (
	// 400 errors
	// user parameters

	// ErrorBadRequestMsg service error message
	ErrorBadRequestMsg = "error_bad_request"
	// ErrorNoCookieMsg service error message
	ErrorNoCookieMsg = "error_no_cookie"
	// ErrorInvalidRequestMsg service error message
	ErrorInvalidRequestMsg = "error_invalid_request"

	// 401 errors

	// ErrorUnauthorizedMsg service error message
	ErrorUnauthorizedMsg = "error_user_unauthorized"
	// ErrorTokenInvalidMsg service error message
	ErrorTokenInvalidMsg = "error_token_invalid"
	// ErrorGenerateJWTTokenMsg service error message
	ErrorGenerateJWTTokenMsg = "error_generate_jwt_token"

	// 500 errors
	// server errors

	// ErrorServerStartFailMsg service error message
	ErrorServerStartFailMsg = "error_server_start_fail"
	// ErrorUserserviceConnectionMsg service error message
	ErrorUserserviceConnectionMsg = "error_userservice_connection"
	// ErrorGrpcClientStartFailMsg service error message
	ErrorGrpcClientStartFailMsg = "error_grpc_client_start_fail"
	// ErrorJaegerInitMsg service error message
	ErrorJaegerInitMsg = "error_jaeger_init"
	// ErrorJWTSignatureInvalidMsg service error message
	ErrorJWTSignatureInvalidMsg = "error_jwt_signature_invalid"
	// ErrorUnexpectedJWTErr service error message
	ErrorUnexpectedJWTErr = "error_unexpected_jwt_err"
)
