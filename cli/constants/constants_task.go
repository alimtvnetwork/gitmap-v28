package constants

// gitmap:cmd top-level
// Task CLI commands.
const (
	CmdTask      = "task"
	CmdTaskAlias = "tk"
)

// gitmap:cmd top-level
// Task subcommands.
const (
	CmdTaskCreate = "create" // gitmap:cmd skip
	CmdTaskList   = "list"   // gitmap:cmd skip
	CmdTaskRun    = "run"    // gitmap:cmd skip
	CmdTaskShow   = "show"   // gitmap:cmd skip
	CmdTaskDelete = "delete" // gitmap:cmd skip
)

// Task help text.
const HelpTask = "  task (tk) <sub>     Manage file-sync watch tasks"

// Task defaults.
const (
	TaskDefaultInterval = 5
	TaskMinInterval     = 2
	TaskMaxGoroutines   = 64
	TaskCopyBufferSize  = 32768
	TasksFileName       = "tasks.json"
)

// Task file path within .gitmap.
const TasksFilePath = GitMapDir + "/" + TasksFileName

// Task flag names.
const (
	FlagTaskSrc      = "src"
	FlagTaskDest     = "dest"
	FlagTaskInterval = "interval"
	FlagTaskVerbose  = "verbose"
	FlagTaskDryRun   = "dry-run"
)

// Task flag descriptions.
const (
	FlagDescTaskSrc      = "Source directory path"
	FlagDescTaskDest     = "Destination directory path"
	FlagDescTaskInterval = "Sync interval in seconds (minimum 2)"
	FlagDescTaskVerbose  = "Show detailed sync output"
	FlagDescTaskDryRun   = "Preview sync actions without copying"
)

// Task terminal messages.
const (
	MsgTaskCreated    = "Task '%s' created.\n"
	MsgTaskDeleted    = "Task '%s' deleted.\n"
	MsgTaskRunning    = "Task '%s' running — syncing every %ds (Ctrl+C to stop)\n"
	MsgTaskSynced     = "Synced: %s\n"
	MsgTaskUpToDate   = "All files up to date.\n"
	MsgTaskDrySync    = "[dry-run] Would sync: %s\n"
	MsgTaskListHeader = "Tasks:\n"
	MsgTaskListRow    = "  %-20s %s → %s\n"
	MsgTaskListEmpty  = "No tasks defined. Use 'gitmap task create' to add one.\n"
	MsgTaskShowFmt    = "Name:     %s\nSource:   %s\nDest:     %s\n"
	MsgTaskStopped    = "\nTask '%s' stopped.\n"
)

// Task error messages.
const (
	ErrTaskNameRequired  = "Task name is required."
	ErrTaskSrcRequired   = "Source directory (--src) is required."
	ErrTaskDestRequired  = "Destination directory (--dest) is required."
	ErrTaskNotFound      = "Task '%s' not found.\n"
	ErrTaskAlreadyExists = "Task '%s' already exists.\n"
	ErrTaskSrcNotExist   = "Error: source directory does not exist at %s (operation: resolve, reason: file does not exist)\n"
	ErrTaskDestCreate    = "Error: failed to create destination directory at %s: %v (operation: mkdir)\n"
	ErrTaskLoadFile      = "Error: failed to load tasks file at %s: %v (operation: read)\n"
	ErrTaskSaveFile      = "Error: failed to save tasks file at %s: %v (operation: write)\n"
	ErrTaskSyncFailed    = "Sync failed for %s: %v\n"
	ErrTaskSubcommand    = "Unknown task subcommand: %s\n"
)
