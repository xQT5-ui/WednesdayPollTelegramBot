package lib

import "time"

func DayOfWeek() int {
	return int(time.Now().Weekday())
}
