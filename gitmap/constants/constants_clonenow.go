package constants

// Constants for `gitmap clone-now <file>` (v3.161.0+).
//
// `clone-now` re-runs `git clone` against scan output: it consumes the
// JSON / CSV / text artifacts produced by `gitmap scan` (the same
// files written under `.gitmap/output/`) and re-creates each repo at
// its recorded relative path using the user-selected URL mode.
//
// Why a separate command from clone-from?
//
//   - clone-from is plan-driven (user-authored row schema: url, dest,
//     branch, depth) — the file describes intent.
//   - clone-now is round-trip-driven: input is gitmap's own scan
//     output, and we honor the recorded RelativePath verbatim so the
//     destination tree is byte-identical to the original layout.
//
// All user-facing strings live here per the no-magic-strings rule.

// CLI surface — canonical name + backward-compat aliases.
//
// CmdCloneReclone ("reclone") is the CANONICAL verb as of v3.x.
// It was renamed from clone-now so the CLI vocabulary makes the
// split between the two clone families unambiguous:
//
//   - `gitmap clone <url>`   — fresh clone from a URL.
//   - `gitmap reclone <file>` — RE-clone from a scan artifact
//     (round-trip the recorded RelativePath layout).
//
// All four legacy spellings are kept as aliases so existing scripts
// keep working forever:
//   - clone-now / cnow  — original name (v3.161.0).
//   - relclone / rc     — earlier explicit "re-clone" alias.
//
// The dispatcher in rootcore.go binds every spelling to runCloneNow.
// The completion generator picks them up via the
// `// gitmap:cmd top-level` marker on this const block.
const (
	CmdCloneReclone      = "reclone"
	CmdCloneRecloneAlias = "rec"
	CmdCloneNow          = "clone-now"
	CmdCloneNowAlias     = "cnow"
	CmdCloneRel          = "relclone"
	CmdCloneRelAlias     = "rc"
)

