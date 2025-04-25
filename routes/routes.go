package routes

import "github.com/gin-gonic/gin"

func RegisterAPIRoutes(rg *gin.RouterGroup) {
	// 所有模块都注册到这里
	RegisterUserRoutes(rg)
	//RegisterTrackerRoutes(rg)
	RegisterTorrentRoutes(rg)
	// RegisterArticleRoutes(rg)
	// RegisterCommentRoutes(rg)
}
