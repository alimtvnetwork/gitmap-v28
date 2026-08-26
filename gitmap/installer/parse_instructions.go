// Package installer — parse_instructions.go parses install instructions.
package installer

import (
	"encoding/json"
	"strings"
)

// Instruction represents a step or shell command in the installation process.
type Instruction struct {
	Command     string `json:"command"`
	Description string `json:"description,omitempty"`
}

// ParseInstructions deserializes instructions from json or parses newline-separated commands.
func ParseInstructions(raw string) []Instruction {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}

	var parsed []Instruction
	if err := json.Unmarshal([]byte(trimmed), &parsed); err == nil && len(parsed) > 0 {
		return parsed
	}

	lines := strings.Split(trimmed, "\n")
	for _, l := range lines {
		cmd := strings.TrimSpace(l)
		if cmd != "" {
			parsed = append(parsed, Instruction{Command: cmd})
		}
	}
	return parsed
}
