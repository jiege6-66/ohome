package model

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UserMediaHistory struct {
	CommonModel
	UserID          uint              `json:"userId" gorm:"not null;index:idx_user_media_history_recent,priority:1"`
	ApplicationType string            `json:"applicationType" gorm:"size:32;not null;index:idx_user_media_history_recent,priority:2"`
	FolderPath      string            `json:"folderPath" gorm:"type:text;not null"`
	ItemTitle       string            `json:"itemTitle" gorm:"type:text;not null"`
	ItemPath        string            `json:"itemPath" gorm:"type:text"`
	PositionMs      int               `json:"positionMs" gorm:"not null;default:0"`
	DurationMs      int               `json:"durationMs" gorm:"not null;default:0"`
	CoverURL        string            `json:"coverUrl" gorm:"type:text"`
	Extra           datatypes.JSONMap `json:"extra" gorm:"type:json"`
	LastPlayedAt    time.Time         `json:"lastPlayedAt" gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_user_media_history_recent,priority:3,sort:desc"`
	UniqueKey       string            `json:"-" gorm:"size:40;uniqueIndex:uk_user_media_history_unique_key"`
}

func (UserMediaHistory) TableName() string {
	return "user_media_history"
}

func (h *UserMediaHistory) NormalizePlaybackFields() {
	if h == nil {
		return
	}

	h.ApplicationType = strings.TrimSpace(h.ApplicationType)
	h.FolderPath = normalizeMediaPath(h.FolderPath)
	if h.FolderPath == "" {
		h.FolderPath = "/"
	}
	h.ItemPath = normalizeMediaPath(h.ItemPath)
	if h.ItemPath == "" && h.Extra != nil {
		if value, ok := h.Extra["itemPath"]; ok {
			if path, ok := value.(string); ok {
				h.ItemPath = normalizeMediaPath(path)
			}
		}
	}
	if h.ItemPath == "" {
		h.ItemPath = joinMediaPath(h.FolderPath, h.ItemTitle)
	}

	if strings.TrimSpace(h.ItemTitle) == "" {
		h.ItemTitle = titleFromMediaPath(h.ItemPath)
	}

	if h.ItemPath != "" {
		if h.Extra == nil {
			h.Extra = datatypes.JSONMap{}
		}
		h.Extra["itemPath"] = h.ItemPath
	}
}

func (h *UserMediaHistory) PrepareForSave() {
	if h == nil {
		return
	}

	h.NormalizePlaybackFields()
	h.UniqueKey = buildUserMediaHistoryUniqueKey(h.UserID, h.ApplicationType, h.FolderPath)
}

func (h *UserMediaHistory) BeforeSave(_ *gorm.DB) error {
	h.PrepareForSave()
	return nil
}

func buildUserMediaHistoryUniqueKey(userID uint, applicationType, folderPath string) string {
	raw := fmt.Sprintf("%d|%s|%s", userID, strings.TrimSpace(applicationType), normalizeMediaPath(folderPath))
	sum := sha1.Sum([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func normalizeMediaPath(value string) string {
	normalized := strings.TrimSpace(strings.ReplaceAll(value, "\\", "/"))
	if normalized == "" {
		return ""
	}

	rawParts := strings.Split(normalized, "/")
	parts := make([]string, 0, len(rawParts))
	for _, part := range rawParts {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		parts = append(parts, trimmed)
	}
	if len(parts) == 0 {
		if strings.HasPrefix(normalized, "/") {
			return "/"
		}
		return ""
	}
	return "/" + strings.Join(parts, "/")
}

func joinMediaPath(folderPath, itemName string) string {
	folder := strings.TrimSpace(folderPath)
	item := strings.TrimSpace(itemName)
	if folder == "" {
		return normalizeMediaPath(item)
	}
	if item == "" {
		return normalizeMediaPath(folder)
	}
	if strings.HasSuffix(folder, "/") {
		return normalizeMediaPath(folder + item)
	}
	return normalizeMediaPath(folder + "/" + item)
}

func titleFromMediaPath(path string) string {
	normalized := normalizeMediaPath(path)
	if normalized == "" {
		return ""
	}
	parts := strings.Split(strings.Trim(normalized, "/"), "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
