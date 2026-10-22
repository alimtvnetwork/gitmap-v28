// Package uuidutil provides RFC 4122 compliant Version 4 UUID generation and canonical formatting utilities.
package uuidutil

import (
	"crypto/rand"
	"encoding/hex"
)

const (
	// UUIDLength represents the canonical length of a hyphenated UUID string (36 characters).
	UUIDLength = 36

	// RawUUIDByteCount represents the byte length of a raw binary UUID (16 bytes = 128 bits).
	RawUUIDByteCount = 16
)

// FormatUUIDBytes formats a 16-byte raw UUID slice into the standard 36-character
// lowercase hexadecimal hyphen-separated RFC 4122 representation (8-4-4-4-12 format).
func FormatUUIDBytes(rawUUIDBytes []byte) string {
	formattedUUIDBuffer := make([]byte, UUIDLength)

	hex.Encode(formattedUUIDBuffer[0:8], rawUUIDBytes[0:4])
	formattedUUIDBuffer[8] = '-'
	hex.Encode(formattedUUIDBuffer[9:13], rawUUIDBytes[4:6])
	formattedUUIDBuffer[13] = '-'
	hex.Encode(formattedUUIDBuffer[14:18], rawUUIDBytes[6:8])
	formattedUUIDBuffer[18] = '-'
	hex.Encode(formattedUUIDBuffer[19:23], rawUUIDBytes[8:10])
	formattedUUIDBuffer[23] = '-'
	hex.Encode(formattedUUIDBuffer[24:36], rawUUIDBytes[10:16])

	return string(formattedUUIDBuffer)
}

// GenerateUUID generates a cryptographically secure RFC 4122 Version 4 random UUID string.
// It sets the version bits (bits 4-7 of the 7th byte to 0100) and the variant bits
// (bits 6-7 of the 9th byte to 10) according to RFC 4122 section 4.4.
func GenerateUUID() string {
	randomEntropyBytes := make([]byte, RawUUIDByteCount)
	_, _ = rand.Read(randomEntropyBytes)

	// Set Version 4 bits: 0100 xxxx
	randomEntropyBytes[6] = (randomEntropyBytes[6] & 0x0f) | 0x40

	// Set Variant 1 bits (RFC 4122): 10xx xxxx
	randomEntropyBytes[8] = (randomEntropyBytes[8] & 0x3f) | 0x80

	return FormatUUIDBytes(randomEntropyBytes)
}

// New is a convenience alias for GenerateUUID returning a new Version 4 UUID string.
func New() string {
	return GenerateUUID()
}
