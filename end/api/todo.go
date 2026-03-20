package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	BaseApi
}

func NewTodoApi() Todo {
	return Todo{BaseApi: NewBaseApi()}
}

func (a *Todo) GetList(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var listDTO dto.TodoListDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &listDTO}).GetErrors(); err != nil {
		return
	}

	records, total, err := todoService.GetList(&listDTO, loginUser.ID)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(gin.H{
		"records": records,
		"total":   total,
	}, c)
}

func (a *Todo) Add(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var createDTO dto.TodoItemCreateDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &createDTO}).GetErrors(); err != nil {
		return
	}

	item, err := todoService.Create(&createDTO, loginUser)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(item, c)
}

func (a *Todo) Update(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var updateDTO dto.TodoItemUpdateDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &updateDTO}).GetErrors(); err != nil {
		return
	}

	item, err := todoService.Update(&updateDTO, loginUser)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(item, c)
}

func (a *Todo) UpdateStatus(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var statusDTO dto.TodoItemStatusDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &statusDTO}).GetErrors(); err != nil {
		return
	}

	item, err := todoService.UpdateStatus(&statusDTO, loginUser)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(item, c)
}

func (a *Todo) Delete(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}

	if err := todoService.Delete(&idDTO, loginUser); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithMessage("删除成功", c)
}

func (a *Todo) Reorder(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	var reorderDTO dto.TodoReorderDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &reorderDTO}).GetErrors(); err != nil {
		return
	}

	if err := todoService.Reorder(&reorderDTO, loginUser); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.Ok(c)
}
