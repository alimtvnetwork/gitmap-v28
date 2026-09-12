package constants

// VersionProbe operations (v3.8.0+, Phase 2.3).
//
// The probe inspects a repo's remote to discover the next available
// version tag. Hybrid strategy: try `git ls-remote` against the HEAD
// first (cheap, network-only), and only fall back to a `--depth 1
// --filter=blob:none` clone when ls-remote returns nothing usable.
//
// Results land in the VersionProbe table. The "Method" column records
// which strategy succeeded ("ls-remote" or "shallow-clone"); "Error"
// captures the failure reason when IsAvailable = 0 so operators can
// debug a probe without re-running it.

// Probe method tokens (stored in VersionProbe.Method).
const (
	ProbeMethodLsRemote     = "ls-remote"
	ProbeMethodShallowClone = "shallow-clone"
	ProbeMethodNone         = "none"
)

// SQL: insert a new probe row.
const SQLInsertVersionProbe = `INSERT INTO VersionProbe
	(RepoId, NextVersionTag, NextVersionNum, Method, IsAvailable, Error)
	VALUES (?, ?, ?, ?, ?, ?)`

// SQL: latest probe per repo.
const SQLSelectLatestVersionProbe = `SELECT VersionProbeId, RepoId, ProbedAt,
		NextVersionTag, NextVersionNum, Method, IsAvailable, Error
	FROM VersionProbe WHERE RepoId = ?
	ORDER BY ProbedAt DESC, VersionProbeId DESC LIMIT 1`

// SQL: bulk-tag every repo whose AbsolutePath was just scanned with the
// active ScanFolderId. Path list is interpolated as `?,?,?,...` because
// SQLite has no array binding.
const SQLTagReposByScanFolderTpl = `UPDATE Repo SET ScanFolderId = ? WHERE AbsolutePath IN (%s)`

// VersionProbe error/message strings.
const (
	ErrProbeOpenDB       = "version probe: failed to open database: %v"
	ErrProbeMissingURL   = "version probe: repo %q has no clone URL"
	ErrProbeLsRemoteFail = "ls-remote failed: %v"
	ErrProbeCloneFail    = "shallow clone failed: %v"
	ErrProbeRecord       = "version probe: failed to record result for repo %d: %v"
	ErrProbeNoRepo       = "version probe: no repo found at %q"
	ErrProbeTagFail      = "scan: failed to tag repos with scan folder %d: %v"
)

// VersionProbe user-facing CLI strings.
const (
	MsgProbeStartFmt    = "→ Probing %d repo(s)...\n"
	MsgProbeOkFmt       = "  ✓ %s → %s (method=%s)\n"
	MsgProbeNoneFmt     = "  · %s → no new version (method=%s)\n"
	MsgProbeFailFmt     = "  ✗ %s → %s\n"
	MsgProbeDoneFmt     = "✓ Probe complete: %d available, %d unchanged, %d failed.\n"
	MsgProbeUsageHeader = "Usage: gitmap probe [<repo-path>|--all] [--json] [--probe-workers N] [--probe-depth N]"
	MsgProbeNoTargets   = "No repos to probe. Pass a path or --all.\n"
)

// VersionProbe CLI tokens.
//
// `--workers` is the legacy spelling kept as a deprecated alias for
// `--probe-workers`; the foreground probe accepts both, but the new
// name is preferred (and is the same string scan accepts, which keeps
// muscle memory consistent across the two commands).
const (
	ProbeFlagAll          = "--all"
	ProbeFlagJSON         = "--json"
	ProbeFlagWorkers      = "--workers"
	ProbeFlagProbeWorkers = "--probe-workers"
	ProbeFlagDepth        = "--probe-depth"
	// ProbeFlagOutput selects an alternative human-output mode.
	// Currently only "terminal" is recognized — it emits the
	// standardized RepoTermBlock per repo (branch, from, to,
	// command) so the per-probe summary matches the shape used
	// by scan, clone-from, and clone-next. Mutually exclusive
	// with --json (terminal beats JSON when both are passed
	// because explicit terminal opts in to the human format).
	ProbeFlagOutput = "--output"
)

// Foreground probe pool sizing (v3.134.0+) and clone-depth (v3.135.0+).
//
// `gitmap probe` runs a small capped worker pool so a probe of N repos
// completes in ~N/2 round-trips instead of N. The cap is intentionally
// tight — git hosting providers (GitHub in particular) throttle bursts
// of unauthenticated ls-remote calls, and going above 3 workers
// produces more 429s than throughput gains. The default of 2 is the
// sweet spot for laptops on residential bandwidth.
//
// ProbeDefaultDepth is the `--depth N` value passed to the shallow-clone
// fallback. 1 is enough for the vast majority of repos because the
// fallback only runs when ls-remote returned nothing — at which point
// even one commit of history is sufficient to surface the tag list.
// Bump it (e.g. via `--probe-depth 25`) for repos whose latest tag
// lives further back than the first refs page exposes.
const (
	ProbeDefaultWorkers = 2
	ProbeMaxWorkers     = 3
	ProbeDefaultDepth   = 1
)

