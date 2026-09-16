package cmdpull

import (
	"fmt"
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
	badge := FormatStepBadge(p.activeStep, p.isSafe)
	desc := p.formatActiveDescription()

	return fmt.Sprintf("%s %s%s", badge, p.activeRepo, desc)
}

func (p *PullProgressBar) formatActiveWorkersList() string {
	workerIDs := p.sortedWorkerIDs()
	if len(workerIDs) == 0 {
		return ""
	}
	if len(workerIDs) == 1 {
		w := p.activeWorkers[workerIDs[0]]
		return fmt.Sprintf("%s %s%s", FormatStepBadge(w.Step, p.isSafe), w.RepoName, formatWorkerSlotDesc(w))
	}
	var parts []string
	for i := 0; i < len(workerIDs) && i < 2; i++ {
		w := p.activeWorkers[workerIDs[i]]
		parts = append(parts, fmt.Sprintf("%s %s", FormatStepBadge(w.Step, p.isSafe), w.RepoName))
	}
	summary := strings.Join(parts, ", ")
	if len(workerIDs) > 2 {
		summary += fmt.Sprintf(" (+%d more)", len(workerIDs)-2)
	}

	return summary
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
