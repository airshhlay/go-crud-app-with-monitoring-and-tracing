package constants

// service code starts with 2<error group><error category><number>
const (
	// 400 errors
	// user parameters
	ErrorUserAlreadyExists = 240011
	ErrorUserDoesNotExist  = 240012
	ErrorUserPassword      = 240013

	// 500 errors
	// server errors
	ErrorServerStartFail = 250021
	// database errors
	ErrorDatabase           = 250011
	ErrorDatabaseInsert     = 250012
	ErrorDatabaseQuery      = 250013
	ErrorDatabaseConnection = 250014

	// encryption errors
	ErrorPasswordEncryption = 250021

	// typecasting
	ErrorTypecast = 250031

	// prometheus
	ErrorPromInitCustomMetrics = 250041
)
