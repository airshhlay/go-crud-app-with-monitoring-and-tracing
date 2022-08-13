package constants

// service code starts with 3<error group><error category><number>
const (
	// 400 errors
	// user parameters

	// ErrorItemInFavourites service error code
	ErrorItemInFavourites = 340011

	// 500 errors
	// server errors

	// ErrorServerStartFail service error code
	ErrorServerStartFail = 350021
	// database errors

	// ErrorDatabase service error code
	ErrorDatabase = 350011
	// ErrorDatabaseInsert service error code
	ErrorDatabaseInsert = 350012
	// ErrorDatabaseQuery service error code
	ErrorDatabaseQuery = 350013
	// ErrorDatabaseConnection service error code
	ErrorDatabaseConnection = 350014
	// ErrorDatabaseDelete service error code
	ErrorDatabaseDelete = 350015

	// ErrorRedis service error code
	ErrorRedis = 350021
	// ErrorRedisConnection service error code
	ErrorRedisConnection = 350022
	// ErrorRedisGet service error code
	ErrorRedisGet = 350023
	// ErrorRedisSet service error code
	ErrorRedisSet = 350024

	// external calls

	// ErrorExternalAPICall service error code
	ErrorExternalAPICall = 350031
	// ErrorExternalShopeeAPICall service error code
	ErrorExternalShopeeAPICall = 350032

	// marshalling and unmarshalling

	// ErrorMarshal service error code
	ErrorMarshal = 350041
	// ErrorUnmarshal service error code
	ErrorUnmarshal = 350042

	// typecasting

	// ErrorTypecast service error code
	ErrorTypecast = 350051
)
