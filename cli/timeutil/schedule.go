package timeutil

import "time"

// IsWeekend reports whether the given time falls on Saturday or Sunday.
func IsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

// IsWeekday reports whether the given time falls on Monday through Friday.
func IsWeekday(t time.Time) bool {
	w := t.Weekday()

	return w >= time.Monday && w <= time.Friday
}

// IsMorningAllowed reports whether the time's hour is before cutoffHour.
func IsMorningAllowed(t time.Time, cutoffHour int) bool {
	return t.Hour() < cutoffHour
}

// IsCommitWindowAllowed returns true if weekend, or if weekday and hour < 8.
func IsCommitWindowAllowed(t time.Time, loc *time.Location) bool {
	target := resolveTimeLocation(t, loc)
	if IsWeekend(target) {
		return true
	}

	return IsWeekday(target) && IsMorningAllowed(target, 8)
}

// CalculateNextWindow calculates the next valid slot (07:30 AM next weekday morning, or weekend).
func CalculateNextWindow(t time.Time, loc *time.Location) time.Time {
	target := resolveTimeLocation(t, loc)
	if IsWeekend(target) {
		return target
	}
	if IsWeekday(target) && IsMorningAllowed(target, 8) {
		return target
	}

	return nextValidSlot(target, target.Location())
}

// SimulateMorningWindow returns a compliant 07:30 AM timestamp for the day in loc.
func SimulateMorningWindow(t time.Time, loc *time.Location) time.Time {
	target := resolveTimeLocation(t, loc)

	return morningSlot(target, target.Location())
}

func resolveTimeLocation(t time.Time, loc *time.Location) time.Time {
	if loc == nil {
		return t.Local()
	}

	return t.In(loc)
}

func nextValidSlot(target time.Time, loc *time.Location) time.Time {
	nextDay := target.AddDate(0, 0, 1)
	if IsWeekend(nextDay) {
		return weekendSlot(nextDay, loc)
	}

	return morningSlot(nextDay, loc)
}

func weekendSlot(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
}

func morningSlot(t time.Time, loc *time.Location) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 7, 30, 0, 0, loc)
}
