package task

import (
	"ohome/global"
	"ohome/service"
	"time"
)

var quarkAutoSaveTaskService service.QuarkAutoSaveTaskService

// 用 channel 作为轻量“TryLock”，避免定时器重入并发叠加
var quarkAutoSaveSem = make(chan struct{}, 1)

func StartQuarkAutoSaveScheduler() {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			<-ticker.C
			runQuarkAutoSaveOnce()
		}
	}()
}

func runQuarkAutoSaveOnce() {
	select {
	case quarkAutoSaveSem <- struct{}{}:
		defer func() { <-quarkAutoSaveSem }()
	default:
		return
	}

	tasks, err := quarkAutoSaveTaskService.GetEnabledTasks()
	if err != nil {
		if global.Logger != nil {
			global.Logger.Errorf("Quark AutoSave Scheduler List Error: %s", err.Error())
		}
		return
	}

	now := time.Now()
	for _, t := range tasks {
		if !quarkAutoSaveTaskService.ShouldRunNow(now, t) {
			continue
		}

		if global.Logger != nil {
			global.Logger.Infof("Quark AutoSave Task Submit Start: #%d %s", t.ID, t.TaskName)
		}
		transferTask, err := quarkAutoSaveTaskService.SubmitScheduled(t)
		if err != nil {
			if global.Logger != nil {
				global.Logger.Errorf("Quark AutoSave Task Submit Error: #%d %s", t.ID, err.Error())
			}
		} else {
			if global.Logger != nil {
				global.Logger.Infof("Quark AutoSave Task Submitted: #%d transferTask=%d %s", t.ID, transferTask.ID, t.TaskName)
			}
		}
	}
}
