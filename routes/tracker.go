package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
)

func RegisterTrackerRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/tracker")
	{
		user.GET("/announce", controllers.Announce)
		user.GET("/scrape", controllers.Scrape)
		// 可以继续加更多接口
	}
}