// Flag names + descriptions. Long-form only; short flags are
// reserved for very-frequent operations per the project convention.
const (
	FlagCloneNowExecute     = "execute"
	FlagDescCloneNowExecute = "Actually run git clone (default: dry-run only)"
	FlagCloneNowQuiet       = "quiet"
	FlagDescCloneNowQuiet   = "Suppress per-row progress lines (summary still prints)"
	// FlagCloneNowMode picks which URL column to use when the input
	// supplies both. Values: "https" (default) | "ssh". When the
	// requested mode is missing on a given row we fall back to the
	// other one rather than skipping the row -- the user's intent is
	// "clone these repos now", not "clone only the ones that have
	// the preferred URL shape".
	FlagCloneNowMode     = "mode"
	FlagDescCloneNowMode = "URL mode to clone with: 'https' (default) or 'ssh'. " +
		"Falls back to the other mode if the preferred URL is missing on a row."
	// FlagCloneNowFormat lets the caller force the input format when
	// the file extension is missing or wrong (e.g., `repos.out`).
	// Values: "" (auto from extension) | "json" | "csv" | "text".
	FlagCloneNowFormat     = "format"
	FlagDescCloneNowFormat = "Force input format: '' (auto from extension), 'json', 'csv', or 'text'."
	// FlagCloneNowCwd sets the working directory clones run in.
	// Empty (default) = the current process cwd. Honored as-is so
	// that scripts can re-create a tree under a fresh root with
	// `gitmap clone-now scan.json --cwd ./mirror --execute`.
	FlagCloneNowCwd     = "cwd"
	FlagDescCloneNowCwd = "Working directory for git clone (default: current dir)."
	// FlagCloneNowOnExists controls re-clone behavior when the
	// destination already contains a git repo. Default "skip" keeps
	// the historical no-op behavior. "update" runs fetch + checkout
	// without destroying local work. "force" removes the directory
	// and re-clones from scratch -- destructive, opt-in only.
	FlagCloneNowOnExists     = "on-exists"
	FlagDescCloneNowOnExists = "Behavior when target already exists: " +
		"'skip' (default, no-op when repo+branch match), " +
		"'update' (fetch + checkout to align with the planned URL/branch), " +
		"'force' (remove target and re-clone -- destructive)."
	// FlagCloneNowManifest is the unified, explicit-path alias for the
	// positional <file> argument. Accepts a JSON or CSV produced by
	// `gitmap scan`; format is auto-detected from the extension (use
	// --format to override). Precedence:
	//   1. --manifest <path>     (explicit, highest priority)
	//   2. positional <file>     (legacy positional form)
	//   3. auto-pickup           (./.gitmap/output/gitmap.{json,csv})
	// Passing BOTH --manifest and a positional file is a usage error
	// (exit 2) so the run is unambiguous.
	FlagCloneNowManifest     = "manifest"
	FlagDescCloneNowManifest = "Path to a scan artifact (JSON or CSV) to consume. " +
		"Format auto-detected from the extension. Equivalent to the positional " +
		"<file> argument; when omitted, auto-pickup under ./.gitmap/output/ is used."
	// FlagCloneNowScanRoot redirects auto-pickup to probe
	// `<scan-root>/.gitmap/output/` instead of `./.gitmap/output/`.
	// Useful when running `reclone` from a different directory than
	// the one originally scanned (e.g. CI scripts, scheduled jobs).
	// Ignored when --manifest or a positional <file> is supplied
	// (those are explicit paths and don't need a root to resolve).
	FlagCloneNowScanRoot     = "scan-root"
	FlagDescCloneNowScanRoot = "Directory to auto-pickup the scan manifest from " +
		"(probes <scan-root>/.gitmap/output/gitmap.{json,csv}). " +
		"Defaults to the current directory. Ignored when --manifest or a positional <file> is given."
	// FlagCloneNowYes bypasses the pre-flight existing-destinations
	// confirmation prompt. The prompt only fires on --execute when
	// at least one row's RelativePath already exists on disk; the
	// per-row --on-exists policy still decides what actually happens
	// to those existing dirs (skip / update / force). --yes is the
	// scripting / CI escape hatch and is mandatory in non-TTY runs
	// because there's no stdin to read a confirmation from.
	FlagCloneNowYes     = "yes"
	FlagDescCloneNowYes = "Skip the pre-flight confirmation when destination folders already exist " +
		"(required for non-interactive / CI runs). The --on-exists policy still applies per row."
	// FlagCloneNowNoSummary suppresses the pre-execute summary
	// (totals + destination-folder tree) printed before the safety
	// prompt. Useful for terse CI logs where the dry-run preview
	// has already been printed in a previous step.
	FlagCloneNowNoSummary     = "no-summary"
	FlagDescCloneNowNoSummary = "Suppress the pre-execute summary " +
		"(row totals + destination folder tree) shown before the safety prompt."
)

// CloneNowSummaryTreeLimit caps how many destination paths are
// rendered in the folder-layout tree to keep terminal output
// scannable on large round-trips. The total row count is always
// shown in the header so the user knows the full impact.
const CloneNowSummaryTreeLimit = 40

// CloneNowConfirmYes is the only stdin response that proceeds with
// --execute when destinations already exist. Anything else (including
// the empty default) aborts with exit code 2. Stable so shell
// scripts piping `yes` keep working: `yes | gitmap reclone --execute`.
const CloneNowConfirmYes = "y"

// CloneNowExistingPreviewLimit caps how many existing destinations
// are listed in the prompt to keep terminal output scannable. The
// total count is always shown so the user knows the full impact.
const CloneNowExistingPreviewLimit = 10

// CloneNowExitConfirmAborted is the exit code used when the user
// declines the prompt OR when --execute is passed in a non-TTY
// context with existing destinations and no --yes. Distinct from
// the per-row failure exit (1) so wrappers can tell "user said no"
// apart from "git clone failed".
const CloneNowExitConfirmAborted = 2

// CloneNowExitManifestInvalid is the exit code used when the
// manifest parses cleanly but one or more rows fail semantic
// validation (missing repo name, unusable URL, absent / absolute /
// traversal RelativePath). Same numeric code as the bad-flag exit
// (2) because both represent "you fed me invalid input" rather
// than a runtime clone failure (1). Aliased to a named constant
// so the validator's intent is self-documenting at the call site.
const CloneNowExitManifestInvalid = 2

