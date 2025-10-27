package utils

import "time"

func getBangkokTimezone() (*time.Location, error) {
	return time.LoadLocation("Asia/Bangkok")
}

func BangkokNow() time.Time {
	location, _ := getBangkokTimezone()
	return time.Now().In(location)
}

func ToStringDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05+07:00")
}

func ToStringDate(t time.Time) string {
	return t.Format("2006-01-02+07:00")
}
