package model

import "time"

// QuarkAutoSaveTask 夸克网盘自动转存任务（精简版：仅支持 daily/weekly + 执行时间）
type QuarkAutoSaveTask struct {
	CommonModel

	TaskName string `json:"taskName" gorm:"size:100;not null"`
	ShareURL string `json:"shareUrl" gorm:"size:1000;not null"`
	SavePath string `json:"savePath" gorm:"size:255;not null"`

	// 同步规则：daily / weekly
	ScheduleType string `json:"scheduleType" gorm:"size:10;not null;default:''"`
	// 执行时间：HH:mm
	RunTime string `json:"runTime" gorm:"size:5;not null;default:''"`
	// weekly 模式下：1-7 逗号分隔（周一=1，周日=7）
	RunWeek string `json:"runWeek" gorm:"size:20"`

	Enabled bool `json:"enabled" gorm:"not null;default:true"`

	// 用于避免同一分钟重复执行
	LastRunAt *time.Time `json:"lastRunAt"`
}

func (QuarkAutoSaveTask) TableName() string { return "sys_quark_auto_save_task" }
