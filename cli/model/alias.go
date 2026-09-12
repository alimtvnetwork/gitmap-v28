// Package model defines the core data structures for gitmap.
package model

// Alias links a short name to a repository for quick access.
type Alias struct {
	ID        int64  `json:"id"`
	Alias     string `json:"alias"`
	RepoID    int64  `json:"repoId"`
	CreatedAt string `json:"createdAt"`
}
