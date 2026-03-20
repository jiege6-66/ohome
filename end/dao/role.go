package dao

import (
	"ohome/global"
	"ohome/model"
)

type RoleDao struct {
	BaseDao
}

func (d *RoleDao) GetByCode(code string) (model.Role, error) {
	var role model.Role
	err := global.DB.Where("code = ?", code).First(&role).Error
	return role, err
}
