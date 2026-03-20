package model

import "time"

const (
	DropsScopeShared   = "shared"
	DropsScopePersonal = "personal"

	DropsItemCategoryKitchen  = "kitchen"
	DropsItemCategoryFood     = "food"
	DropsItemCategoryMedicine = "medicine"
	DropsItemCategoryClothing = "clothing"
	DropsItemCategoryOther    = "other"

	DropsCalendarSolar = "solar"
	DropsCalendarLunar = "lunar"

	DropsEventTypeBirthday    = "birthday"
	DropsEventTypeAnniversary = "anniversary"
	DropsEventTypeCustom      = "custom"

	DropsBizTypeItem  = "item"
	DropsBizTypeEvent = "event"
)

type DropsItem struct {
	CommonModel
	ScopeType    string           `json:"scopeType" gorm:"size:20;not null;default:'shared';index"`
	OwnerUserID  uint             `json:"ownerUserId" gorm:"not null;default:0;index"`
	CreatedBy    uint             `json:"createdBy" gorm:"not null;default:0"`
	UpdatedBy    uint             `json:"updatedBy" gorm:"not null;default:0"`
	Category     string           `json:"category" gorm:"size:32;not null;index"`
	Name         string           `json:"name" gorm:"size:120;not null;index"`
	Location     string           `json:"location" gorm:"size:255;index"`
	ExpireAt     *time.Time       `json:"expireAt" gorm:"index"`
	Remark       string           `json:"remark" gorm:"size:500"`
	ReminderDays string           `json:"-" gorm:"size:64;not null;default:'7,3,1,0'"`
	Enabled      bool             `json:"enabled" gorm:"not null;default:true;index"`
	CoverURL     string           `json:"coverUrl" gorm:"type:text"`
	PhotoCount   int              `json:"photoCount" gorm:"not null;default:0"`
	Photos       []DropsItemPhoto `json:"photos,omitempty" gorm:"foreignKey:ItemID"`
}

func (DropsItem) TableName() string { return "drops_item" }

type DropsItemPhoto struct {
	CommonModel
	ItemID   uint   `json:"itemId" gorm:"not null;index"`
	Sort     int    `json:"sort" gorm:"not null;default:1"`
	FileName string `json:"fileName" gorm:"size:255;not null"`
	FilePath string `json:"filePath" gorm:"type:text;not null"`
	URL      string `json:"url" gorm:"type:text;not null"`
	Size     int64  `json:"size" gorm:"not null;default:0"`
	IsCover  bool   `json:"isCover" gorm:"not null;default:false"`
}

func (DropsItemPhoto) TableName() string { return "drops_item_photo" }

type DropsEvent struct {
	CommonModel
	ScopeType    string     `json:"scopeType" gorm:"size:20;not null;default:'shared';index"`
	OwnerUserID  uint       `json:"ownerUserId" gorm:"not null;default:0;index"`
	CreatedBy    uint       `json:"createdBy" gorm:"not null;default:0"`
	UpdatedBy    uint       `json:"updatedBy" gorm:"not null;default:0"`
	Title        string     `json:"title" gorm:"size:120;not null;index"`
	EventType    string     `json:"eventType" gorm:"size:32;not null;index"`
	CalendarType string     `json:"calendarType" gorm:"size:16;not null;default:'solar';index"`
	EventYear    int        `json:"eventYear" gorm:"not null;default:0"`
	EventMonth   int        `json:"eventMonth" gorm:"not null;default:1"`
	EventDay     int        `json:"eventDay" gorm:"not null;default:1"`
	IsLeapMonth  bool       `json:"isLeapMonth" gorm:"not null;default:false"`
	RepeatYearly bool       `json:"repeatYearly" gorm:"not null;default:true"`
	Remark       string     `json:"remark" gorm:"size:500"`
	ReminderDays string     `json:"-" gorm:"size:64;not null;default:'7,3,1,0'"`
	Enabled      bool       `json:"enabled" gorm:"not null;default:true;index"`
	NextOccurAt  *time.Time `json:"nextOccurAt,omitempty" gorm:"-"`
}

func (DropsEvent) TableName() string { return "drops_event" }
