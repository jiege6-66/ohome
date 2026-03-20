package dto

import "ohome/model"

// DictDataListDTO 配置列表相关DTO
type DictDataListDTO struct {
	Label    string `json:"label"`
	Status   string `json:"status"`
	DictType string `json:"dictType"`
	Paginate
}

type DictDataUpdateDTO struct {
	ID        uint   `json:"id" form:"id"`
	Sort      int    `json:"sort" form:"sort"`
	Label     string `json:"label" form:"label"`
	LabelEn   string `json:"labelEn" form:"labelEn"`
	Value     string `json:"value" form:"value"`
	DictType  string `json:"dictType" form:"dictType"`
	IsDefault string `json:"isDefault" form:"isDefault"`
	Status    string `json:"status" form:"status"`
	Remark    string `json:"remark" form:"remark"`
}

func (d *DictDataUpdateDTO) ConvertToModel(iDictData *model.DictData) {
	iDictData.ID = d.ID
	iDictData.Sort = d.Sort
	iDictData.LabelEn = d.LabelEn
	iDictData.Label = d.Label
	iDictData.Value = d.Value
	iDictData.DictType = d.DictType
	iDictData.IsDefault = d.IsDefault
	iDictData.Status = d.Status
	iDictData.Remark = d.Remark
}
