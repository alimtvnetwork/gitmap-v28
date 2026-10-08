// Package ctxmodetype — Windows shell context-menu entry mode.
//
// Migrated from cli/constants/constants_installctx.go (spec 243, Wave D).
// The serialized values ("terminal", "silent", "prefill") are functional:
// they are baked into Windows registry entries and parsed back from user
// configuration, so they are preserved verbatim per the enum pattern's
// protocol-driven exception.
package ctxmodetype

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Variant represents a Windows shell context-menu entry mode.
type Variant byte

const (
	// Invalid is the zero value (always first).
	Invalid Variant = iota

	// Terminal runs pwsh -NoExit; the window stays open after the command.
	Terminal

	// Silent runs pwsh -WindowStyle Hidden; output goes via the notifier.
	Silent

	// Prefill runs pwsh -NoExit and writes "gitmap " at the prompt without
	// running a command.
	Prefill
)

// variantLabels maps each variant to its serialized form. The values are
// functional (registry/config surface), so they keep their lowercase form
// per the pattern's protocol-driven exception.
var variantLabels = [...]string{
	Invalid:  "Invalid",
	Terminal: "terminal",
	Silent:   "silent",
	Prefill:  "prefill",
}

// String returns the serialized string representation.
func (v Variant) String() string {
	if v.IsInvalid() {
		return variantLabels[Invalid]
	}

	return variantLabels[v]
}

// Label delegates to String.
func (v Variant) Label() string {
	return v.String()
}

// IsValid reports whether the variant is a real mode.
func (v Variant) IsValid() bool {
	return v > Invalid && v < Variant(len(variantLabels))
}

// IsInvalid reports whether the variant is the zero value or out of range.
func (v Variant) IsInvalid() bool {
	return v <= Invalid || v >= Variant(len(variantLabels))
}

// IsTerminal reports whether the variant is Terminal.
func (v Variant) IsTerminal() bool {
	return v == Terminal
}

// IsSilent reports whether the variant is Silent.
func (v Variant) IsSilent() bool {
	return v == Silent
}

// IsPrefill reports whether the variant is Prefill.
func (v Variant) IsPrefill() bool {
	return v == Prefill
}

// All returns all valid variants.
func All() []Variant {
	result := make([]Variant, 0, len(variantLabels)-1)
	for i := 1; i < len(variantLabels); i++ {
		result = append(result, Variant(i))
	}

	return result
}

// ByIndex returns the variant for an index, or Invalid when out of range.
func ByIndex(i int) Variant {
	isOutOfRange := i < 0 || i >= len(variantLabels)
	if isOutOfRange {
		return Invalid
	}

	return Variant(i)
}

// Parse parses a serialized mode using case-insensitive matching.
func Parse(s string) (Variant, error) {
	trimmed := strings.TrimSpace(s)
	for i, label := range variantLabels {
		if strings.EqualFold(label, trimmed) {
			return Variant(i), nil
		}
	}

	return Invalid, fmt.Errorf("invalid ctxmode: %q", s)
}

// Values returns all valid serialized values.
func Values() []string {
	result := make([]string, 0, len(variantLabels)-1)
	for _, s := range variantLabels[1:] {
		result = append(result, s)
	}

	return result
}

// MarshalJSON implements json.Marshaler.
func (v Variant) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler.
func (v *Variant) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	parsed, err := Parse(s)
	if err != nil {
		return err
	}

	*v = parsed

	return nil
}
