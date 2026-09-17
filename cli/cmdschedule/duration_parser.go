package cmdschedule

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// DurationResult represents a single reusable result envelope wrapping time.Duration.
type DurationResult = result.Result[time.Duration]

var (
	colonPattern  = regexp.MustCompile(`^(\d+):(\d+)(?::(\d+))?\s*(hr|h|m|s)?$`)
	dayPattern    = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:days|day|d)$`)
	hourPattern   = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:hours|hour|hrs|hr|h)$`)
	minutePattern = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:minutes|minute|mins|min|m)$`)
	secondPattern = regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*(?:seconds|second|secs|sec|s)$`)
)

// ParseScheduleDuration parses flexible duration strings into a DurationResult.
func ParseScheduleDuration(raw string) DurationResult {
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	if isZeroDuration(cleaned) {
		return result.Ok(time.Duration(0))
	}
	if res, isMatched := matchParsedDuration(cleaned); isMatched {
		return res
	}
	return parseFallbackDuration(cleaned)
}

func matchParsedDuration(s string) (DurationResult, bool) {
	if res, isParsed := tryParseNumericSeconds(s); isParsed {
		return res, true
	}
	if res, isParsed := tryParseColonDuration(s); isParsed {
		return res, true
	}
	return tryParseUnitDuration(s)
}

func isZeroDuration(s string) bool {
	return s == "" || s == "0" || s == "0s" || s == "now"
}

func tryParseNumericSeconds(s string) (DurationResult, bool) {
	secs, err := strconv.ParseInt(s, 10, 64)
	if err != nil || secs < 0 {
		return DurationResult{}, false
	}
	return result.Ok(time.Duration(secs) * time.Second), true
}

func tryParseColonDuration(s string) (DurationResult, bool) {
	m := colonPattern.FindStringSubmatch(s)
	if len(m) == 0 {
		return DurationResult{}, false
	}
	dur := calculateColonDuration(m[1], m[2], m[3], m[4])
	return result.Ok(dur), true
}

func calculateColonDuration(p1, p2, p3, unit string) time.Duration {
	v1, _ := strconv.ParseInt(p1, 10, 64)
	v2, _ := strconv.ParseInt(p2, 10, 64)
	if p3 != "" {
		v3, _ := strconv.ParseInt(p3, 10, 64)
		return time.Duration(v1)*time.Hour + time.Duration(v2)*time.Minute + time.Duration(v3)*time.Second
	}
	if unit == "m" {
		return time.Duration(v1)*time.Minute + time.Duration(v2)*time.Second
	}
	return time.Duration(v1)*time.Hour + time.Duration(v2)*time.Minute
}

func tryParseUnitDuration(s string) (DurationResult, bool) {
	if res, isDay := matchRegexMultiplier(s, dayPattern, 24*time.Hour); isDay {
		return res, true
	}
	if res, isHour := matchRegexMultiplier(s, hourPattern, time.Hour); isHour {
		return res, true
	}
	if res, isMin := matchRegexMultiplier(s, minutePattern, time.Minute); isMin {
		return res, true
	}
	return matchRegexMultiplier(s, secondPattern, time.Second)
}

func matchRegexMultiplier(s string, re *regexp.Regexp, unit time.Duration) (DurationResult, bool) {
	m := re.FindStringSubmatch(s)
	if len(m) == 0 {
		return DurationResult{}, false
	}
	val, err := strconv.ParseFloat(m[1], 64)
	if err != nil || val < 0 {
		return DurationResult{}, false
	}
	dur := time.Duration(val * float64(unit))
	return result.Ok(dur), true
}

func parseFallbackDuration(s string) DurationResult {
	if dur, err := time.ParseDuration(s); err == nil && dur >= 0 {
		return result.Ok(dur)
	}
	return makeInvalidDurationError(s)
}

func makeInvalidDurationError(s string) DurationResult {
	appErr := apperror.NewWithDetails(
		"ParseScheduleDuration",
		"E_INVALID_DURATION",
		fmt.Sprintf("invalid schedule duration %q", s),
		"cmdschedule",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"raw": s},
	)
	return result.Fail[time.Duration](appErr)
}
