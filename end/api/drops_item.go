package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type DropsItem struct {
	BaseApi
}

func NewDropsItemApi() DropsItem {
	return DropsItem{BaseApi: NewBaseApi()}
}

func (a *DropsItem) GetOverview(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	loc, _ := loadDropsLocation()
	overview, err := dropsItemService.GetOverview(loginUser, loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(overview, c)
}

func (a *DropsItem) GetList(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var listDTO dto.DropsItemListDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &listDTO}).GetErrors(); err != nil {
		return
	}
	items, total, err := dropsItemService.GetList(&listDTO, loginUser, false)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(gin.H{"records": items, "total": total}, c)
}

func (a *DropsItem) GetByID(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	item, err := dropsItemService.GetByID(idDTO.ID, loginUser)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(item, c)
}

func (a *DropsItem) Create(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var createDTO dto.DropsItemUpsertDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &createDTO}).GetErrors(); err != nil {
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		utils.FailWithMessage("请上传物资照片", c)
		return
	}
	loc, _ := loadDropsLocation()
	item, err := dropsItemService.Create(loginUser, &createDTO, form.File["photos"], loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(item, c)
}

func (a *DropsItem) Update(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var updateDTO dto.DropsItemUpsertDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &updateDTO}).GetErrors(); err != nil {
		return
	}
	loc, _ := loadDropsLocation()
	item, err := dropsItemService.Update(updateDTO.ID, loginUser, &updateDTO, loc)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(item, c)
}

func (a *DropsItem) Delete(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	if err := dropsItemService.Delete(idDTO.ID, loginUser); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功", c)
}

func (a *DropsItem) AddPhotos(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var idDTO dto.CommonIDDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &idDTO}).GetErrors(); err != nil {
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		utils.FailWithMessage("请上传物资照片", c)
		return
	}
	item, err := dropsItemService.AddPhotos(idDTO.ID, loginUser, form.File["photos"])
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(item, c)
}

func (a *DropsItem) DeletePhoto(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var photoDTO dto.DropsPhotoDeleteDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &photoDTO}).GetErrors(); err != nil {
		return
	}
	if err := dropsItemService.DeletePhoto(photoDTO.ID, photoDTO.PhotoID, loginUser); err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.Ok(c)
}

func (a *DropsItem) SuggestLocations(c *gin.Context) {
	loginUser, err := getLoginUser(c)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	var suggestDTO dto.DropsLocationSuggestDTO
	if err := a.Request(RequestOptions{Ctx: c, DTO: &suggestDTO}).GetErrors(); err != nil {
		return
	}
	locations, err := dropsItemService.SuggestLocations(suggestDTO.Keyword, loginUser)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithData(locations, c)
}
