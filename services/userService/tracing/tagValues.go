package tracing

const (

	// tag values

	// DatabaseTypeSQL for <db.type>
	DatabaseTypeSQL = "sql"

	// ComponentGrpc for <component>
	ComponentGrpc = "gRPC"
	// ComponentMySQL for <component>
	ComponentMySQL = "mysql"

	// PeerServiceMySQL for <peer.service>
	PeerServiceMySQL = "mysql"

	// Span kind

	// SpanKindClient for <span.kind>
	SpanKindClient = "client"
	// SpanKindServer for <span.kind>
	SpanKindServer = "server"
)
