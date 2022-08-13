package routes

import (
	config "gateway/config"
	controllers "gateway/controllers"

	"github.com/gin-gonic/gin"
)

// ItemServiceRoutes define routes used by the item service.
func ItemServiceRoutes(g *gin.RouterGroup, controller *controllers.ItemServiceController, apis *config.ItemServiceAPIs) {
	g.POST(apis.AddFav.Endpoint, controller.AddFavHandler)
	g.GET(apis.GetFavList.Endpoint, controller.GetFavListHandler)
	g.DELETE(apis.DeleteFav.Endpoint, controller.DeleteFavHandler)
}
