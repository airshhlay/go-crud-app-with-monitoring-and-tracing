package response

type LoginRes struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}

type SignupRes struct {
	ErrorCode int32  `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`
}
