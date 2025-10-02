package timeutil

import "time"

func BangkokNow() time.Time {
	location, _ := time.LoadLocation("Asia/Bangkok")

	return time.Now().In(location)
}
