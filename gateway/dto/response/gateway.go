package response

// GatewayResponse defines the response sent by the gateway should errors occur at the gateway.
type GatewayResponse struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}
