package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitFileRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		fileApi := api.NewFileApi()

		rgAuthFile := rgAuth.Group("/file")
		{
			rgAuthFile.POST("/upload", fileApi.Upload)
			rgAuthFile.GET("/:id", fileApi.GetFileById)
			rgAuthFile.POST("/list", fileApi.GetFileList)
			rgAuthFile.PUT("/:id", fileApi.UpdateFile)
			rgAuthFile.DELETE("/:id", fileApi.DeleteFile)
		}
	})
}
