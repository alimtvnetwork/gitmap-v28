// Package cmdclone — clone_auth_types.go defines types and errors for clone authentication.
package cmdclone

import "errors"

// ErrAuthSkipped is returned when the user chooses to skip authenticating a repository.
var ErrAuthSkipped = errors.New("authentication skipped by user")

// AuthMethod represents the user-selected authentication method.
type AuthMethod int

const (
	AuthMethodToken AuthMethod = 1
	AuthMethodBrowser AuthMethod = 2
	AuthMethodSkip AuthMethod = 3
)

// AuthPromptResult encapsulates the outcome of an interactive authentication prompt.
type AuthPromptResult struct {
	Token      string
	IsReused   bool
	Method     AuthMethod
}
