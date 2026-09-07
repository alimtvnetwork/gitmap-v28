package power

import (
	"strconv"
	"strings"
)

// ParseGnomeSeconds parses gsettings output (e.g. "uint32 300\n" or "300\n") to seconds.
func ParseGnomeSeconds(output string) int {
	clean := strings.TrimSpace(output)
	clean = strings.TrimPrefix(clean, "uint32")
	clean = strings.TrimSpace(clean)

	val, err := strconv.Atoi(clean)
	if err != nil || val < 0 {
		return 0
	}

	return val
}

// ParseGnomeTimeoutMinutes parses gsettings output to minutes.
func ParseGnomeTimeoutMinutes(output string) int {
	sec := ParseGnomeSeconds(output)

	return ConvertSecondsToMinutes(sec)
}

// ParseGnomeBoolean parses gsettings boolean output (e.g. "true", "false").
func ParseGnomeBoolean(output string) bool {
	clean := strings.TrimSpace(strings.ToLower(output))

	return clean == "true"
}
