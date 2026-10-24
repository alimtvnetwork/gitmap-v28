package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// OSTUIViewModeType specifies the current mode of the OS TUI dashboard.
type OSTUIViewModeType string

const (
	OSTUIViewModeNavType  OSTUIViewModeType = "nav"
	OSTUIViewModeExecType OSTUIViewModeType = "exec"
	OSTUIViewModeDoneType OSTUIViewModeType = "done"
)

// IsValid reports whether the TUI view mode matches a known variant.
func (m OSTUIViewModeType) IsValid() bool {
	switch m {
	case OSTUIViewModeNavType, OSTUIViewModeExecType, OSTUIViewModeDoneType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the TUI view mode.
func (m OSTUIViewModeType) String() string {
	return string(m)
}

// ParseOSTUIViewModeType parses a raw string into a validated OSTUIViewModeType.
func ParseOSTUIViewModeType(s string) (OSTUIViewModeType, *apperror.AppError) {
	val := OSTUIViewModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseOSTUIViewModeType", "E1000", fmt.Sprintf("invalid OS TUI view mode: %q", s)).WithContext("value", s)
}

// AllOSTUIViewModeTypes returns all valid OSTUIViewModeType variants.
func AllOSTUIViewModeTypes() []OSTUIViewModeType {
	return []OSTUIViewModeType{
		OSTUIViewModeNavType,
		OSTUIViewModeExecType,
		OSTUIViewModeDoneType,
	}
}

// IsNav reports whether the TUI view mode is nav.
func (m OSTUIViewModeType) IsNav() bool {
	return m == OSTUIViewModeNavType
}

// IsExec reports whether the TUI view mode is exec.
func (m OSTUIViewModeType) IsExec() bool {
	return m == OSTUIViewModeExecType
}

// IsDone reports whether the TUI view mode is done.
func (m OSTUIViewModeType) IsDone() bool {
	return m == OSTUIViewModeDoneType
}

// HasViewMode reports whether the TUI view mode is non-empty and valid.
func (m OSTUIViewModeType) HasViewMode() bool {
	return m.IsValid()
}
