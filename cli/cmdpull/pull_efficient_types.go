package cmdpull

import "github.com/alimtvnetwork/gitmap-v28/cli/model"

// InactiveRepoDetail holds metadata for a repository skipped due to inactivity.
type InactiveRepoDetail struct {
	RepoName string
	RepoPath string
	Reason   string
}

// EfficientPullPartition holds separated active scan records and inactive details.
type EfficientPullPartition struct {
	ActiveRecords []model.ScanRecord
	InactiveRepos []InactiveRepoDetail
}

// EfficientPullOptions contains runtime flags for efficient pull execution.
type EfficientPullOptions struct {
	IsTableMode  bool
	InvokedAlias string
	IsShortForm  bool
	UseSSH       bool
	UseHTTPS     bool
	Args         []string
}
