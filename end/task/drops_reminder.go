package task

import (
	"ohome/global"
	"ohome/model"
	"ohome/service"
	"time"
)

var (
	dropsItemService  service.DropsItemService
	dropsEventService service.DropsEventService
	appMessageService service.AppMessageService
)

func StartDropsReminderScheduler() {
	runDropsReminderOnce()
	go func() {
		for {
			next := nextFullHour(time.Now())
			timer := time.NewTimer(time.Until(next))
			<-timer.C
			runDropsReminderOnce()
		}
	}()
}

func nextFullHour(now time.Time) time.Time {
	loc := now.Location()
	base := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	if !base.After(now) {
		base = base.Add(time.Hour)
	}
	return base
}

func runDropsReminderOnce() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.Local
	}
	now := time.Now().In(loc)
	itemMessages, err := dropsItemService.BuildReminderMessages(now, loc)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Errorf("Drops Item Reminder Error: %s", err.Error())
		}
		return
	}
	eventMessages, err := dropsEventService.BuildReminderMessages(now, loc)
	if err != nil {
		if global.Logger != nil {
			global.Logger.Errorf("Drops Event Reminder Error: %s", err.Error())
		}
		return
	}

	allMessages := make([]model.AppMessage, 0, len(itemMessages)+len(eventMessages))
	allMessages = append(allMessages, itemMessages...)
	allMessages = append(allMessages, eventMessages...)
	if len(allMessages) == 0 {
		return
	}
	if err := appMessageService.SaveMessages(allMessages); err != nil && global.Logger != nil {
		global.Logger.Errorf("Drops Reminder Save Error: %s", err.Error())
	}
}
