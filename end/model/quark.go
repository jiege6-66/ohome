package model

import (
	"time"

	"gorm.io/gorm"
)

// QuarkConfig 表示 Quark 存储配置
type QuarkConfig struct {
	Application string `json:"application" gorm:"primaryKey;size:20;not null;default:'other'"`
	RootPath    string `json:"rootPath" gorm:"size:255"`
	Remark      string `json:"remark" gorm:"size:255"`

	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (QuarkConfig) TableName() string { return "sys_quark_config" }
