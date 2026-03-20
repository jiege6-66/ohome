package dto

import (
	"ohome/model"
	"time"

	"gorm.io/datatypes"
)

// UserMediaHistoryCreateDTO defines the payload for creating media history records
type UserMediaHistoryCreateDTO struct {
	UserID          uint           `json:"userId" form:"userId" binding:"required"`
	ApplicationType string         `json:"applicationType" form:"applicationType" binding:"required"`
	FolderPath      string         `json:"folderPath" form:"folderPath" binding:"required"`
	ItemTitle       string         `json:"itemTitle" form:"itemTitle" binding:"required"`
	ItemPath        string         `json:"itemPath" form:"itemPath"`
	PositionMs      int            `json:"positionMs" form:"positionMs"`
	DurationMs      int            `json:"durationMs" form:"durationMs"`
	CoverURL        string         `json:"coverUrl" form:"coverUrl"`
	Extra           map[string]any `json:"extra" form:"extra"`
	LastPlayedAt    string         `json:"lastPlayedAt" form:"lastPlayedAt"`
}

func (d UserMediaHistoryCreateDTO) FillModel(m *model.UserMediaHistory, lastPlayedAt time.Time) {
	m.UserID = d.UserID
	m.ApplicationType = d.ApplicationType
	m.FolderPath = d.FolderPath
	m.ItemTitle = d.ItemTitle
	m.ItemPath = d.ItemPath
	m.PositionMs = d.PositionMs
	m.DurationMs = d.DurationMs
	m.CoverURL = d.CoverURL
	if d.Extra != nil {
		m.Extra = datatypes.JSONMap(d.Extra)
	} else {
		m.Extra = nil
	}
	m.LastPlayedAt = lastPlayedAt
}

// UserMediaHistoryUpdateDTO defines mutable fields for an existing record
type UserMediaHistoryUpdateDTO struct {
	// 注意：BaseApi.Request 会先执行 ShouldBind，再执行 ShouldBindUri。
	// 如果这里加 binding:"required"，在仅通过路径传 id 的 PUT 请求里会先报校验错误。
	ID              uint            `json:"id" form:"id" uri:"id"`
	UserID          *uint           `json:"userId" form:"userId"`
	ApplicationType *string         `json:"applicationType" form:"applicationType"`
	FolderPath      *string         `json:"folderPath" form:"folderPath"`
	ItemTitle       *string         `json:"itemTitle" form:"itemTitle"`
	ItemPath        *string         `json:"itemPath" form:"itemPath"`
	PositionMs      *int            `json:"positionMs" form:"positionMs"`
	DurationMs      *int            `json:"durationMs" form:"durationMs"`
	CoverURL        *string         `json:"coverUrl" form:"coverUrl"`
	Extra           *map[string]any `json:"extra" form:"extra"`
	LastPlayedAt    *string         `json:"lastPlayedAt" form:"lastPlayedAt"`
}

func (d UserMediaHistoryUpdateDTO) ApplyToModel(m *model.UserMediaHistory, lastPlayedAt *time.Time) {
	if d.UserID != nil {
		m.UserID = *d.UserID
	}
	if d.ApplicationType != nil {
		m.ApplicationType = *d.ApplicationType
	}
	if d.FolderPath != nil {
		m.FolderPath = *d.FolderPath
	}
	if d.ItemTitle != nil {
		m.ItemTitle = *d.ItemTitle
	}
	if d.ItemPath != nil {
		m.ItemPath = *d.ItemPath
	}
	if d.PositionMs != nil {
		m.PositionMs = *d.PositionMs
	}
	if d.DurationMs != nil {
		m.DurationMs = *d.DurationMs
	}
	if d.CoverURL != nil {
		m.CoverURL = *d.CoverURL
	}
	if d.Extra != nil {
		if *d.Extra == nil {
			m.Extra = nil
		} else {
			m.Extra = datatypes.JSONMap(*d.Extra)
		}
	}
	if lastPlayedAt != nil {
		m.LastPlayedAt = *lastPlayedAt
	}
}

// UserMediaHistoryListDTO carries pagination/query filters
type UserMediaHistoryListDTO struct {
	UserID          uint   `json:"userId" form:"userId"`
	ApplicationType string `json:"applicationType" form:"applicationType"`
	FolderPath      string `json:"folderPath" form:"folderPath"`
	ItemTitle       string `json:"itemTitle" form:"itemTitle"`
	Paginate
}

type UserMediaHistoryRecentDTO struct {
	UserID          uint   `json:"userId" form:"userId" binding:"required"`
	ApplicationType string `json:"applicationType" form:"applicationType"`
}

type UserMediaHistoryFolderDTO struct {
	UserID          uint   `json:"userId" form:"userId" binding:"required"`
	ApplicationType string `json:"applicationType" form:"applicationType" binding:"required"`
	FolderPath      string `json:"folderPath" form:"folderPath" binding:"required"`
}
