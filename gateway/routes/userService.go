package routes

import (
	config "gateway/config"
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

func UserServiceRoutes(g *gin.RouterGroup, controller *controllers.UserServiceController, apis *config.UserServiceApis) {
	g.POST(apis.Signup.Endpoint, controller.SignupHandler)
	g.POST(apis.Login.Endpoint, controller.LoginHandler)
}
