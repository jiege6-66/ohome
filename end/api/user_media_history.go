package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type UserMediaHistoryApi struct {
	BaseApi
}

func NewUserMediaHistoryApi() UserMediaHistoryApi {
	return UserMediaHistoryApi{
		BaseApi: NewBaseApi(),
	}
}

func (a *UserMediaHistoryApi) GetUserMediaHistoryList(c *gin.Context) {
	var listDTO dto.UserMediaHistoryListDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &listDTO}).GetErrors(); err != nil {
		return
	}

	records, total, err := userMediaHistoryService.GetList(&listDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(gin.H{
		"records": records,
		"total":   total,
	}, c)
}

func (a *UserMediaHistoryApi) GetUserMediaHistoryByID(c *gin.Context) {
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}

	record, err := userMediaHistoryService.GetByID(&idDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(record, c)
}

func (a *UserMediaHistoryApi) CreateUserMediaHistory(c *gin.Context) {
	var createDTO dto.UserMediaHistoryCreateDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &createDTO}).GetErrors(); err != nil {
		return
	}

	record, err := userMediaHistoryService.Create(&createDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(record, c)
}

func (a *UserMediaHistoryApi) UpdateUserMediaHistory(c *gin.Context) {
	var updateDTO dto.UserMediaHistoryUpdateDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &updateDTO}).GetErrors(); err != nil {
		return
	}

	if err := userMediaHistoryService.Update(&updateDTO); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.Ok(c)
}

func (a *UserMediaHistoryApi) DeleteUserMediaHistory(c *gin.Context) {
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}

	if err := userMediaHistoryService.Delete(&idDTO); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithMessage("删除成功", c)
}

func (a *UserMediaHistoryApi) GetRecentUserMediaHistory(c *gin.Context) {
	var recentDTO dto.UserMediaHistoryRecentDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &recentDTO}).GetErrors(); err != nil {
		return
	}

	record, err := userMediaHistoryService.GetRecent(&recentDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(record, c)
}

func (a *UserMediaHistoryApi) GetUserMediaHistoryByFolder(c *gin.Context) {
	var folderDTO dto.UserMediaHistoryFolderDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &folderDTO}).GetErrors(); err != nil {
		return
	}

	record, err := userMediaHistoryService.GetByFolder(&folderDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(record, c)
}
