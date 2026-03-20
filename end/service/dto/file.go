package dto

import "ohome/model"

// FileUploadDTO 上传文件DTO
type FileUploadDTO struct {
	Type        string `json:"type" form:"type" binding:"required" required_err:"文件类型不能为空"`
	Description string `json:"description" form:"description"`
}

// FileAddDTO 添加File实体类DTO
type FileAddDTO struct {
	Name         string `json:"name" gorm:"size:20;not null"`
	Type         string `json:"type" gorm:"size:10;not null"`
	Size         int64  `json:"size" gorm:"size:11;not null"`
	Url          string `json:"url" gorm:"size:255;not null"`
	Description  string `json:"description" gorm:"size:255"`
	UploaderId   uint   `json:"uploaderId" gorm:"size:11;not null"`
	UploaderName string `json:"uploaderName" gorm:"size:50;not null"`
}

func (d FileAddDTO) ConvertToModel(m *model.File) {
	m.UploaderId = d.UploaderId
	m.UploaderName = d.UploaderName
	m.Name = d.Name
	m.Type = d.Type
	m.Description = d.Description
	m.Url = d.Url
	m.Size = d.Size
}

// FileListDTO 文件列表查询 DTO
type FileListDTO struct {
	Name       string `json:"name" form:"name"`
	Type       string `json:"type" form:"type"`
	Status     string `json:"status" form:"status"`
	UploaderId uint   `json:"uploaderId" form:"uploaderId"`
	Paginate
}

// FileUpdateDTO 更新文件 DTO
type FileUpdateDTO struct {
	// 注意：BaseApi.Request 会先执行 ShouldBind，再执行 ShouldBindUri。
	// 如果这里加 binding:"required"，在仅通过路径传 id 的 PUT 请求里会先报校验错误。
	ID           uint    `json:"id" form:"id" uri:"id"`
	Name         *string `json:"name" form:"name"`
	Type         *string `json:"type" form:"type"`
	Status       *string `json:"status" form:"status"`
	UploaderName *string `json:"uploaderName" form:"uploaderName"`
	Description  *string `json:"description" form:"description"`
}

func (d FileUpdateDTO) ApplyToModel(m *model.File) {
	m.ID = d.ID

	if d.Name != nil {
		m.Name = *d.Name
	}
	if d.Type != nil {
		m.Type = *d.Type
	}
	if d.Status != nil {
		m.Status = *d.Status
	}
	if d.UploaderName != nil {
		m.UploaderName = *d.UploaderName
	}
	if d.Description != nil {
		m.Description = *d.Description
	}
}
