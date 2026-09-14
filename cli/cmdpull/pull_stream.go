package cmdpull

import (
	"regexp"
	"strconv"
	"strings"
)

// PullProgressEvent encapsulates a parsed git progress milestone.
type PullProgressEvent struct {
	Phase       string
	Percent     int
	Step        PullStepType
	Description string
	CommitRange string
}

var (
	percentRegex  = regexp.MustCompile(`(\d+)%`)
	updatingRegex = regexp.MustCompile(`Updating\s+([0-9a-fA-F]+\.\.[0-9a-fA-F]+)`)
)

// ParseGitProgressLine parses a single progress line emitted by git.
func ParseGitProgressLine(line string) (PullProgressEvent, bool) {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return PullProgressEvent{}, false
	}
	event, isProgress := parseGitObjectProgress(trimmed)
	if isProgress {
		return event, true
	}

	return parseGitRefUpdate(trimmed)
}

func parseGitObjectProgress(line string) (PullProgressEvent, bool) {
	lower := strings.ToLower(line)
	phase, step := matchProgressPhase(lower)
	if phase == "" {
		return PullProgressEvent{}, false
	}
	percent := extractProgressPercent(line)
	desc := formatProgressDesc(phase, percent)

	return PullProgressEvent{
		Phase:       phase,
		Percent:     percent,
		Step:        step,
		Description: desc,
	}, true
}

func matchProgressPhase(lower string) (string, PullStepType) {
	if strings.Contains(lower, "counting objects") {
		return "Counting", PullStepTypeFetching
	}
	if strings.Contains(lower, "compressing objects") {
		return "Compressing", PullStepTypeFetching
	}
	if strings.Contains(lower, "receiving objects") {
		return "Receiving", PullStepTypeFetching
	}
	if strings.Contains(lower, "resolving deltas") {
		return "Resolving deltas", PullStepTypeFetching
	}

	return "", PullStepTypePending
}

func extractProgressPercent(line string) int {
	matches := percentRegex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return 0
	}
	val, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return val
}

func formatProgressDesc(phase string, percent int) string {
	if percent > 0 {
		return phase + " (" + strconv.Itoa(percent) + "%)"
	}

	return phase
}

func parseGitRefUpdate(line string) (PullProgressEvent, bool) {
	matches := updatingRegex.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return PullProgressEvent{
			Phase:       "Fast-forward",
			Percent:     90,
			Step:        PullStepTypeFastForward,
			Description: "Fast-forward " + matches[1],
			CommitRange: matches[1],
		}, true
	}
	if strings.Contains(strings.ToLower(line), "fast-forward") {
		return PullProgressEvent{
			Phase:       "Fast-forward",
			Percent:     95,
			Step:        PullStepTypeFastForward,
			Description: "Fast-forward",
		}, true
	}

	return PullProgressEvent{}, false
}

// MakePullProgressHandler returns a callback that parses git lines and forwards them to the progress bar.
func MakePullProgressHandler(bar *PullProgressBar, workerID int, repoName string) func(string) {
	return func(line string) {
		if bar == nil {
			return
		}
		event, hasEvent := ParseGitProgressLine(line)
		if !hasEvent {
			return
		}
		dispatchProgressEvent(bar, workerID, repoName, event)
	}
}

func dispatchProgressEvent(bar *PullProgressBar, workerID int, repoName string, event PullProgressEvent) {
	isWorkerActive := workerID >= 0
	if isWorkerActive {
		bar.UpdateWorkerProgress(workerID, event.Step, event.Description, event.Percent)
		return
	}
	bar.UpdateRepoStep(repoName, event.Step, event.Description)
	bar.SetActivePercent(event.Percent)
}
