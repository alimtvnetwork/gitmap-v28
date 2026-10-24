package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// FileEventType identifies file mutation categories.
type FileEventType string

const (
	EventCreate FileEventType = "create"
	EventModify FileEventType = "modify"
	EventDelete FileEventType = "delete"
)

// IsValid reports whether the value matches a known variant.
func (f FileEventType) IsValid() bool {
	switch f {
	case EventCreate, EventModify, EventDelete:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (f FileEventType) String() string {
	return string(f)
}

// ParseFileEventType parses a raw string into a validated FileEventType.
func ParseFileEventType(s string) (FileEventType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllFileEventTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "eventcreate":
		return EventCreate, nil
	case "create":
		return EventCreate, nil
	case "eventmodify":
		return EventModify, nil
	case "modify":
		return EventModify, nil
	case "eventdelete":
		return EventDelete, nil
	case "delete":
		return EventDelete, nil
	}
	return "", apperror.NewValidation("ParseFileEventType", "E1000", fmt.Sprintf("invalid fileeventtype: %q", s)).WithContext("value", s)
}

// AllFileEventTypes returns all valid variants.
func AllFileEventTypes() []FileEventType {
	return []FileEventType{
		EventCreate,
		EventModify,
		EventDelete,
	}
}

// IsCreate reports whether the value is EventCreate.
func (f FileEventType) IsCreate() bool {
	return f == EventCreate
}

// IsModify reports whether the value is EventModify.
func (f FileEventType) IsModify() bool {
	return f == EventModify
}

// IsDelete reports whether the value is EventDelete.
func (f FileEventType) IsDelete() bool {
	return f == EventDelete
}

// HasEvent reports whether the value is non-empty and valid.
func (f FileEventType) HasEvent() bool {
	return f.IsValid()
}
