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

// GitProfileConfig holds all configured Git profiles and active defaults.
type GitProfileConfig struct {
	Profiles  []GitProfile `json:"profiles"`
	Active    string       `json:"active"`
	Default   string       `json:"default"`
	UpdatedAt time.Time    `json:"updatedAt"`
}
