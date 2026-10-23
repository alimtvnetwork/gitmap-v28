package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CommitHistoryRecord models a historical commit entry with AI provenance.
type CommitHistoryRecord struct {
	CommitHistoryId    int64  `json:"commitHistoryId"`
	CommitSha          string `json:"commitSha"`
	ShortSha           string `json:"shortSha"`
	RepoName           string `json:"repoName"`
	RepoSlug           string `json:"repoSlug"`
	RepoRoot           string `json:"repoRoot"`
	BranchName         string `json:"branchName"`
	AuthorName         string `json:"authorName"`
	AuthorEmail        string `json:"authorEmail"`
	CommitMessage      string `json:"commitMessage"`
	CommittedAtUnix    int64  `json:"committedAtUnix"`
	AiModelName        string `json:"aiModelName"`
	AiToolName         string `json:"aiToolName"`
	AgentId            string `json:"agentId"`
	HeatmapSessionId   string `json:"heatmapSessionId"`
	HeatmapActionId    int64  `json:"heatmapActionId"`
	TaskSlug           string `json:"taskSlug"`
	SpecPath           string `json:"specPath"`
	PlansCount         int    `json:"plansCount"`
	SubtasksCount      int    `json:"subtasksCount"`
	IsAiCommitted      bool   `json:"isAiCommitted"`
	IsCurrentMachineAi bool   `json:"isCurrentMachineAi"`
	Notes              string `json:"notes"`
	FilesChanged       int    `json:"filesChanged"`
	LinesAdded         int    `json:"linesAdded"`
	LinesDeleted       int    `json:"linesDeleted"`
	HasPushed          bool   `json:"hasPushed"`
	CreatedAtUnix      int64  `json:"createdAtUnix"`
}

// CommitFileTouch models an individual file modified within a commit.
type CommitFileTouch struct {
	CommitFileTouchId int64  `json:"commitFileTouchId"`
	CommitSha         string `json:"commitSha"`
	FilePath          string `json:"filePath"`
	ChangeType        string `json:"changeType"`
	LinesAdded        int    `json:"linesAdded"`
	LinesDeleted      int    `json:"linesDeleted"`
	CreatedAtUnix     int64  `json:"createdAtUnix"`
}

// CommitHistoryFilter encapsulates filtering criteria for commit history queries.
type CommitHistoryFilter struct {
	Limit    int
	Offset   int
	RepoSlug string
	TaskSlug string
	AgentId  string
	IsAiOnly bool
}

// CommitHistoryQueryParam is an alias for CommitHistoryFilter.
type CommitHistoryQueryParam = CommitHistoryFilter

// Result wrapper type aliases.
type CommitHistoryRecordWrap = result.Wrap[CommitHistoryRecord]
type CommitHistorySliceWrap = result.WrapSlice[CommitHistoryRecord]
type CommitFileTouchSliceWrap = result.WrapSlice[CommitFileTouch]
