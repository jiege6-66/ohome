package api

import (
	"net/http"

	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type PansouApi struct {
	BaseApi
}

func NewPansouApi() PansouApi {
	return PansouApi{BaseApi: NewBaseApi()}
}

func (a *PansouApi) Search(c *gin.Context) {
	var req dto.PansouSearchReq

	if c.Request.Method == http.MethodGet {
		var q dto.PansouSearchQueryDTO
		if err := a.Request(RequestOptions{Ctx: c, DTO: &q}).GetErrors(); err != nil {
			return
		}
		parsed, err := q.ToReq()
		if err != nil {
			utils.FailWithMessage(err.Error(), c)
			return
		}
		req = parsed
	} else {
		if err := a.Request(RequestOptions{Ctx: c, DTO: &req}).GetErrors(); err != nil {
			return
		}
	}

	data, message, err := pansouService.Search(c.Request.Context(), &req)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	if message != "" {
		utils.OkWithDetailed(data, message, c)
		return
	}
	utils.OkWithData(data, c)
}
