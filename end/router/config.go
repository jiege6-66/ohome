package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitConfigRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		configApi := api.NewConfigApi()

		rgAuthConfig := rgAuth.Group("/config")
		{
			rgAuthConfig.GET("/:id", configApi.GetConfigById)
			rgAuthConfig.POST("/list", configApi.GetConfigList)
			rgAuthConfig.PUT("/add", configApi.AddOrUpdateConfig)
			rgAuthConfig.DELETE("/:id", configApi.DeleteConfig)
		}

	})
}
