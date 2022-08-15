package tracing

const (

	// tag values

	// DatabaseTypeSQL for <db.type>
	DatabaseTypeSQL = "sql"
	// DatabaseTypeRedis for <db.type>
	DatabaseTypeRedis = "redis"

	// ComponentServer for <component>
	ComponentServer = "itemService.server"
	// ComponentDB for <component>
	ComponentDB = "itemService.db"
	// ComponentExternal for <component>
	ComponentExternal = "itemService.external"

	// PeerServiceUserService for <peer.service>
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
