package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// TargetType defines the resolution mode for history purging.
type TargetType string

const (
	TargetTypeFolder TargetType = "folder"
	TargetTypeFile   TargetType = "file"
	TargetTypeCommit TargetType = "commit"
)

// IsValid reports whether the value matches a known variant.
func (t TargetType) IsValid() bool {
	switch t {
	case TargetTypeFolder, TargetTypeFile, TargetTypeCommit:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (t TargetType) String() string {
	return string(t)
}

// ParseTargetType parses a raw string into a validated TargetType.
func ParseTargetType(s string) (TargetType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllTargetTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "targettypefolder":
		return TargetTypeFolder, nil
	case "folder":
		return TargetTypeFolder, nil
	case "targettypefile":
		return TargetTypeFile, nil
	case "file":
		return TargetTypeFile, nil
	case "targettypecommit":
		return TargetTypeCommit, nil
	case "commit":
		return TargetTypeCommit, nil
	}
	return "", apperror.NewValidation("ParseTargetType", "E1000", fmt.Sprintf("invalid targettype: %q", s)).WithContext("value", s)
}

// AllTargetTypes returns all valid variants.
func AllTargetTypes() []TargetType {
	return []TargetType{
		TargetTypeFolder,
		TargetTypeFile,
		TargetTypeCommit,
	}
}

// IsFolder reports whether the value is TargetTypeFolder.
func (t TargetType) IsFolder() bool {
	return t == TargetTypeFolder
}

// IsFile reports whether the value is TargetTypeFile.
func (t TargetType) IsFile() bool {
	return t == TargetTypeFile
}

// IsCommit reports whether the value is TargetTypeCommit.
func (t TargetType) IsCommit() bool {
	return t == TargetTypeCommit
}

// HasTarget reports whether the value is non-empty and valid.
func (t TargetType) HasTarget() bool {
	return t.IsValid()
}
