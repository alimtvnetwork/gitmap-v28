package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// FileTouchType represents the categorical operation applied to a touched file.
type FileTouchType string

const (
	FileTouchRead    FileTouchType = "READ"
	FileTouchEdit    FileTouchType = "EDIT"
	FileTouchCreated FileTouchType = "CREATED"
)

// IsValid reports whether the touch type is a recognized enumeration value.
func (t FileTouchType) IsValid() bool {
	switch t {
	case FileTouchRead, FileTouchEdit, FileTouchCreated:
		return true
	default:
		return false
	}
}

// String returns the string representation of the touch type.
func (t FileTouchType) String() string {
	return string(t)
}

// ParseFileTouchType parses a raw string into a validated FileTouchType.
func ParseFileTouchType(s string) (FileTouchType, *apperror.AppError) {
	val := FileTouchType(strings.ToUpper(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseFileTouchType", "E1000", fmt.Sprintf("invalid file touch type: %q", s)).WithContext("value", s)
}

// AllFileTouchTypes returns all valid FileTouchType variants.
func AllFileTouchTypes() []FileTouchType {
	return []FileTouchType{
		FileTouchRead,
		FileTouchEdit,
		FileTouchCreated,
	}
}

// IsRead reports whether the touch type is a read.
func (t FileTouchType) IsRead() bool {
	return t == FileTouchRead
}

// IsEdit reports whether the touch type is an edit.
func (t FileTouchType) IsEdit() bool {
	return t == FileTouchEdit
}

// IsCreated reports whether the touch type is a creation.
func (t FileTouchType) IsCreated() bool {
	return t == FileTouchCreated
}
