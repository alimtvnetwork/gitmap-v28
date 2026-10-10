// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// NodePendingRepoItem represents a single repository's pending state on a remote node.
type NodePendingRepoItem struct {
	RepoName           string   `json:"repo_name"`
	RelativePath       string   `json:"relative_path"`
	Branch             string   `json:"branch"`
	IsClean            bool     `json:"is_clean"`
	HasChanges         bool     `json:"has_changes"`
	DirtyFileCount     int      `json:"dirty_file_count"`
	HasUnpushedCommits bool     `json:"has_unpushed_commits"`
	UnpushedCount      int      `json:"unpushed_count"`
	PriorityScore      int      `json:"priority_score"`
	ChangedFiles       []string `json:"changed_files,omitempty"`
}

// NodePendingCommitsResult represents the pending commits discovery outcome for a single node.
type NodePendingCommitsResult struct {
	NodeAlias       string                `json:"node_alias"`
	Host            string                `json:"host"`
	OS              string                `json:"os"`
	IsNodeOnline    bool                  `json:"is_node_online"`
	IsSuccess       bool                  `json:"is_success"`
	Latency         time.Duration         `json:"latency"`
	LatencyMs       int64                 `json:"latency_ms"`
	PendingRepos    []NodePendingRepoItem `json:"pending_repos,omitempty"`
	TotalDirtyFiles int                   `json:"total_dirty_files"`
	TotalUnpushed   int                   `json:"total_unpushed"`
	ErrorMessage    string                `json:"error_message,omitempty"`
}

// NodesPendingCommitsAggregated is the top-level envelope for gitmap nodes pc --json.
type NodesPendingCommitsAggregated struct {
	Timestamp          string                     `json:"timestamp"`
	TotalNodesQueried  int                        `json:"total_nodes_queried"`
	OnlineNodesCount   int                        `json:"online_nodes_count"`
	TotalPendingRepos  int                        `json:"total_pending_repos"`
	TotalDirtyFiles    int                        `json:"total_dirty_files"`
	TotalUnpushedCount int                        `json:"total_unpushed_count"`
	SortMode           string                     `json:"sort_mode"`
	Results            []NodePendingCommitsResult `json:"results"`
}

// CommitActionKind defines semantic action types.
type CommitActionKind string

const (
	CommitActionStandard CommitActionKind = "commits"
	CommitActionFeature  CommitActionKind = "cpf"
	CommitActionBug      CommitActionKind = "cpb"
	CommitActionRelease  CommitActionKind = "cpr"
	CommitActionFix      CommitActionKind = "commit-fix"
)

// RemoteRepoCommitOutcome captures the commit result for an individual repository on a node.
type RemoteRepoCommitOutcome struct {
	RepoName     string `json:"repo_name"`
	Branch       string `json:"branch"`
	CommitHash   string `json:"commit_hash,omitempty"`
	CommitMsg    string `json:"commit_msg"`
	IsSuccess    bool   `json:"is_success"`
	IsPushed     bool   `json:"is_pushed"`
	HasChanges   bool   `json:"has_changes"`
	FilesChanged int    `json:"files_changed"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// NodeCommitResult captures the fleet commit dispatch outcome for a single node.
type NodeCommitResult struct {
	NodeAlias      string                    `json:"node_alias"`
	Host           string                    `json:"host"`
	OS             string                    `json:"os"`
	Action         string                    `json:"action"`
	IsNodeOnline   bool                      `json:"is_node_online"`
	IsSuccess      bool                      `json:"is_success"`
	IsDryRun       bool                      `json:"is_dry_run"`
	Latency        time.Duration             `json:"latency"`
	LatencyMs      int64                     `json:"latency_ms"`
	Repositories   []RemoteRepoCommitOutcome `json:"repositories,omitempty"`
	CommittedCount int                       `json:"committed_count"`
	ErrorMessage   string                    `json:"error_message,omitempty"`
}

// NodeCommitExecutionResult is an alias for NodeCommitResult.
type NodeCommitExecutionResult = NodeCommitResult

// NodesCommitsAggregated is the top-level envelope for gitmap nodes commits/cpf/... --json.
type NodesCommitsAggregated struct {
	Timestamp         string             `json:"timestamp"`
	Action            string             `json:"action"`
	TargetRepo        string             `json:"target_repo"`
	CommitMessage     string             `json:"commit_message"`
	IsDryRun          bool               `json:"is_dry_run"`
	TotalNodesTarget  int                `json:"total_nodes_target"`
	SuccessfulNodes   int                `json:"successful_nodes"`
	TotalReposUpdated int                `json:"total_repos_updated"`
	Results           []NodeCommitResult `json:"results"`
}

// NodesCommitSummary is an alias for NodesCommitsAggregated.
type NodesCommitSummary = NodesCommitsAggregated

// NodesPendingCommitsOptions configures fleet-wide pending commits discovery.
type NodesPendingCommitsOptions struct {
	FilterOpts  NodeFilterOptions
	SortMode    string
	IsDetail    bool
	IsJSON      bool
	IsDirtyOnly bool
}

// NodesCommitsOptions configures fleet-wide commit delegation.
type NodesCommitsOptions struct {
	FilterOpts    NodeFilterOptions
	Action        string
	TargetScope   string
	CommitMessage string
	IsDryRun      bool
	IsPushEnabled bool
	IsJSON        bool
}

// remoteSSHExecutor defines an interface for executing remote shell commands over SSH.
type remoteSSHExecutor interface {
	Execute(conn db.SSHConnection, command string, shell string) (string, error)
}

type sshProductionExecutor struct{}

func (e *sshProductionExecutor) Execute(conn db.SSHConnection, command string, shell string) (string, error) {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return "", errConnect
	}

	defer client.Close()

	out, errRun := secrets.RunCommand(client, command, shell)
	if errRun != nil && shell == "bash" && isBashMissingError(out, errRun) {
		out, errRun = secrets.RunCommand(client, command, "sh")
	}

	return out, errRun
}

var currentSSHExecutor remoteSSHExecutor = &sshProductionExecutor{}
