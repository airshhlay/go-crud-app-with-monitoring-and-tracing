package response

// LoginRes defines the response sent back to the client by the gateway. It removes the userID from the initial response received from the user service.
type LoginRes struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}
