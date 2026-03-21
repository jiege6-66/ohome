package task

import (
	"ohome/service"
	"time"
)

func StartDropsReminderScheduler() {
	service.TriggerDropsReminderRescan(nil)
	go func() {
		for {
			next := nextFullHour(time.Now())
			timer := time.NewTimer(time.Until(next))
			<-timer.C
			service.TriggerDropsReminderRescan(nil)
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
