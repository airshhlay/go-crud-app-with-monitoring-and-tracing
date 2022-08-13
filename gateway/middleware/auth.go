package middleware

import (
	config "gateway/config"
	constants "gateway/constants"
	res "gateway/dto/response"
	metrics "gateway/metrics"
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

// Claims is a struct that will be encoded to a JWT.
// jwt.StandardClaims is added as an embedded type, to provide fields like expiry time.
type Claims struct {
	UserID string `json:"userID"`
	jwt.StandardClaims
}

// Authenticate middleware is called on relevant routes to retrieve the token cookie attached with the request and validate it using the jwt key.
// It
func Authenticate(secret string, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		validationSuccess := false // label used for metrics
		// observe latency
		timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
			metrics.AuthenticateDuration.WithLabelValues(strconv.FormatBool(validationSuccess)).Observe(v)
		}))
		defer func() {
			timer.ObserveDuration()
		}()

		cookie, err := c.Request.Cookie(constants.Token)
		if err != nil {
			if err == http.ErrNoCookie {
				// no cookie sent with the request
				logger.Info(constants.InfoNoCookieReceived, zap.Error(err))
				c.IndentedJSON(http.StatusUnauthorized, res.GatewayResponse{ErrorCode: constants.ErrorNoCookie})
				c.Abort()
				return
			}
			// other problem with the request
			c.IndentedJSON(http.StatusBadRequest, res.GatewayResponse{ErrorCode: constants.ErrorBadRequest})
			c.Abort()
			return
		}
		tokenString := cookie.Value
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				logger.Error(constants.ErrorJWTSignatureInvalidMsg, zap.Error(err))
				c.IndentedJSON(http.StatusUnauthorized, constants.Unauthorized)
				c.Abort()
				return
			}
			logger.Error(constants.ErrorUnexpectedJWTErr, zap.Error(err))
			c.IndentedJSON(http.StatusBadRequest, constants.BadRequest)
			c.Abort()
			return
		}

		validationSuccess = true

		if !token.Valid {
			logger.Info(constants.InfoInvalidTokenReceived, zap.String(constants.Token, tokenString))
			c.IndentedJSON(http.StatusUnauthorized, constants.Unauthorized)
			c.Abort()
			return
		}

		c.Set(constants.UserID, claims.UserID)
		c.Next()
	}
}

// CORSMiddleware enables cross origin resource sharing.
func CORSMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, v := range config.AllowedOrigins {
			c.Writer.Header().Set("Access-Control-Allow-Origin", v)
		}
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:80")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set(
			"Access-Control-Expose-Headers",
			"Set-Cookie",
		)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
