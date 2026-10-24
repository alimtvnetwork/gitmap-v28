package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// CommandType represents the type of command broadcasted to clients.
type CommandType string

const (
	// CommandPullAll is the command to tell clients to pull all local repositories.
	CommandPullAll CommandType = "PullAll"
)

// IsValid reports whether the command matches a known broadcast command.
func (c CommandType) IsValid() bool {
	switch c {
	case CommandPullAll:
		return true
	default:
		return false
	}
}

// String returns the string representation of the command.
func (c CommandType) String() string {
	return string(c)
}

// ParseCommandType parses a raw string into a validated CommandType.
func ParseCommandType(s string) (CommandType, *apperror.AppError) {
	val := CommandType(strings.TrimSpace(s))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseCommandType", "E1000", fmt.Sprintf("invalid cluster command: %q", s)).WithContext("value", s)
}

// AllCommandTypes returns all valid CommandType variants.
func AllCommandTypes() []CommandType {
	return []CommandType{CommandPullAll}
}

// IsPullAll reports whether the command is pull-all.
func (c CommandType) IsPullAll() bool {
	return c == CommandPullAll
}

// HasCommand reports whether the command is non-empty and valid.
func (c CommandType) HasCommand() bool {
	return c.IsValid()
}
