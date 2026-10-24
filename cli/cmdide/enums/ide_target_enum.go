package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// IDETargetType represents a supported IDE or editor target.
type IDETargetType string

const (
	TargetVSCode        IDETargetType = "vscode"
	TargetCursor        IDETargetType = "cursor"
	TargetAntigravity   IDETargetType = "antigravity"
	TargetDesktop       IDETargetType = "desktop"
	TargetAll           IDETargetType = "all"
	TargetGitHubDesktop IDETargetType = "github-desktop"
)

// IsValid reports whether the value matches a known variant.
func (i IDETargetType) IsValid() bool {
	switch i {
	case TargetVSCode, TargetCursor, TargetAntigravity, TargetDesktop, TargetAll, TargetGitHubDesktop:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (i IDETargetType) String() string {
	return string(i)
}

// ParseIDETargetType parses a raw string into a validated IDETargetType.
func ParseIDETargetType(s string) (IDETargetType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllIDETargetTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "targetvscode":
		return TargetVSCode, nil
	case "vscode":
		return TargetVSCode, nil
	case "targetcursor":
		return TargetCursor, nil
	case "cursor":
		return TargetCursor, nil
	case "targetantigravity":
		return TargetAntigravity, nil
	case "antigravity":
		return TargetAntigravity, nil
	case "targetdesktop":
		return TargetDesktop, nil
	case "desktop":
		return TargetDesktop, nil
	case "targetall":
		return TargetAll, nil
	case "all":
		return TargetAll, nil
	case "targetgithubdesktop":
		return TargetGitHubDesktop, nil
	case "githubdesktop":
		return TargetGitHubDesktop, nil
	}
	return "", apperror.NewValidation("ParseIDETargetType", "E1000", fmt.Sprintf("invalid idetargettype: %q", s)).WithContext("value", s)
}

// AllIDETargetTypes returns all valid variants.
func AllIDETargetTypes() []IDETargetType {
	return []IDETargetType{
		TargetVSCode,
		TargetCursor,
		TargetAntigravity,
		TargetDesktop,
		TargetAll,
		TargetGitHubDesktop,
	}
}

// IsVSCode reports whether the value is TargetVSCode.
func (i IDETargetType) IsVSCode() bool {
	return i == TargetVSCode
}

// IsCursor reports whether the value is TargetCursor.
func (i IDETargetType) IsCursor() bool {
	return i == TargetCursor
}

// IsAntigravity reports whether the value is TargetAntigravity.
func (i IDETargetType) IsAntigravity() bool {
	return i == TargetAntigravity
}

// IsDesktop reports whether the value is TargetDesktop.
func (i IDETargetType) IsDesktop() bool {
	return i == TargetDesktop
}

// IsAll reports whether the value is TargetAll.
func (i IDETargetType) IsAll() bool {
	return i == TargetAll
}

// IsGitHubDesktop reports whether the value is TargetGitHubDesktop.
func (i IDETargetType) IsGitHubDesktop() bool {
	return i == TargetGitHubDesktop
}

// HasTarget reports whether the value is non-empty and valid.
func (i IDETargetType) HasTarget() bool {
	return i.IsValid()
}