// Manifest validation messages. All printed to stderr by the
// pre-flight validator in reclone_validate.go. Phrases are stable
// so tests + shell scripts can grep them, and short so the row
// table stays scannable in an 80-column terminal.
const (
	// %s = manifest path, %d = bad-row count, %d = total rows.
	MsgRecloneValidateHeaderFmt = "reclone: manifest validation failed for %s\n" +
		"  %d issue(s) across %d row(s):\n"
	// %d = 1-based row index, %s = repo name (or "<unnamed>"),
	// %s = dest (or "<empty>"), %s = reason phrase.
	MsgRecloneValidateRowFmt = "  - row %d  repo=%s  dest=%s  -- %s\n"
	// Footer printed once after all per-row lines. Tells the user
	// the run was aborted before any side effects occurred.
	MsgRecloneValidateFooter = "reclone: aborted; fix the manifest and re-run " +
		"(no clones were attempted)\n"
	// Per-row reason phrases. Stable identifiers — do not rephrase
	// without bumping the test fixtures and CHANGELOG.
	MsgRecloneValidateMissingRepoName = "missing RepoName"
	MsgRecloneValidateNoURL           = "no HTTPSUrl or SSHUrl set"
	MsgRecloneValidateMalformedURL    = "URL is not a valid git URL " +
		"(expected scheme://host/path or user@host:path)"
	MsgRecloneValidateMissingDest   = "missing RelativePath"
	MsgRecloneValidateAbsoluteDest  = "RelativePath must be relative, not absolute"
	MsgRecloneValidateTraversalDest = "RelativePath escapes the working dir via '..'"
	// Placeholders for empty fields in the report table.
	MsgRecloneValidateUnnamedRepo = "<unnamed>"
	MsgRecloneValidateEmptyDest   = "<empty>"
)

// On-exists policy enum strings. Stable: surfaced in --on-exists,
// the dry-run header, and the per-row Result.Detail field. Renaming
// is a breaking change for shell scripts that grep these values.
const (
	CloneNowOnExistsSkip   = "skip"
	CloneNowOnExistsUpdate = "update"
	CloneNowOnExistsForce  = "force"
)

// Mode enum strings. Stable: surfaced in the dry-run header and the
// per-row progress lines.
const (
	CloneNowModeHTTPS = "https"
	CloneNowModeSSH   = "ssh"
)

// Format enum strings. The "auto" empty value is intentionally not
// exported because callers detect it via len(format)==0.
const (
	CloneNowFormatJSON = "json"
	CloneNowFormatCSV  = "csv"
	CloneNowFormatText = "text"
)

// Status enum strings. Mirrors clone-from for cross-tool grep-ability:
// downstream pipelines that already filter on "ok"/"skipped"/"failed"
// keep working without a status-name translation table.
const (
	CloneNowStatusOK      = "ok"
	CloneNowStatusSkipped = "skipped"
	CloneNowStatusFailed  = "failed"
)

// CloneNowErrTrimLimit caps the per-row stderr summary length so the
// summary table stays scannable in an 80-column terminal. Full stderr
// remains in the user's scrollback because we use CombinedOutput.
const CloneNowErrTrimLimit = 80

