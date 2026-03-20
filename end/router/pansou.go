package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitPansouRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		pansouApi := api.NewPansouApi()

		rg := rgAuth.Group("/pansou")
		{
			rg.GET("/search", pansouApi.Search)
			rg.POST("/search", pansouApi.Search)
		}
	})
}
