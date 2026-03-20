package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitUserMediaHistoryRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		historyApi := api.NewUserMediaHistoryApi()

		rgHistory := rgAuth.Group("/userMediaHistory")
		{
			rgHistory.POST("/list", historyApi.GetUserMediaHistoryList)
			rgHistory.POST("/byFolder", historyApi.GetUserMediaHistoryByFolder)
			rgHistory.POST("/recent", historyApi.GetRecentUserMediaHistory)
			rgHistory.GET("/:id", historyApi.GetUserMediaHistoryByID)
			rgHistory.POST("", historyApi.CreateUserMediaHistory)
			rgHistory.PUT("/:id", historyApi.UpdateUserMediaHistory)
			rgHistory.DELETE("/:id", historyApi.DeleteUserMediaHistory)
		}
	})
}
