// Package installer — versioning.go provides semantic version calculation helpers.
package installer

import (
	"fmt"
	"strconv"
	"strings"
)

// NextSemanticVersion increments the patch component of a semantic version string.
func NextSemanticVersion(current string) string {
	raw := strings.TrimSpace(current)
	hasPrefix := strings.HasPrefix(strings.ToLower(raw), "v")
	vStr := strings.TrimPrefix(strings.ToLower(raw), "v")

	parts := strings.Split(vStr, ".")
	if len(parts) != 3 {
		return defaultFallbackVersion(hasPrefix)
	}

	patch, err := strconv.Atoi(parts[2])
	if err != nil {
		patch = 0
	}
	patch++

	next := fmt.Sprintf("%s.%s.%d", parts[0], parts[1], patch)
	if hasPrefix {
		return "v" + next
	}
	return next
}

func defaultFallbackVersion(hasPrefix bool) string {
	if hasPrefix {
		return "v1.0.1"
	}
	return "1.0.1"
}
