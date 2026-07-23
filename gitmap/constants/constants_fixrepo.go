package constants

// Fix-repo exit codes. These match fix-repo.ps1 and fix-repo.sh 1:1
// so CI scripts that branch on the exit code keep working when they
// switch from invoking the script to invoking the binary.
const (
	FixRepoExitOk              = 0
	FixRepoExitNotARepo        = 2
	FixRepoExitNoRemote        = 3
	FixRepoExitNoVersionSuffix = 4
	FixRepoExitBadVersion      = 5
	FixRepoExitBadFlag         = 6
	FixRepoExitWriteFailed     = 7
	FixRepoExitBadConfig       = 8

	// FixRepoExitTestsFailed is the exit code emitted by --strict when
	// `go test` on the touched packages reports failure. Distinct from
	// 7 (write-failed) so CI can branch on "the rewrite produced
	// semantically broken code" vs "the file system rejected the
	// write". Picked as 9 (next free code) so the existing 0/2..8
	// matrix stays stable for downstream parsers; the help block
	// documents the new value verbatim.
	FixRepoExitTestsFailed = 9
)

// Fix-repo flag names. Both GNU long-form (`--dry-run`) and the
// PowerShell single-dash forms (`-DryRun`) are accepted as aliases.
const (
	FixRepoFlagAll     = "all"
	FixRepoFlagDryRun  = "dry-run"
	FixRepoFlagVerbose = "verbose"
	FixRepoFlagConfig  = "config"

	// FixRepoFlagStrict gates the post-rewrite `go test` step that
	// catches semantic desyncs the byte-level rewriter cannot see
	// (e.g. a v9→v10 bump that desynced a hard-coded sibling literal
	// from the same file's `-vN` token — exactly the failure mode
	// closed by v4.12.0). Off by default so `gitmap fix-repo` stays a
	// pure rewriter for non-Go repos and for users without a Go
	// toolchain on PATH; opt-in via --strict / -Strict for Go repos
	// in CI.
	FixRepoFlagStrict     = "strict"
	FixRepoModeFlag2      = "-2"
	FixRepoModeFlag3      = "-3"
	FixRepoModeFlag5      = "-5"
	FixRepoConfigFileName = "fix-repo.config.json"
)

// FixRepoFlagRestrict (v5.39.0+) narrows the rewrite scope. The only
// currently-defined value is `no-version` (alias `nv`), which
// suppresses the v1→v2 bare-base sweep so ONLY `{base}-vN` tokens are
// rewritten. Long form `--restrict`, single-dash `-restrict`, and
// short alias `-r` are all accepted.
const (
	FixRepoFlagRestrict           = "restrict"
	FixRepoFlagRestrictShort      = "r"
	FixRepoRestrictNoVersion      = "no-version"
	FixRepoRestrictNoVersionShort = "nv"
)

// Fix-repo defaults.
const (
	FixRepoDefaultSpan    = 2
	FixRepoMaxFileBytes   = int64(5 * 1024 * 1024)
	FixRepoBinarySniffMax = 8192
)

