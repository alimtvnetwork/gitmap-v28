package constants

// Constants for `gitmap clone-pick <url> <paths>` (spec 100, v3.153.0+).
//
// `clone-pick` performs a partial, sparse-checkout-based clone of a
// subset of a git repository into the current working directory. Unlike
// `clone-from` (manifest-driven) or `clone-now` (round-trip from scan
// output), `clone-pick` operates on a single repo + a comma-separated
// list of repo-relative paths.
//
// Persistence: every successful run writes one row to the
// CloneInteractiveSelection table (see constants_clonepick_store.go) so
// the same selection can be replayed later via --replay <id|name>.
//
// Naming note: the short alias is `cpk`, NOT `ci` (which collides with
// CI/CD muscle memory). See spec/01-app/100-clone-pick.md §3.

// gitmap:cmd top-level
// CLI surface. Both names are referenced by rootcore.go's dispatcher
// and indirectly by the AST parity test (which requires every Cmd*
// constant declared under a `gitmap:cmd top-level` block to also
// appear in topLevelCmds()).
const (
	CmdClonePick      = "clone-pick"
	CmdClonePickAlias = "cpk"
)

// Flag names + descriptions. Long-form only; short flags are
// reserved for very-frequent operations per the project convention.
const (
	FlagClonePickAsk     = "ask"
	FlagDescClonePickAsk = "Open an interactive tree picker before fetching. " +
		"Paths supplied on the command line are pre-checked."

	FlagClonePickName     = "name"
	FlagDescClonePickName = "Optional human label saved with this selection. " +
		"Reuse later with --replay <name>."

	FlagClonePickReplay     = "replay"
	FlagDescClonePickReplay = "Re-run a previously saved selection by id or name. " +
		"When set, <paths> is not required."

	FlagClonePickBranch     = "branch"
	FlagDescClonePickBranch = "Branch to check out (passed to git clone --branch). " +
		"Empty (default) lets git use the remote's HEAD."

	FlagClonePickMode     = "mode"
	FlagDescClonePickMode = "URL form to use when <repo-url> is shorthand " +
		"(owner/repo): 'https' (default) or 'ssh'."

	FlagClonePickDepth     = "depth"
	FlagDescClonePickDepth = "Shallow clone depth. 0 = full history (default: 1)."

	FlagClonePickCone     = "cone"
	FlagDescClonePickCone = "Use sparse-checkout cone mode (faster, folder-only). " +
		"Auto-disabled when any path looks like a glob or a file."

	FlagClonePickDest     = "dest"
	FlagDescClonePickDest = "Destination directory. Created if missing. " +
		"Default: '.' (current directory)."

	FlagClonePickKeepGit     = "keep-git"
	FlagDescClonePickKeepGit = "Leave the .git directory in <dest>. " +
		"Set --keep-git=false for a files-only checkout."

	FlagClonePickDryRun     = "dry-run"
	FlagDescClonePickDryRun = "Print the plan and the git commands without executing. " +
		"No DB write happens in dry-run mode."

	FlagClonePickQuiet     = "quiet"
	FlagDescClonePickQuiet = "Suppress per-step progress on stderr."

	FlagClonePickForce     = "force"
	FlagDescClonePickForce = "Allow non-empty <dest>. Default refuses to clobber."
)

// Mode enum. Stable: surfaced in dry-run header + DB rows.
const (
	ClonePickModeHTTPS = "https"
	ClonePickModeSSH   = "ssh"
)

// ClonePickAutoExclude lists folder names auto-greyed in the picker
// and pre-unchecked even when matched by a glob. Override per-project
// via the `clonePick.autoExclude` config array.
//
// Conservative defaults: things every developer wants out of
// sparse-checkouts almost always (build outputs, vendored deps,
// language caches). The list is small on purpose -- false positives
// here are silent data loss for the user.
var ClonePickAutoExclude = []string{
	".git",
	"node_modules",
	"vendor",
	"dist",
	"build",
	"__pycache__",
	".venv",
	"target", // Rust / Java
}

