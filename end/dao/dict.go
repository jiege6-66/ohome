package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
)

type DictDao struct {
	BaseDao
}

func (d DictDao) GetDictById(id uint) (model.DictType, error) {
	var iDictType model.DictType
	err := global.DB.First(&iDictType, id).Error
	return iDictType, err
}

func (d DictDao) GetDictTypeList(iDictTypeListDTO *dto.DictTypeListDTO) ([]model.DictType, int64, error) {
	var giDictTypeList []model.DictType
	var nTotal int64

	query := global.DB.Model(&model.DictType{}).
		Scopes(Paginate(iDictTypeListDTO.Paginate))

	if iDictTypeListDTO.Name != "" {
		query = query.Where("name like ?", "%"+iDictTypeListDTO.Name+"%")
	}

	if iDictTypeListDTO.Key != "" {
		query = query.Where("`key` = ?", iDictTypeListDTO.Key)
	}

	if iDictTypeListDTO.Status != "" {
		query = query.Where("`status` = ?", iDictTypeListDTO.Status)
	}

	query.
		Order("id asc").
		Find(&giDictTypeList).
		Offset(-1).Limit(-1).
		Count(&nTotal)

	return giDictTypeList, nTotal, query.Error
}

func (d DictDao) GetDataTypeList(iDictDataListDTO *dto.DictDataListDTO) ([]model.DictData, int64, error) {
	var giDictDataList []model.DictData
	var nTotal int64

	query := global.DB.Model(&model.DictData{}).
		Scopes(Paginate(iDictDataListDTO.Paginate))

	if iDictDataListDTO.Label != "" {
		query = query.Where("label like ?", "%"+iDictDataListDTO.Label+"%")
	}

	if iDictDataListDTO.Status != "" {
		query = query.Where("`status` = ?", iDictDataListDTO.Status)
	}

	if iDictDataListDTO.DictType != "" {
		query = query.Where("`dict_type` = ?", iDictDataListDTO.DictType)
	}

	query.
		Order("sort asc").
		Order("id asc").
		Find(&giDictDataList).
		Offset(-1).Limit(-1).
		Count(&nTotal)

	return giDictDataList, nTotal, query.Error
}

func (d DictDao) AddDictType(iDictTypeUpdateDTO *dto.DictTypeUpdateDTO) error {
	var iDictType model.DictType

	if iDictTypeUpdateDTO.ID != 0 {
		global.DB.First(&iDictType, iDictTypeUpdateDTO.ID)
	}

	iDictTypeUpdateDTO.ConvertToModel(&iDictType)

	return global.DB.Save(&iDictType).Error
}

func (d DictDao) DeleteDictType(id uint) error {
	return global.DB.Delete(&model.DictType{}, id).Error
}

func (d DictDao) GetSysDictList() ([]model.DictType, []model.DictData, error) {
	var allDictType []model.DictType
	var allDictItem []model.DictData

	err := global.DB.Find(&allDictType).Error
	err = global.DB.Find(&allDictItem).Error

	return allDictType, allDictItem, err
}

func (d DictDao) AddDictData(updateDTO *dto.DictDataUpdateDTO) error {
	var iDictData model.DictData

	if updateDTO.ID != 0 {
		global.DB.First(&iDictData, updateDTO.ID)
	}

	updateDTO.ConvertToModel(&iDictData)

	return global.DB.Save(&iDictData).Error
}

func (d DictDao) DeleteDictData(id uint) error {
	return global.DB.Delete(&model.DictData{}, id).Error
}

func (d DictDao) HasDictValue(dictType string, value string) (bool, error) {
	var total int64
	err := global.DB.Model(&model.DictData{}).
		Where("dict_type = ? AND value = ?", dictType, value).
		Count(&total).Error

	return total > 0, err
}
