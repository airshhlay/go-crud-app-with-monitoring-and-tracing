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
	// ErrorItemserviceConnectionMsg service error message
	ErrorItemserviceConnectionMsg = "error_itemservice_connection"
	// ErrorNoUserIDReturnedMsg service error message
	ErrorNoUserIDReturnedMsg = "error_no_userid_returned"
	// ErrorGrpcClientStartFailMsg service error message
	ErrorGrpcClientStartFailMsg = "error_grpc_client_start_fail"
	// ErrorJaegerInitMsg service error message
	ErrorJaegerInitMsg = "error_jaeger_init"
	// ErrorJWTSignatureInvalidMsg service error message
	ErrorJWTSignatureInvalidMsg = "error_jwt_signature_invalid"
	// ErrorUnexpectedJWTErr service error message
	ErrorUnexpectedJWTErr = "error_unexpected_jwt_err"
	// ErrorGetUserIDFromTokenMsg service error message
	ErrorGetUserIDFromTokenMsg = "error_get_userid_from_token"
	// ErrorNoUserIDInTokenMsg service error message
	ErrorNoUserIDInTokenMsg = "error_no_userid_in_token"
	// ErrorParseIntMsg service error message
	ErrorParseIntMsg = "error_parse_int"
	// ErrorTypeAssertionMsg service error message
	ErrorTypeAssertionMsg = "error_type_assertion"
	// ErrorCreateGRPCChannelMsg service error message
	ErrorCreateGRPCChannelMsg = "error_create_grpc_channel"
)
