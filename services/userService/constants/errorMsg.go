package constants

const (
	// ErrorDatabaseInsertMsg for database insert failures
	ErrorDatabaseInsertMsg = "error_database_insert_failure"
	// ErrorDatabaseQueryMsg for database query failures
	ErrorDatabaseQueryMsg = "error_database_query_failure"
	// ErrorDatabaseConnectionMsg for database connection errors
	ErrorDatabaseConnectionMsg = "error_database_connection_failure"
	// ErrorPasswordEncryptionMsg for bcrypt encryption errors
	ErrorPasswordEncryptionMsg = "error_password_encryption"
	// ErrorServerStartFailMsg for when the grpc server fails to start
	ErrorServerStartFailMsg = "error_server_start_fail"
	// ErrorPromInitCustomMetricsMsg for when prometheus fails to initialise custom metrics
	ErrorPromInitCustomMetricsMsg = "error_prom_init_custom_metrics"
	// ErrorPromHTTPServerMsg for when the http server for prometheus metrics fails to start
	ErrorPromHTTPServerMsg = "error_prom_http_server"
	// ErrorUserPasswordMSg for errors with the user password
	ErrorUserPasswordMsg = "error_user_password"
	// ErrorTypecastMsg for errors typecasting error to customErr
	ErrorTypecastMsg = "error_typecast"
)
