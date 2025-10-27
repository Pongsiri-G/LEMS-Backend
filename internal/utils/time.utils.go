package utils

import "time"

func getBangkokTimezone() (*time.Location, error) {
	return time.LoadLocation("Asia/Bangkok")
}

func BangkokNow() time.Time {
	// location, _ := getBangkokTimezone()
	return time.Now().UTC()
}

func ToStringDateTime(t time.Time) string {
	return t.Add(7 * time.Hour).Format("2006-01-02 15:04:05")
}

func ToStringDate(t time.Time) string {
	return t.Add(7 * time.Hour).Format("2006-01-02")
}
