package tracing

const (

	// tag values

	// Component values

	// ComponentGrpc for <component>
	ComponentGrpc = "grpc"

	// PeerServiceUserService for <peer.service>
	PeerServiceUserService = "userservice"
	// PeerServiceItemService for <peer.service>
	PeerServiceItemService = "itemservice"

	// Span kind

	// SpanKindClient for <span.kind>
	SpanKindClient = "client"
	// SpanKindServer for <span.kind>
	SpanKindServer = "server"
)
