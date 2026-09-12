// Package model — history.go defines the command history audit record.
package model

// CommandHistoryRecord represents a single CLI command execution.
type CommandHistoryRecord struct {
	ID         int64  `json:"id"`
	Command    string `json:"command"`
	Alias      string `json:"alias,omitempty"`
	Args       string `json:"args,omitempty"`
	Flags      string `json:"flags,omitempty"`
	StartedAt  string `json:"startedAt"`
	FinishedAt string `json:"finishedAt,omitempty"`
	DurationMs int64  `json:"durationMs"`
	ExitCode   int    `json:"exitCode"`
	Summary    string `json:"summary,omitempty"`
	RepoCount  int    `json:"repoCount"`
	CreatedAt  string `json:"createdAt,omitempty"`
}
