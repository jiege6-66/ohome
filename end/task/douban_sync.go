package task

import (
	"context"
	"ohome/global"
	"ohome/service"
	"time"
)

var doubanService service.DoubanService

func StartDoubanSyncScheduler() {
	go func() {
		for {
			next := nextNoon(time.Now())
			timer := time.NewTimer(time.Until(next))
			<-timer.C

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
			if global.Logger != nil {
				global.Logger.Infof("Douban Sync Scheduler Start: %s", time.Now().Format(time.RFC3339))
			}
			if err := doubanService.SyncAllCache(ctx); err != nil && global.Logger != nil {
				global.Logger.Errorf("Douban Sync Scheduler Error: %s", err.Error())
			}
			if global.Logger != nil {
				global.Logger.Infof("Douban Sync Scheduler Done: %s", time.Now().Format(time.RFC3339))
			}
			cancel()
		}
	}()
}

func nextNoon(now time.Time) time.Time {
	loc := now.Location()
	noon := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, loc)
	if !noon.After(now) {
		noon = noon.Add(24 * time.Hour)
	}
	return noon
}
