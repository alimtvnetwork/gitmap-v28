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

const (
	DefaultBotProfileID  = "autobot"
	DefaultBotName       = "gitmap automation"
	DefaultLegacyBotName = "AI automation bot"
	DefaultBotEmail      = "local@local"
	DefaultBotRole       = "automation agent"
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
	HasName      bool   `json:"hasName"`
	HasEmail     bool   `json:"hasEmail"`
}

// GitIdentityOptions configures Git identity setting and unsetting operations.
type GitIdentityOptions struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	IsGlobal      bool   `json:"isGlobal"`
	IsUnsetEmail  bool   `json:"isUnsetEmail"`
	HasEmptyEmail bool   `json:"hasEmptyEmail"`
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
