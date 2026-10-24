package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// BorderStyleType defines the visual framing style for rendered terminal tables.
type BorderStyleType string

const (
	// BorderNone indicates no border framing (standard spaced table).
	BorderNone BorderStyleType = "none"
	// BorderRounded uses curved UTF-8 box-drawing characters (╭, ╮, ╯, ╰).
	BorderRounded BorderStyleType = "rounded"
	// BorderSharp uses classic rectangular box-drawing characters (┌, ┐, ┘, └).
	BorderSharp BorderStyleType = "sharp"
	// BorderMinimal uses understated separator lines without outer vertical bars.
	BorderMinimal BorderStyleType = "minimal"
)

// IsValid reports whether the value matches a known variant.
func (b BorderStyleType) IsValid() bool {
	switch b {
	case BorderNone, BorderRounded, BorderSharp, BorderMinimal:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (b BorderStyleType) String() string {
	return string(b)
}

// ParseBorderStyleType parses a raw string into a validated BorderStyleType.
func ParseBorderStyleType(s string) (BorderStyleType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllBorderStyleTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "bordernone":
		return BorderNone, nil
	case "none":
		return BorderNone, nil
	case "borderrounded":
		return BorderRounded, nil
	case "rounded":
		return BorderRounded, nil
	case "bordersharp":
		return BorderSharp, nil
	case "sharp":
		return BorderSharp, nil
	case "borderminimal":
		return BorderMinimal, nil
	case "minimal":
		return BorderMinimal, nil
	}
	return "", apperror.NewValidation("ParseBorderStyleType", "E1000", fmt.Sprintf("invalid borderstyletype: %q", s)).WithContext("value", s)
}

// AllBorderStyleTypes returns all valid variants.
func AllBorderStyleTypes() []BorderStyleType {
	return []BorderStyleType{
		BorderNone,
		BorderRounded,
		BorderSharp,
		BorderMinimal,
	}
}

// IsNone reports whether the value is BorderNone.
func (b BorderStyleType) IsNone() bool {
	return b == BorderNone
}

// IsRounded reports whether the value is BorderRounded.
func (b BorderStyleType) IsRounded() bool {
	return b == BorderRounded
}

// IsSharp reports whether the value is BorderSharp.
func (b BorderStyleType) IsSharp() bool {
	return b == BorderSharp
}

// IsMinimal reports whether the value is BorderMinimal.
func (b BorderStyleType) IsMinimal() bool {
	return b == BorderMinimal
}

// HasStyle reports whether the value is non-empty and valid.
func (b BorderStyleType) HasStyle() bool {
	return b.IsValid()
}
