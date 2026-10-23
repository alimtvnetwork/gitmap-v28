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

// Pull task execution status enums.
const (
	PullTaskStatusPending   = "pending"
	PullTaskStatusRunning   = "running"
	PullTaskStatusCompleted = "completed"
	PullTaskStatusFailed    = "failed"
	PullTaskStatusCanceled  = "canceled"
)

// Pull task repo progress status enums.
const (
	PullTaskRepoStatusPending = "pending"
	PullTaskRepoStatusPulling = "pulling"
	PullTaskRepoStatusSuccess = "success"
	PullTaskRepoStatusSkipped = "skipped"
	PullTaskRepoStatusFailed  = "failed"
)

// Pull task settings and operational tuning constants.
const (
	SettingPullRetentionDays       = "pull_retention_days"
	DefaultPullRetentionDays       = 7
	PullTaskHeartbeatIntervalSec   = 5
	PullTaskHeartbeatExpirySec     = 30
	PullTaskStaleLockThresholdSec  = 60
	PullTaskWatcherPollIntervalSec = 10
	PullTaskDefaultBaselineEtaSec  = 0.85
	PullTaskRollingEtaSampleLimit  = 8
	PullTaskRollingEtaMinSamples   = 2
	SQLPragmaBusyTimeout10s        = "PRAGMA busy_timeout = 10000"
)
