package utils

import "time"

func InTimeSpan(start, end, check time.Time) bool {
	//case end does not expire
	if end.IsZero() {
		return true
	}
	return check.After(start) && check.Before(end)
}
