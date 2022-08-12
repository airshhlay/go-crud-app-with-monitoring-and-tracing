package controllers

import (
	client "gateway/client"
	config "gateway/config"
	constants "gateway/constants"
	req "gateway/dto/request"
	res "gateway/dto/response"
	metrics "gateway/metrics"
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

type UserServiceController struct {
	config *config.UserServiceConfig
	logger *zap.Logger
	client *client.UserServiceClient
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func NewUserServiceController(config *config.UserServiceConfig, logger *zap.Logger, client *client.UserServiceClient) *UserServiceController {
	return &UserServiceController{
		config,
		logger,
		client,
	}
}

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
			constants.ERROR_BAD_REQUEST_MSG,
			zap.Error(err),
		)
		u.removeCookie(c, token)
		errorCodeInt = constants.ERROR_BAD_REQUEST
		c.JSON(200, res.GatewayResponse{ErrorCode: errorCodeInt})
		return
	}
	u.logger.Info(
		"info_request",
		zap.Any("username", loginReq.Username),
	)

	clientLoginReq := &proto.LoginReq{
		Username: loginReq.Username,
		Password: loginReq.Password,
	}

	// err if issue with client connection
	clientLoginRes, err := u.client.Login(c, clientLoginReq)
	if err != nil {
		u.removeCookie(c, token)
		errorCodeInt = constants.ERROR_USERSERVICE_CONNECTION
		c.JSON(200, res.GatewayResponse{ErrorCode: errorCodeInt})
		return
	}
	if clientLoginRes.ErrorCode != -1 {
		// remove any credentials if there is a login error
		u.removeCookie(c, token)
	}

	if clientLoginRes.UserId != 0 {
		// no error occured
		tokenString, expirationTime, err := u.generateToken(clientLoginRes.UserId)
		if err != nil {
			u.logger.Error(
				constants.ERROR_GENERATE_JWT_TOKEN_MSG,
				zap.Error(err),
			)
			u.removeCookie(c, token)
			errorCodeInt = constants.ERROR_GENERATE_JWT_TOKEN
			c.IndentedJSON(200, res.GatewayResponse{ErrorCode: errorCodeInt})
		} else {
			// set jwt token in cookie
			http.SetCookie(
				c.Writer, &http.Cookie{
					Name:     "token",
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
	u.logger.Info(
		"response",
		zap.Any("res", clientLoginRes),
	)
	errorCodeInt = loginRes.ErrorCode

	// errorCode, if any
	errorCodeStr = strconv.Itoa(int(errorCodeInt))
	c.IndentedJSON(200, loginRes)
}

func (u *UserServiceController) generateToken(userId int64) (string, time.Time, error) {
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		UserId: strconv.FormatInt(userId, 10),
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(u.config.Secret))
	return tokenString, expirationTime, err
}

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
	u.removeCookie(c, token)

	var signupReq req.SignupReq
	err := c.BindJSON(&signupReq)
	if err != nil {
		u.logger.Info(
			constants.ERROR_BAD_REQUEST_MSG,
			zap.Error(err),
		)
		u.removeCookie(c, token)
		errorCodeInt = constants.ERROR_BAD_REQUEST
		c.JSON(200, res.GatewayResponse{ErrorCode: constants.ERROR_BAD_REQUEST})
		return
	}
	u.logger.Info(
		"info_request",
		zap.Any("request", signupReq),
	)

	clientSignupReq := &proto.SignupReq{
		Username: signupReq.Username,
		Password: signupReq.Password,
	}
	// err if issue with client connection
	clientLoginRes, err := u.client.Signup(c, clientSignupReq)
	if err != nil {
		errorCodeInt = constants.ERROR_USERSERVICE_CONNECTION
		u.removeCookie(c, token)
		c.JSON(200, res.GatewayResponse{ErrorCode: errorCodeInt})
		return
	}
	if clientLoginRes.ErrorCode != 0 {
		u.removeCookie(c, token)
	}

	errorCodeInt = clientLoginRes.ErrorCode
	// errorCode, if any
	errorCodeStr = strconv.Itoa(int(errorCodeInt))
	c.JSON(200, clientLoginRes)
}

func (u *UserServiceController) removeCookie(c *gin.Context, cookieName string) {
	// set jwt token in cookie
	http.SetCookie(
		c.Writer, &http.Cookie{
			Name:   cookieName,
			MaxAge: -1,
		},
	)
}
