package model

import "time"

type TodoItem struct {
	CommonModel
	OwnerUserID uint       `json:"ownerUserId" gorm:"not null;default:0;index"`
	CreatedBy   uint       `json:"createdBy" gorm:"not null;default:0"`
	UpdatedBy   uint       `json:"updatedBy" gorm:"not null;default:0"`
	Title       string     `json:"title" gorm:"size:255;not null;default:''"`
	SortOrder   int64      `json:"sortOrder" gorm:"not null;default:0;index"`
	Completed   bool       `json:"completed" gorm:"not null;default:false;index"`
	CompletedAt *time.Time `json:"completedAt" gorm:"index"`
}

func (TodoItem) TableName() string { return "sys_todo_item" }
