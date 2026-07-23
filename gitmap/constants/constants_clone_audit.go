package constants

// Constants for the read-only `gitmap clone --audit` mode. The audit
// command never executes git or writes to disk; every printable line
// flows through one of the format strings below so users can grep and
// downstream tools can parse stable text.
//
// See: spec/01-app/05-cloner.md, gitmap/cloner/audit.go.

// CloneFlagAudit is the long-form flag name (`--audit`) that switches
// `gitmap clone` into the read-only planner.
const CloneFlagAudit = "audit"

// FlagDescCloneAudit is the help text shown by `gitmap help clone`.
const FlagDescCloneAudit = "Validate planned git clone commands and print a diff-style summary; never executes."

// Audit output format strings. Layout decisions:
//   - Header carries source path, target dir, and total record count.
//   - Each row begins with a single-character diff marker so the output
//     is visually scannable and easy to grep ("^+ " for new clones,
//     "^! " for invalid records, etc.).
//   - The optional command line is indented two spaces to differentiate
//     it from the row header.
//   - Summary always prints, even when every count is zero.
const (
	MsgCloneAuditHeader  = "clone audit: source=%s target=%s records=%d\n"
	MsgCloneAuditRow     = "  %s %s %s (%s) — %s\n"
	MsgCloneAuditCmd     = "      $ %s\n"
	MsgCloneAuditSummary = "audit summary: +clone=%d ~pull=%d =cached=%d ?conflict=%d !invalid=%d\n"
)

// Error format used when the audit cannot load the source manifest. The
// runtime falls back to os.Exit(1) after printing this — partial output
// would be misleading in audit mode.
const ErrCloneAuditLoad = "clone --audit: could not load %q: %v\n"

// ErrCloneAuditDirectURL is printed when the user combines `--audit` with
// a direct git URL. Audit only makes sense against a manifest file
// because it needs a list of records to validate.
const ErrCloneAuditDirectURL = "clone --audit: requires a manifest file (json|csv|text|path), not a direct git URL\n"
