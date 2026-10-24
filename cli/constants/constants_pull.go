package constants

// Pull-specific flags + messages added in Phase 2.5 (v3.10.0).
//
// `--parallel <N>` runs the pull batch through a worker pool of width N.
// `--only-available` intersects the target set with `gitmap find-next`
// before pulling, so we skip every repo whose latest VersionProbe row
// reports no new tag — turns "pull all" into "pull what's actually new".

// CLI flag descriptions (used by flag.NewFlagSet).
const (
	FlagDescPullParallel      = "Run up to N pulls concurrently (0 = auto / NumCPU, 1 = sequential)"
	FlagDescPullOnlyAvailable = "Skip repos whose latest probe reports no new tag (run `gitmap probe --all` first)"
	FlagDescPullRaw           = "Stream raw git output directly instead of using progress bar"
	FlagDescPullStatus        = "Render full post-pull repository status table"
	FlagDescPullJSON          = "Output batch pull summary as JSON"
	FlagDescPullProbe         = "Probe GitHub account and workspace for companion repositories (repo-secrets & repo-cache) before pulling"
)

// Progress bar visual glyphs.
const (
	ProgressBarFilledRich = "█"
	ProgressBarEmptyRich  = "░"
	ProgressBarFilledSafe = "="
	ProgressBarEmptySafe  = "-"
)

// User-facing pull messages (Phase 2.5).
const (
	MsgPullNoAvailable     = "No repos with available updates. Run `gitmap probe --all` first.\n"
	WarnPullFilterFallback = "  ⚠ --only-available filter unavailable (probe DB unreadable); falling back to full target set"
)

// Pull Task table names.
const (
	TablePullTask       = "PullTask"
	TablePullTaskRepo   = "PullTaskRepo"
	TablePullEtaHistory = "PullEtaHistory"
)

// PullTaskStatusType represents overall execution status of a batch pull task.
type PullTaskStatusType string

const (
	PullTaskStatusPending   PullTaskStatusType = "pending"
	PullTaskStatusRunning   PullTaskStatusType = "running"
	PullTaskStatusCompleted PullTaskStatusType = "completed"
	PullTaskStatusFailed    PullTaskStatusType = "failed"
	PullTaskStatusCanceled  PullTaskStatusType = "canceled"
)

// IsCompleted reports whether the task reached completion.
func (s PullTaskStatusType) IsCompleted() bool {
	return s == PullTaskStatusCompleted
}

// IsFailed reports whether the task failed.
func (s PullTaskStatusType) IsFailed() bool {
	return s == PullTaskStatusFailed
}

// IsCanceled reports whether the task was canceled.
func (s PullTaskStatusType) IsCanceled() bool {
	return s == PullTaskStatusCanceled
}

// PullTaskRepoStatusType represents the progress status of an individual repo in a pull task.
type PullTaskRepoStatusType string

const (
	PullTaskRepoStatusPending PullTaskRepoStatusType = "pending"
	PullTaskRepoStatusPulling PullTaskRepoStatusType = "pulling"
	PullTaskRepoStatusSuccess PullTaskRepoStatusType = "success"
	PullTaskRepoStatusSkipped PullTaskRepoStatusType = "skipped"
	PullTaskRepoStatusFailed  PullTaskRepoStatusType = "failed"
)

// IsPulling reports whether the repo is actively pulling.
func (s PullTaskRepoStatusType) IsPulling() bool {
	return s == PullTaskRepoStatusPulling
}

// IsSuccess reports whether the repo pulled successfully.
func (s PullTaskRepoStatusType) IsSuccess() bool {
	return s == PullTaskRepoStatusSuccess
}

// IsFailed reports whether the repo pull failed.
func (s PullTaskRepoStatusType) IsFailed() bool {
	return s == PullTaskRepoStatusFailed
}

// PullChangesStatusType represents the outcome status of a repo pull operation.
type PullChangesStatusType string

const (
	PullChangesStatusActive      PullChangesStatusType = "active"
	PullChangesStatusUpToDate    PullChangesStatusType = "up-to-date"
	PullChangesStatusSynced      PullChangesStatusType = "synced"
	PullChangesStatusSuccess     PullChangesStatusType = "success"
	PullChangesStatusUpdated     PullChangesStatusType = "updated"
	PullChangesStatusFastForward PullChangesStatusType = "fast-forward"
	PullChangesStatusFailed      PullChangesStatusType = "failed"
	PullChangesStatusClean       PullChangesStatusType = "clean"
)

// IsActive reports whether the status denotes active / up-to-date.
func (s PullChangesStatusType) IsActive() bool {
	lower := string(s)
	return s == PullChangesStatusActive || s == PullChangesStatusUpToDate || s == PullChangesStatusSynced || lower == "up_to_date" || lower == "active"
}

// IsUpdated reports whether the status denotes a successful update.
func (s PullChangesStatusType) IsUpdated() bool {
	lower := string(s)
	return s == PullChangesStatusSuccess || s == PullChangesStatusUpdated || s == PullChangesStatusFastForward || lower == "ok" || lower == "fast_forward"
}

// IsFailed reports whether the status denotes a failed pull.
func (s PullChangesStatusType) IsFailed() bool {
	lower := string(s)
	return s == PullChangesStatusFailed || lower == "fail" || lower == "error"
}

// Pull task settings and operational tuning constants.
const (
	PullTaskHeartbeatIntervalSec   = 5
	PullTaskHeartbeatExpirySec     = 30
	PullTaskStaleLockThresholdSec  = 60
	PullTaskWatcherPollIntervalSec = 10
	PullTaskDefaultBaselineEtaSec  = 0.85
	PullTaskRollingEtaSampleLimit  = 8
	PullTaskRollingEtaMinSamples   = 2
	SQLPragmaBusyTimeout10s        = "PRAGMA busy_timeout = 10000"
)
