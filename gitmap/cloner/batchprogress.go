package cloner

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// FailureRecord stores details about a single failed batch item.
type FailureRecord struct {
	Name  string
	Error string
}

// BatchProgress tracks progress for any batch operation (pull, exec, status).
type BatchProgress struct {
	total      int
	current    int
	start      time.Time
	quiet      bool
	succeeded  int
	failed     int
	skipped    int
	operation  string
	failures   []FailureRecord
	stopOnFail bool
	stopped    bool
	lastName   string
}

// NewBatchProgress creates a progress tracker for a named operation.
func NewBatchProgress(total int, operation string, quiet bool) *BatchProgress {
	return &BatchProgress{
		total:     total,
		start:     time.Now(),
		quiet:     quiet,
		operation: operation,
	}
}

// SetStopOnFail enables early termination after the first failure.
func (p *BatchProgress) SetStopOnFail(v bool) { p.stopOnFail = v }

// Stopped returns true if the batch was halted due to --stop-on-fail.
func (p *BatchProgress) Stopped() bool { return p.stopped }

// BeginItem prints progress for starting an item.
func (p *BatchProgress) BeginItem(name string) {
	p.current++
	p.lastName = name
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressBeginFmt, p.current, p.total, name)
}

// Succeed marks an item as successful.
func (p *BatchProgress) Succeed() {
	p.succeeded++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressDoneFmt, formatDuration(time.Since(p.start)))
}

// Fail marks an item as failed.
func (p *BatchProgress) Fail() {
	p.failed++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressFailFmt)
}

// FailWithError marks an item as failed and records the error detail.
func (p *BatchProgress) FailWithError(name, errMsg string) {
	p.failed++
	p.failures = append(p.failures, FailureRecord{Name: name, Error: errMsg})
	if p.stopOnFail {
		p.stopped = true
	}
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressFailFmt)
}

// Skip marks an item as skipped (e.g., missing directory).
func (p *BatchProgress) Skip() {
	p.skipped++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.BatchProgressSkipFmt)
}

// PrintSummary prints the final summary.
func (p *BatchProgress) PrintSummary() {
	if p.quiet {
		return
	}

	elapsed := formatDuration(time.Since(p.start))
	fmt.Fprintf(os.Stderr, constants.BatchProgressSummaryFmt,
		p.operation, p.current, p.total, elapsed)
	fmt.Fprintf(os.Stderr, constants.BatchProgressDetailFmt,
		p.succeeded, p.failed, p.skipped)

	if p.stopped {
		fmt.Fprintf(os.Stderr, constants.BatchStoppedMsg)
	}
}

// Succeeded returns the success count.
func (p *BatchProgress) Succeeded() int { return p.succeeded }

// Failed returns the failure count.
func (p *BatchProgress) Failed() int { return p.failed }

// Skipped returns the skip count.
func (p *BatchProgress) Skipped() int { return p.skipped }

// Failures returns all recorded failure details.
func (p *BatchProgress) Failures() []FailureRecord { return p.failures }

// HasFailures returns true if any items failed.
func (p *BatchProgress) HasFailures() bool { return len(p.failures) > 0 }
