package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.GET("/profile", controllers.Profile)
		user.POST("/logout", controllers.Logout)
		// 可以继续加更多接口
	}
}
