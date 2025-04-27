package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
)

func RegisterTorrentRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/torrent")
	{
		user.GET("/list", controllers.TorrentList)
		user.POST("/info", controllers.TorrenInfo)
		user.GET("/download", controllers.TorrenDownload)
		user.POST("/upload", controllers.TorrentUpload)
		user.GET("/count", controllers.TorrentsCount)
		user.POST("/delete", controllers.TorrentDelete)

		// 可以继续加更多接口
	}
}
