package service

import (
	"ohome/model"
	"ohome/service/dto"
)

type ConfigService struct {
	BaseService
}

func (s *ConfigService) GetConfigById(iCommonIDDTO *dto.CommonIDDTO) (model.Config, error) {
	return configDao.GetFormById(iCommonIDDTO.ID)
}

func (s *ConfigService) GetConfigList(iConfigListDTO *dto.ConfigListDTO) ([]model.Config, int64, error) {
	return configDao.GetConfigList(iConfigListDTO)
}

func (s *ConfigService) AddOrUpdateConfig(d *dto.ConfigUpdateDTO) error {
	return configDao.AddOrUpdateConfig(d)
}

func (s *ConfigService) DeleteConfigById(iCommonIDDTO *dto.CommonIDDTO) error {
	return configDao.DeleteConfig(iCommonIDDTO.ID)
}

func (s *ConfigService) GetConfigsByKeys(keys []string) ([]model.Config, error) {
	return configDao.GetByKeys(keys)
}
