package cloner

import (
	"fmt"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

// FailureRecord stores details about a single failed batch item.
type FailureRecord struct {
	Name  string
	Error string
}

// BatchProgress tracks progress for any batch operation (pull, exec, status).
// It has been upgraded to use a bounded worker pool dynamic spinner UI.
type BatchProgress struct {
	mu         sync.Mutex
	total      int
	current    int
	start      time.Time
	quiet      bool
	succeeded  int
	failed     int
	skipped    int
	upToDate   int
	operation  string
	failures   []FailureRecord
	stopOnFail bool
	stopped    bool

	multi    *pterm.MultiPrinter
	spinners map[string]*pterm.SpinnerPrinter
}

// NewBatchProgress creates a progress tracker for a named operation.
func NewBatchProgress(total int, operation string, quiet bool) *BatchProgress {
	p := &BatchProgress{
		total:     total,
		start:     time.Now(),
		quiet:     quiet,
		operation: operation,
		spinners:  make(map[string]*pterm.SpinnerPrinter),
	}
	if !quiet {
		p.multi = &pterm.DefaultMultiPrinter
		p.multi.Start()
		pterm.Info.Printf("[gitmap] Processing %d repositories for %s...\n", total, operation)
	}
	return p
}

// SetStopOnFail enables early termination after the first failure.
func (p *BatchProgress) SetStopOnFail(v bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.stopOnFail = v
}

// Stopped returns true if the batch was halted due to --stop-on-fail.
func (p *BatchProgress) Stopped() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.stopped
}

// BeginItem prints progress for starting an item.
func (p *BatchProgress) BeginItem(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current++
	if p.quiet {
		return
	}

	spinner, _ := pterm.DefaultSpinner.WithWriter(p.multi.NewWriter()).Start(fmt.Sprintf("%s %s...", p.operation, name))
	p.spinners[name] = spinner
}

// Succeed marks an item as successful.
func (p *BatchProgress) Succeed(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.succeeded++
	if p.quiet {
		return
	}

	if spinner, ok := p.spinners[name]; ok {
		spinner.Success(fmt.Sprintf("%s %s in %s", name, p.operation, formatDuration(time.Since(p.start))))
		delete(p.spinners, name)
	}
}

// UpToDate marks an item as successful but up-to-date (no changes pulled).
func (p *BatchProgress) UpToDate(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.upToDate++
	p.succeeded++ // up-to-date counts as success
	if p.quiet {
		return
	}

	if spinner, ok := p.spinners[name]; ok {
		spinner.SuccessPrinter = &pterm.PrefixPrinter{
			MessageStyle: &pterm.Style{pterm.FgDarkGray},
			Prefix: pterm.Prefix{
				Text:  "~",
				Style: &pterm.Style{pterm.FgDarkGray},
			},
		}
		spinner.Success(fmt.Sprintf("%s (up to date)", name))
		delete(p.spinners, name)
	}
}

// Fail marks an item as failed.
func (p *BatchProgress) Fail(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failed++
	if p.quiet {
		return
	}

	if spinner, ok := p.spinners[name]; ok {
		spinner.Fail(fmt.Sprintf("%s failed", name))
		delete(p.spinners, name)
	}
}

// FailWithError marks an item as failed and records the error detail.
func (p *BatchProgress) FailWithError(name, errMsg string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failed++
	p.failures = append(p.failures, FailureRecord{Name: name, Error: errMsg})
	if p.stopOnFail {
		p.stopped = true
	}
	if p.quiet {
		return
	}

	if spinner, ok := p.spinners[name]; ok {
		spinner.Fail(fmt.Sprintf("%s failed", name))
		delete(p.spinners, name)
	}
}

// Skip marks an item as skipped (e.g., missing directory).
func (p *BatchProgress) Skip(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.skipped++
	if p.quiet {
		return
	}

	if spinner, ok := p.spinners[name]; ok {
		spinner.Warning(fmt.Sprintf("%s skipped", name))
		delete(p.spinners, name)
	}
}

// PrintSummary prints the final summary.
func (p *BatchProgress) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.quiet {
		return
	}

	if p.multi != nil {
		p.multi.Stop()
	}

	elapsed := formatDuration(time.Since(p.start))
	fmt.Println()
	pterm.DefaultBox.WithTitle("Summary").Println(
		fmt.Sprintf("Completed: %d, Skipped: %d, Failed: %d\nTime: %s",
			p.succeeded, p.skipped, p.failed, elapsed))
}

// Succeeded returns the success count.
func (p *BatchProgress) Succeeded() int { p.mu.Lock(); defer p.mu.Unlock(); return p.succeeded }

// Failed returns the failure count.
func (p *BatchProgress) Failed() int { p.mu.Lock(); defer p.mu.Unlock(); return p.failed }

// Skipped returns the skip count.
func (p *BatchProgress) Skipped() int { p.mu.Lock(); defer p.mu.Unlock(); return p.skipped }

// Failures returns all recorded failure details.
func (p *BatchProgress) Failures() []FailureRecord {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.failures
}

// HasFailures returns true if any items failed.
func (p *BatchProgress) HasFailures() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.failures) > 0
}
