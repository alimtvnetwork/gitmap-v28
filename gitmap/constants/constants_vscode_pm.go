package constants

// VS Code Project Manager (alefragnani.project-manager) sync constants.
//
// Path resolution: discover the VS Code USER-DATA root per OS first, then
// append the relative tail. Never hardcode the full path.
//
//	Windows : %APPDATA%\Code           (fallback %USERPROFILE%\AppData\Roaming\Code)
//	macOS   : $HOME/Library/Application Support/Code
//	Linux   : $XDG_CONFIG_HOME/Code   (fallback $HOME/.config/Code)
//
// Final path = <userDataRoot>/User/globalStorage/alefragnani.project-manager/projects.json
//
// See: spec/01-vscode-project-manager-sync/README.md

// User-data root segments per OS.
const (
	VSCodeUserDataRootDirName   = "Code"
	VSCodeUserDataMacRel        = "Library/Application Support/Code"
	VSCodeUserDataLinuxFallback = ".config/Code"
	VSCodeEnvAppData            = "APPDATA"
	VSCodeEnvUserProfile        = "USERPROFILE"
	VSCodeEnvHome               = "HOME"
	VSCodeEnvXDGConfigHome      = "XDG_CONFIG_HOME"
	VSCodeUserProfileAppDataRel = "AppData/Roaming/Code"
)

// Relative tail under the user-data root (constant across all OSes).
const (
	VSCodePMUserDir            = "User"
	VSCodePMGlobalStorageDir   = "globalStorage"
	VSCodePMExtensionDir       = "alefragnani.project-manager"
	VSCodePMProjectsFile       = "projects.json"
	VSCodePMProjectsTempSuffix = ".tmp"
	VSCodePMJSONIndent         = "\t"
)

// Default field values gitmap writes when inserting a NEW projects.json entry.
// Existing entries' values are preserved across re-syncs.
const (
	VSCodePMDefaultEnabled = true
	VSCodePMDefaultProfile = ""
)

// CLI flag for opting out of the automatic sync during scan.
const (
	FlagNoVSCodeSync     = "no-vscode-sync"
	FlagDescNoVSCodeSync = "skip syncing scanned repos into VS Code Project Manager projects.json"
)

// Global kill switch for the VS Code Project Manager sync. When the
// flag is passed at any position on the gitmap command line, OR when
// the env var is set to "1", every clone/scan/reclone helper that
// would otherwise call vscodepm.Sync short-circuits and prints
// MsgVSCodePMSyncDisabled instead. This is the "fully disable" lever
// for users who never want gitmap to touch projects.json — distinct
// from the per-command `--no-vscode-sync` opt-out.
const (
	FlagVSCodeSyncDisabled     = "vscode-sync-disabled"
	FlagDescVSCodeSyncDisabled = "fully disable VS Code Project Manager projects.json sync for this and all sub-invocations (sets GITMAP_VSCODE_SYNC_DISABLED=1)"
	EnvVSCodeSyncDisabled      = "GITMAP_VSCODE_SYNC_DISABLED"
	EnvVSCodeSyncDisabledOn    = "1"
	MsgVSCodePMSyncDisabled    = "  • VS Code Project Manager sync disabled (--vscode-sync-disabled / GITMAP_VSCODE_SYNC_DISABLED=1)\n"
)

// Debug-paths process-wide switch. The CLI flag (FlagDebugPaths in
// constants_cli.go) sets this env var to "1" for the current process
// so every code path that calls canonicalizePMPath — across every
// clone variant — emits the same trace line without N flag-plumbing
// edits. Reading from an env var (instead of a struct field) keeps
// the helper signature stable and lets future callers opt in by
// setting the env var directly in CI.
const (
	EnvDebugPaths      = "GITMAP_DEBUG_PATHS"
	EnvDebugPathsOn    = "1"
	MsgDebugPathsTrace = "[debug-paths] in=%q clean=%q resolved=%q\n"
)

// Error messages (Code Red zero-swallow policy).
const (
	ErrVSCodePMUserDataNotFound = "vscode: user data directory not found at %q (is VS Code installed?)\n"
	ErrVSCodePMExtDirMissing    = "vscode: project-manager extension dir not found at %q (open VS Code, install the alefragnani.project-manager extension, then retry)\n"
	ErrVSCodePMReadFailed       = "vscode: failed to read %s: %v\n"
	ErrVSCodePMParseFailed      = "vscode: %s is not valid JSON: %v (left untouched)\n"
	ErrVSCodePMWriteTempFailed  = "vscode: failed to write temp %s: %v\n"
	ErrVSCodePMRenameFailed     = "vscode: failed to commit %s: %v\n"
	ErrVSCodePMNoUserDataEnv    = "vscode: cannot determine user-data directory (no APPDATA / USERPROFILE / HOME env)\n"
	// ErrVSCodePMSyncBadMode is emitted when --mode is passed an
	// unknown literal. Args: bad-value, then the three accepted
	// literals (union, replace, intersection).
	ErrVSCodePMSyncBadMode = "vscode-pm-sync: unknown --mode %q (accepted: %s | %s | %s)"
)

// User-facing messages.
const (
	MsgVSCodePMSectionHeader = "  → VS Code Project Manager: %s\n"
	MsgVSCodePMSyncSummary   = "  ✓ projects.json synced: %d added, %d updated, %d unchanged (%d total)\n"
	MsgVSCodePMSyncSkipped   = "  • VS Code Project Manager sync skipped (--no-vscode-sync)\n"
	MsgVSCodePMRenamed       = "  ✓ projects.json: renamed %q -> %q\n"
	MsgVSCodePMRenameNoMatch = "  • projects.json: no entry matched %q (skipped rename)\n"

	// Diagnostic messages used by `gitmap vscode-pm-path` (v3.41.0+).
	MsgVSCodePMPathRootMissing = "vscode: user-data directory not found (is VS Code installed? checked APPDATA / HOME / XDG_CONFIG_HOME)"
	MsgVSCodePMPathExtMissing  = "vscode: project-manager extension storage dir not found near %s (open VS Code, install the alefragnani.project-manager extension, then retry)\n"

	// vscode-pm-sync (v4.36.0+) — full re-tag of every projects.json entry.
	MsgVSCodePMSyncStart     = "→ vscode-pm-sync: re-tagging projects.json entries at %s\n"
	MsgVSCodePMSyncDryRun    = "  • dry-run: %d entries scanned, %d would change (no write)\n"
	MsgVSCodePMSyncEntryStat = "  ✓ scanned %d entries, %d skipped (rootPath missing on disk)\n"
	MsgVSCodePMSyncEmptyFile = "  • projects.json contains 0 entries — nothing to re-tag\n"
)
