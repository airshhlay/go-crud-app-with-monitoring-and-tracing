package routes

import (
	config "gateway/config"
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

// UserServiceRoutes defines routes used by the user service.
func UserServiceRoutes(g *gin.RouterGroup, controller *controllers.UserServiceController, apis *config.UserServiceAPIs) {
	g.POST(apis.Signup.Endpoint, controller.SignupHandler)
	g.POST(apis.Login.Endpoint, controller.LoginHandler)
}
