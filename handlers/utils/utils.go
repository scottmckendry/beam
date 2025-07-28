package utils

import (
	"time"
)

// NextBillingDate calculates the next billing date for a subscription given start date and cadence.
func NextBillingDate(start time.Time, cadence string, end *time.Time) time.Time {
	now := time.Now()
	d := start

	var add func(time.Time) time.Time
	switch cadence {
	case "monthly":
		add = func(t time.Time) time.Time { return t.AddDate(0, 1, 0) }
	case "yearly":
		add = func(t time.Time) time.Time { return t.AddDate(1, 0, 0) }
	default:
		add = func(t time.Time) time.Time { return t.AddDate(0, 1, 0) }
	}

	for d.Before(now) {
		d = add(d)
	}

	if end != nil && !end.IsZero() && d.After(*end) {
		return time.Time{} // zero value, means no next billing
	}
	return d
}
