package timeutil

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	_ "time/tzdata"
)

// ResolveLocation resolves IANA names, "UTC", "Local", and numeric offsets (+08:00, -0400).
func ResolveLocation(raw string) (*time.Location, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" || strings.EqualFold(trimmed, "local") {
		return time.Local, nil
	}
	if strings.EqualFold(trimmed, "utc") {
		return time.UTC, nil
	}
	if hasOffsetSign(trimmed) {
		return parseNumericOffset(trimmed)
	}

	return time.LoadLocation(trimmed)
}

func hasOffsetSign(s string) bool {
	return strings.HasPrefix(s, "+") || strings.HasPrefix(s, "-")
}

func parseNumericOffset(raw string) (*time.Location, error) {
	isNegative := strings.HasPrefix(raw, "-")
	numPart := strings.TrimPrefix(strings.TrimPrefix(raw, "+"), "-")
	hours, mins, err := splitOffsetHoursMins(numPart)
	if err != nil {
		return nil, err
	}
	sec := calculateOffsetSeconds(hours, mins, isNegative)

	return time.FixedZone(raw, sec), nil
}

func splitOffsetHoursMins(numPart string) (int, int, error) {
	if strings.Contains(numPart, ":") {
		return parseColonOffset(numPart)
	}

	return parsePackedOffset(numPart)
}

func parseColonOffset(s string) (int, int, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid offset format: %s", s)
	}

	return parseHourMinStrings(parts[0], parts[1])
}

func parsePackedOffset(s string) (int, int, error) {
	if len(s) == 4 {
		return parseHourMinStrings(s[:2], s[2:])
	}
	if len(s) == 1 || len(s) == 2 {
		return parseHourMinStrings(s, "0")
	}

	return 0, 0, fmt.Errorf("invalid packed offset: %s", s)
}

func parseHourMinStrings(hourStr, minStr string) (int, int, error) {
	h, errH := strconv.Atoi(hourStr)
	m, errM := strconv.Atoi(minStr)
	if errH != nil || errM != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, 0, fmt.Errorf("invalid hour/min in offset: %s:%s", hourStr, minStr)
	}

	return h, m, nil
}

func calculateOffsetSeconds(hours, mins int, isNegative bool) int {
	totalSec := hours*3600 + mins*60
	if isNegative {
		return -totalSec
	}

	return totalSec
}
