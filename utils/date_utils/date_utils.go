package date_utils

import "time"

func GetTime() time.Time {
	return time.Now().UTC()
}

func GetFormattedTime() string {
	time := GetTime()
	return time.Format("2006-02-01T15:04:05Z")
}

func GetDbFormattedTime() string {
	time := GetTime()
	return time.Format("2006-01-02 15:04:05")
}
