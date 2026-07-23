package constants

// Clone progress format strings.
const (
	ProgressBeginFmt   = "[%3d/%d]  Cloning %s ..."
	ProgressDoneFmt    = " done (%s)\n"
	ProgressSkipFmt    = " skipped (cached)\n"
	ProgressFailFmt    = " FAILED\n"
	ProgressSummaryFmt = "\nClone complete: %d/%d repos in %s\n"
	ProgressDetailFmt  = "  Cloned: %d | Pulled: %d | Skipped: %d | Failed: %d\n"
)

// Batch progress format strings (generic operations).
const (
	BatchProgressBeginFmt   = "[%3d/%d]  %s ..."
	BatchProgressDoneFmt    = " done (%s)\n"
	BatchProgressFailFmt    = " FAILED\n"
	BatchProgressSkipFmt    = " skipped\n"
	BatchProgressSummaryFmt = "\n%s complete: %d/%d in %s\n"
	BatchProgressDetailFmt  = "  Succeeded: %d | Failed: %d | Skipped: %d\n"
	BatchStoppedMsg         = "  ⚠ Halted early (--stop-on-fail)\n"
)

// Batch failure report format strings.
const (
	BatchFailureHeader    = "  ── Failed Items ──"
	BatchFailureEntryFmt  = "  %d. %s: %s\n"
	BatchFailureFooterFmt = "  ── %d failure(s) total ──\n"
	ExitPartialFailure    = 3
)

// Batch flag constants.
const (
	FlagStopOnFail     = "stop-on-fail"
	FlagDescStopOnFail = "Stop batch operation after first failure"
)

// Clone shorthands — short aliases for `gitmap clone <source>` that
// expand to the default scan output files (json/csv/text).
const (
	ShorthandJSON = "json"
	ShorthandCSV  = "csv"
	ShorthandText = "text"
)

// Clone SSH host-key setup diagnostics.
const ErrCloneSSHEnvFmt = "clone: failed to configure SSH host-key acceptance: %v\n"

// Multi-URL clone messages (spec/01-app/104-clone-multi.md).
const (
	MsgCloneInvalidURLFmt    = "  ⚠ Skipping invalid URL: %s\n"
	MsgCloneSummaryMultiFmt  = "\n  Multi-clone summary: %d succeeded, %d failed (of %d URLs)\n"
	MsgCloneRegisteredInline = "  ✓ Registered with GitHub Desktop: %s\n"
	MsgCloneMultiBegin       = "\n  Cloning %d repositories...\n"
	MsgCloneMultiItem        = "\n  [%d/%d] %s\n"
	ErrCloneAllInvalid       = "  ✗ All URLs were invalid — nothing to clone\n"
	ErrCloneMultiFailedFmt   = "  ✗ [%d/%d] %s failed: %v\n"
)

// Multi-clone exit codes.
const (
	ExitCloneMultiPartialFail = 1
	ExitCloneMultiAllInvalid  = 3
)

// Stale-binary detection — fired when executeDirectClone is called with a
// folder name that itself parses as a URL. That shape is impossible in
// current source (multi-URL routing in runClone catches it), so when it
// happens it almost always means the user is running a deployed binary
// that pre-dates v3.80.0's multi-URL fix. We refuse to build the broken
// `D:\...\https:\github.com\...` path and tell the user exactly why.
const ErrCloneStaleBinaryFolderURL = "" +
	"  ✗ Refusing to clone: the folder name resolved to a URL (%q).\n" +
	"    This means your installed gitmap binary is older than v3.80.0\n" +
	"    (current source: v%s). The multi-URL clone fix is not present\n" +
	"    in the binary on your PATH.\n\n" +
	"    To fix:\n" +
	"      1. gitmap doctor                # confirm the active binary version\n" +
	"      2. gitmap update                # rebuild + redeploy from current source\n" +
	"      3. open a NEW terminal so PATH refreshes\n" +
	"      4. gitmap pending clear --yes   # drop any orphaned pending rows\n" +
	"      5. retry the clone — `gitmap clone <url1> <url2>` works either\n" +
	"         space-separated or comma-separated in PowerShell and bash.\n"
