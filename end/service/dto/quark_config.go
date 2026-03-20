package dto

import "ohome/model"

type QuarkConfigListDTO struct {
	Application string `json:"application" form:"application"`
	Paginate
}

type QuarkConfigCreateDTO struct {
	Application string `json:"application" form:"application" binding:"required" required_err:"应用不能为空"`
	RootPath    string `json:"rootPath" form:"rootPath"`
	Remark      string `json:"remark" form:"remark"`
}

func (d *QuarkConfigCreateDTO) ConvertToModel(m *model.QuarkConfig) {
	if d.Application != "" {
		m.Application = d.Application
	} else {
		m.Application = "other"
	}
	m.RootPath = d.RootPath
	m.Remark = d.Remark
}

type QuarkConfigUpdateDTO struct {
	QuarkApplicationDTO
	RootPath *string `json:"rootPath" form:"rootPath"`
	Remark   *string `json:"remark" form:"remark"`
}

func (d *QuarkConfigUpdateDTO) ApplyToModel(m *model.QuarkConfig) {
	if d.RootPath != nil {
		m.RootPath = *d.RootPath
	}
	if d.Remark != nil {
		m.Remark = *d.Remark
	}
}
