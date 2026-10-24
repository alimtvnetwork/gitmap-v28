package enums

// AuthMethodType represents the user-selected authentication method.
type AuthMethodType int

const (
	AuthMethodToken   AuthMethodType = 1
	AuthMethodBrowser AuthMethodType = 2
	AuthMethodSkip    AuthMethodType = 3
)

// IsValid reports whether the authentication method is recognized.
func (a AuthMethodType) IsValid() bool {
	return a >= AuthMethodToken && a <= AuthMethodSkip
}

// String returns the string representation.
func (a AuthMethodType) String() string {
	switch a {
	case AuthMethodToken:
		return "token"
	case AuthMethodBrowser:
		return "browser"
	case AuthMethodSkip:
		return "skip"
	default:
		return "unknown"
	}
}

// IsToken reports whether token auth was chosen.
func (a AuthMethodType) IsToken() bool {
	return a == AuthMethodToken
}

// IsBrowser reports whether browser auth was chosen.
func (a AuthMethodType) IsBrowser() bool {
	return a == AuthMethodBrowser
}

// IsSkip reports whether authentication was skipped.
func (a AuthMethodType) IsSkip() bool {
	return a == AuthMethodSkip
}
