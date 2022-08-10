package constants

// service code starts with 3<error group><error category><number>
const (
	// 400 errors
	// user parameters
	ERROR_ITEM_IN_FAVOURITES = 340011

	// 500 errors
	// server errors
	ERROR_SERVER_START_FAIL = 350021
	// database errors
	ERROR_DATABASE            = 350011
	ERROR_DATABASE_INSERT     = 350012
	ERROR_DATABASE_QUERY      = 350013
	ERROR_DATABASE_CONNECTION = 350014
	ERROR_DATABASE_DELETE     = 350015

	ERROR_REDIS            = 350021
	ERROR_REDIS_CONNECTION = 350022
	ERROR_REDIS_GET        = 350023
	ERROR_REDIS_SET        = 350024

	// external calls
	ERROR_EXTERNAL_API_CALL        = 350031
	ERROR_EXTERNAL_SHOPEE_API_CALL = 350032

	// marshalling and unmarshalling
	ERROR_MARSHAL   = 350041
	ERROR_UNMARSHAL = 350042

	// typecasting
	ERROR_TYPECAST = 350051
)
