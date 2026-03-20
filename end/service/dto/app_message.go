package dto

type AppMessageListDTO struct {
	Source   string `json:"source" form:"source"`
	ReadOnly *bool  `json:"readOnly" form:"readOnly"`
	Paginate
}

type AppMessageReadDTO struct {
	ID uint `json:"id" form:"id"`
}
