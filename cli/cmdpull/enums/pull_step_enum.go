package enums

// PullStepType represents the lifecycle step or final state of a pull operation.
type PullStepType byte

const (
	PullStepInvalid PullStepType = iota
	PullStepPending
	PullStepInspecting
	PullStepFetching
	PullStepMerging
	PullStepUpToDate
	PullStepFastForward
	PullStepConflict
	PullStepError
	PullStepSkipped
)

var stepNames = [...]string{
	PullStepInvalid:     "INVALID",
	PullStepPending:     "PENDING",
	PullStepInspecting:  "INSPECTING",
	PullStepFetching:    "FETCHING",
	PullStepMerging:     "MERGING",
	PullStepUpToDate:    "UP_TO_DATE",
	PullStepFastForward: "FAST_FORWARD",
	PullStepConflict:    "CONFLICT",
	PullStepError:       "ERROR",
	PullStepSkipped:     "SKIPPED",
}

var stepBadges = [...]string{
	PullStepInvalid:     "",
	PullStepPending:     "·",
	PullStepInspecting:  "🔍",
	PullStepFetching:    "↓",
	PullStepMerging:     "🔀",
	PullStepUpToDate:    "✔",
	PullStepFastForward: "⚡",
	PullStepConflict:    "✖",
	PullStepError:       "✖",
	PullStepSkipped:     "↷",
}

// IsValid reports whether the pull step is a recognized active or final step.
func (s PullStepType) IsValid() bool {
	return s >= PullStepPending && s <= PullStepSkipped
}

// String returns the uppercase identifier representation of the step.
func (s PullStepType) String() string {
	if int(s) >= len(stepNames) {
		return "INVALID"
	}

	return stepNames[s]
}

// Badge returns the rich terminal glyph badge representing the step.
func (s PullStepType) Badge() string {
	if int(s) >= len(stepBadges) {
		return ""
	}

	return stepBadges[s]
}

// IsComplete reports whether the step represents a terminal state.
func (s PullStepType) IsComplete() bool {
	return s == PullStepUpToDate || s == PullStepFastForward || s == PullStepConflict || s == PullStepError || s == PullStepSkipped
}

// HasFailed reports whether the step denotes a pull failure or merge conflict.
func (s PullStepType) HasFailed() bool {
	return s == PullStepError || s == PullStepConflict
}
