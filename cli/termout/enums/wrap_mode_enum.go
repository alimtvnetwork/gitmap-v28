package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// WrapModeType defines cell text wrapping behavior for terminal tables.
type WrapModeType string

const (
	// WrapClip truncates cell text at column width without appending ellipsis.
	WrapClip WrapModeType = "clip"
	// WrapEllipsis truncates cell text exceeding column width with an ellipsis indicator.
	WrapEllipsis WrapModeType = "ellipsis"
	// WrapWord wraps text softly at whitespace boundaries across multiple lines.
	WrapWord WrapModeType = "word"
)

// IsValid reports whether the value matches a known variant.
func (w WrapModeType) IsValid() bool {
	switch w {
	case WrapClip, WrapEllipsis, WrapWord:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (w WrapModeType) String() string {
	return string(w)
}

// ParseWrapModeType parses a raw string into a validated WrapModeType.
func ParseWrapModeType(s string) (WrapModeType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllWrapModeTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "wrapclip":
		return WrapClip, nil
	case "clip":
		return WrapClip, nil
	case "wrapellipsis":
		return WrapEllipsis, nil
	case "ellipsis":
		return WrapEllipsis, nil
	case "wrapword":
		return WrapWord, nil
	case "word":
		return WrapWord, nil
	}
	return "", apperror.NewValidation("ParseWrapModeType", "E1000", fmt.Sprintf("invalid wrapmodetype: %q", s)).WithContext("value", s)
}

// AllWrapModeTypes returns all valid variants.
func AllWrapModeTypes() []WrapModeType {
	return []WrapModeType{
		WrapClip,
		WrapEllipsis,
		WrapWord,
	}
}

// IsClip reports whether the value is WrapClip.
func (w WrapModeType) IsClip() bool {
	return w == WrapClip
}

// IsEllipsis reports whether the value is WrapEllipsis.
func (w WrapModeType) IsEllipsis() bool {
	return w == WrapEllipsis
}

// IsWord reports whether the value is WrapWord.
func (w WrapModeType) IsWord() bool {
	return w == WrapWord
}

// HasMode reports whether the value is non-empty and valid.
func (w WrapModeType) HasMode() bool {
	return w.IsValid()
}
