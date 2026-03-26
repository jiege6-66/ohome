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
	keysToInvalidate := make([]string, 0, 2)
	if d.ID != 0 {
		current, err := configDao.GetFormById(d.ID)
		if err == nil {
			keysToInvalidate = append(keysToInvalidate, current.Key)
		}
	}
	keysToInvalidate = append(keysToInvalidate, d.Key)

	if err := configDao.AddOrUpdateConfig(d); err != nil {
		return err
	}
	invalidateRuntimeConfigByKeys(keysToInvalidate...)
	return nil
}

func (s *ConfigService) DeleteConfigById(iCommonIDDTO *dto.CommonIDDTO) error {
	current, err := configDao.GetFormById(iCommonIDDTO.ID)
	if err != nil {
		return err
	}
	if err := configDao.DeleteConfig(iCommonIDDTO.ID); err != nil {
		return err
	}
	invalidateRuntimeConfigByKeys(current.Key)
	return nil
}

func (s *ConfigService) GetConfigsByKeys(keys []string) ([]model.Config, error) {
	return configDao.GetByKeys(keys)
}

func invalidateRuntimeConfigByKeys(keys ...string) {
	for _, key := range keys {
		if isQuarkStreamConfigKey(key) {
			invalidateQuarkStreamConfigCache()
			return
		}
	}
}
