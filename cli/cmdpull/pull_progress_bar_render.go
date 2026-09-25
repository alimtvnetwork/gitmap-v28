package cmdpull

import (
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	coloredBar := p.colorizeBar(bar, percent)
	spinner := p.currentSpinner()
	badge := p.colorizeBadge(FormatStepBadge(p.activeStep, p.isSafe), p.activeStep)
	stepLabel := fmt.Sprintf("[Step %d/%d] %s", p.subStepIndex, p.subStepTotal, StepMilestoneTitle(p.subStepIndex, p.activeStep))
	elapsed := fmt.Sprintf("%s%.1fs%s", constants.ColorDim, time.Since(p.startTime).Seconds(), constants.ColorReset)
	desc := p.formatActiveDescription()

	line := fmt.Sprintf("  %s %s %3d%% %s | %s %s%s (%s)", spinner, coloredBar, percent, stepLabel, badge, p.activeRepo, desc, elapsed)
	line = clampLineToTermWidth(line)
	fmt.Fprintf(p.out, "\r\033[K%s", line)
}

func (p *PullProgressBar) renderMultiRepoTTY() {
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 24, p.isSafe)
	coloredBar := p.colorizeBar(bar, percent)
	spinner := p.currentSpinner()
	counter := fmt.Sprintf("%s%3d%%%s (%d/%d repos)", constants.ColorBold, percent, constants.ColorReset, p.completed, p.total)
	elapsed := fmt.Sprintf("%s%.1fs%s", constants.ColorDim, time.Since(p.startTime).Seconds(), constants.ColorReset)
	summary := p.formatWorkersSummary()

	line := fmt.Sprintf("  %s %s %s | %s (%s)", spinner, coloredBar, counter, summary, elapsed)
	line = clampLineToTermWidth(line)
	fmt.Fprintf(p.out, "\r\033[K%s", line)
}

func (p *PullProgressBar) colorizeBar(bar string, percent int) string {
	if percent >= 100 {
		return constants.ColorGreen + bar + constants.ColorReset
	}

	return constants.ColorCyan + bar + constants.ColorReset
}

func (p *PullProgressBar) colorizeBadge(badge string, step PullStepType) string {
	switch step {
	case PullStepTypeFetching:
		return constants.ColorCyan + badge + constants.ColorReset
	case PullStepTypeMerging:
		return constants.ColorMagenta + badge + constants.ColorReset
	case PullStepTypeUpToDate:
		return constants.ColorGreen + badge + constants.ColorReset
	case PullStepTypeFastForward:
		return constants.ColorYellow + badge + constants.ColorReset
	case PullStepTypeConflict, PullStepTypeError:
		return constants.ColorRed + badge + constants.ColorReset
	case PullStepTypeInspecting:
		return constants.ColorCyan + badge + constants.ColorReset
	case PullStepTypePending, PullStepTypeSkipped:
		return badge
	default:
		return badge
	}
}

func clampLineToTermWidth(line string) string {
	termWidth := detectTerminalWidth()
	maxLen := termWidth - 1
	if maxLen < 30 {
		maxLen = 79
	}
	plain := stripANSI(line)
	if len(plain) <= maxLen {
		return line
	}
	if maxLen <= 3 {
		return plain[:maxLen]
	}

	return plain[:maxLen-3] + "..."
}

func (p *PullProgressBar) formatActiveDescription() string {
	if p.activeDesc == "" {
		return ""
	}

	return ": " + p.activeDesc
}
