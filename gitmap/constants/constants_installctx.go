package constants

// gitmap install ctx — Windows right-click context menu (v1).
// Spec: spec/04-generic-cli/30-install-ctx.md.

// Tool name for the install dispatcher.
const (
	ToolCtx     = "ctx"
	ToolCtxDesc = "Add gitmap to the OS right-click context menu (Windows / macOS / Linux)"
)

// Top-level cascade label and registry root key names. The two roots
// give us the menu both on folder backgrounds (clicking inside a
// folder) and on folder items (right-clicking the folder itself).
const (
	CtxRootKeyBackground = `HKCU\Software\Classes\Directory\Background\shell\gitmap`
	CtxRootKeyDirectory  = `HKCU\Software\Classes\Directory\shell\gitmap`
	CtxRootMUIVerb       = "gitmap"
)

// CtxMode controls how an entry is wired into the registry.
type CtxMode = string

// CtxMode values.
const (
	CtxModeTerminal CtxMode = "terminal" // pwsh -NoExit, command runs and window stays open
	CtxModeSilent   CtxMode = "silent"   // pwsh -WindowStyle Hidden, output via notifier
	CtxModePrefill  CtxMode = "prefill"  // pwsh -NoExit + writes "gitmap " prompt, no command run
)

// CtxIconExeToken is the placeholder swapped for the resolved gitmap
// binary path when an entry's Icon field is rendered into a registry
// value. Lets ctxMenu() declare Icons without knowing the exe path.
const CtxIconExeToken = "{exe}"

// Common per-entry Icon values for Windows context-menu entries.
// Format: "<path>,<index>" — Windows shell uses the indexed icon
// resource inside the binary/DLL. Index 0 picks the default icon.
const (
	CtxIconGitmap = CtxIconExeToken + ",0"
	CtxIconPwsh   = `C:\Program Files\PowerShell\7\pwsh.exe,0`
)

// User-facing labels and exec messages.
const (
	MsgCtxInstallStart    = "  Adding gitmap to Windows context menu...\n"
	MsgCtxInstallDone     = "  ✓ gitmap context menu installed (%d/%d registry keys).\n"
	MsgCtxUninstallStart  = "  Removing gitmap from Windows context menu...\n"
	MsgCtxUninstallDone   = "  ✓ gitmap context menu removed (%d/%d registry keys).\n"
	MsgCtxRegFail         = "  ! Registry command failed: %v\n"
	MsgCtxOSUnsupported   = "  Error: ctx is not supported on this OS yet (current OS: %s). Supported: windows, darwin, linux.\n"
	MsgCtxOpenTerminalLbl = "Open terminal here"
	MsgCtxDocsLbl         = "Docs"
)

// Raw git context-menu entries. These bypass gitmap and shell out to
// `git` directly so users can inspect history/diff/log on the clicked
// folder. Operate on cwd; for file-scoped views the user can pass the
// path interactively from the opened terminal.
const (
	CtxExeGit          = "git"
	CtxGitHistoryLabel = "History (git log graph)"
	CtxGitDiffLabel    = "Diff (git diff)"
	CtxGitLogLabel     = "Log (git log)"
	CtxGitStatusLabel  = "Status (git status)"
)

// Raw git argument vectors. Kept here to avoid magic strings in
// installctxentries.go (constants-only policy).
var (
	CtxGitHistoryArgs = []string{"log", "--oneline", "--graph", "--decorate", "--all", "-n", "100"}
	CtxGitDiffArgs    = []string{"diff", "--stat", "HEAD"}
	CtxGitLogArgs     = []string{"log", "-n", "30", "--pretty=format:%h %ad %s", "--date=short"}
	CtxGitStatusArgs  = []string{"status"}
)
