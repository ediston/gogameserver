package util

import (
   "time"
)

func CurrentDate() string {
    return GetDate(0)
}

func GetDate(daysOld int) string {
    t := time.Now().AddDate(0, 0,-daysOld)
	return t.Format("2006-01-02")
}