package utils

import "time"

func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

func TimeInt64Format(sec int64) string {
	return time.Unix(sec, 0).Format("2006-01-02 15:04:05")
}

func TimeString2Time(value string) time.Time {
	d, _ := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
	return d
}
