package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitDropsEventRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		_ = rgPublic
		eventApi := api.NewDropsEventApi()
		rg := rgAuth.Group("")
		{
			rg.POST("/dropsEvent/list", eventApi.GetList)
			rg.GET("/dropsEvent/:id", eventApi.GetByID)
			rg.POST("/dropsEvent/add", eventApi.AddOrUpdate)
			rg.PUT("/dropsEvent/add", eventApi.AddOrUpdate)
			rg.PUT("/dropsEvent/:id", eventApi.AddOrUpdate)
			rg.DELETE("/dropsEvent/:id", eventApi.Delete)
		}
	})
}
