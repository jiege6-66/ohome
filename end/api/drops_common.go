package api

import "time"

func loadDropsLocation() (*time.Location, error) {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return time.Local, nil
	}
	return loc, nil
}
