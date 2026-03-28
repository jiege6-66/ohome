package model

import (
	"ohome/utils"

	"gorm.io/gorm"
)

type User struct {
	CommonModel
	Name     string `json:"name" gorm:"size:64;not null;uniqueIndex:uk_sys_user_name"`
	RealName string `json:"realName" gorm:"size:128;"`
	Avatar   string `json:"avatar" gorm:"size:256"`
	Password string `json:"password" gorm:"size:128;not null;"`
	RoleID   uint   `json:"roleId" gorm:"not null;index:idx_sys_user_role_id"`
	Role     Role   `json:"role" gorm:"foreignKey:RoleID"`
}

func (m *User) BeforeCreate(orm *gorm.DB) error {
	return m.Encrypt()
}

func (m *User) Encrypt() error {
	stHash, err := utils.Encrypt(m.Password)
	if err == nil {
		m.Password = stHash
	}

	return err
}

// LoginUser ===============================================================================
// = 用户登录信息
type LoginUser struct {
	ID       uint
	Name     string
	RoleID   uint
	RoleCode string
}

func (m LoginUser) IsSuperAdmin() bool {
	return m.RoleCode == RoleCodeSuperAdmin
}
