package middleware

import (
	"fmt"
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

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

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

		cookie, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				// no cookie sent with the request
				logger.Info(constants.ERROR_NO_COOKIE_MSG, zap.Error(err))
				c.IndentedJSON(http.StatusUnauthorized, res.GatewayResponse{ErrorCode: constants.ERROR_NO_COOKIE})
				c.Abort()
				return
			}
			// other problem with the request
			c.IndentedJSON(http.StatusBadRequest, res.GatewayResponse{ErrorCode: constants.ERROR_BAD_REQUEST})
			c.Abort()
			return
		}
		tokenString := cookie.Value
		logger.Info("info_received_token", zap.String("token", tokenString))
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				fmt.Println("jwt signature invalid")
				fmt.Println(err)
				c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
				c.Abort()
				return
			}
			fmt.Println("another error occured: ", err)
			c.IndentedJSON(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}

		validationSuccess = true

		if !token.Valid {
			fmt.Println("token invalid")
			c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}

func CORSMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		for _, v := range config.AllowedOrigins {
			c.Writer.Header().Set("Access-Control-Allow-Origin", v)
		}
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:80")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set(
			"Access-Control-Expose-Headers",
			"Set-Cookie",
		)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
