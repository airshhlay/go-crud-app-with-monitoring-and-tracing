package middleware

import (
	"fmt"
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
				c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}
		tokenString := cookie.Value
		claims := &Claims{}
		fmt.Println(tokenString)

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
				c.Abort()
				return
			}
			c.IndentedJSON(http.StatusBadRequest, "bad request")
			c.Abort()
			return
		}

		if !token.Valid {
			c.IndentedJSON(http.StatusUnauthorized, "unauthorised")
			c.Abort()
			return
		}

		c.Set("userId", claims.UserId)
		c.Next()
	}
}
