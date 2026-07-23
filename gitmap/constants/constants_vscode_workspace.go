package constants

// VS Code multi-root workspace (`*.code-workspace`) constants.
//
// A `.code-workspace` file is a small JSON document VS Code opens via
// File → Open Workspace from File…. Schema (only the parts gitmap
// emits — extensions/settings are kept as `{}` to match a fresh
// VS Code "Save Workspace As…" file byte-for-byte where possible):
//
//	{
//	  "folders": [ { "name": "<repo>", "path": "<abs-or-rel-path>" } ],
//	  "settings": {}
//	}
//
// Source-of-truth is the same Repo table that drives the Project
// Manager sync, so the workspace file always mirrors the latest
// `gitmap scan` output without a separate plumbing path.
const (
	VSCodeWorkspaceDefaultFilename = "gitmap.code-workspace"
	VSCodeWorkspaceJSONIndent      = "\t"
	VSCodeWorkspaceTempSuffix      = ".tmp"
)

// CLI flag IDs for `gitmap vscode-workspace`.
const (
	FlagVSCodeWorkspaceOut            = "out"
	FlagDescVSCodeWorkspaceOut        = "output `.code-workspace` file path (default: ./" + VSCodeWorkspaceDefaultFilename + ")"
	FlagVSCodeWorkspaceRelative       = "relative"
	FlagDescVSCodeWorkspaceRelative   = "emit folder paths relative to the workspace file's directory"
	FlagVSCodeWorkspaceTag            = "tag"
	FlagDescVSCodeWorkspaceTag        = "include only repos whose auto-detected tag set contains the given tag (e.g. go, node, git)"
	FlagVSCodeWorkspaceRootSubdir     = "root-subdir"
	FlagDescVSCodeWorkspaceRootSubdir = "add `<repo>/<subdir>` as the workspace folder instead of the repo root; repos without that subdir are skipped"
)

// User-facing messages.
const (
	MsgVSCodeWorkspaceWritten    = "  ✓ wrote %s with %d folder(s)\n"
	MsgVSCodeWorkspaceEmpty      = "  • no tracked repos — workspace not written (run `gitmap scan` first)\n"
	MsgVSCodeWorkspaceSubdirSkip = "  • skipped %s: subdir %q not found\n"
)

// Error templates (Code Red zero-swallow policy).
const (
	ErrVSCodeWorkspaceDBOpen     = "vscode-workspace: failed to open repo database: %v\n"
	ErrVSCodeWorkspaceDBList     = "vscode-workspace: failed to list repos: %v\n"
	ErrVSCodeWorkspaceWriteTemp  = "vscode-workspace: failed to write temp file %q: %v\n"
	ErrVSCodeWorkspaceRename     = "vscode-workspace: failed to commit %q: %v\n"
	ErrVSCodeWorkspaceRelativize = "vscode-workspace: failed to relativize %q against %q: %v\n"
)
