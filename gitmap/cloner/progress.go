package cloner

import (
	"fmt"

	"sync"
	"time"


	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/pterm/pterm"
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
	quiet    bool
	cloned   int
	pulled   int
	skipped  int
	failed   int
	multi    *pterm.MultiPrinter
	spinners map[string]*pterm.SpinnerPrinter
}

// NewProgress creates a progress tracker.
func NewProgress(total int, quiet bool) *Progress {
	p := &Progress{
		total:    total,
		start:    time.Now(),
		quiet:    quiet,
		spinners: make(map[string]*pterm.SpinnerPrinter),
	}
	if !quiet {
		p.multi = &pterm.DefaultMultiPrinter
		p.multi.Start()
		pterm.Info.Printf("[gitmap] Processing %d repositories...\n", total)
	}
	return p
}

// Begin prints the starting line for a repo.
func (p *Progress) Begin(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current++
	if p.quiet {
		return
	}

	spinner, _ := pterm.DefaultSpinner.WithWriter(p.multi.NewWriter()).Start(fmt.Sprintf("Processing %s...", name))
	p.spinners[name] = spinner
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

	name := repoDisplayName(result.Record)
	if spinner, ok := p.spinners[name]; ok {
		elapsed := time.Since(p.start)
		if pulled {
			spinner.Success(fmt.Sprintf("%s updated (pull) in %s", name, formatDuration(elapsed)))
		} else {
			spinner.Success(fmt.Sprintf("%s cloned in %s", name, formatDuration(elapsed)))
		}
		delete(p.spinners, name)
	}
}

// Skip marks a repo as skipped because it was already up to date.
func (p *Progress) Skip(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.skipped++
	if p.quiet {
		return
	}

	name := repoDisplayName(result.Record)
	if spinner, ok := p.spinners[name]; ok {
		spinner.Warning(fmt.Sprintf("%s skipped (existing/up-to-date)", name))
		delete(p.spinners, name)
	}
}

// Fail marks a repo as failed.
func (p *Progress) Fail(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failed++
	if p.quiet {
		return
	}

	name := repoDisplayName(result.Record)
	if spinner, ok := p.spinners[name]; ok {
		spinner.Fail(fmt.Sprintf("%s failed", name))
		delete(p.spinners, name)
	}
}

// PrintSummary prints the final summary line.
func (p *Progress) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.quiet {
		return
	}

	if p.multi != nil {
		p.multi.Stop()
	}

	elapsed := time.Since(p.start)
	fmt.Println()
	pterm.DefaultBox.WithTitle("Summary").Println(
		fmt.Sprintf("Completed: %d, Skipped: %d, Failed: %d\nTime: %s",
			p.cloned+p.pulled, p.skipped, p.failed, formatDuration(elapsed)))
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
