package constants

// gitmap install ctx — Windows right-click context menu (v1).
// Spec: 02-spec/04-generic-cli/30-install-ctx.md.

// Tool name for the install dispatcher.
const (
	ToolCtx            = "ctx"
	ToolContextMenu    = "context-menu"
	ToolContextMenuAlt = "contextmenu"
	ToolCtxDesc        = "Add gitmap to the OS right-click context menu (Windows / macOS / Linux)"
)

// gitmap:cmd top-level
const (
	CmdContextMenu      = "context-menu" // gitmap:cmd skip
	CmdContextMenuAlias = "contextmenu"  // gitmap:cmd skip
)

// Top-level cascade label and registry root key names. The three roots
// give us the menu on files (right-clicking any file item), on folder
// backgrounds (clicking inside a folder) and on folder items (right-clicking
// the folder itself).
const (
	CtxRootKeyBackground = `HKCU\Software\Classes\Directory\Background\shell\gitmap`
	CtxRootKeyDirectory  = `HKCU\Software\Classes\Directory\shell\gitmap`
	CtxRootKeyFile       = `HKCU\Software\Classes\*\shell\gitmap`
	CtxRootMUIVerb       = "gitmap"
)

// CtxMode moved to cli/enums/ctxmodetype (spec 243, Wave D) — see the
// ctxmodetype.VariantType enum. The constants surface keeps no alias; use the
// enum package directly.

// CtxIconExeToken is the placeholder swapped for the resolved gitmap
// binary path when an entry's Icon field is rendered into a registry
// value. Lets ctxMenu() declare Icons without knowing the exe path.
const CtxIconExeToken = "{exe}"

// Common per-entry Icon values for Windows context-menu entries.
// Format: "<path>,<index>" — Windows shell uses the indexed icon
// resource inside the binary/DLL. Index 0 picks the default icon.
const (
	CtxIconGitmap        = CtxIconExeToken + ",0"
	CtxIconPwsh          = `C:\Program Files\PowerShell\7\pwsh.exe,0`
	CtxIconShowLog       = "shell32.dll,264"
	CtxIconDiff          = "shell32.dll,246"
	CtxIconBranchCompare = "shell32.dll,133"
	CtxIconCommit        = "shell32.dll,297"
	CtxIconPull          = "shell32.dll,146"
	CtxIconPush          = "shell32.dll,147"
	CtxIconScan          = "shell32.dll,22"
	CtxIconClone         = "shell32.dll,134"
	CtxIconRelease       = "shell32.dll,265"
	CtxIconRepos         = "shell32.dll,4"
	CtxIconWebUI         = "shell32.dll,14"
	CtxIconSettings      = "shell32.dll,315"
	CtxIconTerminal      = "shell32.dll,220"
	CtxIconDocs          = "shell32.dll,266"
	CtxIconHelp          = "shell32.dll,23"
	CtxIconInstall       = "shell32.dll,260"
	CtxIconMacro         = "shell32.dll,250"
	CtxIconVisibility    = "shell32.dll,174"
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
