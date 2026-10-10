// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"time"
)

// NewCommandEntry represents a cataloged command entry.
type NewCommandEntry struct {
	Name        string   `json:"name"`
	Alias       string   `json:"alias,omitempty"`
	Version     string   `json:"version"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	Example     string   `json:"example"`
	Flags       []string `json:"flags,omitempty"`
}

// NewCommandsPayload encapsulates structured output for new commands.
type NewCommandsPayload struct {
	Timestamp        time.Time         `json:"timestamp"`
	TotalCommands    int               `json:"totalCommands"`
	FilteredCommands int               `json:"filteredCommands"`
	Limit            int               `json:"limit"`
	Category         string            `json:"category,omitempty"`
	Filter           string            `json:"filter,omitempty"`
	Commands         []NewCommandEntry `json:"commands"`
}

// NewCommandsOptions encapsulates parsed CLI flags.
type NewCommandsOptions struct {
	Limit    int
	Category string
	Filter   string
	IsJSON   bool
	IsHelp   bool
}

// DefaultNewCommandsOptions returns options with default values.
func DefaultNewCommandsOptions() NewCommandsOptions {
	return NewCommandsOptions{
		Limit:    100,
		Category: "",
		Filter:   "",
		IsJSON:   false,
		IsHelp:   false,
	}
}
