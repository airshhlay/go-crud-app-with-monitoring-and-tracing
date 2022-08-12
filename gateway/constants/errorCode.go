package constants

// service code starts with 2<error group><error category><number>
const (
	// 400 errors
	// user parameters
	ERROR_BAD_REQUEST     = 140011
	ERROR_NO_COOKIE       = 140012
	ERROR_INVALID_REQUEST = 140013

	// 401 errors
	ERROR_UNAUTHORIZED          = 140111
	ERROR_TOKEN_INVALID         = 140112
	ERROR_JWT_SIGNATURE_INVALID = 140113

	// 500 errors
	// server errors
	ERROR_USERSERVICE_CONNECTION = 150011
	// database errors
	ERROR_DATABASE            = 150011
	ERROR_DATABASE_INSERT     = 150012
	ERROR_DATABASE_QUERY      = 150013
	ERROR_DATABASE_CONNECTION = 150014

	// encryption errors
	ERROR_PASSWORD_ENCRYPTION = 150021

	// typecasting
	ERROR_TYPECAST = 150031

	// prometheus
	ERROR_PROM_INIT_CUSTOM_METRICS = 150041

	// token errors
	ERROR_GENERATE_JWT_TOKEN = 150051
)
