// Package model — installer.go defines the installer script and version records.
package model

// InstallerScript represents an installer script record stored in SQLite.
type InstallerScript struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	TargetOS     string `json:"targetOs"`
	Version      string `json:"version"`
	Instructions string `json:"instructions"`
	CreatedAt    string `json:"createdAt"`
	UpdatedAt    string `json:"updatedAt"`
}

// InstallerVersion represents a versioned snapshot of an installer script.
type InstallerVersion struct {
	ID           int64  `json:"id"`
	ScriptID     int64  `json:"scriptId"`
	Slug         string `json:"slug"`
	Version      string `json:"version"`
	TargetOS     string `json:"targetOs"`
	Instructions string `json:"instructions"`
	CreatedAt    string `json:"createdAt"`
}
