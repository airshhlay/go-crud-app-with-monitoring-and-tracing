package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Went through authenticate middleware")
	}
}
