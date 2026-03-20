package model

type Config struct {
	CommonModel
	Name string `json:"name" gorm:"size:100;not null"`
	Key  string `json:"key" gorm:"size:100;not null"`
	// Value 可能存放较长内容（例如 quark cookies），使用 text 避免长度不足
	Value  string `json:"value" gorm:"type:text;not null"`
	IsLock string `json:"isLock" gorm:"size:1;not null"`
	Remark string `json:"remark" gorm:"size:255"`
}
