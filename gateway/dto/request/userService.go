package request

// LoginReq defines the expected incoming request body to Login
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// SignupReq defines the expected incoming request body to Signup
type SignupReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
