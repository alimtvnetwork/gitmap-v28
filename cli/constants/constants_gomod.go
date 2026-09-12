package constants

// gitmap:cmd top-level
// GoMod CLI commands.
const (
	CmdGoMod      = "gomod"
	CmdGoModAlias = "gm"
)

// GoMod help text.
const (
	HelpGoMod      = "  gomod (gm) <path>   Rename Go module path across repo with branch safety"
	HelpGoModFlags = "GoMod flags:"
	HelpGoModDry   = "  --dry-run           Preview changes without modifying files or branches"
	HelpGoModNoMrg = "  --no-merge          Commit on feature branch but do not merge back"
	HelpGoModNoTdy = "  --no-tidy           Skip go mod tidy after replacement"
	HelpGoModVerb  = "  --verbose           Print each file path as it is modified"
	HelpGoModExt   = "  --ext <exts>        Comma-separated extensions to filter (e.g. *.go,*.md); default: all files"
)

// GoMod flag names.
const (
	FlagGoModDryRun  = "dry-run"
	FlagGoModNoMerge = "no-merge"
	FlagGoModNoTidy  = "no-tidy"
	FlagGoModExt     = "ext"
)

// GoMod flag descriptions.
const (
	FlagDescGoModDryRun  = "Preview changes without modifying files or branches"
	FlagDescGoModNoMerge = "Commit on feature branch but do not merge back"
	FlagDescGoModNoTidy  = "Skip go mod tidy after replacement"
	FlagDescGoModExt     = "Comma-separated file extensions to filter (e.g. *.go,*.md); default: all files"
)

// GoMod file and directory constants.
const (
	GoModFile       = "go.mod"
	GoModModuleLine = "module "
	GoFileExt       = ".go"
)

// GoMod excluded directories during file walk.
var GoModExcludeDirs = []string{".git", "vendor", "node_modules"}

// GoMod branch prefixes.
const (
	GoModFeaturePrefix = "feature/replace-"
	GoModBackupPrefix  = "backup/before-replace-"
)

// GoMod terminal messages.
const (
	MsgGoModSummary       = "✔ Module path renamed\n"
	MsgGoModOld           = "  Old: %s\n"
	MsgGoModNew           = "  New: %s\n"
	MsgGoModFiles         = "  Files updated: %d\n"
	MsgGoModBackupBranch  = "  Backup branch: %s\n"
	MsgGoModFeatureBranch = "  Feature branch: %s\n"
	MsgGoModMergedInto    = "  Merged into: %s\n"
	MsgGoModLeftOn        = "  Left on branch: %s\n"
	MsgGoModVerboseFile   = "  replaced: %s\n"
	MsgGoModDryHeader     = "gomod (dry-run): would rename module path\n"
	MsgGoModDryOld        = "  Old: %s\n"
	MsgGoModDryNew        = "  New: %s\n"
	MsgGoModDryFiles      = "  Files that would change: %d\n"
	MsgGoModDryFile       = "  %s\n"
	MsgGoModNoImports     = "Warning: no files found containing the old path to replace (only go.mod updated)\n"
	MsgGoModTidyWarn      = "Warning: go mod tidy failed: %v (continuing)\n"
	MsgGoModNothingRename = "module path is already %s, nothing to rename\n"
)

// GoMod error messages.
const (
	ErrGoModUsage         = "usage: gitmap gomod <new-module-path> [--ext *.go,*.md] [--dry-run] [--no-merge] [--no-tidy] [--verbose]\n"
	ErrGoModNoFile        = "error: go.mod not found in current directory\n"
	ErrGoModNoModule      = "error: no module directive found in go.mod\n"
	ErrGoModNotRepo       = "error: not inside a git repository\n"
	ErrGoModDirtyTree     = "error: working tree has uncommitted changes, commit or stash first\n"
	ErrGoModBranchExists  = "error: branch %s already exists, aborting\n"
	ErrGoModMergeConflict = "error: merge conflict — resolve manually on %s\n"
	ErrGoModReadFailed    = "error: failed to read %s: %v (operation: read)\n"
	ErrGoModWriteFailed   = "error: failed to write %s: %v (operation: write)\n"
	ErrGoModCommitFailed  = "error: git commit failed: %v\n"
)

// GoMod git arguments.
const (
	GitAdd         = "add"
	GitAddAll      = "-A"
	GitCommit      = "commit"
	GitCommitMsg   = "-m"
	GitMerge       = "merge"
	GitMergeNoFF   = "--no-ff"
	GitStatusShort = "--porcelain"
	GitStatus      = "status"
)

// GoMod commit message format.
const GoModCommitMsgFmt = "refactor: rename go module path\n\nOld: %s\nNew: %s\n\nReplaced module directive in go.mod and all matching paths\nacross %d files."

// GoMod merge message format.
const GoModMergeMsgFmt = "merge: module rename to %s"
