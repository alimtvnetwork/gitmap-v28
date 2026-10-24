package cmdpull

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull/enums"
)

// PullStepType represents the lifecycle step or final state of a pull operation.
type PullStepType = enums.PullStepType

const (
	PullStepTypePending     = enums.PullStepPending
	PullStepTypeInspecting  = enums.PullStepInspecting
	PullStepTypeFetching    = enums.PullStepFetching
	PullStepTypeMerging     = enums.PullStepMerging
	PullStepTypeUpToDate    = enums.PullStepUpToDate
	PullStepTypeFastForward = enums.PullStepFastForward
	PullStepTypeConflict    = enums.PullStepConflict
	PullStepTypeError       = enums.PullStepError
	PullStepTypeSkipped     = enums.PullStepSkipped
)

// Backward compatibility aliases.
const (
	PullStepInvalid     = enums.PullStepInvalid
	PullStepPending     = enums.PullStepPending
	PullStepInspecting  = enums.PullStepInspecting
	PullStepFetching    = enums.PullStepFetching
	PullStepMerging     = enums.PullStepMerging
	PullStepUpToDate    = enums.PullStepUpToDate
	PullStepFastForward = enums.PullStepFastForward
	PullStepConflict    = enums.PullStepConflict
	PullStepError       = enums.PullStepError
	PullStepSkipped     = enums.PullStepSkipped
)

// PullRepoState holds runtime execution and completion state for a single repository pull.
type PullRepoState struct {
	RepoName     string
	RepoPath     string
	Branch       string
	OldSHA       string
	NewSHA       string
	CommitRange  string
	Changes      string
	Step         PullStepType
	Duration     time.Duration
	IsDirty      bool
	IsAutoMerged bool
	ErrorMsg     string
}

var safeBadges = map[PullStepType]string{
	PullStepTypePending:     "[.]",
	PullStepTypeInspecting:  "[?]",
	PullStepTypeFetching:    "[v]",
	PullStepTypeMerging:     "[+]",
	PullStepTypeUpToDate:    "[=]",
	PullStepTypeFastForward: "[^]",
	PullStepTypeConflict:    "[!]",
	PullStepTypeError:       "[X]",
	PullStepTypeSkipped:     "[-]",
}

// FormatStepBadge formats a step indicator badge adhering to safe or rich glyph mode.
func FormatStepBadge(step PullStepType, isSafe bool) string {
	if isSafe {
		return formatStepBadgeSafe(step)
	}

	return formatStepBadgeRich(step)
}

func formatStepBadgeSafe(step PullStepType) string {
	badge, isFound := safeBadges[step]
	if isFound {
		return badge
	}

	return "[?]"
}

func formatStepBadgeRich(step PullStepType) string {
	badge := step.Badge()
	if badge != "" {
		return badge
	}

	return "·"
}

// ParseGitPullOutput classifies pull output, calculating step, commit range, and changes.
func ParseGitPullOutput(output, oldSHA, newSHA string) (PullStepType, string, string) {
	step := detectOutputStep(output, oldSHA, newSHA)
	commitRange := resolveCommitRange(step, oldSHA, newSHA)
	changes := resolveChangesSummary(step, output)

	return step, commitRange, changes
}

func detectOutputFailureStep(lower string) (PullStepType, bool) {
	if strings.Contains(lower, "conflict") || strings.Contains(lower, "automatic merge failed") {
		return PullStepTypeConflict, true
	}
	if strings.Contains(lower, "fatal:") || strings.Contains(lower, "error:") {
		return PullStepTypeError, true
	}

	return PullStepTypePending, false
}

func detectOutputSuccessStep(lower, oldSHA, newSHA string) PullStepType {
	if strings.Contains(lower, "already up to date") || strings.Contains(lower, "already up-to-date") {
		return PullStepTypeUpToDate
	}
	if strings.Contains(lower, "fast-forward") || (oldSHA != "" && newSHA != "" && oldSHA != newSHA) {
		return PullStepTypeFastForward
	}
	if strings.Contains(lower, "merge made by") {
		return PullStepTypeMerging
	}

	return PullStepTypeUpToDate
}

func detectOutputStep(output, oldSHA, newSHA string) PullStepType {
	lower := strings.ToLower(output)
	step, isFailed := detectOutputFailureStep(lower)
	if isFailed {
		return step
	}

	return detectOutputSuccessStep(lower, oldSHA, newSHA)
}

func resolveCommitRange(step PullStepType, oldSHA, newSHA string) string {
	if step == PullStepTypeConflict || step == PullStepTypeError {
		return ""
	}

	return FormatCommitRange(oldSHA, newSHA)
}

func resolveChangesSummary(step PullStepType, output string) string {
	if step == PullStepTypeConflict {
		return "conflict"
	}
	if step == PullStepTypeError {
		return "error"
	}
	if step == PullStepTypeUpToDate {
		return "up-to-date"
	}

	return extractDiffStatSummary(output)
}

func extractDiffStatSummary(output string) string {
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(trimmed, "changed") && strings.Contains(trimmed, "insertion") {
			return FormatDiffStat(trimmed)
		}
	}

	return "updated"
}

// StepMilestonePercent returns the milestone percentage associated with each pull step.
func StepMilestonePercent(step PullStepType) int {
	if step == PullStepTypeFetching {
		return 35
	}
	if step == PullStepTypeMerging {
		return 75
	}
	if step == PullStepTypeFastForward {
		return 90
	}

	return resolveStepMilestoneEdge(step)
}

func resolveStepMilestoneEdge(step PullStepType) int {
	if step == PullStepTypeInspecting {
		return 10
	}
	if step == PullStepTypePending {
		return 0
	}

	return 100
}

// StepSubStepIndex maps step to a 1-based index (1..4) for single-repo progression.
func StepSubStepIndex(step PullStepType) int {
	switch step {
	case PullStepTypePending, PullStepTypeInspecting:
		return 1
	case PullStepTypeFetching:
		return 2
	case PullStepTypeMerging, PullStepTypeFastForward:
		return 3
	case PullStepTypeUpToDate, PullStepTypeConflict, PullStepTypeError, PullStepTypeSkipped:
		return 4
	default:
		return 4
	}
}

// StepMilestoneTitle returns the single-repo milestone description.
func StepMilestoneTitle(index int, step PullStepType) string {
	switch index {
	case 1:
		return "Inspecting"
	case 2:
		return "Fetching remote objects"
	case 3:
		return "Fast-forwarding / Merging"
	default:
		return resolveStepCompletionTitle(step)
	}
}

func resolveStepCompletionTitle(step PullStepType) string {
	if step == PullStepTypeUpToDate {
		return "Complete (up-to-date)"
	}
	if step == PullStepTypeError {
		return "Complete (failed)"
	}

	return "Complete"
}
