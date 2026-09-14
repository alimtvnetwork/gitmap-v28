package cmdschedule

import (
	"strconv"
	"strings"
)

// TryParseRunInterval parses "run <interval-spec>" args and advances index.
func TryParseRunInterval(args []string, idx *int, opts *scheduleAddOpts) bool {
	a := strings.ToLower(args[*idx])
	if a != "run" && a != "--run" {
		return false
	}

	if *idx+1 >= len(args) {
		return true
	}

	*idx++
	next := strings.ToLower(args[*idx])
	if isMatched := matchFixedIntervalKeyword(next, opts); isMatched {
		return true
	}

	return parseParameterizedEvery(args, idx, opts)
}

func matchFixedIntervalKeyword(kw string, opts *scheduleAddOpts) bool {
	if isStartup := matchStartupIntervalKeyword(kw, opts); isStartup {
		return true
	}

	return matchRecurringIntervalKeyword(kw, opts)
}

func matchStartupIntervalKeyword(kw string, opts *scheduleAddOpts) bool {
	switch kw {
	case "startup", "s":
		opts.IsStartup = true
		opts.Interval = "startup"
		return true
	case "startup-once", "so":
		opts.IsStartup = true
		opts.Interval = "startup-once"
		return true
	case "startup-weekly", "sw":
		opts.IsStartup = true
		opts.Interval = "startup-weekly"
		return true
	case "startup-monthly", "sm":
		opts.IsStartup = true
		opts.Interval = "startup-monthly"
		return true
	case "startup-yearly", "sy":
		opts.IsStartup = true
		opts.Interval = "startup-yearly"
		return true
	}
	return false
}

func matchRecurringIntervalKeyword(kw string, opts *scheduleAddOpts) bool {
	switch kw {
	case "daily", "d":
		opts.Interval = "1d"
		return true
	case "weekly", "w":
		opts.Interval = "7d"
		return true
	case "yearly", "y":
		opts.Interval = "365d"
		return true
	case "every-hour", "eh", "h":
		opts.Interval = "1h"
		return true
	}
	return false
}

func parseParameterizedEvery(args []string, idx *int, opts *scheduleAddOpts) bool {
	token := strings.ToLower(args[*idx])
	if token != "every" && token != "e" {
		return false
	}

	num := 1
	unit := "h"
	consumeEveryNumber(args, idx, &num)
	consumeEveryUnit(args, idx, &unit)
	opts.Interval = strconv.Itoa(num) + unit

	return true
}

func consumeEveryNumber(args []string, idx *int, num *int) {
	if *idx+1 >= len(args) {
		return
	}

	parsed, err := strconv.Atoi(args[*idx+1])
	if err == nil && parsed > 0 {
		*num = parsed
		*idx++
	}
}

func consumeEveryUnit(args []string, idx *int, unit *string) {
	if *idx+1 >= len(args) {
		return
	}

	next := strings.ToLower(args[*idx+1])
	canonicalUnit, isRecognized := normalizeTimeUnit(next)
	if isRecognized {
		*unit = canonicalUnit
		*idx++
	}
}

func normalizeTimeUnit(raw string) (string, bool) {
	switch raw {
	case "hours", "hour", "h", "hr":
		return "h", true
	case "daily", "days", "day", "d":
		return "d", true
	case "weekly", "weeks", "week", "w":
		return "w", true
	case "yearly", "years", "year", "y":
		return "y", true
	case "minutes", "minute", "min", "m":
		return "m", true
	case "seconds", "second", "sec", "s":
		return "s", true
	}

	return "", false
}
