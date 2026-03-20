package dao

import (
	"ohome/global"
	"ohome/model"
	"ohome/service/dto"
	"strings"
)

type ConfigDao struct {
	BaseDao
}

func (m *ConfigDao) GetFormById(id uint) (model.Config, error) {
	var iConfig model.Config
	err := global.DB.First(&iConfig, id).Error
	return iConfig, err
}

func (m *ConfigDao) GetByKey(key string) (model.Config, error) {
	var iConfig model.Config
	err := global.DB.Where("`key` = ?", key).First(&iConfig).Error
	return iConfig, err
}

func (m *ConfigDao) GetConfigList(iConfigListDTO *dto.ConfigListDTO) ([]model.Config, int64, error) {

	var giConfigList []model.Config
	var nTotal int64

	query := global.DB.Model(&model.Config{}).
		Scopes(Paginate(iConfigListDTO.Paginate))

	if iConfigListDTO.Name != "" {
		query = query.Where("name like ?", "%"+iConfigListDTO.Name+"%")
	}

	if iConfigListDTO.Key != "" {
		query = query.Where("`key` = ?", iConfigListDTO.Key)
	}

	if len(iConfigListDTO.Keys) > 0 {
		keys := make([]string, 0, len(iConfigListDTO.Keys))
		for _, key := range iConfigListDTO.Keys {
			key = strings.TrimSpace(key)
			if key != "" {
				keys = append(keys, key)
			}
		}
		if len(keys) > 0 {
			query = query.Where("`key` IN ?", keys)
		}
	}

	if iConfigListDTO.IsLock != "" {
		query = query.Where("`is_lock` = ?", iConfigListDTO.IsLock)
	}

	if len(iConfigListDTO.ExcludeKeys) > 0 {
		query = query.Where("`key` NOT IN ?", iConfigListDTO.ExcludeKeys)
	}

	err := query.Find(&giConfigList).Error
	if err != nil {
		return giConfigList, 0, err
	}
	err = query.Offset(-1).Limit(-1).Count(&nTotal).Error

	return giConfigList, nTotal, err

}

func (m *ConfigDao) AddOrUpdateConfig(d *dto.ConfigUpdateDTO) error {
	var iConfig model.Config

	if d.ID != 0 {
		global.DB.First(&iConfig, d.ID)
	}

	d.ConvertToModel(&iConfig)

	return global.DB.Save(&iConfig).Error
}

func (m *ConfigDao) DeleteConfig(id uint) error {
	return global.DB.Delete(&model.Config{}, id).Error
}

func (m *ConfigDao) GetByKeys(keys []string) ([]model.Config, error) {
	normalized := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			normalized = append(normalized, key)
		}
	}
	if len(normalized) == 0 {
		return []model.Config{}, nil
	}

	var configs []model.Config
	err := global.DB.Where("`key` IN ?", normalized).Find(&configs).Error
	return configs, err
}
