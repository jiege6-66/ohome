package dto

type QuarkTransferTaskListDTO struct {
	Status     string `json:"status" form:"status"`
	SourceType string `json:"sourceType" form:"sourceType"`
	Paginate
}
