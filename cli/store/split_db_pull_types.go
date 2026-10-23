package store

import (
	"database/sql"
)

// PullSplitDB manages the dedicated SQLite split database for git pull telemetry.
type PullSplitDB struct {
	conn *sql.DB
	Path string
}

// PullRunRecord represents a master pull execution session in PullRun table.
type PullRunRecord struct {
	PullRunId     int64
	CommandType   string
	WorkingDir    string
	TotalRepos    int
	PulledRepos   int
	SkippedRepos  int
	SuccessCount  int
	FailedCount   int
	IsEfficient   bool
	DurationMs    int64
	GitMapVersion string
	Notes         string
	Comments      string
	CreatedAt     string
}

// PullRepoRunRecord represents a repository result in PullRepoRun table.
type PullRepoRunRecord struct {
	PullRepoRunId     int64
	PullRunId         int64
	RepoPath          string
	RepoName          string
	PullStatus        string
	FilesChanged      int
	LastCommitSha     string
	PreviousCommitSha string
	CommitMessage     string
	CommitAuthor      string
	IsActive          bool
	HasChanges        bool
	DurationMs        int64
	ErrorMessage      string
	Notes             string
	Comments          string
	CreatedAt         string
}

// RepoInactivityStatus encapsulates evaluation result for repository activity.
type RepoInactivityStatus struct {
	RepoPath                 string
	RepoName                 string
	IsInactive               bool
	ConsecutiveZeroRunsCount int
	TotalRunsInspected       int
	LatestPullTimestamp      string
	Reason                   string
}

// PullTaskRecord represents an execution session row in PullTask.
type PullTaskRecord struct {
	PullTaskId    int64  `json:"pullTaskId"`
	TaskUuid      string `json:"taskUuid"`
	CommandType   string `json:"commandType"`
	TargetNode    string `json:"targetNode"`
	MachineId     string `json:"machineId"`
	Status        string `json:"status"`
	Pid           int    `json:"pid"`
	TotalRepos    int    `json:"totalRepos"`
	PulledRepos   int    `json:"pulledRepos"`
	EtaSeconds    int    `json:"etaSeconds"`
	IsActive      bool   `json:"isActive"`
	IsCompleted   bool   `json:"isCompleted"`
	IsFailed      bool   `json:"isFailed"`
	HasActiveLock bool   `json:"hasActiveLock"`
	StartedAt     string `json:"startedAt"`
	CompletedAt   string `json:"completedAt"`
	ErrorMessage  string `json:"errorMessage"`
	HeartbeatAt   string `json:"heartbeatAt"`
	CreatedAt     string `json:"createdAt"`
}

// PullItemHistoryRecord represents a repository-level progress row in PullItemHistory.
// Strictly stores NO file contents or diffs (only metadata: RepoName, RepoPath, CommitHash, Title, CreatedAt, PulledAt).
type PullItemHistoryRecord struct {
	PullItemHistoryId int64  `json:"pullItemHistoryId"`
	PullTaskId        int64  `json:"pullTaskId"`
	RepoName          string `json:"repoName"`
	RepoPath          string `json:"repoPath"`
	CommitHash        string `json:"commitHash"`
	Title             string `json:"title"`
	PullStatus        string `json:"pullStatus"`
	FilesChanged      int    `json:"filesChanged"`
	DurationMs        int64  `json:"durationMs"`
	IsSuccess         bool   `json:"isSuccess"`
	HasChanges        bool   `json:"hasChanges"`
	IsSkipped         bool   `json:"isSkipped"`
	ErrorMessage      string `json:"errorMessage"`
	PulledAt          string `json:"pulledAt"`
	CreatedAt         string `json:"createdAt"`
}

// PullEtaHistoryRecord models a historical duration entry in PullEtaHistory.
type PullEtaHistoryRecord struct {
	PullEtaHistoryId int64   `json:"pullEtaHistoryId"`
	CommandType      string  `json:"commandType"`
	TargetNode       string  `json:"targetNode"`
	TotalRepos       int     `json:"totalRepos"`
	DurationSeconds  float64 `json:"durationSeconds"`
	CreatedAt        string  `json:"createdAt"`
}
