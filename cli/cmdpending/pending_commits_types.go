// Package cmd — pending_commits_types.go defines telemetry payloads and data structures
// for gitmap pending-commits and gitmap sends suites.
package cmdpending

import "time"

// UncommittedChangesSummary aggregates metrics for uncommitted working tree changes.
type UncommittedChangesSummary struct {
	UntrackedFilesCount int `json:"untrackedFilesCount"`
	ModifiedFilesCount  int `json:"modifiedFilesCount"`
	StagedFilesCount    int `json:"stagedFilesCount"`
	TotalUncommitted    int `json:"totalUncommitted"`
}

// UnpushedCommitSummary aggregates metrics for local commits ahead of upstream branch.
type UnpushedCommitSummary struct {
	UnpushedCommitsCount int  `json:"unpushedCommitsCount"`
	HasUpstream          bool `json:"hasUpstream"`
}

// RemediationOption represents an actionable fix for a dirty repository.
type RemediationOption struct {
	OptionNumber int    `json:"optionNumber"`
	Label        string `json:"label"`
	Command      string `json:"command"`
}

// RepoPendingCommitRecord represents the pending changes status of an individual repository.
type RepoPendingCommitRecord struct {
	RepoName             string              `json:"repoName"`
	RelativePath         string              `json:"relativePath"`
	CurrentBranch        string              `json:"currentBranch"`
	Version              string              `json:"version,omitempty"`
	ShortVersionBranch   string              `json:"shortVersionBranch"`
	IsDirty              bool                `json:"isDirty"`
	IsClean              bool                `json:"isClean"`
	HasUncommitted       bool                `json:"hasUncommitted"`
	HasUnpushed          bool                `json:"hasUnpushed"`
	HasUpstream          bool                `json:"hasUpstream"`
	TotalUncommitted     int                 `json:"totalUncommitted"`
	UntrackedFilesCount  int                 `json:"untrackedFilesCount"`
	ModifiedFilesCount   int                 `json:"modifiedFilesCount"`
	StagedFilesCount     int                 `json:"stagedFilesCount"`
	UnpushedCommitsCount int                 `json:"unpushedCommitsCount"`
	PendingFiles         []string            `json:"pendingFiles,omitempty"`
	UnpushedCommitSHAs   []string            `json:"unpushedCommitShas,omitempty"`
	RemediationOptions   []RemediationOption `json:"remediationOptions,omitempty"`
}

// PendingCommitsPayload represents the top-level JSON telemetry for pending-commits.
type PendingCommitsPayload struct {
	Timestamp             time.Time                 `json:"timestamp"`
	TotalReposScanned     int                       `json:"totalReposScanned"`
	TotalDirtyRepos       int                       `json:"totalDirtyRepos"`
	TotalUncommittedFiles int                       `json:"totalUncommittedFiles"`
	TotalUnpushedCommits  int                       `json:"totalUnpushedCommits"`
	SortMode              string                    `json:"sortMode"`
	DetailMode            string                    `json:"detailMode"`
	Repositories          []RepoPendingCommitRecord `json:"repositories"`
}

// PendingCommitsOptions encapsulates CLI execution parameters for pending-commits.
type PendingCommitsOptions struct {
	SortMode    string
	DetailMode  string
	IsSSH       bool
	IsJSON      bool
	IsDirtyOnly bool
	IsAll       bool
	TargetRepo  string
	IsNoCache   bool
	IsRefresh   bool
}

// NodePendingCommitsRecord represents pending commits for a specific cluster SSH node.
type NodePendingCommitsRecord struct {
	NodeAlias    string                 `json:"nodeAlias"`
	Host         string                 `json:"host"`
	OSType       string                 `json:"osType"`
	IsOnline     bool                   `json:"isOnline"`
	IsSuccess    bool                   `json:"isSuccess"`
	LatencyMs    int64                  `json:"latencyMs"`
	ErrorMessage string                 `json:"errorMessage,omitempty"`
	Payload      *PendingCommitsPayload `json:"payload,omitempty"`
}

// NodesPendingCommitsPayload represents aggregated fleet pending commits across SSH cluster nodes.
type NodesPendingCommitsPayload struct {
	Timestamp    time.Time                  `json:"timestamp"`
	TotalNodes   int                        `json:"totalNodes"`
	OnlineNodes  int                        `json:"onlineNodes"`
	FleetResults []NodePendingCommitsRecord `json:"fleetResults"`
}

// RepoSendResultRecord details the commit and push operation for a single repository.
type RepoSendResultRecord struct {
	RepoName     string `json:"repoName"`
	RelativePath string `json:"relativePath"`
	Status       string `json:"status"` // committed | pushed | clean-skipped | dry-run-simulated | failed | already-clean
	Branch       string `json:"branch"`
	HeadSHA      string `json:"headSha,omitempty"`
	FilesStaged  int    `json:"filesStaged"`
	IsSuccess    bool   `json:"isSuccess"`
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// SendsExecutionPayload represents the top-level JSON telemetry for gitmap sends.
type SendsExecutionPayload struct {
	Timestamp          time.Time              `json:"timestamp"`
	Verb               string                 `json:"verb"`
	PrefixApplied      string                 `json:"prefixApplied"`
	RawMessage         string                 `json:"rawMessage"`
	FinalCommitMessage string                 `json:"finalCommitMessage"`
	TargetScope        string                 `json:"targetScope"`
	IsDryRun           bool                   `json:"isDryRun"`
	IsPushed           bool                   `json:"isPushed"`
	TotalProcessed     int                    `json:"totalProcessed"`
	TotalCommitted     int                    `json:"totalCommitted"`
	TotalSkipped       int                    `json:"totalSkipped"`
	Results            []RepoSendResultRecord `json:"results"`
}

// SendsOptions encapsulates CLI execution parameters for gitmap sends.
type SendsOptions struct {
	Verb       string
	Target     string
	RawMessage string
	IsDryRun   bool
	IsPushed   bool
	IsJSON     bool
	IsVerbose  bool
}
