// Package model defines the core data structures for gitmap.
package model

// ReleaseRecord holds release metadata stored in the database.
// v17: RepoId FK to Repo.RepoId (see spec/04-generic-cli/24-release-repo-relationship.md).
type ReleaseRecord struct {
	ID           int64  `json:"id"`
	RepoID       int64  `json:"repoId"`
	Version      string `json:"version"`
	Tag          string `json:"tag"`
	Branch       string `json:"branch"`
	SourceBranch string `json:"sourceBranch"`
	CommitSha    string `json:"commitSha"`
	Changelog    string `json:"changelog"`
	Notes        string `json:"notes"`
	IsDraft      bool   `json:"isDraft"`
	IsPreRelease bool   `json:"isPreRelease"`
	IsLatest     bool   `json:"isLatest"`
	Source       string `json:"source"`
	CreatedAt    string `json:"createdAt"`
}

// Release source values.
const (
	SourceRelease = "release"
	SourceImport  = "import"
	SourceRepo    = "repo"
	SourceTag     = "tag"
)
