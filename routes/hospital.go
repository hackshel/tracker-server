package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
)

func RegisterHospitalRoutes(rg *gin.RouterGroup) {
	user := rg.Group("/hospital")
	{
		user.POST("/list", controllers.ListHospital)
		user.POST("/delete", controllers.DeleteHospital)
		// user.POST("/modify", controllers.ModifyUser)
		user.POST("/add", controllers.AddHospital)
		// 可以继续加更多接口
	}
}
