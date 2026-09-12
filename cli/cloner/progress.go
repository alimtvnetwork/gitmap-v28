package cloner

import (
	"fmt"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// Progress tracks clone operation progress cleanly.
type Progress struct {
	mu      sync.Mutex
	total   int
	current int
	start   time.Time
	isQuiet bool
	cloned  int
	pulled  int
	skipped int
	failed  int
}

// NewProgress creates a progress tracker.
func NewProgress(total int, isQuiet bool) *Progress {
	p := &Progress{
		total:   total,
		start:   time.Now(),
		isQuiet: isQuiet,
	}

	if !isQuiet {
		fmt.Printf("  %s⚡ Parallel clone active: %d repositories%s\n\n",
			constants.ColorCyan, total, constants.ColorReset)
	}

	return p
}

// Begin records a repo processing start.
func (p *Progress) Begin(name string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current++
}

// Done marks a repo as successfully completed.
func (p *Progress) Done(result model.CloneResult, isPulled bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if isPulled {
		p.pulled++
	} else {
		p.cloned++
	}

	if p.isQuiet {
		return
	}

	name := repoDisplayName(result.Record)
	elapsed := time.Since(p.start)
	if isPulled {
		fmt.Printf("  [%2d/%d] 📂 %-32s %s✔ updated (pull) (%s)%s\n",
			p.cloned+p.pulled+p.skipped+p.failed, p.total, name,
			constants.ColorGreen, formatDuration(elapsed), constants.ColorReset)

		return
	}

	fmt.Printf("  [%2d/%d] 📂 %-32s %s✔ cloned (%s)%s\n",
		p.cloned+p.pulled+p.skipped+p.failed, p.total, name,
		constants.ColorGreen, formatDuration(elapsed), constants.ColorReset)
}

// Skip marks a repo as skipped because it was already up to date.
func (p *Progress) Skip(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.skipped++
	if p.isQuiet {
		return
	}

	name := repoDisplayName(result.Record)
	fmt.Printf("  [%2d/%d] 📂 %-32s %s✔ up-to-date (skipped)%s\n",
		p.cloned+p.pulled+p.skipped+p.failed, p.total, name,
		constants.ColorCyan, constants.ColorReset)
}

// Fail marks a repo as failed.
func (p *Progress) Fail(result model.CloneResult) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.failed++
	if p.isQuiet {
		return
	}

	name := repoDisplayName(result.Record)
	fmt.Printf("  [%2d/%d] 📂 %-32s %s✖ failed%s\n",
		p.cloned+p.pulled+p.skipped+p.failed, p.total, name,
		constants.ColorRed, constants.ColorReset)
}

// PrintSummary prints the final summary line.
func (p *Progress) PrintSummary() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isQuiet {
		return
	}

	elapsed := time.Since(p.start)
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %s✔ Clone complete: %d succeeded, %d skipped, %d failed · Elapsed: %s%s\n\n",
		constants.ColorGreen, p.cloned+p.pulled, p.skipped, p.failed, formatDuration(elapsed), constants.ColorReset)
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
