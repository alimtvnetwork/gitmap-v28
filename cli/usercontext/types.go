package usercontext

import "github.com/alimtvnetwork/gitmap-v28/cli/usercontext/enums"

// AuthStatusType defines the authorization state for GitHub CLI and provider tools.
// AuthStatusType aliases the strongly-typed enum from the enums package.
type AuthStatusType = enums.AuthStatusType

const (
	StatusAuthorized      = enums.StatusAuthorized
	StatusUnauthenticated = enums.StatusUnauthenticated
	StatusToolMissing     = enums.StatusToolMissing
)

// AuthStatus is an alias for AuthStatusType for backwards compatibility.
type AuthStatus = enums.AuthStatusType

// GhAuthInfo represents GitHub CLI authorization readout details.
type GhAuthInfo struct {
	Username string     `json:"username"`
	Status   AuthStatus `json:"status"`
	Source   string     `json:"source"`
	Scope    string     `json:"scope,omitempty"`
}

// GitIdentity captures Git author name, email, and configured status.
type GitIdentity struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	IsConfigured bool   `json:"isConfigured"`
}

// UserSummary consolidates user identity and context for display and inspection.
type UserSummary struct {
	GhAuth        GhAuthInfo  `json:"ghAuth"`
	LocalGitUser  GitIdentity `json:"localGitUser"`
	GlobalGitUser GitIdentity `json:"globalGitUser"`
	ActiveProfile string      `json:"activeProfile"`
	BoundProfile  string      `json:"boundProfile"`
	IsInsideRepo  bool        `json:"isInsideRepo"`
}

// IsAuthorized reports whether the GitHub authorization status is active.
func (a GhAuthInfo) IsAuthorized() bool {
	return a.Status == StatusAuthorized
}
