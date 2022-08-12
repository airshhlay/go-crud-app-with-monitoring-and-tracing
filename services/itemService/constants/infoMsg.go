package constants

const (
	INFO_SERVER_START_MSG = "info_server_start_success"

	// database
	INFO_DATABASE_QUERY           = "info_db_query"
	INFO_DATABASE_QUERY_ROWS      = "info_db_query_rows"
	INFO_DATABASE_INSERT          = "info_db_insert"
	INFO_DATABASE_CONNECT_SUCCESS = "info_db_connect_success"
	INFO_DATABASE_DELETE          = "info_db_delete"

	INFO_REDIS_CONNECT_SUCCESS = "info_redis_connect_success"
	INFO_REDIS_SET             = "info_redis_set"
	INFO_REDIS_GET             = "info_redis_get"
	INFO_REDIS_NOT_FOUND       = "info_redis_item_not_found"

	INFO_EXTERNAL_API_CALL = "info_External_api_call"

	// queries
	INFO_ITEM_NOT_IN_FAVOURITES = "info_item_not_in_favourites"
	INFO_ITEM_IN_FAVOURITES     = "info_item_in_favourites"

	INFO_PROM_SERVER_START_MSG = "info_prom_sever_start"
)