// User-facing messages. Trailing newlines are baked in so call sites
// don't need to remember them.
const (
	// %s = source path, %s = format, %s = mode, %d = row count.
	MsgCloneNowDryHeader = "gitmap clone-now: dry-run\n" +
		"source: %s (%s, mode=%s)\n" +
		"%d row(s) -- pass --execute to actually clone\n\n"
	// %d ok, %d skipped, %d failed, %d total.
	MsgCloneNowSummaryHeader = "\ngitmap clone-now: %d ok, %d skipped, %d failed (%d total)\n"
	MsgCloneNowDestExists    = "dest exists"
	MsgCloneNowMissingArg    = "reclone: <file> argument is required and " +
		"no scan artifact was found under ./.gitmap/output/ " +
		"(looked for gitmap.json then gitmap.csv). " +
		"Run `gitmap scan` first, or pass an explicit path " +
		"(e.g. reclone .gitmap/output/gitmap.json)."
	// %s = scan-root value. Variant of MsgCloneNowMissingArg used
	// when the user explicitly requested auto-pickup from a custom
	// root via --scan-root and that root yielded no artifact —
	// echoing the path back makes the typo / wrong-dir case obvious.
	MsgCloneNowMissingArgScanRoot = "reclone: no scan artifact was found under " +
		"%s/.gitmap/output/ (looked for gitmap.json then gitmap.csv). " +
		"Run `gitmap scan` against that root first, or pass --manifest / a positional <file>."
	// %s = auto-discovered manifest path. Printed to stderr when
	// reclone is invoked with no <file> arg and a scan artifact is
	// found in the conventional location. Lets users see exactly
	// which file fed the run instead of guessing.
	MsgCloneNowAutoPickup = "reclone: using scan artifact %s (auto-discovered; pass an explicit path to override)\n"
	// %s = positional file, %s = --manifest value. Printed when the
	// caller supplies BOTH forms — refusing is safer than silently
	// preferring one and having the run consume the wrong artifact.
	MsgCloneNowManifestConflict = "reclone: cannot combine positional <file> %q with --manifest %q; pass only one\n"
	// %d = total existing dirs, %s = on-exists policy. Header
	// printed before the bullet list of existing destinations.
	// %s = source path, %s = format, %s = mode, %s = on-exists,
	// %s = resolved cwd. Top banner of the pre-execute summary.
	MsgCloneNowSummaryHeaderFmt = "\nreclone: pre-execute summary\n" +
		"  source     : %s (%s)\n" +
		"  mode       : %s\n" +
		"  on-exists  : %s\n" +
		"  cwd        : %s\n"
	// %d = total rows, %d = new dirs, %d = existing dirs.
	MsgCloneNowSummaryCountsFmt = "  rows       : %d total (%d new, %d already exist)\n"
	// Section title for the destination folder tree.
	MsgCloneNowSummaryTreeTitle = "  destinations:\n"
	// %s = tree-formatted line (already includes the leading
	// indent + branch glyph). One line per visible destination.
	MsgCloneNowSummaryTreeLineFmt = "    %s\n"
	// %d = number of dirs not shown.
	MsgCloneNowSummaryTreeTruncFmt = "    ... and %d more\n"
	MsgCloneNowConfirmHeader       = "reclone: %d destination folder(s) already exist on disk " +
		"(--on-exists=%s will be applied to each):\n"
	// %s = relative path. One bullet per existing destination, up
	// to CloneNowExistingPreviewLimit rows.
	MsgCloneNowConfirmBullet = "  - %s\n"
	// %d = number of dirs not shown.
	MsgCloneNowConfirmTruncated = "  ... and %d more\n"
	// Final prompt line. The trailing space (no newline) keeps the
	// cursor on the prompt line for the user's response.
	MsgCloneNowConfirmPrompt = "Proceed with `git clone` against these destinations? [y/N]: "
	// Printed (stderr) when the user declines. Exit code follows.
	MsgCloneNowConfirmAborted = "reclone: aborted by user; no clones were performed\n"
	// Printed (stderr) when --execute lands in a non-TTY context
	// with existing destinations and no --yes. Tells the user
	// exactly which flag would unblock the run.
	MsgCloneNowConfirmNonTTY = "reclone: refusing to proceed -- destinations already exist " +
		"and stdin is not a TTY; pass --yes to confirm non-interactively\n"
	MsgCloneNowNoURL = "no url for selected mode"
	// Idempotency / re-clone messages. Each lands in Result.Detail
	// so the per-row summary tells the user exactly which branch
	// of the on-exists policy fired. The mismatch + fail messages
	// take printf args -- documented per-line.
	MsgCloneNowAlreadyMatches = "already matches (url + branch)"
	MsgCloneNowNotARepo       = "dest exists but is not a git repo"
	// %s = current remote, %s = expected remote.
	MsgCloneNowURLMismatch = "skipped: remote url differs (have=%s, want=%s)"
	// %s = current branch, %s = expected branch.
	MsgCloneNowBranchMismatch = "skipped: branch differs (have=%s, want=%s)"
	MsgCloneNowUpdated        = "updated (fetch + checkout)"
	MsgCloneNowForceRecloned  = "force-recloned (previous dir removed)"
	// %s = path, %v = err.
	MsgCloneNowForceRemoveFail = "force: remove %s: %v"
	// %s = trimmed git stderr.
	MsgCloneNowFetchFail = "update: git fetch failed: %s"
	// %s = branch, %s = trimmed git stderr.
	MsgCloneNowCheckoutFail = "update: git checkout %s failed: %s"
)

