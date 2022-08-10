package routes

import (
	config "gateway/config"
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

func ItemServiceRoutes(g *gin.RouterGroup, controller *controllers.ItemServiceController, apis *config.ItemServiceApis) {
	g.POST(apis.AddFav.Endpoint, controller.AddFavHandler)
	g.GET(apis.GetFavList.Endpoint, controller.GetFavListHandler)
	g.DELETE(apis.DeleteFav.Endpoint, controller.DeleteFavHandler)
}
