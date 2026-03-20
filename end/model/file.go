package model

type File struct {
	CommonModel
	Name         string `json:"name" gorm:"size:20;not null"`
	Type         string `json:"type" gorm:"size:10;not null"`
	Size         int64  `json:"size" gorm:"size:11;not null"`
	Status       string `json:"status" gorm:"size:1"`
	Url          string `json:"url" gorm:"size:255;not null"`
	UploaderId   uint   `json:"uploaderId" gorm:"size:11;not null"`
	UploaderName string `json:"uploaderName" gorm:"size:50;not null"`
	Description  string `json:"description" gorm:"size:255"`
}
