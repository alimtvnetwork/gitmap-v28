package cmdpullerror

import (
	"time"
)

// PullErrorOptions configures the execution of the pull-error inspection command.
type PullErrorOptions struct {
	IsJSON   bool
	IsSSH    bool
	IsHelp   bool
	IsClear  bool
	RepoSlug string
	Limit    int
}

// PullErrorCardDetails holds formatted fields for displaying an error diagnostic card.
type PullErrorCardDetails struct {
	ErrorID        string    `json:"error_id"`
	RepoSlug       string    `json:"repo_slug"`
	RepoPath       string    `json:"repo_path"`
	NodeID         string    `json:"node_id"`
	NodeVersion    string    `json:"node_version"`
	ErrorType      string    `json:"error_type"`
	ErrorText      string    `json:"error_text"`
	StackTrace     string    `json:"stack_trace"`
	RemediationCmd string    `json:"remediation_cmd"`
	CreatedAt      time.Time `json:"created_at"`
}
