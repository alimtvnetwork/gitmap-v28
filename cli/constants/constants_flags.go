package constants

const (
	FlagCloneYes      = "yes"
	FlagCloneYesShort = "y"
	FlagDescCloneYes  = "Auto-accept new SSH host keys for SSH clone URLs " +
		"(sets StrictHostKeyChecking=accept-new; changed keys still fail)"

	FlagVSCodePMSyncDryRun     = "dry-run"
	FlagDescVSCodePMSyncDryRun = "Preview projects.json changes without writing"
	// FlagVSCodePMSyncProjectsJSON overrides the resolved
	// projects.json location for a single run. Useful when VS Code
	// is installed in a non-standard location (portable build,
	// remote-dev container, second user profile) or when scripting
	// against a backup copy. Path must be absolute. When omitted the
	// command falls back to the default discovery chain.
	FlagVSCodePMSyncProjectsJSON     = "projects-json"
	FlagDescVSCodePMSyncProjectsJSON = "Absolute path to projects.json (overrides VS Code user-data discovery)"
	// FlagVSCodePMSyncTag REPLACES the auto-detected tag set for
	// every entry with the user-supplied list. Repeatable AND
	// accepts a comma-separated value: `--tag a --tag b,c` produces
	// {a,b,c}. Differs from the global `--vscode-tag` (which only
	// ADDS to detection); use this when you want exact control over
	// what lands in projects.json. The "gitmap" brand tag is NOT
	// auto-prepended in this mode — pass it explicitly if you want
	// to keep it.
	FlagVSCodePMSyncTag     = "tag"
	FlagDescVSCodePMSyncTag = "Replace detected tags with this name (repeatable; accepts comma-list). Use instead of --vscode-tag for exact control"
	// FlagVSCodePMSyncMode picks the tag-merge strategy used when the
	// detected set is reconciled with the existing on-disk tags.
	// Values: union (default, additive), replace (detector wins),
	// intersection (set-AND + brand pin — only tags present in BOTH
	// sources survive, but the "gitmap" brand tag is always preserved
	// so re-syncs never silently strip our own marker).
	FlagVSCodePMSyncMode     = "mode"
	FlagDescVSCodePMSyncMode = "Tag merge strategy: union (default) | replace | intersection"
	// VSCodePMSyncMode* are the canonical CLI string values for the
	// --mode flag. Centralized so the parser, the helptext, and the
	// MergeMode.String() implementation can never drift.
	VSCodePMSyncModeUnion        = "union"
	VSCodePMSyncModeReplace      = "replace"
	VSCodePMSyncModeIntersection = "intersection"
)

// Audit-legacy defaults + flag names. Patterns are comma-separated regexes.
// Defaults catch the gitmap-v28/v6/v7 → v8 rename regression class. // gitmap-legacy-ref-allow
const (
	FlagAuditLegacyPatterns     = "patterns"
	FlagDescAuditLegacyPatterns = "Comma-separated regex patterns to flag (default: gitmap-v[567])"
	FlagAuditLegacyPath         = "path"
	FlagDescAuditLegacyPath     = "Root directory to scan (default: current working directory)"
	FlagAuditLegacyJSON         = "json"
	FlagDescAuditLegacyJSON     = "Emit machine-readable JSON instead of human text"
	FlagAuditLegacyReport       = "report"
	FlagDescAuditLegacyReport   = "Write a Markdown audit report to this path (use '' for default .gitmap/audit-legacy-report.md)"
	FlagAuditLegacyDiffs        = "diffs"
	FlagDescAuditLegacyDiffs    = "Also write per-file unified diffs (legacy → v8) into the report's diffs/ folder and link them"
	DefaultAuditLegacyPatterns  = "gitmap-v[567]\\b"
	DefaultAuditLegacyReport    = ".gitmap/audit-legacy-report.md"
	DefaultAuditLegacyDiffsDir  = "diffs"
	DefaultAuditLegacyReplace   = "gitmap-v28"
	MsgAuditLegacyClean         = "  ✓ audit-legacy: no forbidden legacy strings found in %s\n"
	MsgAuditLegacyHeader        = "  ✗ audit-legacy: %d match(es) across %d file(s) for pattern(s) %q\n"
	MsgAuditLegacyHit           = "    %s:%d: %s\n"
	MsgAuditLegacyReportWrote   = "  📄 audit-legacy: report written to %s\n"
	MsgAuditLegacyDiffsWrote    = "  📎 audit-legacy: %d per-file diff(s) written to %s\n"
	ErrAuditLegacyWalk          = "audit-legacy: walk failed for %s: %v\n"
	ErrAuditLegacyRegex         = "audit-legacy: invalid regex %q: %v\n"
	ErrAuditLegacyReportWrite   = "audit-legacy: failed to write report to %s: %v\n"
	ErrAuditLegacyDiffWrite     = "audit-legacy: failed to write diff for %s: %v\n"
)

