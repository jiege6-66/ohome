package router

import (
	"ohome/api"

	"github.com/gin-gonic/gin"
)

func InitQuarkLoginRoutes() {
	RegisterRouter(func(rgPublic *gin.RouterGroup, rgAuth *gin.RouterGroup) {
		quarkLoginApi := api.NewQuarkLoginApi()

		rgQuarkLogin := rgAuth.Group("/quarkLogin")
		{
			rgQuarkLogin.GET("/tv/status", quarkLoginApi.GetQuarkTVStatus)
			rgQuarkLogin.POST("/tv/start", quarkLoginApi.StartQuarkTVLogin)
			rgQuarkLogin.POST("/tv/poll", quarkLoginApi.PollQuarkTVLogin)
		}
	})
}
