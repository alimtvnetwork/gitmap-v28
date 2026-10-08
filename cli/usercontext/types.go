package usercontext

// AuthStatus defines the authorization state for GitHub CLI and provider tools.
type AuthStatus string

const (
	StatusAuthorized      AuthStatus = "authorized"
	StatusUnauthenticated AuthStatus = "unauthenticated"
	StatusToolMissing     AuthStatus = "tool_missing"
)

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
