// Package model — installer.go defines the installer script and version records.
package model

// OSScript represents script instructions and runtime engine for a specific operating system.
type OSScript struct {
	Runtime      string `json:"runtime,omitempty"`   // e.g., "bash", "powershell", "python", "node"
	Instructions string `json:"instructions"`        // script commands or JSON payload
	VerifyCmd    string `json:"verifyCmd,omitempty"` // post-installation health check command
}

// InstallerScript represents an installer script record stored in SQLite.
type InstallerScript struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	Slug         string              `json:"slug"`
	Description  string              `json:"description"`
	TargetOS     string              `json:"targetOs"`
	Version      string              `json:"version"`
	OrderMode    string              `json:"orderMode,omitempty"` // "unix-first", "os-first", "os-only", "fallback"
	Instructions string              `json:"instructions"`
	Scripts      map[string]OSScript `json:"scripts,omitempty"` // keyed by OS ("ubuntu", "win", "arch", "unix", etc.)
	CreatedAt    string              `json:"createdAt"`
	UpdatedAt    string              `json:"updatedAt"`
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
