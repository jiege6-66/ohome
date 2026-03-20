package dto

import "ohome/model"

// DictTypeListDTO 配置列表相关DTO
type DictTypeListDTO struct {
	Name   string `json:"name"`
	Key    string `json:"key"`
	Status string `json:"status"`
	Paginate
}

type DictTypeUpdateDTO struct {
	ID     uint   `json:"id" form:"id"`
	Name   string `json:"name" form:"name"`
	Key    string `json:"key" form:"key"`
	Status string `json:"status" form:"status"`
	Remark string `json:"remark" form:"remark"`
}

func (d *DictTypeUpdateDTO) ConvertToModel(iDictType *model.DictType) {
	iDictType.ID = d.ID
	iDictType.Name = d.Name
	iDictType.Key = d.Key
	iDictType.Status = d.Status
	iDictType.Remark = d.Remark
}
