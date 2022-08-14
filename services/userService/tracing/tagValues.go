package tracing

const (

	// tag values

	// DatabaseTypeSQL for <db.type>
	DatabaseTypeSQL = "sql"

	// ComponentServer for <component>
	ComponentServer = "userService.server"
	// ComponentDB for <component>
	ComponentDB = "userService.db"

	// PeerServiceUserService for <peer.service>
	PeerServiceMySQL = "mysql"

	// Span kind

	// SpanKindClient for <span.kind>
	SpanKindClient = "client"
	// SpanKindServer for <span.kind>
	SpanKindServer = "server"
)
