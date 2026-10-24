package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strconv"
	"strings"
)

// ChoiceType is the outcome of resolving one conflict.
type ChoiceType int

const (
	// ChoiceLeft writes LEFT's version onto the destination side.
	ChoiceLeft ChoiceType = 0
	// ChoiceRight writes RIGHT's version onto the destination side.
	ChoiceRight ChoiceType = 1
	// ChoiceSkip leaves both sides untouched.
	ChoiceSkip ChoiceType = 2
	// ChoiceQuit aborts the run; partial changes are kept.
	ChoiceQuit ChoiceType = 3
)

// IsValid reports whether the value matches a known variant.
func (c ChoiceType) IsValid() bool {
	switch c {
	case ChoiceLeft, ChoiceRight, ChoiceSkip, ChoiceQuit:
		return true
	default:
		return false
	}
}

// String returns the numeric representation of the value.
func (c ChoiceType) String() string {
	return strconv.Itoa(int(c))
}

// ParseChoiceType parses a raw string into a validated ChoiceType.
func ParseChoiceType(s string) (ChoiceType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	if n, err := strconv.Atoi(trimmed); err == nil {
		val := ChoiceType(n)
		if val.IsValid() {
			return val, nil
		}
	}
	for _, val := range AllChoiceTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "choiceleft":
		return ChoiceLeft, nil
	case "left":
		return ChoiceLeft, nil
	case "choiceright":
		return ChoiceRight, nil
	case "right":
		return ChoiceRight, nil
	case "choiceskip":
		return ChoiceSkip, nil
	case "skip":
		return ChoiceSkip, nil
	case "choicequit":
		return ChoiceQuit, nil
	case "quit":
		return ChoiceQuit, nil
	}
	return 0, apperror.NewValidation("ParseChoiceType", "E1000", fmt.Sprintf("invalid choicetype: %q", s)).WithContext("value", s)
}

// AllChoiceTypes returns all valid variants.
func AllChoiceTypes() []ChoiceType {
	return []ChoiceType{
		ChoiceLeft,
		ChoiceRight,
		ChoiceSkip,
		ChoiceQuit,
	}
}

// IsLeft reports whether the value is ChoiceLeft.
func (c ChoiceType) IsLeft() bool {
	return c == ChoiceLeft
}

// IsRight reports whether the value is ChoiceRight.
func (c ChoiceType) IsRight() bool {
	return c == ChoiceRight
}

// IsSkip reports whether the value is ChoiceSkip.
func (c ChoiceType) IsSkip() bool {
	return c == ChoiceSkip
}

// IsQuit reports whether the value is ChoiceQuit.
func (c ChoiceType) IsQuit() bool {
	return c == ChoiceQuit
}

// HasChoice reports whether the value is non-empty and valid.
func (c ChoiceType) HasChoice() bool {
	return c.IsValid()
}
