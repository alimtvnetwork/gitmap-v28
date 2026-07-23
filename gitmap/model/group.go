// Package model defines the core data structures for gitmap.
package model

// Group represents a named collection of repositories.
type Group struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	CreatedAt   string `json:"createdAt"`
}

// GroupRepo links a group to a repository.
type GroupRepo struct {
	GroupID int64 `json:"groupId"`
	RepoID  int64 `json:"repoId"`
}
