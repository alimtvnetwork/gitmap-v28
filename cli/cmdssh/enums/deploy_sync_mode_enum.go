package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// DeploySyncModeType defines bidirectional synchronization conflict behavior.
type DeploySyncModeType string

const (
	// SyncModeNone prompts interactively on conflict.
	SyncModeNone DeploySyncModeType = ""
	// SyncModeOverwrite forces overwriting the remote file unconditionally.
	SyncModeOverwrite DeploySyncModeType = "overwrite"
	// SyncModeSkip skips transfer if remote file exists.
	SyncModeSkip DeploySyncModeType = "skip"
	// SyncModeSync synchronizes bidirectionally based on latest mtime.
	SyncModeSync DeploySyncModeType = "sync"
	// SyncModeSyncRight pushes local to remote only if local is newer.
	SyncModeSyncRight DeploySyncModeType = "sync-right"
	// SyncModeSyncLeft pulls remote to local only if remote is newer.
	SyncModeSyncLeft DeploySyncModeType = "sync-left"
)

// IsValid reports whether the mode matches a known sync mode.
func (m DeploySyncModeType) IsValid() bool {
	switch m {
	case SyncModeNone, SyncModeOverwrite, SyncModeSkip,
		SyncModeSync, SyncModeSyncRight, SyncModeSyncLeft:
		return true
	default:
		return false
	}
}

// String returns the string representation of the mode.
func (m DeploySyncModeType) String() string {
	return string(m)
}

// ParseDeploySyncModeType parses a raw string into a validated DeploySyncModeType.
func ParseDeploySyncModeType(s string) (DeploySyncModeType, *apperror.AppError) {
	val := DeploySyncModeType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseDeploySyncModeType", "E1000", fmt.Sprintf("invalid deploy sync mode: %q", s)).WithContext("value", s)
}

// AllDeploySyncModeTypes returns all valid DeploySyncModeType variants.
func AllDeploySyncModeTypes() []DeploySyncModeType {
	return []DeploySyncModeType{
		SyncModeNone,
		SyncModeOverwrite,
		SyncModeSkip,
		SyncModeSync,
		SyncModeSyncRight,
		SyncModeSyncLeft,
	}
}

// IsNone reports whether the mode prompts interactively on conflict.
func (m DeploySyncModeType) IsNone() bool {
	return m == SyncModeNone
}

// IsOverwrite reports whether the mode forces overwriting the remote file.
func (m DeploySyncModeType) IsOverwrite() bool {
	return m == SyncModeOverwrite
}

// IsSkip reports whether the mode skips transfer if the remote file exists.
func (m DeploySyncModeType) IsSkip() bool {
	return m == SyncModeSkip
}

// IsSync reports whether the mode synchronizes bidirectionally.
func (m DeploySyncModeType) IsSync() bool {
	return m == SyncModeSync
}

// IsSyncRight reports whether the mode pushes local to remote only.
func (m DeploySyncModeType) IsSyncRight() bool {
	return m == SyncModeSyncRight
}

// IsSyncLeft reports whether the mode pulls remote to local only.
func (m DeploySyncModeType) IsSyncLeft() bool {
	return m == SyncModeSyncLeft
}

// HasSyncMode reports whether the mode is non-empty and valid.
func (m DeploySyncModeType) HasSyncMode() bool {
	return m != SyncModeNone && m.IsValid()
}