// Probe worker- and depth-flag messages.
const (
	ErrProbeWorkersValue   = "version probe: --probe-workers requires a positive integer, got %q"
	ErrProbeWorkersMissing = "version probe: --probe-workers requires a value"
	ErrProbeDepthValue     = "version probe: --probe-depth requires a positive integer, got %q"
	ErrProbeDepthMissing   = "version probe: --probe-depth requires a value"
	MsgProbeWorkersClamped = "  · --probe-workers %d exceeds cap, clamping to %d\n"
	MsgProbeWorkersAlias   = "  · --workers is deprecated; use --probe-workers instead\n"
)

// Background probe tuning for `gitmap scan` (v3.123.0+).
//
// When scan finds a small repo set we eagerly kick off a probe pass in
// the background so the next `gitmap find-next` call already has fresh
// data. The defaults are intentionally gentle: 3 workers max so we do
// not hammer GitHub's rate limit, and the auto-trigger only fires for
// scans of <50 repos so a directory full of vendored sources doesn't
// suddenly fan out 500 ls-remote calls. Power users can override or
// disable any of this with flags below.
const (
	// ScanProbeFlagDisable disables the background probe entirely.
	ScanProbeFlagDisable = "no-probe"
	// ScanProbeFlagNoWait makes scan return immediately after kicking
	// off probes (they keep running in background until completion or
	// process exit).
	ScanProbeFlagNoWait = "no-probe-wait"
	// ScanProbeFlagConcurrency sets the worker-pool size for the
	// background probe runner. Deprecated alias for ScanProbeFlagProbeWorkers.
	ScanProbeFlagConcurrency = "probe-concurrency"
	// ScanProbeFlagProbeWorkers is the unified spelling shared with
	// `gitmap probe --probe-workers`. Same semantics as
	// ScanProbeFlagConcurrency; either is accepted.
	ScanProbeFlagProbeWorkers = "probe-workers"
	// ScanProbeFlagProbeDepth controls the `git clone --depth N`
	// passed to the shallow-clone fallback (v3.135.0+).
	ScanProbeFlagProbeDepth = "probe-depth"

	// ScanProbeDefaultConcurrency caps the background pool. Three
	// workers is the documented sweet spot: parallel enough to
	// finish 50 probes in ~ten seconds, low enough that GitHub's
	// abuse detection doesn't kick in.
	ScanProbeDefaultConcurrency = 3
	// ScanProbeAutoTriggerCeiling is the repo-count threshold under
	// which background probing fires automatically. Above it the
	// user must opt in by passing --probe-workers explicitly.
	ScanProbeAutoTriggerCeiling = 50

	// FlagDescScanProbeDisable, FlagDescScanProbeNoWait,
	// FlagDescScanProbeConcurrency, FlagDescScanProbeProbeWorkers, and
	// FlagDescScanProbeProbeDepth are the help strings shown by
	// `gitmap help scan`.
	FlagDescScanProbeDisable      = "Skip the background version probe entirely (offline / air-gapped runs)"
	FlagDescScanProbeNoWait       = "Return as soon as scan finishes; let probes keep running in the background"
	FlagDescScanProbeConcurrency  = "(deprecated, use --probe-workers) Worker count for the background probe pool"
	FlagDescScanProbeProbeWorkers = "Worker count for the background probe pool (default 3, 0 = disable)"
	FlagDescScanProbeProbeDepth   = "Shallow-clone depth used by the probe fallback (default 1)"

	// MsgScanProbeConcurrencyAlias warns when the deprecated
	// --probe-concurrency spelling is used at the scan command line.
	MsgScanProbeConcurrencyAlias = "  · --probe-concurrency is deprecated; use --probe-workers instead\n"
)

// Background probe runtime messages.
const (
	// MsgScanProbeStartFmt fires once when scan kicks off the
	// background runner. The %d is the number of repos queued.
	MsgScanProbeStartFmt = "  ↪ background probe queued for %d repo(s) (workers=%d)\n"
	// MsgScanProbeWaitingFmt is printed at scan-end while we block
	// on the runner's Wait. Includes the count remaining so users
	// can gauge how long they'll be waiting.
	MsgScanProbeWaitingFmt = "  ⏳ Waiting for background probes to finish (%d remaining)...\n"
	// MsgScanProbeDoneFmt prints the per-bucket tally once Wait
	// returns. Mirrors the foreground probe summary line.
	MsgScanProbeDoneFmt = "  ✓ Background probe done: %d available, %d unchanged, %d failed.\n"
	// MsgScanProbeDetached prints when --no-probe-wait was passed
	// and we're returning before the pool drains.
	MsgScanProbeDetached = "  ↪ Background probe detached; results will land in the DB asynchronously.\n"
	// MsgScanProbeSkippedAutoFmt prints when the auto-trigger
	// declined to start (repo count above the ceiling) so users
	// understand why no probe ran.
	MsgScanProbeSkippedAutoFmt = "  · Background probe skipped: %d repos exceeds auto-trigger ceiling (%d). Pass --probe-workers to force.\n"
)
