package power

import (
	"bufio"
	"strconv"
	"strings"
)

// ParsePmsetSettings converts output of `pmset -g` into Settings.
func ParsePmsetSettings(output string) Settings {
	displayMin := ParsePmsetValue(output, "displaysleep")
	sleepMin := ParsePmsetValue(output, "sleep")
	isNever := ParsePmsetIsNeverSleep(displayMin, sleepMin)

	return Settings{
		Platform:              "darwin",
		DisplayTimeoutMinutes: displayMin,
		SleepTimeoutMinutes:   sleepMin,
		IsNeverSleep:          isNever,
		Source:                "pmset",
	}
}

// ParsePmsetValue searches lines for `<key> <int>`.
func ParsePmsetValue(output, key string) int {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		val, hasMatch := extractPmsetLineValue(scanner.Text(), key)
		if hasMatch {
			return val
		}
	}

	return 0
}

func extractPmsetLineValue(line, key string) (int, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, key) {
		return 0, false
	}

	fields := strings.Fields(trimmed)
	if len(fields) < 2 || fields[0] != key {
		return 0, false
	}

	val, err := strconv.Atoi(fields[1])
	if err != nil {
		return 0, false
	}

	return val, true
}

// ParsePmsetIsNeverSleep checks whether display and system sleep are 0.
func ParsePmsetIsNeverSleep(displayMinutes, sleepMinutes int) bool {
	isNever := displayMinutes == 0 && sleepMinutes == 0

	return isNever
}
