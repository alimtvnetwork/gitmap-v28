package cmdpull

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ColorizeVisualBar wraps visual bar in bright green when 100% complete, or cyan while in progress.
func ColorizeVisualBar(bar string, percent int) string {
	if percent >= 100 {
		return constants.ColorGreen + bar + constants.ColorReset
	}

	return constants.ColorCyan + bar + constants.ColorReset
}

func (p *PullProgressBar) renderNonTTY() {
	if isIntermediateGitMilestone(p.activeDesc) {
		return
	}
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 20, p.isSafe)
	counter := fmt.Sprintf("%3d%% (%d/%d repos)", percent, p.completed, p.total)
	badge := p.colorizeBadge(FormatStepBadge(p.activeStep, p.isSafe), p.activeStep)
	desc := p.formatActiveDescription()

	fmt.Fprintf(p.out, "  %s %s | %s %s%s\n", p.colorizeBar(bar, percent), counter, badge, p.activeRepo, desc)
}

func isIntermediateGitMilestone(desc string) bool {
	lower := strings.ToLower(desc)
	return strings.Contains(lower, "counting") ||
		strings.Contains(lower, "compressing") ||
		strings.Contains(lower, "receiving") ||
		strings.Contains(lower, "resolving")
}

func (p *PullProgressBar) resolveSingleRepoPercent() int {
	if p.completed >= 1 {
		return 100
	}
	if p.activePercent > 0 {
		return p.activePercent
	}

	return StepMilestonePercent(p.activeStep)
}

func calcProgressPercent(completed, total int) int {
	if total <= 0 || completed <= 0 {
		return 0
	}
	if completed >= total {
		return 100
	}

	return (completed * 100) / total
}

// FormatVisualBar returns ASCII or block visual bar.
func FormatVisualBar(completed, total, barWidth int, isSafe bool) string {
	if barWidth <= 0 {
		barWidth = 20
	}
	filledLen := calcFilledBarLength(completed, total, barWidth)
	emptyLen := barWidth - filledLen
	fillChar, emptyChar := resolveBarChars(isSafe)

	return fmt.Sprintf("[%s%s]", strings.Repeat(fillChar, filledLen), strings.Repeat(emptyChar, emptyLen))
}

func calcFilledBarLength(completed, total, barWidth int) int {
	if total <= 0 || completed <= 0 {
		return 0
	}
	if completed >= total {
		return barWidth
	}

	return (completed * barWidth) / total
}

func resolveBarChars(isSafe bool) (string, string) {
	if isSafe {
		return constants.ProgressBarFilledSafe, constants.ProgressBarEmptySafe
	}

	return constants.ProgressBarFilledRich, constants.ProgressBarEmptyRich
}
