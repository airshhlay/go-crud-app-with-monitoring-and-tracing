package routes

import (
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

func ItemServiceRoutes(g *gin.RouterGroup) {
	g.GET("/ping", controllers.Ping)
	g.POST("/add", controllers.AddItemHandler)
	g.GET("/get/list", controllers.GetListHandler)
	g.DELETE("/delete", controllers.DeleteItemHandler)
}
