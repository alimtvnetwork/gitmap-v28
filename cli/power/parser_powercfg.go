package power

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePowercfgSettingIndex extracts the AC power setting value in seconds.
func ParsePowercfgSettingIndex(output, alias string) (int, error) {
	lines := strings.Split(output, "\n")
	foundAlias := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "GUID Alias: "+alias) {
			foundAlias = true
			continue
		}

		if foundAlias && strings.HasPrefix(trimmed, "Current AC Power Setting Index:") {
			return extractSecondsFromLine(trimmed)
		}
	}

	return 0, fmt.Errorf("alias %s not found in powercfg output", alias)
}

func extractSecondsFromLine(line string) (int, error) {
	parts := strings.Split(line, ":")
	if len(parts) < 2 {
		return 0, fmt.Errorf("malformed index line: %s", line)
	}

	hexStr := strings.TrimSpace(parts[1])
	hexStr = strings.TrimPrefix(hexStr, "0x")
	val, err := strconv.ParseInt(hexStr, 16, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse index %q: %w", hexStr, err)
	}

	return int(val), nil
}

// ConvertSecondsToMinutes converts seconds to minutes, rounding up sub-minute values.
func ConvertSecondsToMinutes(seconds int) int {
	if seconds <= 0 {
		return 0
	}

	if seconds < 60 {
		return 1
	}

	return seconds / 60
}
