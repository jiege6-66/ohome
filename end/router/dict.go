package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitDictRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		dictApi := api.NewDictApi()

		rgPublicDict := rgPublic.Group("/dict")
		{
			rgPublicDict.GET("/:id", dictApi.GetDictById)
			rgPublicDict.POST("/list", dictApi.GetDictTypeList)
			rgPublicDict.POST("/data_list", dictApi.GetDictDataList)
		}

		rgAuthDict := rgAuth.Group("/dict")
		{
			rgAuthDict.POST("/add", dictApi.UpdateDictType)
			rgAuthDict.DELETE("/:id", dictApi.DeleteDictType)
			rgAuthDict.GET("/getSysDictList", dictApi.GetSysDictList)
		}

		rgAuthDictData := rgAuth.Group("/dictData")
		{
			rgAuthDictData.POST("/add", dictApi.UpdateDictItem)
			rgAuthDictData.DELETE("/:id", dictApi.DeleteDictData)
		}

	})
}
