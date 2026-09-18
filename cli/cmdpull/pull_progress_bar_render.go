package cmdpull

import (
	"fmt"
	"time"
)

// Render redraws the bar immediately.
func (p *PullProgressBar) Render() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.renderLocked()
}

func (p *PullProgressBar) renderLocked() {
	if p.isQuiet {
		return
	}
	if p.isTTY {
		p.renderTTY()
		return
	}
	p.renderNonTTY()
}

func (p *PullProgressBar) currentSpinner() string {
	frame := brailleFrames[p.spinnerIdx%len(brailleFrames)]
	if p.isSafe {
		frame = asciiFrames[p.spinnerIdx%len(asciiFrames)]
	}

	return fmt.Sprintf("%s%s%s", constants.ColorCyan, frame, constants.ColorReset)
}

func (p *PullProgressBar) renderTTY() {
	if p.total <= 1 {
		p.renderSingleRepoTTY()
		return
	}
	p.renderMultiRepoTTY()
}

func (p *PullProgressBar) renderSingleRepoTTY() {
	percent := p.resolveSingleRepoPercent()
	bar := FormatVisualBar(percent, 100, 20, p.isSafe)
	spinner := p.currentSpinner()
	badge := FormatStepBadge(p.activeStep, p.isSafe)
	stepLabel := fmt.Sprintf("[Step %d/%d] %s", p.subStepIndex, p.subStepTotal, StepMilestoneTitle(p.subStepIndex, p.activeStep))
	elapsed := fmt.Sprintf("%.1fs", time.Since(p.startTime).Seconds())
	desc := p.formatActiveDescription()

	line := fmt.Sprintf("%s %s %3d%% %s | %s %s%s (%s)", spinner, bar, percent, stepLabel, badge, p.activeRepo, desc, elapsed)
	line = clampLineToTermWidth(line)
	fmt.Fprintf(p.out, "\r\033[K%s", line)
}

func (p *PullProgressBar) renderMultiRepoTTY() {
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 20, p.isSafe)
	spinner := p.currentSpinner()
	counter := fmt.Sprintf("%3d%% (%d/%d repos)", percent, p.completed, p.total)
	elapsed := fmt.Sprintf("%.1fs", time.Since(p.startTime).Seconds())
	summary := p.formatWorkersSummary()

	line := fmt.Sprintf("%s %s %s | %s (%s)", spinner, bar, counter, summary, elapsed)
	line = clampLineToTermWidth(line)
	fmt.Fprintf(p.out, "\r\033[K%s", line)
}

func clampLineToTermWidth(line string) string {
	termWidth := detectTerminalWidth()
	maxLen := termWidth - 1
	if maxLen < 30 {
		maxLen = 79
	}
	if len(line) <= maxLen {
		return line
	}
	if maxLen <= 3 {
		return line[:maxLen]
	}

	return line[:maxLen-3] + "..."
}

func (p *PullProgressBar) formatActiveDescription() string {
	if p.activeDesc == "" {
		return ""
	}

	return ": " + p.activeDesc
}
