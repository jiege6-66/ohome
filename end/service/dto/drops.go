package dto

import "ohome/model"

type DropsItemListDTO struct {
	ScopeType string `json:"scopeType" form:"scopeType"`
	Category  string `json:"category" form:"category"`
	Keyword   string `json:"keyword" form:"keyword"`
	Paginate
}

type DropsItemUpsertDTO struct {
	ID        uint   `json:"id" form:"id" uri:"id"`
	ScopeType string `json:"scopeType" form:"scopeType"`
	Category  string `json:"category" form:"category" binding:"required" required_err:"物资分类不能为空"`
	Name      string `json:"name" form:"name" binding:"required" required_err:"物资名称不能为空"`
	Location  string `json:"location" form:"location"`
	ExpireAt  string `json:"expireAt" form:"expireAt"`
	Remark    string `json:"remark" form:"remark"`
	Enabled   *bool  `json:"enabled" form:"enabled"`
}

func (d *DropsItemUpsertDTO) ApplyToModel(item *model.DropsItem) {
	item.ScopeType = d.ScopeType
	item.Category = d.Category
	item.Name = d.Name
	item.Location = d.Location
	item.Remark = d.Remark
}

type DropsPhotoDeleteDTO struct {
	ID      uint `json:"id" form:"id" uri:"id"`
	PhotoID uint `json:"photoId" form:"photoId" uri:"photoId"`
}

type DropsLocationSuggestDTO struct {
	Keyword string `json:"keyword" form:"keyword"`
}

type DropsEventListDTO struct {
	ScopeType string `json:"scopeType" form:"scopeType"`
	EventType string `json:"eventType" form:"eventType"`
	Month     int    `json:"month" form:"month"`
	Keyword   string `json:"keyword" form:"keyword"`
	Paginate
}

type DropsEventUpsertDTO struct {
	ID           uint   `json:"id" form:"id" uri:"id"`
	ScopeType    string `json:"scopeType" form:"scopeType"`
	Title        string `json:"title" form:"title" binding:"required" required_err:"重要日期标题不能为空"`
	EventType    string `json:"eventType" form:"eventType" binding:"required" required_err:"重要日期类型不能为空"`
	CalendarType string `json:"calendarType" form:"calendarType"`
	EventYear    int    `json:"eventYear" form:"eventYear"`
	EventMonth   int    `json:"eventMonth" form:"eventMonth"`
	EventDay     int    `json:"eventDay" form:"eventDay"`
	IsLeapMonth  bool   `json:"isLeapMonth" form:"isLeapMonth"`
	RepeatYearly *bool  `json:"repeatYearly" form:"repeatYearly"`
	Remark       string `json:"remark" form:"remark"`
	Enabled      *bool  `json:"enabled" form:"enabled"`
}

func (d *DropsEventUpsertDTO) ApplyToModel(event *model.DropsEvent) {
	event.ScopeType = d.ScopeType
	event.Title = d.Title
	event.EventType = d.EventType
	event.CalendarType = d.CalendarType
	event.EventYear = d.EventYear
	event.EventMonth = d.EventMonth
	event.EventDay = d.EventDay
	event.IsLeapMonth = d.IsLeapMonth
	event.Remark = d.Remark
}
