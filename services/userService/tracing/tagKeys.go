package tracing

const (
	// tag keys

	// Component tag
	Component = "component"
	// DatabaseInstance tag
	DatabaseInstance = "db.instance"
	// DatabaseStatement tag
	DatabaseStatement = "db.statement"
	// DatabaseType tag
	DatabaseType = "db.type"
	// DatabaseUser tag
	DatabaseUser = "db.user"
	// Error tag
	Error = "http.error"
	// SpanKind tag
	SpanKind = "span.kind"
	// PeerService tag
	PeerService = "peer.service"
	// PeerHostname tag
	PeerHostname = "peer.hostname"
	// PeerPort tag
	PeerPort = "peer.port"

	// custom tag keys

	// ServiceErrorCode custom tag key
	ServiceErrorCode = "service.errorCode"
	// ServiceErrorMsg custom tag key
	ServiceErrorMsg = "service.errorMsg"

	// log fields

	// Message log field
	Message = "message"
	// Event log field
	Event = "event"
)
