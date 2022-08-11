package middleware

import (
	"fmt"
	config "gateway/config"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func Authenticate(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Went through authenticate middleware")
		cookie, err := c.Request.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				fmt.Println("No cookie")
				c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
				c.Abort()
				return
			}
			fmt.Println(err)
			c.IndentedJSON(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}
		tokenString := cookie.Value
		// tokenString, _ := c.Cookie("token")
		fmt.Println("token: ", tokenString)
		claims := &Claims{}
		fmt.Println(tokenString)

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
