// Package cmdclone — clone_auth_types.go defines types and errors for clone authentication.
package cmdclone

import (
	"errors"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone/enums"
)

// ErrAuthSkipped is returned when the user chooses to skip authenticating a repository.
var ErrAuthSkipped = errors.New("authentication skipped by user")

// AuthMethod represents the user-selected authentication method.
type AuthMethod = enums.AuthMethodType

const (
	AuthMethodToken   = enums.AuthMethodToken
	AuthMethodBrowser = enums.AuthMethodBrowser
	AuthMethodSkip    = enums.AuthMethodSkip
)

// AuthPromptResult encapsulates the outcome of an interactive authentication prompt.
type AuthPromptResult struct {
	Token    string
	IsReused bool
	Method   AuthMethod
}

// IsSSHCloneURL checks if rawURL is an SSH clone URL (shorthand git@ or ssh://).
func IsSSHCloneURL(url string) bool {
	return isSSHCloneURL(url)
}
