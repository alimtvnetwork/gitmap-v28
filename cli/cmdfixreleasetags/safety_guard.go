// Package cmdfixreleasetags provides safety guard evaluation to prevent
// accidental deletion of active, latest healthy, or in-flight releases.
package cmdfixreleasetags

import (
	"strings"
	"time"
)

// NormalizeVersionTag strips whitespace and optional leading 'v' or 'V' prefixes.
func NormalizeVersionTag(versionStr string) string {
	trimmed := strings.TrimSpace(versionStr)
	trimmed = strings.TrimPrefix(trimmed, "v")
	trimmed = strings.TrimPrefix(trimmed, "V")

	return trimmed
}

// IsActiveVersion checks if a tag matches the currently running binary version.
func IsActiveVersion(tag, activeVersion string) bool {
	normTag := NormalizeVersionTag(tag)
	normActive := NormalizeVersionTag(activeVersion)
	if normTag == "" || normActive == "" {
		return false
	}

	return normTag == normActive
}

// IsLatestHealthyTag checks if a tag matches the latest verified healthy release.
func IsLatestHealthyTag(tag, latestHealthyTag string) bool {
	normTag := NormalizeVersionTag(tag)
	normLatest := NormalizeVersionTag(latestHealthyTag)
	if normTag == "" || normLatest == "" {
		return false
	}

	return normTag == normLatest
}

// isRunInProgress checks if an individual CI run is active.
func isRunInProgress(run CIWorkflowRunInfo) bool {
	if run.IsInProgress {
		return true
	}

	status := strings.ToLower(run.Status)
	if status == "in_progress" || status == "queued" || status == "waiting" {
		return true
	}

	return false
}

// isRunWithinGrace checks if a run was created within the grace window.
func isRunWithinGrace(run CIWorkflowRunInfo, now time.Time, graceWindow time.Duration) bool {
	if run.CreatedAt.IsZero() || graceWindow <= 0 {
		return false
	}

	diff := now.Sub(run.CreatedAt)
	if diff >= 0 && diff < graceWindow {
		return true
	}

	return false
}

// HasActiveWorkflowRun evaluates if any CI run is active or within grace period.
func HasActiveWorkflowRun(runs []CIWorkflowRunInfo, now time.Time, graceWindow time.Duration) bool {
	for _, run := range runs {
		if isRunInProgress(run) {
			return true
		}

		if isRunWithinGrace(run, now, graceWindow) {
			return true
		}
	}

	return false
}

// EvaluateTagSafety checks all three safety invariants and returns status and eligibility.
func EvaluateTagSafety(tag string, activeVersion string, latestHealthyTag string, runs []CIWorkflowRunInfo, now time.Time, graceWindow time.Duration) (ProtectionStatus, bool) {
	if IsActiveVersion(tag, activeVersion) {
		return StatusProtectedActive, false
	}

	if IsLatestHealthyTag(tag, latestHealthyTag) {
		return StatusProtectedLatest, false
	}

	if HasActiveWorkflowRun(runs, now, graceWindow) {
		return StatusProtectedGrace, false
	}

	return StatusEligible, true
}

// EvaluateSafety evaluates safety for a given audit record and mutates protection flags.
func EvaluateSafety(record *ReleaseTagAuditRecord, currentVersion string, latestHealthyTag string, graceDuration time.Duration) ProtectionStatus {
	if record == nil {
		return StatusEligible
	}

	now := time.Now()
	status, isEligible := EvaluateTagSafety(record.Tag, currentVersion, latestHealthyTag, record.WorkflowRuns, now, graceDuration)

	record.ProtectionStatus = status
	record.IsProtected = !isEligible
	if record.IsProtected {
		record.IsEligibleForDeletion = false
	}

	return status
}
