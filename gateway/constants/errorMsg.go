package constants

const (
	// 400 errors
	// user parameters
	ERROR_BAD_REQUEST_MSG     = "error_bad_request"
	ERROR_NO_COOKIE_MSG       = "error_no_cookie"
	ERROR_INVALID_REQUEST_MSG = "error_invalid_request"

	// 401 errors
	ERROR_UNAUTHORIZED_MSG          = "error_user_unauthorized"
	ERROR_TOKEN_INVALID_MSG         = "error_token_invalid"
	ERROR_JWT_SIGNATURE_INVALID_MSG = "error_signature_invalid"
	ERROR_GENERATE_JWT_TOKEN_MSG    = "error_generate_jwt_token"

	// 500 errors
	// server errors
	ERROR_SERVER_START_FAIL_MSG      = "error_server_start_fail"
	ERROR_USERSERVICE_CONNECTION_MSG = "error_userservice_connection"
)
