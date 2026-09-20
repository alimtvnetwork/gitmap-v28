// Package prdb provides SQLite split-database storage and operations for the GitMap PR subsystem (Spec 129).
package prdb

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PullRequestRecord encapsulates a pull request database record (Spec 129).
type PullRequestRecord struct {
	PullRequestId  int64  `json:"pullRequestId"`
	PrNumber       int    `json:"prNumber"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	SourceBranch   string `json:"sourceBranch"`
	TargetBranch   string `json:"targetBranch"`
	Status         string `json:"status"`
	MergeCommitSha string `json:"mergeCommitSha,omitempty"`
	CreatedAt      int64  `json:"createdAt"`
	MergedAt       int64  `json:"mergedAt,omitempty"`
	ClosedAt       int64  `json:"closedAt,omitempty"`
	Notes          string `json:"notes,omitempty"`
	Comments       string `json:"comments,omitempty"`
	UpdatedAt      int64  `json:"updatedAt"`
}

// PrReleaseRecord encapsulates a release tag database record linked to a pull request.
type PrReleaseRecord struct {
	PrReleaseId   int64  `json:"prReleaseId"`
	PullRequestId int64  `json:"pullRequestId"`
	ReleaseTag    string `json:"releaseTag"`
	CommitSha     string `json:"commitSha"`
	Notes         string `json:"notes,omitempty"`
	Comments      string `json:"comments,omitempty"`
	CreatedAt     int64  `json:"createdAt"`
}

// PrBranchRecord encapsulates a feature or pull request branch database record.
type PrBranchRecord struct {
	PrBranchId int64  `json:"prBranchId"`
	BranchName string `json:"branchName"`
	BranchType string `json:"branchType"`
	IsMerged   bool   `json:"isMerged"`
	IsDeleted  bool   `json:"isDeleted"`
	CreatedAt  int64  `json:"createdAt"`
	MergedAt   int64  `json:"mergedAt,omitempty"`
	DeletedAt  int64  `json:"deletedAt,omitempty"`
	Notes      string `json:"notes,omitempty"`
	Comments   string `json:"comments,omitempty"`
	UpdatedAt  int64  `json:"updatedAt"`
}

// Monadic Result wrappers for PR database operations.
type (
	PullRequestResult  = result.Result[*PullRequestRecord]
	PrBranchListResult = result.Result[[]PrBranchRecord]
	PrReleaseResult    = result.Result[*PrReleaseRecord]
)
