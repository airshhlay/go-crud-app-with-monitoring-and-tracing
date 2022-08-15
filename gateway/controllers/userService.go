package controllers

import (
	client "gateway/client"
	config "gateway/config"
	constants "gateway/constants"
	req "gateway/dto/request"
	res "gateway/dto/response"
	metrics "gateway/metrics"
	"gateway/middleware"
	proto "gateway/proto"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	token = "token"
)

// UserServiceController is called to handle incoming HTTP requests directed to the user service.
type UserServiceController struct {
	config *config.UserServiceConfig
	logger *zap.Logger
	client *client.UserServiceClient
}

// NewUserServiceController returns a UserServiceController.
func NewUserServiceController(config *config.UserServiceConfig, logger *zap.Logger, client *client.UserServiceClient) *UserServiceController {
	return &UserServiceController{
		config,
		logger,
		client,
	}
}

// LoginHandler handles requests to the /user/login endpoint.
func (u *UserServiceController) LoginHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int32

	// observe request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(u.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	}))

	defer func() {
		timer.ObserveDuration()
	}()

	var loginReq req.LoginReq
	err := c.BindJSON(&loginReq)
	if err != nil {
		u.logger.Error(
			constants.ErrorBadRequestMsg,
			zap.Error(err),
		)
		u.removeCookie(c, constants.Token)
		errorCodeInt = constants.ErrorBadRequest
		c.JSON(200, res.GatewayResponse{ErrorCode: errorCodeInt})
		return
	}
	u.logger.Info(
		constants.InfoUserServiceRequest,
		zap.Any(constants.Username, loginReq.Username),
	)

	// construct the request to be made as a grpc client to user service
	clientLoginReq := &proto.LoginReq{
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	// call user service
	clientLoginRes, err := u.client.Login(c, clientLoginReq)
	if err != nil {
		u.removeCookie(c, constants.Token)
		errorCodeInt = constants.ErrorUserserviceConnection
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorUserserviceConnection})
		return
	}
	if clientLoginRes.ErrorCode != -1 {
		// remove any credentials if there is a login error
		u.removeCookie(c, constants.Token)
	}

	if clientLoginRes.UserID != 0 {
		// no error occured
		tokenString, expirationTime, err := u.generateToken(clientLoginRes.UserID)
		if err != nil {
			u.logger.Error(
				constants.ErrorGenerateJWTTokenMsg,
				zap.Error(err),
			)
			u.removeCookie(c, constants.Token)
			errorCodeInt = constants.ErrorGenerateJWTToken
			c.IndentedJSON(200, res.GatewayResponse{ErrorCode: constants.ErrorGenerateJWTToken})
		} else {
			// set jwt token in cookie
			http.SetCookie(
				c.Writer, &http.Cookie{
					Name:     constants.Token,
					Value:    tokenString,
					Expires:  expirationTime,
					HttpOnly: true,
					Path:     "/",
				},
			)
		}
	}

	loginRes := res.LoginRes{
		ErrorCode: clientLoginRes.ErrorCode,
		ErrorMsg:  clientLoginRes.ErrorMsg,
	}
	errorCodeInt = loginRes.ErrorCode

	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(errorCodeInt))
	c.IndentedJSON(200, loginRes)
}

// generateToken is a helper function to generate the JWT token for an authenticated user's session.
func (u *UserServiceController) generateToken(userID int64) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &middleware.Claims{
		UserID: strconv.FormatInt(userID, 10),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.config.Secret))
	return tokenString, expirationTime, err
}

// SignupHandler handles incoming requests to the /user/signup endpoint.
func (u *UserServiceController) SignupHandler(c *gin.Context) {
	var errorCodeStr string
	var errorCodeInt int32

	// observe request latency
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		metrics.RequestLatency.WithLabelValues(u.config.Label, c.Request.URL.Path, errorCodeStr).Observe(v)
	}))
	defer func() {
		timer.ObserveDuration()
	}()

	// remove any token cookies
	u.removeCookie(c, constants.Token)

	var signupReq req.SignupReq
	err := c.BindJSON(&signupReq)
	if err != nil {
		u.logger.Info(
			constants.ErrorBadRequestMsg,
			zap.Error(err),
		)
		u.removeCookie(c, constants.Token)
		errorCodeInt = constants.ErrorBadRequest
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorBadRequest})
		return
	}
	u.logger.Info(
		constants.InfoUserServiceRequest,
		zap.Any(constants.Request, signupReq),
	)

	// construct the request to be made as a grpc client to user service
	clientSignupReq := &proto.SignupReq{
		Username: signupReq.Username,
		Password: signupReq.Password,
	}
	// call user service
	clientLoginRes, err := u.client.Signup(c, clientSignupReq)
	if err != nil {
		errorCodeInt = constants.ErrorUserserviceConnection
		u.removeCookie(c, constants.Token)
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ErrorUserserviceConnection})
		return
	}
	if clientLoginRes.ErrorCode != 0 {
		u.removeCookie(c, constants.Token)
	}

	errorCodeInt = clientLoginRes.ErrorCode
	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(errorCodeInt))
	c.JSON(200, clientLoginRes)
}

// removeCookie is a helper function to remove the http cookie with cookieName from the client side.
func (u *UserServiceController) removeCookie(c *gin.Context, cookieName string) {
	// set jwt token in cookie
	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:   cookieName,
			MaxAge: -1,
		},
	)
}
