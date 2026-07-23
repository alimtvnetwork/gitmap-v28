package formatter

// Shared warning emitter used by every record-aware writer (WriteJSON,
// WriteCSV) before serializing. The default destination is os.Stderr,
// but tests can swap it via SetValidationSink to capture warnings.

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// validationSink is the io.Writer that receives per-issue warning lines.
// Guarded by sinkMu so concurrent writers + tests can swap it safely.
var (
	validationSink io.Writer = os.Stderr
	sinkMu         sync.RWMutex
)

// SetValidationSink redirects validation warnings to w. Pass os.Stderr to
// restore the default. Returns the previous sink so tests can defer-restore.
func SetValidationSink(w io.Writer) io.Writer {
	sinkMu.Lock()
	defer sinkMu.Unlock()
	prev := validationSink
	validationSink = w

	return prev
}

// emitValidationWarnings runs the validator over records and writes one
// `gitmap: validation: <issue>` line per finding to the active sink.
// Returns the number of issues found so the caller can include the
// count in a post-write summary line. Never returns an error — by
// policy the write proceeds regardless.
func emitValidationWarnings(records []model.ScanRecord) int {
	issues := ValidateRecords(records)
	if len(issues) == 0 {
		return 0
	}

	w := activeSink()
	for _, issue := range issues {
		fmt.Fprintf(w, "gitmap: validation: %s\n", issue.String())
	}

	return len(issues)
}

// emitWriteSummary prints a one-line tally after a successful write so
// users see the outcome without having to count stderr lines. Format:
//
//	gitmap: <format>: wrote N record(s), M validation issue(s)
//
// Always goes to the same sink as the per-issue warnings so test capture
// stays consistent.
func emitWriteSummary(format string, recordCount, issueCount int) {
	fmt.Fprintf(activeSink(), "gitmap: %s: wrote %d record(s), %d validation issue(s)\n",
		format, recordCount, issueCount)
}

// activeSink returns the current validation sink under a read lock.
// Extracted so both emitters share the same lookup path.
func activeSink() io.Writer {
	sinkMu.RLock()
	defer sinkMu.RUnlock()

	return validationSink
}