// Regoldens flag names + descriptions. The two env vars driven by
// this command (GITMAP_UPDATE_GOLDEN, GITMAP_ALLOW_GOLDEN_UPDATE)
// are the canonical two-key safety gate from the goldenguard
// package; the values themselves live there to keep one source of
// truth (do NOT duplicate "1" here).
const (
	FlagRegoldensPattern         = "pattern"
	FlagDescRegoldensPattern     = "Test name pattern passed to `go test -run` (required)"
	FlagRegoldensPackage         = "package"
	FlagDescRegoldensPackage     = "Go package selector to test (default: ./...)"
	FlagRegoldensSkipVerify      = "skip-verify"
	FlagDescRegoldensSkipVerify  = "Skip the second verification pass that re-runs without the gate env vars"
	FlagRegoldensDryRun          = "dry-run"
	FlagDescRegoldensDryRun      = "Print the `go test` invocations that would run, then exit 0"
	FlagRegoldensDiff            = "diff"
	FlagDescRegoldensDiff        = "After pass 1, print a summary of touched testdata/ golden files. Modes: 'short' = one terse line per file (status + path + +/- counts); 'full' = adds rename details and aggregate totals. Usage: --diff=short or --diff=full"
	FlagRegoldensDeterminism     = "determinism"
	FlagDescRegoldensDeterminism = "Run a determinism pre-check pass BEFORE pass 1: invokes `go test` with the per-test trigger but WITHOUT the allow-update gate, so any writer using goldenguard.AllowUpdateAfterDeterminism is exercised but no fixtures are written. Fails fast (exit 1) on the first non-deterministic writer; otherwise pass 1 proceeds normally."
	RegoldensDiffModeShort       = "short"
	RegoldensDiffModeFull        = "full"
	RegoldensDefaultPackageGlob  = "./..."
	// RegoldensNonDetMarker is the substring goldenguard prints when
	// a writer fails AssertWriterDeterministic. The pre-check pass
	// scans `go test` output for this marker to distinguish a real
	// non-determinism failure from the EXPECTED gate-block failure.
	// Keep in sync with goldenguard/determinism.go nonDeterministicWriterMsgFmt.
	RegoldensNonDetMarker = "is non-deterministic"
)

