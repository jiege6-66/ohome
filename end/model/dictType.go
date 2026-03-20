package model

type DictType struct {
	CommonModel
	Name   string `json:"name" gorm:"size:100;not null"`
	Key    string `json:"key" gorm:"size:100;not null"`
	Status string `json:"status" gorm:"size:1"`
	Remark string `json:"remark" gorm:"size:500"`
}