// ClonePickPathMaxBytes caps a single path entry at git's documented
// sparse-checkout pattern length limit. Anything longer is rejected
// up-front with ErrClonePickPathTooLong.
const ClonePickPathMaxBytes = 4096

// User-facing messages. Trailing newlines are baked in so call sites
// don't need to remember them.
const (
	// %s = repo url, %s = dest, %s = mode, %s = branch-or-(default),
	// %d = depth, %s = "cone"/"non-cone", %d = path count.
	MsgClonePickDryHeader = "gitmap clone-pick: dry-run\n" +
		"repo:   %s\n" +
		"dest:   %s\n" +
		"mode:   %s   branch: %s   depth: %d   sparse: %s\n" +
		"%d path(s) -- pass without --dry-run to actually clone\n\n"

	// %s = canonical repo id, %d = selection id, %s = name (or "(unnamed)").
	MsgClonePickSaved = "saved selection #%d for %s (%s)\n"

	// %s = canonical repo id, %d = selection id, %s = name.
	MsgClonePickReplayed = "replayed selection #%d for %s (%s)\n"

	MsgClonePickMissingURL   = "clone-pick: <repo-url> argument is required"
	MsgClonePickMissingPaths = "clone-pick: <paths> argument is required " +
		"(or pass --replay <id|name>)"

	// %s = bad path.
	MsgClonePickPathEmpty     = "clone-pick: empty path entry in --paths list"
	MsgClonePickPathAbsolute  = "clone-pick: absolute paths are not allowed (%s)"
	MsgClonePickPathTraversal = "clone-pick: '..' path traversal not allowed (%s)"
	MsgClonePickPathTooLong   = "clone-pick: path exceeds 4096 bytes (%s)"

	MsgClonePickDestDirty = "clone-pick: <dest> is not empty (use --force to override)"

	// %s = lookup key.
	MsgClonePickReplayNotFound  = "clone-pick: --replay: no saved selection matches %q"
	MsgClonePickReplayAmbiguous = "clone-pick: --replay: multiple selections match %q; " +
		"use --replay <id> to disambiguate"

	MsgClonePickUserCancelled = "clone-pick: canceled by user"
)

// Errors. printf-style verbs documented inline.
const (
	// %s = bad value.
	ErrClonePickBadMode = "clone-pick: --mode must be 'https' or 'ssh', got %q"
	// %d = bad value.
	ErrClonePickBadDepth = "clone-pick: --depth must be >= 0, got %d"
	// %v = err.
	ErrClonePickGitClone       = "clone-pick: git clone failed: %v"
	ErrClonePickGitSparseInit  = "clone-pick: git sparse-checkout init failed: %v"
	ErrClonePickGitSparseSet   = "clone-pick: git sparse-checkout set failed: %v"
	ErrClonePickGitCheckout    = "clone-pick: git checkout failed: %v"
	ErrClonePickGitLsTree      = "clone-pick: git ls-tree failed: %v"
	ErrClonePickFsCreateDest   = "clone-pick: create dest dir: %v"
	ErrClonePickFsRemoveDotGit = "clone-pick: remove .git: %v"
	ErrClonePickDBInsert       = "clone-pick: save selection: %v"
	ErrClonePickDBLookup       = "clone-pick: lookup selection: %v"
	ErrClonePickPickerLaunch   = "clone-pick: launch picker: %v"
	// %v = err. Emitted when promoting the picker's pre-cloned
	// metadata dir into DestDir fails (rename + cross-fs copy
	// fallback both errored).
	ErrClonePickPromoteSrc = "clone-pick: promote pre-cloned source: %v"
)

// User-cancel exit code. 130 mirrors the SIGINT convention so shell
// scripts can branch on `$? -eq 130` cleanly.
const ClonePickExitUserCancel = 130
