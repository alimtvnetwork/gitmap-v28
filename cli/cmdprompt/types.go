// Package cmdprompt — types.go centralizes all prompt structs, options, and Result aliases.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// PromptTemplate represents a structured AI prompt with metadata and body.
type PromptTemplate struct {
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Variables   []string `json:"variables"`
	Body        string   `json:"body"`
	FilePath    string   `json:"filePath,omitempty"`
}

// PromptInstallOptions parses CLI flags and targets for ct install-prompts.
type PromptInstallOptions struct {
	Targets  []string
	Exclude  string
	IsDryRun bool
	IsAll    bool
	Action   string
}

// PromptStatusTableLayout calculates column dimensions for prompt status table.
type PromptStatusTableLayout struct {
	MaxRepo    int
	MaxStatus  int
	MaxVersion int
	MaxDate    int
}

// PromptTargetSliceResult wraps a slice of repository directory targets.
type PromptTargetSliceResult = result.ResultSlice[string]
