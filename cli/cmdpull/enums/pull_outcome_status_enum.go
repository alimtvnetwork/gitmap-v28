package enums

// PullOutcomeStatusType represents the outcome status of a repo pull operation.
type PullOutcomeStatusType string

const (
	StatusUpToDate    PullOutcomeStatusType = "up-to-date"
	StatusUpdated     PullOutcomeStatusType = "updated"
	StatusFastForward PullOutcomeStatusType = "fast-forward"
	StatusMerged      PullOutcomeStatusType = "auto-merged"
	StatusDirty       PullOutcomeStatusType = "dirty"
	StatusFailed      PullOutcomeStatusType = "failed"
	StatusClean       PullOutcomeStatusType = "clean"
	StatusLocalOnly   PullOutcomeStatusType = "local-only"
)

// IsValid reports whether the status matches a recognized outcome variant.
func (s PullOutcomeStatusType) IsValid() bool {
	switch s {
	case StatusUpToDate, StatusUpdated, StatusFastForward, StatusMerged,
		StatusDirty, StatusFailed, StatusClean, StatusLocalOnly:
		return true
	default:
		return false
	}
}

// String returns the string representation of the outcome status.
func (s PullOutcomeStatusType) String() string {
	return string(s)
}

// IsDirty reports whether the status denotes a dirty working directory.
func (s PullOutcomeStatusType) IsDirty() bool {
	return s == StatusDirty
}

// IsClean reports whether the status denotes a clean repository.
func (s PullOutcomeStatusType) IsClean() bool {
	return s == StatusClean
}
