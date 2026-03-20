package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitDiscoveryRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		discoveryApi := api.NewDiscoveryApi()

		rgPublic.GET("/discovery", func(context *gin.Context) {
			discoveryApi.GetInfo(context)
		})
	})
}
