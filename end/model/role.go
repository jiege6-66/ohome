package model

import "strings"

const (
	RoleCodeSuperAdmin = "super_admin"
	RoleCodeUser       = "user"
)

type Role struct {
	CommonModel
	Name   string `json:"name" gorm:"size:32;not null;uniqueIndex:uk_sys_role_name"`
	Code   string `json:"code" gorm:"size:32;not null;uniqueIndex:uk_sys_role_code"`
	Remark string `json:"remark" gorm:"size:255"`
}

func (r Role) IsSuperAdmin() bool {
	return strings.TrimSpace(r.Code) == RoleCodeSuperAdmin
}
