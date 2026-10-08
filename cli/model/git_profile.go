package model

import "time"

// GitProfile represents an authenticated Git provider account or organization.
type GitProfile struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Provider   string    `json:"provider"` // "github", "gitlab"
	Type       string    `json:"type"`     // "user", "organization"
	Email      string    `json:"email,omitempty"`
	AuthMethod string    `json:"authMethod"` // "gh-cli", "ssh", "token", "glab-cli"
	IsDefault  bool      `json:"isDefault"`
	UsageCount int       `json:"usageCount"`
	LastUsedAt time.Time `json:"lastUsedAt"`
}

// ProjectBinding associates a repository folder or slug with a Git profile alias.
type ProjectBinding struct {
	RepoSlug     string    `json:"repoSlug"`
	ProfileAlias string    `json:"profileAlias"`
	BoundAt      time.Time `json:"boundAt"`
}

// GitProfileConfig holds all configured Git profiles, bindings, and active defaults.
type GitProfileConfig struct {
	Profiles        []GitProfile              `json:"profiles"`
	Active          string                    `json:"active"`
	Default         string                    `json:"default"`
	ProjectBindings map[string]ProjectBinding `json:"projectBindings,omitempty"`
	UpdatedAt       time.Time                 `json:"updatedAt"`
}
