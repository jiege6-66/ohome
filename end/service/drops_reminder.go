package service

import (
	"ohome/global"
	"ohome/model"
	"sync"
	"time"
)

var dropsReminderRunMu sync.Mutex

func RunDropsReminderOnce(loc *time.Location) error {
	dropsReminderRunMu.Lock()
	defer dropsReminderRunMu.Unlock()

	resolvedLoc := resolveDropsReminderLocation(loc)
	now := time.Now().In(resolvedLoc)

	itemMessages, err := (&DropsItemService{}).BuildReminderMessages(now, resolvedLoc)
	if err != nil {
		return err
	}
	eventMessages, err := (&DropsEventService{}).BuildReminderMessages(now, resolvedLoc)
	if err != nil {
		return err
	}
	if len(itemMessages) == 0 && len(eventMessages) == 0 {
		return nil
	}

	allMessages := make([]model.AppMessage, 0, len(itemMessages)+len(eventMessages))
	allMessages = append(allMessages, itemMessages...)
	allMessages = append(allMessages, eventMessages...)
	return (&AppMessageService{}).SaveMessages(allMessages)
}

func TriggerDropsReminderRescan(loc *time.Location) {
	if err := RunDropsReminderOnce(loc); err != nil && global.Logger != nil {
		global.Logger.Errorf("Drops Reminder Trigger Error: %s", err.Error())
	}
}

func resolveDropsReminderLocation(loc *time.Location) *time.Location {
	if loc != nil {
		return loc
	}
	resolved, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Local
	}
	return resolved
}
