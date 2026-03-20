package model

type DictData struct {
	CommonModel
	Sort      int    `json:"sort" gorm:"size:4"`
	Label     string `json:"label" gorm:"size:100;not null"`
	LabelEn   string `json:"labelEn" gorm:"size:100;not null"`
	Value     string `json:"value" gorm:"size:100;not null"`
	DictType  string `json:"dictType" gorm:"size:11;not null"`
	IsDefault string `json:"isDefault" gorm:"size:1"`
	Status    string `json:"status" gorm:"size:1"`
	Remark    string `json:"remark" gorm:"size:500"`
}
