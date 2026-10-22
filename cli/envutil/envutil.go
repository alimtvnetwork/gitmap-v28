package envutil

import (
	"os"
	"strings"
)

// GetTrimmed retrieves the environment variable for key with surrounding whitespace trimmed.
func GetTrimmed(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

// GetOrDefault retrieves the trimmed environment variable for key, or fallback if empty.
func GetOrDefault(key, fallback string) string {
	if val := GetTrimmed(key); val != "" {
		return val
	}

	return fallback
}
