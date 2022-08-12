package response

type GatewayResponse struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}
