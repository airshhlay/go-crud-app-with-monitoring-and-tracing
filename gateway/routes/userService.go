package routes

import (
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

func UserServiceRoutes(g *gin.RouterGroup) {
	g.GET("/ping", controllers.Ping)
	g.POST("/signup", controllers.SignupHandler)
	g.POST("/login", controllers.LoginHandler)
}
