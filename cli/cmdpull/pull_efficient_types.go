package cmdpull

import "github.com/alimtvnetwork/gitmap-v28/cli/model"

// InactiveRepoDetail holds metadata for a repository skipped due to inactivity.
type InactiveRepoDetail struct {
	RepoName string `json:"repoName"`
	RepoPath string `json:"repoPath"`
	Reason   string `json:"reason"`
}

// EfficientPullPartition holds separated active scan records and inactive details.
type EfficientPullPartition struct {
	ActiveRecords []model.ScanRecord
	InactiveRepos []InactiveRepoDetail
}

// EfficientPullOptions contains runtime flags for efficient pull execution.
type EfficientPullOptions struct {
	IsTableMode  bool
	IsJSON       bool
	InvokedAlias string
	IsShortForm  bool
	UseSSH       bool
	UseHTTPS     bool
	TargetSSH    string
	Args         []string
}

// PullRepoSummaryItem holds the outcome of pulling an individual repository.
type PullRepoSummaryItem struct {
	RepoName string `json:"repoName"`
	Status   string `json:"status"`
	Changes  string `json:"changes"`
}

// PullEfficientSummary represents the serialized JSON summary of an efficient pull operation.
type PullEfficientSummary struct {
	Node          string                `json:"node,omitempty"`
	Total         int                   `json:"total"`
	ActiveCount   int                   `json:"activeCount"`
	InactiveCount int                   `json:"inactiveCount"`
	States        []PullRepoSummaryItem `json:"states"`
	Inactive      []InactiveRepoDetail  `json:"inactive"`
	DurationMs    int64                 `json:"durationMs"`
}
