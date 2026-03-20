package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitQuarkTransferTaskRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		_ = rgPublic
		taskApi := api.NewQuarkTransferTaskApi()

		rg := rgAuth.Group("/quarkTransferTask")
		{
			rg.POST("/list", taskApi.GetList)
			rg.DELETE("/:id", taskApi.Delete)
		}
	})
}
