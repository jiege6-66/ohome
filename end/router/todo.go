package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitTodoRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		_ = rgPublic
		todoApi := api.NewTodoApi()

		rg := rgAuth.Group("/todo")
		{
			rg.POST("/list", todoApi.GetList)
			rg.POST("/add", todoApi.Add)
			rg.PUT("/reorder", todoApi.Reorder)
			rg.PUT("/:id", todoApi.Update)
			rg.PUT("/:id/status", todoApi.UpdateStatus)
			rg.DELETE("/:id", todoApi.Delete)
		}
	})
}
