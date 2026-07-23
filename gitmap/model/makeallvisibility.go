// Package model — makeallvisibility.go: row structs for the bulk
// wildcard visibility audit tables (MakeAllVisibilityRun +
// MakeAllVisibilityResult). Mirrors the schema in
// constants/constants_visibility_store.go.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §plan steps 17-18.
package model

// MakeAllVisibilityRunRecord is one row of MakeAllVisibilityRun.
// Counts are mutated in place by the apply loop and persisted via
// FinalizeMakeAllVisibilityRun on completion.
type MakeAllVisibilityRunRecord struct {
	ID               int64  `json:"id"`
	CommandKind      string `json:"commandKind"`
	TargetVisibility string `json:"targetVisibility"`
	Provider         string `json:"provider"`
	Owner            string `json:"owner"`
	TargetRaw        string `json:"targetRaw"`
	PatternList      string `json:"patternList"`
	YesFlag          bool   `json:"yesFlag"`
	VerboseFlag      bool   `json:"verboseFlag"`
	OwnerRepoTotal   int    `json:"ownerRepoTotal"`
	MatchedCount     int    `json:"matchedCount"`
	ExcludedCount    int    `json:"excludedCount"`
	OkCount          int    `json:"okCount"`
	SkippedCount     int    `json:"skippedCount"`
	FailedCount      int    `json:"failedCount"`
	ExitCode         int    `json:"exitCode"`
	StartedAt        string `json:"startedAt"`
	FinishedAt       string `json:"finishedAt,omitempty"`
}

// MakeAllVisibilityResultRecord is one row of MakeAllVisibilityResult.
// Status flows Pending → (Ok | Skipped | Failed | Excluded) — see
// constants.ResultStatus*. PrevVisibility is captured during the read
// phase; NewVisibility during the verify phase.
type MakeAllVisibilityResultRecord struct {
	ID             int64  `json:"id"`
	RunID          int64  `json:"runId"`
	RepoName       string `json:"repoName"`
	MatchedPattern string `json:"matchedPattern"`
	Status         string `json:"status"`
	PrevVisibility string `json:"prevVisibility,omitempty"`
	NewVisibility  string `json:"newVisibility,omitempty"`
	FailureMessage string `json:"failureMessage,omitempty"`
	StartedAt      string `json:"startedAt,omitempty"`
	FinishedAt     string `json:"finishedAt,omitempty"`
	DurationMs     int64  `json:"durationMs"`
}
