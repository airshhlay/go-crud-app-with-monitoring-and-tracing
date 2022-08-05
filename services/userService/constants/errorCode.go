package constants

// service code starts with 2<error group><error category><number>
const (
	// 400 errors
	// user parameters
	ERROR_USER_ALREADY_EXISTS = 240011
	ERROR_USER_DOES_NOT_EXIST = 240012
	ERROR_USER_PASSWORD       = 240013

	// 500 errors
	// server errors
	ERROR_SERVER_START_FAIL = 250021
	// database errors
	ERROR_DATABASE            = 250011
	ERROR_DATABASE_INSERT     = 250012
	ERROR_DATABASE_QUERY      = 250013
	ERROR_DATABASE_CONNECTION = 250014

	// encryption errors
	ERROR_PASSWORD_ENCRYPTION = 250021
)
