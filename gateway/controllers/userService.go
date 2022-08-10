package controllers

import (
	client "gateway/client"
	config "gateway/config"
	req "gateway/dto/request"
	res "gateway/dto/response"
	proto "gateway/proto"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

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
	var loginReq req.LoginReq
	err := c.BindJSON(&loginReq)
	if err != nil {
		u.logger.Error(
			"error_request_binding",
			zap.Error(err),
		)
		u.removeCookie(c, token)
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_login_request",
		})
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
		c.JSON(500, "server_error")
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
				"error_jwt_token",
				zap.Error(err),
			)
			u.removeCookie(c, token)
			c.IndentedJSON(500, "server_error")
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
	var signupReq req.SignupReq
	err := c.BindJSON(&signupReq)
	if err != nil {
		u.logger.Error(
			"error_request_binding",
			zap.Error(err),
		)
		u.removeCookie(c, token)
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_signup_request",
		})
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
		u.removeCookie(c, token)
		c.JSON(500, "server_error")
		return
	}
	if clientLoginRes.ErrorCode != 0 {
		u.removeCookie(c, token)
	}

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
