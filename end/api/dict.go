package api

import (
	"ohome/service/dto"
	"ohome/utils"

	"github.com/gin-gonic/gin"
)

type Dict struct {
	BaseApi
}

func NewDictApi() Dict {
	return Dict{
		BaseApi: NewBaseApi(),
	}
}

func (dict *Dict) GetDictById(c *gin.Context) {
	var iCommonIDDTO dto.CommonIDDTO
	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iCommonIDDTO}).GetErrors(); err != nil {
		return
	}

	iDictTypes, err := dictService.GetDictById(&iCommonIDDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(iDictTypes, c)
}

func (dict *Dict) GetDictTypeList(c *gin.Context) {
	var iDictTypeListDTO dto.DictTypeListDTO
	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iDictTypeListDTO}).GetErrors(); err != nil {
		return
	}

	giDictTypeList, nTotal, err := dictService.GetDictTypeList(&iDictTypeListDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(gin.H{
		"records": giDictTypeList,
		"total":   nTotal,
	}, c)
}

func (dict *Dict) GetDictDataList(c *gin.Context) {
	var iDictDataListDTO dto.DictDataListDTO
	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iDictDataListDTO}).GetErrors(); err != nil {
		return
	}

	giDictDataListDTO, nTotal, err := dictService.GetDataDataList(&iDictDataListDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(gin.H{
		"records": giDictDataListDTO,
		"total":   nTotal,
	}, c)
}

func (dict *Dict) UpdateDictType(c *gin.Context) {
	var iDictTypeUpdateDTO dto.DictTypeUpdateDTO

	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iDictTypeUpdateDTO}).GetErrors(); err != nil {
		return
	}

	err := dictService.AddDictType(&iDictTypeUpdateDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.Ok(c)
}

func (dict *Dict) DeleteDictType(c *gin.Context) {
	var iCommonIDDTO dto.CommonIDDTO
	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iCommonIDDTO}).GetErrors(); err != nil {
		return
	}

	err := dictService.DeleteDictTypeById(&iCommonIDDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功!", c)
}

func (dict *Dict) GetSysDictList(c *gin.Context) {
	dictType, dictData, err := dictService.GetSysDictList()
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.OkWithData(gin.H{
		"dictType": dictType,
		"dictData": dictData,
	}, c)
}

func (dict *Dict) UpdateDictItem(c *gin.Context) {
	var iDictDataUpdateDTO dto.DictDataUpdateDTO

	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iDictDataUpdateDTO}).GetErrors(); err != nil {
		return
	}

	err := dictService.AddDictData(&iDictDataUpdateDTO)

	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}

	utils.Ok(c)
}

func (dict *Dict) DeleteDictData(c *gin.Context) {
	var iCommonIDDTO dto.CommonIDDTO
	if err := dict.Request(RequestOptions{Ctx: c, DTO: &iCommonIDDTO}).GetErrors(); err != nil {
		return
	}

	err := dictService.DeleteDictDataById(&iCommonIDDTO)
	if err != nil {
		utils.FailWithMessage(err.Error(), c)
		return
	}
	utils.OkWithMessage("删除成功!", c)
}
