package constants

import "strings"

// PinKeywordType represents prefix keywords allowed in pin/unpin commands.
type PinKeywordType string

const (
	PinKeywordInstaller PinKeywordType = "installer"
	PinKeywordObject    PinKeywordType = "object"
	PinKeywordVersion   PinKeywordType = "version"
)

// IsValid reports whether k is a valid pin keyword.
func (k PinKeywordType) IsValid() bool {
	switch k {
	case PinKeywordInstaller, PinKeywordObject, PinKeywordVersion:
		return true
	default:
		return false
	}
}

// String returns the string representation.
func (k PinKeywordType) String() string {
	return string(k)
}

// IsPinKeyword reports whether the given token matches any PinKeywordType.
func IsPinKeyword(s string) bool {
	return PinKeywordType(strings.ToLower(s)).IsValid()
}
