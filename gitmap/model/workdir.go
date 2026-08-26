// Package model — workdir.go defines the registered work directory record.
package model

// WorkDir represents a registered parent workspace path.
type WorkDir struct {
	ID           int64  `json:"id"`
	AbsolutePath string `json:"absolutePath"`
	Label        string `json:"label"`
	IsDefault    bool   `json:"isDefault"`
	RepoCount    int    `json:"repoCount,omitempty"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}
