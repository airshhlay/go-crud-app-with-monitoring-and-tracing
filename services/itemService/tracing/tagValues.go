package tracing

const (

	// tag values

	// DatabaseTypeSQL for <db.type>
	DatabaseTypeSQL = "sql"
	// DatabaseTypeRedis for <db.type>
	DatabaseTypeRedis = "redis"

	// ComponentServer for <component>
	ComponentServer = "itemService.server"
	// ComponentMySQL for <component>
	ComponentMySQL = "mysql"

	// ComponentRedis for <component>
	ComponentRedis = "redis"
	// ComponentHTTP for <component>
	ComponentHTTP = "http"

	// PeerServiceMySQL for <peer.service>
	PeerServiceMySQL = "mysql"
	// PeerServiceRedis for <peer.service>
	PeerServiceRedis = "redis"
	// Span kind

	// SpanKindClient for <span.kind>
	SpanKindClient = "client"
	// SpanKindServer for <span.kind>
	SpanKindServer = "server"

	// DatabaseStatementRedisSet for <db.statement> with format specifiers for key and value
	DatabaseStatementRedisSet = "SET %s %s"
	// DatabaseStatementRedisGet for <db.statement> with format specifiers for key and value
	DatabaseStatementRedisGet = "GET %s"
)
