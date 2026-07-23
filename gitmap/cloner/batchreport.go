package cloner

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// PrintFailureReport outputs a detailed list of failed items after a batch
// operation. Call after PrintSummary when HasFailures() is true.
func (p *BatchProgress) PrintFailureReport() {
	if len(p.failures) == 0 {
		return
	}

	fmt.Fprintf(os.Stderr, "\n%s\n", constants.BatchFailureHeader)
	for i, f := range p.failures {
		fmt.Fprintf(os.Stderr, constants.BatchFailureEntryFmt, i+1, f.Name, f.Error)
	}
	fmt.Fprintf(os.Stderr, constants.BatchFailureFooterFmt, len(p.failures))
}

// ExitCodeForBatch returns 0 if all items succeeded, or the partial-failure
// exit code if any items failed. Use with os.Exit() in the calling command.
func (p *BatchProgress) ExitCodeForBatch() int {
	if p.failed > 0 {
		return constants.ExitPartialFailure
	}

	return 0
}
