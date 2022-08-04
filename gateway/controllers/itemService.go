package controllers

import (
	"github.com/gin-gonic/gin"
)

func AddItemHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Add item",
	})
}

func DeleteItemHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Delete item",
	})
}

func GetListHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Get list",
	})
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{"message": "ping"})
}
