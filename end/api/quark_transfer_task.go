package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type QuarkTransferTask struct {
	BaseApi
}

func NewQuarkTransferTaskApi() QuarkTransferTask {
	return QuarkTransferTask{BaseApi: NewBaseApi()}
}

func (a *QuarkTransferTask) GetList(c *gin.Context) {
	var listDTO dto.QuarkTransferTaskListDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &listDTO}).GetErrors(); err != nil {
		return
	}
	records, total, err := quarkTransferTaskService.GetList(&listDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(gin.H{"records": records, "total": total}, c)
}

func (a *QuarkTransferTask) Delete(c *gin.Context) {
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	if err := quarkTransferTaskService.DeleteByID(&idDTO); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}
