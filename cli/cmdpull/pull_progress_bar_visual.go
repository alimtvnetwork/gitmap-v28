package cmdpull

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func (p *PullProgressBar) renderNonTTY() {
	percent := calcProgressPercent(p.completed, p.total)
	bar := FormatVisualBar(p.completed, p.total, 20, p.isSafe)
	counter := fmt.Sprintf("%3d%% (%d/%d repos)", percent, p.completed, p.total)
	badge := p.colorizeBadge(FormatStepBadge(p.activeStep, p.isSafe), p.activeStep)
	desc := p.formatActiveDescription()

	fmt.Fprintf(p.out, "  %s %s | %s %s%s\n", p.colorizeBar(bar, percent), counter, badge, p.activeRepo, desc)
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
