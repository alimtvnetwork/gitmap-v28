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
