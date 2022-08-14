package controllers

import (
	"entry-task/gateway/constants"
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

	ot "github.com/opentracing/opentracing-go"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	loginHandler  = "gateway.LoginHandler"
	signupHandler = "gateway.SignupHandler"
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

	// start tracing span from context
	span := ot.SpanFromContext(c.Request.Context())
	u.addSpanTags(span, c)
	defer span.Finish()

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
		errorCodeStr = strconv.Itoa(constants.ErrorInvalidRequest)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorInvalidRequest, constants.ErrorInvalidRequestMsg)
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
	clientLoginRes, err := u.client.Login(c.Request.Context(), clientLoginReq)
	if err != nil {
		u.removeCookie(c, constants.Token)
		errorCodeStr = strconv.Itoa(constants.ErrorUserserviceConnection)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorUserserviceConnection, constants.ErrorUserserviceConnectionMsg)
		return
	}
	if clientLoginRes.ErrorCode != -1 {
		// remove any credentials if there is a login error
		u.removeCookie(c, constants.Token)
	}

	if clientLoginRes.UserID == 0 {
		// userID not created and returned by user service, unexpected error occured
		errorCodeStr = strconv.Itoa(constants.ErrorNoUserIDReturned)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorGenerateJWTToken, constants.ErrorNoUserIDReturnedMsg)
		return
	}

	// a userID was succesfully created by user service
	// generate the JWT token containing the userID
	tokenString, expirationTime, err := u.generateToken(clientLoginRes.UserID)
	if err != nil {
		// error occured during token generation
		u.logger.Error(
			constants.ErrorGenerateJWTTokenMsg,
			zap.Error(err),
		)
		u.removeCookie(c, constants.Token)
		errorCodeStr = strconv.Itoa(constants.ErrorGenerateJWTToken)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorGenerateJWTToken, constants.ErrorGenerateJWTTokenMsg)
		return
	}

	// successful token generation
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

	loginRes := res.LoginRes{
		ErrorCode: clientLoginRes.ErrorCode,
		ErrorMsg:  clientLoginRes.ErrorMsg,
	}

	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(loginRes.ErrorCode))
	// add resulting errorCode to span
	AddErrorTagsToSpan(span, loginRes.ErrorCode, loginRes.ErrorMsg)
	// return response
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
	// start tracing span from context
	span := ot.SpanFromContext(c.Request.Context())
	u.addSpanTags(span, c)
	defer span.Finish()

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
		errorCodeStr = strconv.Itoa(constants.ErrorBadRequest)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorInvalidRequest, constants.ErrorInvalidRequestMsg)
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
	clientLoginRes, err := u.client.Signup(c.Request.Context(), clientSignupReq)
	if err != nil {
		errorCodeStr = strconv.Itoa(constants.ErrorUserserviceConnection)
		u.removeCookie(c, constants.Token)
		// add the resulting error code to the span and send a standard gateway response back to the client
		SendStandardGatewayResponse(c, span, constants.ErrorUserserviceConnection, constants.ErrorUserserviceConnectionMsg)
		return
	}
	if clientLoginRes.ErrorCode != 0 {
		u.removeCookie(c, constants.Token)
	}

	// convert error code to string for metrics
	errorCodeStr = strconv.Itoa(int(clientLoginRes.ErrorCode))

	// add the resulting error code to the span
	AddErrorTagsToSpan(span, clientLoginRes.ErrorCode, clientLoginRes.ErrorMsg)
	// return response
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

func (u *UserServiceController) addSpanTags(span ot.Span, c *gin.Context) {
	// TODO: add additional tags if needed
}
