package store

const (
	InstallStatusRunningType = "running"
	InstallStatusSuccessType = "success"
	InstallStatusFailedType  = "failed"
	InstallStatusSkippedType = "skipped"
)

// InstallLogRecord encapsulates execution telemetry for an installation event.
type InstallLogRecord struct {
	ID           string `json:"id"`
	TargetType   string `json:"targetType"`
	TargetName   string `json:"targetName"`
	Action       string `json:"action"`
	Status       string `json:"status"`
	ExitCode     int    `json:"exitCode"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	LogPath      string `json:"logPath,omitempty"`
	StartedAt    string `json:"startedAt"`
	EndedAt      string `json:"endedAt,omitempty"`
}

// IsSuccess reports whether the record represents a successful installation.
func (r InstallLogRecord) IsSuccess() bool {
	return r.Status == InstallStatusSuccessType || r.Status == "completed"
}

// IsFailed reports whether the record represents a failed installation.
func (r InstallLogRecord) IsFailed() bool {
	return r.Status == InstallStatusFailedType || r.Status == "failure"
}

// IsSkipped reports whether the record represents a skipped installation.
func (r InstallLogRecord) IsSkipped() bool {
	return r.Status == InstallStatusSkippedType
}

// IsRunning reports whether the record represents an in-progress installation.
func (r InstallLogRecord) IsRunning() bool {
	return r.Status == InstallStatusRunningType || r.Status == "started"
}
