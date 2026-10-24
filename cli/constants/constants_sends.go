package constants

// SendResultStatusType represents the outcome of a semantic sends operation on a repository.
type SendResultStatusType string

const (
	SendResultStatusPushed          SendResultStatusType = "pushed"
	SendResultStatusCommitted       SendResultStatusType = "committed"
	SendResultStatusDryRunSimulated SendResultStatusType = "dry-run-simulated"
	SendResultStatusCleanSkipped    SendResultStatusType = "clean-skipped"
	SendResultStatusAlreadyClean    SendResultStatusType = "already-clean"
	SendResultStatusFailed          SendResultStatusType = "failed"
)

// IsSuccess reports whether the status denotes a successful operation.
func (s SendResultStatusType) IsSuccess() bool {
	return s == SendResultStatusPushed || s == SendResultStatusCommitted || s == SendResultStatusDryRunSimulated
}

// IsSkipped reports whether the status denotes a skipped repository.
func (s SendResultStatusType) IsSkipped() bool {
	return s == SendResultStatusCleanSkipped || s == SendResultStatusAlreadyClean
}

// IsFailed reports whether the status denotes a failure.
func (s SendResultStatusType) IsFailed() bool {
	return s == SendResultStatusFailed
}

// String returns the string representation.
func (s SendResultStatusType) String() string {
	return string(s)
}