// Fix-repo user-facing messages and error formats. All literals
// printed by the command live here so the no-magic-strings rule is
// honored and script-vs-binary parity stays explicit.
const (
	FixRepoMsgHeaderFmt   = "fix-repo  base=%s  current=v%d  mode=%s\n"
	FixRepoMsgTargetsFmt  = "targets:  %s\n"
	FixRepoMsgIdentityFmt = "host:     %s  owner=%s\n"
	FixRepoMsgScannedFmt  = "scanned: %d files\n"
	FixRepoMsgChangedFmt  = "changed: %d files (%d replacements)\n"
	FixRepoMsgModeFmt     = "mode:    %s\n"
	FixRepoMsgModified    = "modified: %s (%d replacements)\n"
	// Dry-run preview output (v5.40.0+). One line per file that would
	// be modified, followed by a compact per-rule breakdown so users
	// see exactly which `{base}-vN` targets would fire. The breakdown
	// renders via FixRepoMsgDryRunHit / FixRepoMsgDryRunHitBare.
	FixRepoMsgDryRunPreview = "  [dry-run] would rewrite %s (%d replacements): %s\n"
	FixRepoMsgDryRunHit     = "v%d×%d"
	FixRepoMsgDryRunHitBare = "bare×%d"
	FixRepoMsgNothing       = "fix-repo: nothing to replace\n"
	FixRepoTargetsNone      = "(none)"
	FixRepoModeWrite        = "write"
	FixRepoModeDryRun       = "dry-run"
	FixRepoErrNotARepo      = "fix-repo: ERROR not a git repository (E_NOT_A_REPO)\n"
	FixRepoErrNoRemote      = "fix-repo: ERROR no remote URL found (E_NO_REMOTE)\n"
	FixRepoErrParseURLFmt   = "fix-repo: ERROR cannot parse remote URL %q (E_NO_REMOTE)\n"
	FixRepoErrNoVerSuffFmt  = "fix-repo: ERROR no -vN suffix on repo name %q (E_NO_VERSION_SUFFIX)\n"
	FixRepoErrBadVersion    = "fix-repo: ERROR version <= 0 (E_BAD_VERSION)\n"
	FixRepoErrBadFlagFmt    = "fix-repo: ERROR %s (E_BAD_FLAG)\n"
	FixRepoErrBadConfigFmt  = "fix-repo: ERROR %s (E_BAD_CONFIG)\n"
	FixRepoErrWriteFmt      = "fix-repo: ERROR write failed for %s: %v\n"
	FixRepoMsgGofmtFmt      = "gofmt:   %d .go file(s) reformatted\n"
	FixRepoMsgGofmtBatchFmt = "gofmt:   %d .go file(s) reformatted across %d batch(es)\n"
	FixRepoMsgGofmtSkip     = "gofmt:   skipped (dry-run)\n"
	FixRepoMsgGofmtNoneFmt  = "gofmt:   no .go files modified\n"
	FixRepoErrGofmtFmt      = "fix-repo: ERROR gofmt failed: %v\n%s"
	FixRepoErrGofmtMissing  = "fix-repo: WARN  gofmt not found on PATH; skipping post-rewrite formatting\n"

	// Dry-run batch preview (v6.80.1+). Emitted when --dry-run is
	// passed and there are .go files that would be formatted; lets
	// users on Windows see per-batch cmd-line sizes before shipping
	// the real rewrite.
	FixRepoMsgGofmtDryFmt      = "gofmt (dry-run): would run %d batch(es) across %d file(s), budget=%d\n"
	FixRepoMsgGofmtDryBatchFmt = "  batch %d/%d: files=%d, cmdLen=%d bytes (%d%% of budget)%s\n"
	FixRepoMsgGofmtDryNearTag  = " " + ColorYellow + "NEAR-LIMIT" + ColorReset
	FixRepoMsgGofmtDryOverTag  = " " + ColorRed + "OVER-LIMIT" + ColorReset

	// Verbose progress output (v6.80.1+). Emitted when --verbose is
	// passed. Header is printed once before the loop; per-batch
	// start/done lines fire around each exec.
	FixRepoMsgGofmtVerbHeaderFmt     = "gofmt:   %d batch(es), %d file(s), budget=%d\n"
	FixRepoMsgGofmtVerbBatchStartFmt = "  [%d/%d] formatting %d files (cmdLen=%d)...\n"
	FixRepoMsgGofmtVerbBatchDoneFmt  = "    done in %s; ETA ~%s\n"

	// FixRepoGofmtNearLimitPct is the percent-of-budget threshold at
	// which the dry-run preview tags a batch as NEAR-LIMIT. Batches at
	// or above 100% of budget are tagged OVER-LIMIT (should be rare;
	// chunker only permits over-limit for a single pathological path).
	FixRepoGofmtNearLimitPct = 90

	// FixRepoGofmtMinCmdLen is the floor accepted by
	// --gofmt-max-cmd-len. Values below this cannot fit even one
	// typical repo-relative path plus argv overhead.
	FixRepoGofmtMinCmdLen = 512

	// FixRepoGofmtMaxCmdLen bounds the combined length of file-path
	// arguments passed to a single `gofmt -w` invocation. Windows'
	// CreateProcess caps the command line at 32,767 characters; 30,000
	// leaves headroom for the executable path and the `-w` flag. On a
	// large repo with hundreds of touched .go files under a long
	// absolute path (e.g. D:\wp-work\...\gitmap\...), a single-batch
	// exec overflowed this cap and surfaced as
	// "The filename or extension is too long."
	// See mem://issues/2026-05-01-fixrepo-no-gofmt.md.
	FixRepoGofmtMaxCmdLen = 30000

	// Strict-mode (post-rewrite `go test`) message family. Phrased to
	// mirror the gofmt block so a user reading the trailing summary
	// can tell at a glance which post-step ran, what it did, and
	// (on failure) which exit code maps to which root cause.
	FixRepoMsgStrictSkipDryRun = "strict:  skipped (dry-run)\n"
	FixRepoMsgStrictNoGoFiles  = "strict:  no .go files modified; skipping go test\n"
	FixRepoMsgStrictNoPackages = "strict:  no Go packages derived from modified files; skipping go test\n"
	FixRepoMsgStrictRunFmt     = "strict:  running go test on %d package(s): %s\n"
	FixRepoMsgStrictPassFmt    = "strict:  go test passed (%d package(s))\n"
	FixRepoErrStrictMissing    = "fix-repo: WARN  go not found on PATH; --strict skipped\n"
	FixRepoErrStrictFailFmt    = "fix-repo: ERROR strict mode: go test failed (E_TESTS_FAILED): %v\n%s"
)
