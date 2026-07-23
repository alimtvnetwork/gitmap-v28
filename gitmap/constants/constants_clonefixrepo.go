// Package constants — clone-fix-repo command IDs, messages, and
// exit codes for `gitmap clone-fix-repo` (cfr) and
// `gitmap clone-fix-repo-pub` (cfrp).
//
// These commands chain `clone` → `fix-repo` (and optionally
// `make-public`) into a single invocation. See cmd/clonefixrepo.go.
package constants

// gitmap:cmd top-level
// Clone-fix-repo command IDs and short aliases.
const (
	CmdCloneFixRepo         = "clone-fix-repo"
	CmdCloneFixRepoAlias    = "cfr"
	CmdCloneFixRepoPub      = "clone-fix-repo-pub"
	CmdCloneFixRepoPubAlias = "cfrp"
)

// Clone-fix-repo help-line entries surfaced by `gitmap help`.
const (
	HelpCloneFixRepo    = "  clone-fix-repo (cfr) <url> [folder]      Clone, then run fix-repo --all in the new folder"
	HelpCloneFixRepoPub = "  clone-fix-repo-pub (cfrp) <url> [folder] Clone, fix-repo --all, then make-public --yes"
)

// Clone-fix-repo user-facing messages and errors.
const (
	MsgCloneFixRepoDone        = "clone-fix-repo: pipeline completed in %s\n"
	MsgCloneFixRepoSkipNoVer   = "  fix-repo: skipped (repo %q has no -vN suffix, nothing to rewrite)\n    pass --require-version to fail instead.\n"
	WarnCloneFixRepoRemoteFmt  = "  Warning: could not resolve cloned repo remote from %q: %v\n"
	ErrCloneFixRepoUsage       = "clone-fix-repo: ERROR <url> is required\n  usage: gitmap clone-fix-repo <url> [folder]\n  usage: gitmap clone-fix-repo-pub <url> [folder]\n"
	ErrCloneFixRepoChdirFmt    = "clone-fix-repo: ERROR cannot cd into %q: %v\n"
	ErrCloneFixRepoExecFmt     = "clone-fix-repo: ERROR could not run chained step: %v\n"
	ErrCloneFixRepoNeedVersion = "clone-fix-repo: ERROR --require-version set but repo %q has no -vN suffix\n"
	ErrCloneFixRepoRemoteParse = "unparseable remote URL"

	// MsgCFRFolderTransport fires when cfr/cfrp rewrites the user's
	// positional URL to match the destination folder's existing
	// origin transport. Format: scheme, before, after.
	MsgCFRFolderTransport = "clone-fix-repo: rewriting URL to %s to match existing folder origin: %s → %s\n"
	// WarnCFRFolderTransport surfaces non-fatal transport-detection
	// failures (existing origin unreadable, URL rewrite failed).
	// Format: absPath, reason, err.
	WarnCFRFolderTransport      = "clone-fix-repo: warning: %s: %s: %v\n"
	WarnCFRFolderTransportNoErr = "clone-fix-repo: warning: %s: %s\n"

	// Nested-repo escape (v6.48.0+): when cwd is itself a git repo,
	// cfr/cfrp walks up to the first non-repo ancestor before
	// cloning, so the new tree is never nested inside another repo's
	// git context (which caused `fetch-pack: invalid index-pack
	// output` on Windows).
	MsgCFREscapeNested = "  " + ColorCyan + "↑ cfr: cwd is a git repo — escaping to non-repo ancestor" + ColorReset + "\n" +
		"    " + ColorDim + "from:" + ColorReset + " %s\n" +
		"    " + ColorDim + "  to:" + ColorReset + " %s\n"
	WarnCFREscapeChdir = "clone-fix-repo: warning: could not chdir to %s: %v\n"

	// Parallel comma-separated URL fan-out (v6.54.0+). Each URL is
	// re-execed through the single-URL pipeline so chdir/fix-repo
	// chaining stays isolated per worker.
	MsgCloneFixRepoParallelHeader   = ColorCyan + "▶ clone-fix-repo: " + ColorReset + "%d URL(s) across %d worker(s) [%s]\n"
	MsgCloneFixRepoParallelItem     = ColorDim + "  [%d/%d] start " + ColorReset + "%s\n"
	MsgCloneFixRepoParallelItemOk   = ColorGreen + "  [%d/%d] ✓ " + ColorReset + "%s " + ColorDim + "(%s)" + ColorReset + "\n"
	MsgCloneFixRepoParallelItemFail = ColorRed + "  [%d/%d] ✗ " + ColorReset + "%s " + ColorDim + "(%s)" + ColorReset + ": %v\n"
	MsgCloneFixRepoParallelDoneOk   = ColorGreen + "✔ clone-fix-repo: all %d URL(s) ok" + ColorReset + "\n"
	MsgCloneFixRepoParallelDoneFail = ColorYellow + "⚠ clone-fix-repo: %d ok, %d failed" + ColorReset + "\n"
)

// Clone-fix-repo flags.
const (
	FlagRequireVersion          = "require-version"
	CloneFixRepoDefaultParallel = 8
)

// Clone-fix-repo exit codes.
const (
	ExitCloneFixRepoOk          = 0
	ExitCloneFixRepoBadFlag     = 6
	ExitCloneFixRepoChdir       = 9
	ExitCloneFixRepoChainFailed = 10
)
