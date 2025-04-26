package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
)

func RegisterUserRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/user")
	{
		user.GET("/profile", controllers.Profile)
		//user.POST("/logout", controllers.Logout)
		user.POST("/list", controllers.ListUser)
		user.POST("/delete", controllers.DeleteUser)
		user.POST("/modify", controllers.ModifyUser)
		user.POST("/add", controllers.AddUser)
		// 可以继续加更多接口
	}
}
