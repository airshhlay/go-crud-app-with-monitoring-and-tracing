package constants

const (
	// server error

	// ErrorLoadConfigFailMsg server error message
	ErrorLoadConfigFailMsg = "error_load_config_fail"
	// ErrorServerStartFailMsg server error message
	ErrorServerStartFailMsg = "error_server_start_fail"
	// ErrorPromHTTPServerMsg server error message
	ErrorPromHTTPServerMsg = "error_prom_http_sever"
	// ErrorPromInitCustomMetricsMsg server error message
	ErrorPromInitCustomMetricsMsg = "error_prom_init_custom_metrics"
	// ErrorJaegerInitMsg service error message
	ErrorJaegerInitMsg = "error_jaeger_init"

	// database

	// ErrorDatabaseMsg server error message
	ErrorDatabaseMsg = "error_database_operation"
	// ErrorDatabaseInsertMsg server error message
	ErrorDatabaseInsertMsg = "error_database_insert"
	// ErrorDatabaseQueryMsg server error message
	ErrorDatabaseQueryMsg = "error_database_query"
	// ErrorDatabaseDeleteMsg server error message
	ErrorDatabaseDeleteMsg = "error_database_delete"
	// ErrorDatabaseConnectionMsg server error message
	ErrorDatabaseConnectionMsg = "error_database_connection"
	// ErrorRedisConnectionMsg server error message
	ErrorRedisConnectionMsg = "error_redis_connection"
	// ErrorRedisGetMsg server error message
	ErrorRedisGetMsg = "error_redis_get"
	// ErrorRedisSetMsg server error message
	ErrorRedisSetMsg = "error_redis_set"

	// external calls

	// ErrorExternalAPICallMsg server error message
	ErrorExternalAPICallMsg = "error_external_api_call"
	// ErrorExternalShopeeAPICallMsg server error message
	ErrorExternalShopeeAPICallMsg = "error_external_shopee_api_call"

	// marshalling / unmarshalling

	// ErrorMarshalMsg server error message
	ErrorMarshalMsg = "error_marshal"
	// ErrorUnmarshalMsg server error message
	ErrorUnmarshalMsg = "error_unmarshal"

	// typecasting

	// ErrorTypecastMsg server error message
	ErrorTypecastMsg = "error_typecast"
)