// Regoldens user-facing messages.
const (
	MsgRegoldensPass1Header     = "▸ Pass 1: regenerating goldens (GITMAP_UPDATE_GOLDEN=1 GITMAP_ALLOW_GOLDEN_UPDATE=1)\n"
	MsgRegoldensPass2Header     = "▸ Pass 2: verifying determinism (gate env vars cleared)\n"
	MsgRegoldensSkipVerify      = "▸ Pass 2 skipped (--skip-verify)\n"
	MsgRegoldensSuccess         = "✓ Goldens regenerated and verified for pattern %q in package %q\n"
	MsgRegoldensSuccessNoVeri   = "✓ Goldens regenerated for pattern %q in package %q (verification skipped)\n"
	MsgRegoldensDryRun          = "▸ Dry run — would execute:\n  %s\n  %s\n"
	MsgRegoldensDiffHeader      = "▸ Golden diff summary (testdata/**/*golden* files touched by pass 1) [mode=%s]:\n"
	MsgRegoldensDiffNoChanges   = "  (no testdata/**/*golden* files changed)\n"
	MsgRegoldensDiffLineShort   = "  %s  %s\n"
	MsgRegoldensDiffLineFull    = "  %s  %s  (+%d / -%d)\n"
	MsgRegoldensDiffRenameFull  = "      ↳ renamed from %s\n"
	MsgRegoldensDiffTotals      = "  ─ %d file(s) changed: %d added, %d modified, %d renamed, %d deleted (+%d / -%d total)\n"
	MsgRegoldensDiffSkipped     = "▸ Diff summary skipped: not a git working tree (or `git` not on PATH)\n"
	MsgRegoldensPass2NotRun     = "▣ Pass 2: did NOT run — pass 1 failed (exit %d); fixtures left in partial state, determinism unverified\n"
	MsgRegoldensPass2NotRunSkip = "▣ Pass 2: did NOT run — --skip-verify was set; determinism unverified\n"
	MsgRegoldensPass2Ran        = "▣ Pass 2: ran and PASSED — writer is deterministic; goldens are safe to commit\n"
	MsgRegoldensPrecheckHeader  = "▸ Pre-check: determinism (trigger ON, allow-update gate OFF — no fixture writes)\n"
	MsgRegoldensPrecheckPass    = "✓ Pre-check passed — no non-deterministic writers detected; proceeding to pass 1\n"
	ErrRegoldensMissingPat      = "regoldens: --pattern is required (e.g. --pattern TestCloneFromReportJSON_Golden)"
	ErrRegoldensPass1Failed     = "regoldens: pass 1 (regenerate) failed with exit code %d — fixtures may be partially written; inspect git status before committing"
	ErrRegoldensPass2Failed     = "regoldens: pass 2 (verify) failed with exit code %d — the writer is non-deterministic; fix the writer, do not re-run with --skip-verify"
	ErrRegoldensDiffMode        = "regoldens: --diff=%q is not valid; use --diff=short or --diff=full"
	ErrRegoldensPrecheckFailed  = "regoldens: determinism pre-check FAILED — at least one golden writer produced different bytes across runs; fix the writer (likely culprits: map iteration, time.Now, randomness, locale-dependent formatting). Pass 1 was NOT run; no fixtures were touched."
)

// Usage header — colorful banner shown above every help screen.
// Width is intentionally unboxed because ANSI escapes + dynamic version
// width make fixed-width box alignment unreliable across terminals.

const (
	FlagDebugPaths     = "debug-paths"
	FlagDescDebugPaths = "trace VS Code Project Manager sync rootPath canonicalization (raw → cleaned → resolved) to stderr; useful for diagnosing Windows 8.3 short-name and symlink dedup"
)

// Release flag names — surfaced as Cmd*-style constants so callers
// (e.g. the right-click context menu in installctxentries.go and the
// release parser in release.go) reference the same single source of
// truth rather than literal "--bump" / "bump" strings.
const (
	FlagBump      = "bump"
	FlagBumpDash  = "--bump"
	FlagDryRunDsh = "--dry-run"

	// Cluster command flags
	ClusterFlagExcept         = "except"
	ClusterFlagIP             = "ip"
	ClusterFlagID             = "id"
	ClusterFlagYes            = "yes"
	ClusterFlagYesShort       = "Y"
	ClusterFlagForceLifecycle = "force-lifecycle"
	ClusterFlagNoPreflight    = "no-preflight"
	ClusterFlagVerbose        = "verbose"
	ClusterFlagDryRun         = "dry-run"
)
