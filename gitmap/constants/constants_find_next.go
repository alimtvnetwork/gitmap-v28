package constants

// find-next: list every repo whose latest probe row reports an available
// upgrade (Phase 2.4, v3.9.0).
//
// Joins Repo against the newest VersionProbe per repo (via correlated
// subquery on MAX(ProbedAt)) and filters to IsAvailable=1. Optionally
// scoped to a single ScanFolderId so callers can query "what's new in
// E:\src" without seeing unrelated repos.

// SQL: every repo whose latest VersionProbe row has IsAvailable=1.
// Sort by NextVersionNum DESC so the freshest tags float to the top.
const SQLSelectFindNext = `
SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
       r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes,
       p.NextVersionTag, p.NextVersionNum, p.Method, p.ProbedAt
FROM Repo r
JOIN VersionProbe p ON p.RepoId = r.RepoId
WHERE p.IsAvailable = 1
  AND p.ProbedAt = (
    SELECT MAX(ProbedAt) FROM VersionProbe WHERE RepoId = r.RepoId
  )
ORDER BY p.NextVersionNum DESC, r.Slug ASC`

// SQL: same as above, scoped to a specific ScanFolderId.
const SQLSelectFindNextByScanFolder = `
SELECT r.RepoId, r.Slug, r.RepoName, r.HttpsUrl, r.SshUrl, r.Branch,
       r.RelativePath, r.AbsolutePath, r.CloneInstruction, r.Notes,
       p.NextVersionTag, p.NextVersionNum, p.Method, p.ProbedAt
FROM Repo r
JOIN VersionProbe p ON p.RepoId = r.RepoId
WHERE p.IsAvailable = 1
  AND r.ScanFolderId = ?
  AND p.ProbedAt = (
    SELECT MAX(ProbedAt) FROM VersionProbe WHERE RepoId = r.RepoId
  )
ORDER BY p.NextVersionNum DESC, r.Slug ASC`

// find-next user-facing strings.
const (
	MsgFindNextEmpty     = "No repos with available updates. Run `gitmap probe --all` first.\n"
	MsgFindNextHeaderFmt = "Available updates (%d):\n"
	MsgFindNextRowFmt    = "  %s → %s [method=%s, probed=%s]\n      %s\n"
	MsgFindNextDoneFmt   = "Hint: run `gitmap pull` or `gitmap cn next all` to apply.\n"
	// ErrFindNextQuery is the bare wrap-string used by store-side
	// fmt.Errorf calls. No trailing \n because errors are returned,
	// not printed. The cmd-layer counterpart (ErrFindNextQueryFmt)
	// adds the project's standard "Error: ... (operation: ...,
	// reason: ...)\n" framing for stderr.
	ErrFindNextQuery = "find-next: failed to query: %w"
	// ErrFindNextScanRow — same convention, used by store/find_next.go
	// when a per-row Scan fails.
	ErrFindNextScanRow = "find-next: failed to scan row: %w"
	// ErrFindNextQueryFmt wraps any DB-side failure surfaced by
	// store.DB.FindNext. Trailing \n so callers can Fprintf directly
	// to os.Stderr — matches the ErrScan* family in
	// constants_messages.go.
	ErrFindNextQueryFmt = "Error: find-next failed to query database: %v (operation: select, reason: db error)\n"
	// ErrFindNextScanRowFmt is the stderr-formatted counterpart of
	// ErrFindNextScanRow, kept for symmetry should a future caller
	// need to print a row-scan failure directly.
	ErrFindNextScanRowFmt = "Error: find-next failed to scan row: %v (operation: row-scan, reason: db error)\n"
	// ErrFindNextJSONEncodeFmt fires when json.Encoder.Encode fails
	// while writing the result array to stdout. Vanishingly rare
	// (stdout broken pipe), but still routed through stderr with
	// the standard format so scripts can detect it.
	ErrFindNextJSONEncodeFmt = "Error: find-next failed to encode JSON output: %v (operation: encode, reason: io error)\n"
	MsgFindNextUsageHeader   = "Usage: gitmap find-next [--scan-folder <id>] [--json]"
)

// find-next flag-validation errors. Each one is printed to stderr with
// the usage header before the process exits 2 (the conventional exit
// code for CLI usage errors — distinct from the exit-1 used for I/O
// or DB failures so scripts can branch on the cause).
const (
	// ErrFindNextUnknownFlagFmt fires on tokens like `--jsno` or
	// `--scanfolder` that don't match any known flag. The %q wraps
	// the offending token so embedded whitespace stays visible.
	ErrFindNextUnknownFlagFmt = "find-next: unknown flag %q\n"
	// ErrFindNextUnknownFlagSuggestFmt is the same as above but
	// includes a "did you mean" hint when the unknown token is one
	// edit away from a known flag.
	ErrFindNextUnknownFlagSuggestFmt = "find-next: unknown flag %q (did you mean %q?)\n"
	// ErrFindNextBoolTakesNoValueFmt fires on `--json=true` and
	// friends. --json is a pure boolean flag; accepting `=true`
	// would imply `=false` is also valid, and silently ignoring
	// the value would let `--json=fasle` do the wrong thing.
	ErrFindNextBoolTakesNoValueFmt = "find-next: %s does not take a value (got %q)\n"
	// ErrFindNextMissingValueFmt fires when `--scan-folder` is the
	// last token or is followed by another flag.
	ErrFindNextMissingValueFmt = "find-next: %s requires an integer scan-folder ID\n"
	// ErrFindNextBadIntFmt replaces the previous silent ignore for
	// non-integer scan-folder IDs (e.g. `--scan-folder abc`).
	ErrFindNextBadIntFmt = "find-next: %s expects an integer, got %q\n"
	// ErrFindNextUnexpectedArgFmt fires on bare positional tokens.
	// find-next takes flags only; a stray positional is almost
	// always a quoting bug worth surfacing.
	ErrFindNextUnexpectedArgFmt = "find-next: unexpected positional argument %q\n"
)

// find-next CLI flag tokens.
const (
	FindNextFlagScanFolder = "--scan-folder"
	FindNextFlagJSON       = "--json"
)

// FindNextKnownFlags lists every flag find-next accepts. Used by the
// suggestion engine to compute "did you mean?" hints for typos and by
// the validator to detect unknown tokens. Order is irrelevant.
var FindNextKnownFlags = []string{
	FindNextFlagScanFolder,
	FindNextFlagJSON,
}

// gitmap:cmd top-level
// find-next CLI commands.
const (
	CmdFindNext      = "find-next"
	CmdFindNextAlias = "fn"
)
