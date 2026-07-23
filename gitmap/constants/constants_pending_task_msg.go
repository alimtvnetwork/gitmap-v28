package constants

// Pending task error messages.
const (
	ErrPendingTaskInsert   = "failed to insert pending task: %v (operation: insert)"
	ErrPendingTaskQuery    = "failed to query pending tasks: %v (operation: query)"
	ErrPendingTaskComplete = "failed to complete task: %v (operation: complete)"
	ErrPendingTaskFail     = "failed to update task failure: %v (operation: update)"
	ErrPendingTaskNotFound = "pending task not found: %d\n"
	ErrTaskTypeNotFound    = "task type not found: %s"
	ErrPendingTaskExists   = "pending task already exists for %s at %s (Id %d)\n"
	ErrPendingReplayFailed = "failed to replay command: %v (operation: exec)"
)

// Pending task warning messages.
const (
	WarnPendingDBOpen       = "Warning: could not open DB for task tracking: %v\n"
	WarnPendingTypeLookup   = "Warning: task type lookup failed: %v\n"
	WarnPendingInsertFailed = "Warning: could not record pending task: %v\n"
	WarnPendingCompleteFail = "Warning: could not mark task #%d complete: %v\n"
	WarnPendingFailUpdate   = "Warning: could not update task #%d failure: %v\n"
)

// Pending task failure reasons for FailureReason field.
const (
	ReasonLockScanFailed   = "lock scan failed: %v"
	ReasonNoLockingProcs   = "removal failed, no locking processes found: %v"
	ReasonUserDeclined     = "user declined to terminate locking processes"
	ReasonRetryFailed      = "retry removal failed: %v"
	ReasonReplayFailed     = "command replay failed: %v"
	ReasonTargetNotFound   = "target path does not exist: %s (operation: stat, reason: file does not exist)"
	ReasonWorkDirNotFound  = "working directory does not exist: %s (operation: stat, reason: directory does not exist)"
	ReasonPermissionDenied = "permission denied at path: %s (operation: %s, reason: %v)"
)

// Pending task help text.
const (
	HelpPending   = "  pending              List all pending tasks"
	HelpDoPending = "  do-pending (dp)      Retry pending tasks (all or by ID)"
)

// Pending task terminal messages.
const (
	MsgPendingTaskCreated   = "Task #%d created: %s %s\n"
	MsgPendingTaskCompleted = "Task #%d completed: %s\n"
	MsgPendingTaskFailed    = "Task #%d failed: %s\n"
	MsgPendingListHeader    = "Pending Tasks:\n"
	MsgPendingListRow       = "  #%-6d %-8s %-40s %s\n"
	MsgPendingListEmpty     = "No pending tasks.\n"
	MsgPendingRetryAll      = "Retrying %d pending task(s)...\n"
	MsgPendingRetryOne      = "Retrying task #%d...\n"
	MsgPendingReplaying     = "Replaying: gitmap %s\n"
	MsgPendingSkipNotExist  = "Task #%d skipped: target path no longer exists, marking complete\n"
)

// `gitmap pending clear` messages.
const (
	MsgPendingClearHeader     = "\n  ╔══════════════════════════════════════╗\n  ║       gitmap pending clear           ║\n  ╚══════════════════════════════════════╝\n\n"
	MsgPendingClearMode       = "  → Mode: %s\n"
	MsgPendingClearScanned    = "  → Scanned %d pending task(s)\n"
	MsgPendingClearCandidate  = "  • #%d  type=%-8s reason=%-18s target=%s\n"
	MsgPendingClearDryRun     = "\n  (dry-run) — re-run without --dry-run to actually delete %d task(s).\n\n"
	MsgPendingClearConfirm    = "  → Delete the %d task(s) above? (yes/N): "
	MsgPendingClearAborted    = "  ✗ Cleared canceled by user. No tasks deleted.\n"
	MsgPendingClearDeleted    = "  ✓ Deleted task #%d (%s)\n"
	MsgPendingClearDone       = "\n  ✓ Cleared %d/%d pending task(s).\n\n"
	MsgPendingClearNoMatches  = "  ✓ No %s tasks found — nothing to clear.\n"
	MsgPendingClearReasonURL  = "url-shaped-target"
	MsgPendingClearReasonChar = "illegal-path-char"
	MsgPendingClearReasonOrph = "orphan-target-missing"
	MsgPendingClearReasonAll  = "user-requested"
	MsgPendingClearReasonByID = "id-match"
)

// `gitmap pending clear` errors.
const (
	ErrPendingClearUnknownMode = "Error: unknown clear mode %q (allowed: orphans|illegal|all|<id>)\n"
	ErrPendingClearBadID       = "Error: invalid task id %q (operation: parse, reason: must be a positive integer)\n"
	ErrPendingClearDeleteFail  = "Error: failed to delete task #%d: %v\n"
)
