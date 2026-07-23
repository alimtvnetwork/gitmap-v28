// Package constants — visibility command IDs, flags, messages, and
// exit codes for `gitmap make-public` / `gitmap make-private`.
//
// The two commands are thin wrappers around the host platform's CLI
// (`gh` for GitHub, `glab` for GitLab). They:
//
//  1. Resolve provider + owner/repo from `git remote get-url origin`.
//  2. Read the current visibility via the provider CLI.
//  3. Skip if already in the target state (idempotent).
//  4. Prompt the user when going private → public (skip with --yes).
//  5. Apply, then verify the change took effect.
//
// Spec parity: spec-authoring/23-visibility-change/01-spec.md
// (PowerShell reference: visibility-change.ps1).
package constants

// Visibility command IDs live in constants_cli.go (CmdMakePublic /
// CmdMakePrivate) per the project-wide rule that all CLI tokens are
// centralized there. This file owns everything else (target tokens,
// flags, messages, exit codes).

// Visibility target tokens — what the provider CLI expects, what the
// user can type for the (optional) explicit-target form, and what we
// store/print internally.
const (
	VisibilityPublic  = "public"
	VisibilityPrivate = "private"

	VisShortPub = "pub"
	VisShortPri = "pri"
)

// Visibility flags. --yes skips the private→public confirmation;
// --dry-run prints what would change without invoking the provider
// CLI; --verbose echoes each shell command before running it.
const (
	FlagVisYes     = "yes"
	FlagVisYesAlt  = "y"
	FlagVisDryRun  = "dry-run"
	FlagVisVerbose = "verbose"

	FlagDescVisYes     = "Skip the private→public confirmation prompt (no-op for public→private)."
	FlagDescVisDryRun  = "Print the provider CLI command that would run; do not invoke it."
	FlagDescVisVerbose = "Echo every shell command to stderr before running it."
)

// Provider tokens — match what we detect from the origin URL host.
const (
	ProviderGitHub = "github"
	ProviderGitLab = "gitlab"

	HostGitHub = "github.com"
	HostGitLab = "gitlab.com"

	CLIGitHub = "gh"
	CLIGitLab = "glab"
)

// Visibility help-line entries surfaced by `gitmap help` (Utilities).
const (
	HelpMakePublic  = "  make-public         Make current repo public on GitHub/GitLab (gh/glab required)"
	HelpMakePrivate = "  make-private        Make current repo private on GitHub/GitLab (gh/glab required)"
)

// Visibility user-facing messages.
const (
	MsgVisAlreadyFmt   = "visibility: already %s on %s\n"
	MsgVisChangedFmt   = "visibility: %s → %s on %s (%s)\n"
	MsgVisDryRunFmt    = "[dry-run] visibility: %s → %s on %s (%s)\n"
	MsgVisConfirmFmt   = "Make %s PUBLIC on %s? Type 'yes' to confirm: "
	MsgVisVerboseExec  = "+ %s %s\n"
	MsgVisVerifyOK     = "  ✓ verified: visibility is now %s\n"
	MsgVisLocalSkipFmt = "visibility: skipping local remote %q (no provider CLI applies to file:// / local paths)\n"
)


// Visibility error messages.
const (
	ErrVisNotInRepo       = "visibility: not a git repository\n"
	ErrVisNoOrigin        = "visibility: no `origin` remote configured\n"
	ErrVisBadProviderFmt  = "visibility: unsupported host in %q (only github.com / gitlab.com are supported)\n"
	ErrVisBadSlugFmt      = "visibility: cannot parse owner/repo from %q\n"
	ErrVisCLIMissingFmt   = "visibility: %q not found on PATH (install: https://cli.github.com or https://gitlab.com/gitlab-org/cli)\n"
	ErrVisReadCurrentFmt  = "visibility: cannot read current visibility (auth via `%s auth login`?): %v\n"
	ErrVisConfirmRequired = "visibility: confirmation required (re-run with --yes for non-interactive use)\n"
	ErrVisApplyFailedFmt  = "visibility: apply failed: %v\n%s"
	ErrVisVerifyFailedFmt = "visibility: verification failed — current is %q, expected %q\n"
)

// Visibility exit codes (mirrored from visibility-change.ps1 so wrappers
// and CI can branch on the same numbers).
const (
	ExitVisOK           = 0
	ExitVisNotARepo     = 2
	ExitVisNoOrigin     = 3
	ExitVisBadProvider  = 4
	ExitVisAuthFailed   = 5
	ExitVisBadFlag      = 6
	ExitVisConfirmReq   = 7
	ExitVisVerifyFailed = 8
)

