// Package model — prompt_architect.go defines Prompt Architect version metadata.
package model

// PromptArchitectMetadata holds prompt-architect configuration and version status in version.json.
type PromptArchitectMetadata struct {
	Version     string `json:"version"`
	InstalledAt string `json:"installed_at"`
	Status      string `json:"status"`
	Channel     string `json:"channel,omitempty"`
}
