package pipelinedb

// PipelineCompactErrorRecord stores a compact error diagnostic entry with ok lines removed.
type PipelineCompactErrorRecord struct {
	CompactErrorLogId uint64 `json:"compactErrorLogId"`
	RunId             uint64 `json:"runId"`
	RepoSlug          string `json:"repoSlug"`
	WorkflowName      string `json:"workflowName"`
	StepName          string `json:"stepName"`
	ErrorText         string `json:"errorText"`
	CompactLogs       string `json:"compactLogs,omitempty"`
	FilteredOkCount   int    `json:"filteredOkCount"`
	Notes             string `json:"notes,omitempty"`
	Comments          string `json:"comments,omitempty"`
	CreatedAt         string `json:"createdAt"`
}
