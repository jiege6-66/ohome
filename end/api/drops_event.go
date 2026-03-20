package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type DropsEvent struct {
	BaseApi
}

func NewDropsEventApi() DropsEvent {
	return DropsEvent{BaseApi: NewBaseApi()}
}

func (a *DropsEvent) GetList(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var listDTO dto.DropsEventListDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &listDTO}).GetErrors(); err != nil {
		return
	}
	loc, _ := loadDropsLocation()
	events, total, err := dropsEventService.GetList(&listDTO, loginUser, false, loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(gin.H{"records": events, "total": total}, c)
}

func (a *DropsEvent) GetByID(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	loc, _ := loadDropsLocation()
	event, err := dropsEventService.GetByID(idDTO.ID, loginUser, loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(event, c)
}

func (a *DropsEvent) AddOrUpdate(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var updateDTO dto.DropsEventUpsertDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &updateDTO}).GetErrors(); err != nil {
		return
	}
	loc, _ := loadDropsLocation()
	event, err := dropsEventService.AddOrUpdate(&updateDTO, loginUser, loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(event, c)
}

func (a *DropsEvent) Delete(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	loc, _ := loadDropsLocation()
	if err := dropsEventService.Delete(idDTO.ID, loginUser, loc); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}