// Spec 113 — bulk visibility + cfrp prior-version privatize messages.
const (
	MsgVisBulkHeaderFmt    = "visibility: %s × %d versions of %s on %s\n"
	MsgVisBulkItemFmt      = "  [%d/%d] %s … "
	MsgVisBulkSkipFmt      = "already %s\n"
	MsgVisBulkOKFmt        = "%s → %s\n"
	MsgVisBulkFailFmt      = "FAILED (%v)\n"
	MsgVisBulkDryFmt       = "[dry-run] %s → %s (slug=%s)\n"
	MsgCFRPPriorHeaderFmt  = "\ncfrp: scanning prior versions of %s (≤%d back)…\n"
	MsgCFRPPriorFoundFmt   = "cfrp: %d prior version(s) currently public: %s\n"
	MsgCFRPPriorPromptFmt  = "Privatize all %d prior version(s)? [y/N]: "
	MsgCFRPPriorNoneFound  = "cfrp: no prior public versions found.\n"
	MsgCFRPPriorSkipped    = "cfrp: leaving prior versions unchanged.\n"
	ErrVisBulkBadCountFmt  = "visibility: count must be a positive integer, got %q\n"
	ErrVisBulkRepoParseFmt = "visibility: cannot parse repo identity from %q\n"
)

// CFRPPriorMaxLookback caps the prior-version probe at v(N-5)..v(N-1).
// v6.63.0: narrowed from 15 → 5 per user request (privatize the last
// 5 prior versions automatically when cfrp succeeds).
const CFRPPriorMaxLookback = 5

// Spec 116 — bulk wildcard visibility (make-all-public / make-all-private
// / MAPUB / MAPRI). Owner-only resolver + repo-list pagination cap +
// interactive prompt strings.
const (
	ProviderUnknownReason   = "unknown"
	MsgMakeAllNotImpl       = "make-all-*: handler not yet wired (spec/01-app/116)\n"
	ErrMakeAllResolveFmt    = "make-all-*: cannot resolve owner: %v\n"
	ErrMakeAllMissingArgFmt = "make-all-*: usage: %s <target> <patterns> [-Y|--yes]\n"

	// OwnerRepoListLimit caps `gh/glab repo list --limit`. Owners with
	// more than this many repos will hit a WARNING (see plan step 26).
	OwnerRepoListLimit      = 1000
	WarnOwnerRepoListCapFmt = "make-all-*: WARNING — owner %[2]s returned %[1]d repos (the --limit cap). Repos beyond the cap were NOT enumerated; narrow the patterns or raise the limit.\n"

	// Interactive prompt / table strings (spec §4).
	MsgBulkMatchedHeaderFmt = "Matched %d of %d repos under %s:\n"
	MsgBulkConfirmFmt       = "\nProceed with these %d repo(s)? [y/N/exclude e.g. 1,3-5]: "
	MsgBulkExcludedFmt      = "  → excluded %d, %d remaining. Re-prompting.\n"
	ErrBulkExclusionFmt     = "  ! %v — try again, or type 'y' / 'n'.\n"
	ErrBulkPromptEOF        = "make-all-*: prompt aborted (EOF on stdin)\n"

	// -Y flag + apply-loop strings (plan steps 13-14).
	FlagVisYesUpper       = "Y"
	MsgBulkNoMatches      = "make-all-*: no repos matched the supplied pattern(s); nothing to do.\n"
	MsgBulkApplyHeaderFmt = "\nmake-all-*: applying visibility=%s to %d repo(s) on %s\n"
	MsgBulkApplyItemFmt   = "  [%d/%d] %s … "
	MsgBulkApplySkipFmt   = "already %s\n"
	MsgBulkApplyOKFmt     = "%s → %s\n"
	MsgBulkApplyFailFmt   = "FAILED (%v)\n"
	MsgBulkSummaryFmt     = "make-all-*: %d changed, %d skipped, %d failed of %d\n"
	MsgBulkAborted        = "make-all-*: aborted by user (no changes applied).\n"
)

// Bulk-apply exit codes layered on top of ExitVis*. ExitVisBulkPartial
// signals "at least one repo failed but others succeeded" so CI can
// branch separately from a hard auth failure.
const (
	ExitVisBulkPartial = 9
)

