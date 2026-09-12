// Package model — version_history.go defines the repo version history record.
package model

// RepoVersionHistoryRecord represents a single version transition for a repo.
type RepoVersionHistoryRecord struct {
	ID             int64  `json:"id"`
	RepoID         int64  `json:"repoId"`
	FromVersionTag string `json:"fromVersionTag"`
	FromVersionNum int    `json:"fromVersionNum"`
	ToVersionTag   string `json:"toVersionTag"`
	ToVersionNum   int    `json:"toVersionNum"`
	FlattenedPath  string `json:"flattenedPath,omitempty"`
	CreatedAt      string `json:"createdAt,omitempty"`
}
