package api

import (
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type QuarkLogin struct {
	BaseApi
}

func NewQuarkLoginApi() QuarkLogin {
	return QuarkLogin{
		BaseApi: NewBaseApi(),
	}
}

func (a *QuarkLogin) GetQuarkTVStatus(c *gin.Context) {
	if _, ok := requireSuperAdmin(c); !ok {
		return
	}

	status, err := quarkTVLoginService.GetStatus()
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(status, c)
}

func (a *QuarkLogin) StartQuarkTVLogin(c *gin.Context) {
	if _, ok := requireSuperAdmin(c); !ok {
		return
	}

	result, err := quarkTVLoginService.StartLogin(c.Request.Context())
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(result, c)
}

func (a *QuarkLogin) PollQuarkTVLogin(c *gin.Context) {
	if _, ok := requireSuperAdmin(c); !ok {
		return
	}

	result, err := quarkTVLoginService.PollLogin(c.Request.Context())
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(result, c)
}
