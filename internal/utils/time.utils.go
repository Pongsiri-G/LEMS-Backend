package utils

import "time"

func getBangkokTimezone() (*time.Location, error) {
	return time.LoadLocation("Asia/Bangkok")
}

func BangkokNow() time.Time {
	location, _ := getBangkokTimezone()
	return time.Now().In(location)
}
