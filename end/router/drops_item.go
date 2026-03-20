package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitDropsItemRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		_ = rgPublic
		dropsApi := api.NewDropsItemApi()
		rg := rgAuth.Group("")
		{
			rg.GET("/dropsOverview", dropsApi.GetOverview)
			rg.GET("/dropsLocation/suggestions", dropsApi.SuggestLocations)
			rg.POST("/dropsItem/list", dropsApi.GetList)
			rg.GET("/dropsItem/:id", dropsApi.GetByID)
			rg.POST("/dropsItem/create", dropsApi.Create)
			rg.PUT("/dropsItem/:id", dropsApi.Update)
			rg.DELETE("/dropsItem/:id", dropsApi.Delete)
			rg.POST("/dropsItem/:id/photos", dropsApi.AddPhotos)
			rg.DELETE("/dropsItem/:id/photos/:photoId", dropsApi.DeletePhoto)
		}
	})
}
