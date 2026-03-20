package service

import (
	"ohome/model"
	"ohome/service/dto"
)

type DictService struct {
	BaseService
}

func (s *DictService) GetDictById(iCommonIDDTO *dto.CommonIDDTO) (model.DictType, error) {
	return dictDao.GetDictById(iCommonIDDTO.ID)
}

func (s *DictService) GetDictTypeList(iDictTypeListDTO *dto.DictTypeListDTO) ([]model.DictType, int64, error) {
	return dictDao.GetDictTypeList(iDictTypeListDTO)
}

func (s *DictService) GetDataDataList(iDictDataListDTO *dto.DictDataListDTO) ([]model.DictData, int64, error) {
	return dictDao.GetDataTypeList(iDictDataListDTO)
}

func (s *DictService) AddDictType(iDictTypeUpdateDTO *dto.DictTypeUpdateDTO) error {
	return dictDao.AddDictType(iDictTypeUpdateDTO)
}

func (s *DictService) DeleteDictTypeById(iCommonIDDTO *dto.CommonIDDTO) error {
	return dictDao.DeleteDictType(iCommonIDDTO.ID)
}

func (s *DictService) GetSysDictList() ([]model.DictType, []model.DictData, error) {
	return dictDao.GetSysDictList()
}

func (s *DictService) AddDictData(iDictDataUpdateDTO *dto.DictDataUpdateDTO) error {
	return dictDao.AddDictData(iDictDataUpdateDTO)
}

func (s *DictService) DeleteDictDataById(iCommonIDDTO *dto.CommonIDDTO) error {
	return dictDao.DeleteDictData(iCommonIDDTO.ID)
}
