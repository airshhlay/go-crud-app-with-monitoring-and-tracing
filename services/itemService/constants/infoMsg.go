package constants

const (
	// server

	// InfoConfigLoaded info for logging
	InfoConfigLoaded = "info_config_loaded"
	// InfoPromServerStart info for logging
	InfoPromServerStart = "info_prom_sever_start"
	// InfoHTTPServerStart info for logging
	InfoHTTPServerStart = "info_http_server_start"
	// InfoGRPCServerStart info for logging
	InfoGRPCServerStart = "info_grpc_server_start"

	// database

	// InfoDatabaseQuery info for logging
	InfoDatabaseQuery = "info_db_query"
	// InfoDatabaseQueryRows info for logging
	InfoDatabaseQueryRows = "info_db_query_rows"
	// InfoDatabaseInsert info for logging
	InfoDatabaseInsert = "info_db_insert"
	// InfoDatabaseConnectSuccess info for logging
	InfoDatabaseConnectSuccess = "info_db_connect_success"
	// InfoDatabaseDelete info for logging
	InfoDatabaseDelete = "info_db_delete"

	// InfoRedisConnectSuccess info for logging
	InfoRedisConnectSuccess = "info_redis_connect_success"
	// InfoRedisSet info for logging
	InfoRedisSet = "info_redis_set"
	// InfoRedisGet info for logging
	InfoRedisGet = "info_redis_get"
	// InfoRedisNotFound info for logging
	InfoRedisNotFound = "info_redis_item_not_found"

	// InfoExternalAPICall info for logging
	InfoExternalAPICall = "info_External_api_call"

	// queries

	// InfoFavouriteAdded info for logging
	InfoFavouriteAdded = "info_favourite_added"
	// InfoItemNotInFavourites info for logging
	InfoItemNotInFavourites = "info_item_not_in_favourites"
	// InfoItemInFavourites info for logging
	InfoItemInFavourites = "info_item_in_favourites"
)
