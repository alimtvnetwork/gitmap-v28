// Package utils provides general-purpose shared utilities.
package utils

import "github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"

// GenerateUUID generates a cryptographically secure RFC 4122 Version 4 random UUID string.
func GenerateUUID() string {
	return uuidutil.GenerateUUID()
}

// FormatUUIDBytes formats a 16-byte raw UUID slice into the standard 36-character hyphenated representation.
func FormatUUIDBytes(rawUUIDBytes []byte) string {
	return uuidutil.FormatUUIDBytes(rawUUIDBytes)
}

// NewUUID returns a new Version 4 UUID string.
func NewUUID() string {
	return uuidutil.New()
}
