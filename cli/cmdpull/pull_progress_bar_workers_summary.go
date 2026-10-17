package cmdpull

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

func (p *PullProgressBar) formatWorkersSummary() string {
	if len(p.activeWorkers) == 0 {
		return p.formatFallbackSummary()
	}

	return p.formatActiveWorkersList()
}

func (p *PullProgressBar) formatFallbackSummary() string {
	badge := p.colorizeBadge(FormatStepBadge(p.activeStep, p.isSafe), p.activeStep)
	desc := p.formatActiveDescription()

	return fmt.Sprintf("%s %s%s", badge, cleanActiveRepoName(p.activeRepo), desc)
}

func (p *PullProgressBar) formatActiveWorkersList() string {
	workerIDs := p.sortedWorkerIDs()
	if len(workerIDs) == 0 {
		return ""
	}
	if len(workerIDs) == 1 {
		return p.formatSingleActiveWorker(p.activeWorkers[workerIDs[0]])
	}
	maxWorkers := p.resolveMaxWorkersToDisplay()
	var parts []string
	for i := 0; i < len(workerIDs) && i < maxWorkers; i++ {
		w := p.activeWorkers[workerIDs[i]]
		if w == nil {
			continue
		}
		badge := p.colorizeBadge(FormatStepBadge(w.Step, p.isSafe), w.Step)
		repoName := cleanActiveRepoName(w.RepoName)
		parts = append(parts, fmt.Sprintf("%s %s", badge, repoName))
	}
	summary := strings.Join(parts, ", ")
	if len(workerIDs) > maxWorkers {
		summary += fmt.Sprintf(" (+%d more)", len(workerIDs)-maxWorkers)
	}

	return summary
}

func (p *PullProgressBar) resolveMaxWorkersToDisplay() int {
	if p.isTTY && detectTerminalWidth() >= widePullTableThreshold {
		return 4
	}

	return 2
}

func cleanActiveRepoName(repo string) string {
	clean := strings.TrimSpace(repo)
	if strings.Contains(clean, "/") || strings.Contains(clean, "\\") {
		clean = filepath.Base(filepath.ToSlash(clean))
	}

	return clean
}

func formatWorkerSlotDesc(w *WorkerSlotState) string {
	if w.Description == "" {
		return ""
	}

	return ": " + w.Description
}

func (p *PullProgressBar) sortedWorkerIDs() []int {
	ids := make([]int, 0, len(p.activeWorkers))
	for id := range p.activeWorkers {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	return ids
}

func (p *PullProgressBar) formatSingleActiveWorker(w *WorkerSlotState) string {
	if w == nil {
		return ""
	}

	badge := p.colorizeBadge(FormatStepBadge(w.Step, p.isSafe), w.Step)
	repoName := cleanActiveRepoName(w.RepoName)

	return fmt.Sprintf("%s %s%s", badge, repoName, formatWorkerSlotDesc(w))
}