// Spec 116 §parallel — bulk wildcard visibility parallelism and TTL
// cache for the owner-wide repo enumeration.
const (
	// DefaultBulkParallelism is the worker count used when --parallel
	// is not supplied. Tuned to keep provider API rate-limits happy
	// while still cutting wall-clock by ~6-8x on large owners.
	DefaultBulkParallelism = 8

	// MaxBulkParallelism caps --parallel to prevent accidental
	// thundering-herd against gh / glab.
	MaxBulkParallelism = 32

	// OwnerRepoListCacheTTLSeconds is the default TTL for the
	// SQLite-backed owner repo-list cache. 5 minutes balances
	// freshness vs. avoiding repeated `gh repo list` round-trips
	// when users iterate on patterns. Overridable per-invocation
	// via --cache-ttl=<seconds> (0 disables the cache).
	OwnerRepoListCacheTTLSeconds = 300

	// SettingOwnerRepoListCacheTTL is the Setting key consulted
	// before falling back to the compiled default. Persisted as a
	// string of seconds.
	SettingOwnerRepoListCacheTTL = "owner_repo_list_ttl_seconds"

	// SettingBulkParallelism overrides DefaultBulkParallelism when set.
	SettingBulkParallelism = "bulk_visibility_parallelism"

	// Owner repo-list cache table.
	TableOwnerRepoListCache     = "OwnerRepoListCache"
	SQLCreateOwnerRepoListCache = `CREATE TABLE IF NOT EXISTS OwnerRepoListCache (
		Provider   TEXT NOT NULL,
		Owner      TEXT NOT NULL,
		NamesJson  TEXT NOT NULL,
		FetchedAt  TEXT NOT NULL,
		PRIMARY KEY (Provider, Owner)
	)`

	// Owner repo-name index — one row per repo with BaseName +
	// VersionNumber pre-parsed so make-last-* / fuzzy fallback can
	// query without re-walking the JSON blob. Populated alongside
	// the cache write.
	TableOwnerRepoNameIndex     = "OwnerRepoNameIndex"
	SQLCreateOwnerRepoNameIndex = `CREATE TABLE IF NOT EXISTS OwnerRepoNameIndex (
		Provider       TEXT NOT NULL,
		Owner          TEXT NOT NULL,
		RepoName       TEXT NOT NULL,
		BaseName       TEXT NOT NULL,
		VersionNumber  INTEGER NOT NULL DEFAULT -1,
		FetchedAt      TEXT NOT NULL,
		PRIMARY KEY (Provider, Owner, RepoName)
	)`

	// Log lines for cache + parallel surface.
	MsgBulkCacheHitFmt     = "make-all-*: owner repo list cache HIT (%d repos, age %s)\n"
	MsgBulkCacheMissFmt    = "make-all-*: owner repo list cache MISS — refreshing from %s\n"
	MsgBulkParallelFmt     = "make-all-*: applying with parallelism=%d\n"
	MsgBulkExceptLatest    = "make-all-*: --except-latest active — newest -vN per base group will be flipped to the OPPOSITE visibility\n"
	MsgBulkExceptInvertFmt = "  → except-latest: %s (highest -v%d) will be set to %s\n"
	MsgBulkInvertHeaderFmt = "\nmake-all-*: applying INVERTED visibility=%s to %d latest-version repo(s) on %s\n"

	// Fuzzy fallback surface.
	MsgBulkFuzzyAutoFixFmt = "make-all-*: auto-fix retry with pattern %q\n"
	MsgBulkFuzzyHintHeader = "make-all-*: did you mean one of these?\n"

	// make-last-* surface.
	MsgMakeLastResolvedFmt = "make-last-*: %s → %s (highest -v%d under base %q)\n"
	ErrMakeLastNoBaseFmt   = "Error: no repo matches base %q under %s (operation: make-last-resolve, reason: cache empty or no -vN siblings) — try `gitmap make-all-public %s %s*` to inspect candidates.\n"
	ErrMakeLastMissingArg  = "Error: usage: gitmap %s <owner-or-url> <base>\n"

	// Flags.
	FlagBulkParallel          = "--parallel"
	FlagBulkCacheTTL          = "--cache-ttl"
	FlagBulkExceptLatest      = "--except-latest"
	FlagBulkExceptLatestShort = "-XL"
)

// Spec 116 §preflight + §undo-redo — auth-status gate, drift guard,
// and reverse-loop messages shared by `make-all-*`, `vu`, and `vr`.
const (
	ErrVisAuthStatusFailedFmt = "Error: %s auth status failed: %v (operation: provider-auth-preflight, reason: %s) — run `%[1]s auth login` and retry.\n"
	MsgUndoDriftSkipFmt       = "DRIFT SKIP (current=%q expected=%q; re-run with --force to overwrite)\n"
	MsgUndoForceOverrideFmt   = "[--force] overriding drift guard for %s ... "
	ErrUndoAuditDBOpenFmt     = "%s: audit DB open failed: %v\n"
	MsgUndoReverseHeaderFmt   = "%s: reversing run #%d (%s/%s) — %d repo(s)\n"
	UndoPatternsRawFmt        = "%s:source-run=%d"
)
