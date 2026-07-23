package finalize

import (
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Counters captures the per-run tally rendered by the summary line.
type Counters struct {
	RunId   int64
	Created int
	Skipped int
	Failed  int
}

// Outcome maps the counters to the spec §2.7 exit code:
//   - Failed > 0 AND Created == 0 → SourceUnusable / PartiallyFailed
//     decision is up to the caller; we only resolve the
//     "any-failure" axis.
//   - Failed > 0 AND Created > 0 → PartiallyFailed.
//   - else → Ok.
func Outcome(c Counters) int {
	if c.Failed == 0 {
		return constants.CommitInExitOk
	}
	return constants.CommitInExitPartiallyFailed
}

// PrintSummary writes the canonical summary line to `out` (typically
// os.Stderr). Using the constant format keeps the line stable across
// callers and easy to grep in CI logs.
func PrintSummary(out io.Writer, c Counters) {
	fmt.Fprintf(out, constants.CommitInMsgSummaryLine,
		c.RunId, c.Created, c.Skipped, c.Failed)
}

// PrintDryRunBanner emits a one-line "DRY RUN — no commits created"
// notice so users grepping CI logs cannot mistake a dry run for a real
// run with zero produced commits. Caller responsibility to invoke only
// when --dry-run was actually set.
func PrintDryRunBanner(out io.Writer) {
	fmt.Fprint(out, "commit-in: DRY RUN — no commits were created\n")
}

// CleanupTemp removes the run's temp staging dir unless --keep-temp
// was passed. Errors are logged to stderr but never escalated; temp
// leakage is preferable to masking a real run failure.
func CleanupTemp(tempRunDir string, keepTemp bool) {
	if keepTemp || tempRunDir == "" {
		return
	}
	if err := os.RemoveAll(tempRunDir); err != nil {
		fmt.Fprintf(os.Stderr, "commit-in: cleanup: %v\n", err)
	}
}
