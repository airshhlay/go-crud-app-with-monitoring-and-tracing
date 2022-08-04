package controllers

import (
	"fmt"

	req "gateway/dto/request"

	"github.com/gin-gonic/gin"
)

func LoginHandler(c *gin.Context) {
	fmt.Println("Login handler called")

	var loginRequest req.LoginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		fmt.Println("Error with login request body")
		c.JSON(200, gin.H{
			"errorCode": 400,
			"errorMsg":  "invalid_login_request",
		})
		return
	}

	fmt.Println(loginRequest.Username)
	c.IndentedJSON(200, loginRequest)
}

func SignupHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "signup",
	})
}
