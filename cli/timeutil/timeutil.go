// Package timeutil provides standardized UTC time helpers.
package timeutil

import "time"

// NowUTCISO returns the current UTC timestamp formatted as RFC3339.
func NowUTCISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// NowUTC returns the current UTC time.
func NowUTC() time.Time {
	return time.Now().UTC()
}

// FormatISO formats a timestamp to RFC3339 in UTC.
func FormatISO(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// FormatUTCISO formats a timestamp to RFC3339 in UTC.
func FormatUTCISO(t time.Time) string {
	return FormatISO(t)
}

// ParseISO parses an RFC3339 timestamp string into time.Time.
func ParseISO(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
