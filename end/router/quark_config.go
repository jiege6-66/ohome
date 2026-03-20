package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitQuarkConfigRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		quarkApi := api.NewQuarkConfigApi()

		rgQuark := rgAuth.Group("/quarkConfig")
		{
			rgQuark.POST("/list", quarkApi.GetQuarkConfigList)
			rgQuark.GET("/:application", quarkApi.GetQuarkConfigByApplication)
			rgQuark.POST("", quarkApi.CreateQuarkConfig)
			rgQuark.PUT("/:application", quarkApi.UpdateQuarkConfig)
			rgQuark.DELETE("/:application", quarkApi.DeleteQuarkConfig)
		}
	})
}