// Errors. All use printf-style verbs documented inline.
const (
	// %s = path, %v = err.
	ErrCloneNowAbsPath = "clone-now: resolve path %s: %v"
	ErrCloneNowOpen    = "clone-now: open %s: %v"
	// %v = err.
	ErrCloneNowJSONDecode = "clone-now: decode JSON: %v"
	ErrCloneNowCSVRead    = "clone-now: read CSV: %v"
	ErrCloneNowTextRead   = "clone-now: read text: %v"
	// %s = bad value.
	ErrCloneNowBadMode     = "clone-now: --mode must be 'https' or 'ssh', got %q"
	ErrCloneNowBadFormat   = "clone-now: --format must be 'json', 'csv', or 'text', got %q"
	ErrCloneNowBadOnExists = "clone-now: --on-exists must be 'skip', 'update', or 'force', got %q"
	// %s = file extension (with leading dot, or "" when missing),
	// %s = path. Emitted when auto-detect cannot map the file
	// extension to a supported format. Use --format to override.
	ErrCloneNowUnsupportedExt = "clone-now: unsupported file extension %q for %s; " +
		"supported extensions are .json, .csv, .txt (or pass --format json|csv|text)"
	// %s = path.
	ErrCloneNowEmpty = "clone-now: %s contains zero clonable rows"
	// Schema validation errors. Run BEFORE the tolerant
	// formatter.ParseJSON / formatter.ParseCSV so a typo'd field
	// name or a wrong-shape input fails loudly instead of producing
	// silently-empty Rows downstream. Phrases are stable so tests
	// and shell scripts can grep them.
	//
	// %v = decode error from json.Decoder.
	ErrCloneNowJSONShape = "clone-now: JSON input must be an array of objects: %v"
	// %d = 1-based row index, %s = actual JSON kind ("string",
	// "number", "array", "null", etc.). Emitted when the top-level
	// array is well-formed but one of its elements is not an object.
	ErrCloneNowJSONRowNotObject = "clone-now: JSON row %d must be an object, got %s"
	// %d = 1-based row index, %s = unknown field name, %s = sorted
	// list of all known field names.
	ErrCloneNowUnknownJSONField = "clone-now: JSON row %d has unknown field %q; " +
		"known fields: %s"
	// %d = 1-based row index. Emitted when a row has neither
	// httpsUrl nor sshUrl set to a non-empty string.
	ErrCloneNowMissingURL = "clone-now: JSON row %d has neither httpsUrl nor sshUrl set"
	// Emitted when the CSV file is completely empty (no header).
	ErrCloneNowEmptyCSV = "clone-now: CSV input is empty (no header row found)"
	// %s = unknown column name, %s = sorted list of known fields.
	ErrCloneNowUnknownCSVField = "clone-now: CSV header has unknown column %q; " +
		"known columns: %s"
	// Emitted when the CSV header has neither httpsUrl nor sshUrl.
	ErrCloneNowCSVMissingURLCol = "clone-now: CSV header must include 'httpsUrl' or 'sshUrl'"
	// %d = 1-based data row number (header is NOT counted; data row
	// 1 is the first row after the header — matches what a user
	// inspecting the file in a spreadsheet sees as "row 2"). %v =
	// underlying csv.Reader error.
	ErrCloneNowCSVRowRead = "clone-now: CSV data row %d failed to parse: %v"
	// %d = 1-based data row number. Emitted when a CSV data row has
	// neither httpsUrl nor sshUrl set to a non-empty value.
	ErrCloneNowCSVRowMissingURL = "clone-now: CSV data row %d has neither httpsUrl nor sshUrl set"
	// MsgCloneNowMkdirParentFailFmt is the per-row Detail set when
	// pre-creating the destination's parent directory fails. Mirrors
	// the clonefrom equivalent so summary tables read consistently.
	// %v = err.
	MsgCloneNowMkdirParentFailFmt = "mkdir parent: %v"
	// ErrCloneNowMkdirParent is the standardized stderr log emitted
	// alongside MsgCloneNowMkdirParentFailFmt (Code Red zero-swallow).
	// %s = parent path, %v = err.
	ErrCloneNowMkdirParent = "Error: clone-now: failed to create dest parent at %s: %v " +
		"(check permissions / disk space)\n"
)
