package cloner

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// Progress tracks clone operation progress.
//
// Thread-safety: all counter mutations and stderr writes go through mu so
// concurrent workers in the parallel runner (concurrent.go) cannot
// interleave half-written status lines or corrupt the running totals.
// The sequential runner pays only the cost of an uncontended mutex.
type Progress struct {
	mu      sync.Mutex
	total   int
	current int
	start   time.Time
	quiet   bool
	cloned  int
	pulled  int
	skipped int
	failed  int
}

// NewProgress creates a progress tracker.
func NewProgress(total int, quiet bool) *Progress {
	return &Progress{
		total: total,
		start: time.Now(),
		quiet: quiet,
	}
}

// Begin prints the starting line for a repo.
func (p *Progress) Begin(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ProgressBeginFmt, p.current, p.total, name)
}

// Done marks a repo as successfully completed.
func (p *Progress) Done(result model.CloneResult, pulled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if pulled {
		p.pulled++
	} else {
		p.cloned++
	}

	if p.quiet {
		return
	}

	elapsed := time.Since(p.start)
	fmt.Fprintf(os.Stderr, constants.ProgressDoneFmt, formatDuration(elapsed))
}

// Skip marks a repo as skipped because it was already up to date.
func (p *Progress) Skip(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.skipped++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ProgressSkipFmt)
}

// Fail marks a repo as failed.
func (p *Progress) Fail(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failed++
	if p.quiet {
		return
	}

	fmt.Fprintf(os.Stderr, constants.ProgressFailFmt)
}

// PrintSummary prints the final summary line.
func (p *Progress) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.quiet {
		return
	}

	elapsed := time.Since(p.start)
	fmt.Fprintf(os.Stderr, constants.ProgressSummaryFmt,
		p.current, p.total, formatDuration(elapsed))
	fmt.Fprintf(os.Stderr, constants.ProgressDetailFmt,
		p.cloned, p.pulled, p.skipped, p.failed)
}

// formatDuration returns a human-readable duration string.
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}

	mins := int(d.Minutes())
	secs := int(d.Seconds()) % 60

	return fmt.Sprintf("%dm %ds", mins, secs)
}
