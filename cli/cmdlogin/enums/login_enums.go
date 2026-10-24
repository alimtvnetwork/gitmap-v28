package enums

import (
	"fmt"
	"strings"
)

// LoginAuthMethodType represents supported authentication mechanisms.
type LoginAuthMethodType string

const (
	LoginAuthMethodSSH     LoginAuthMethodType = "ssh"
	LoginAuthMethodToken   LoginAuthMethodType = "token"
	LoginAuthMethodBrowser LoginAuthMethodType = "browser"
	LoginAuthMethodStatus  LoginAuthMethodType = "status"
)

// IsValid reports whether the authentication method is recognized.
func (m LoginAuthMethodType) IsValid() bool {
	switch m {
	case LoginAuthMethodSSH, LoginAuthMethodToken, LoginAuthMethodBrowser, LoginAuthMethodStatus:
		return true
	default:
		return false
	}
}

// String returns the string representation of the auth method.
func (m LoginAuthMethodType) String() string {
	return string(m)
}

// IsSSH reports whether this auth method uses SSH keys.
func (m LoginAuthMethodType) IsSSH() bool {
	return m == LoginAuthMethodSSH
}

// IsToken reports whether this auth method uses a personal access token.
func (m LoginAuthMethodType) IsToken() bool {
	return m == LoginAuthMethodToken
}

// IsBrowser reports whether this auth method uses browser authentication.
func (m LoginAuthMethodType) IsBrowser() bool {
	return m == LoginAuthMethodBrowser
}

// IsStatus reports whether this auth method checks authentication status.
func (m LoginAuthMethodType) IsStatus() bool {
	return m == LoginAuthMethodStatus
}

func normalizeAuthMethod(s string) LoginAuthMethodType {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "browser-flow":
		return LoginAuthMethodBrowser
	case "pat":
		return LoginAuthMethodToken
	default:
		return LoginAuthMethodType(strings.ToLower(strings.TrimSpace(s)))
	}
}

// ParseLoginAuthMethod parses a raw string into a LoginAuthMethodType.
func ParseLoginAuthMethod(s string) (LoginAuthMethodType, error) {
	val := normalizeAuthMethod(s)
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid login auth method: %q", s)
}

// ParseLoginAuthMethodType is an alias for ParseLoginAuthMethod.
func ParseLoginAuthMethodType(s string) (LoginAuthMethodType, error) {
	return ParseLoginAuthMethod(s)
}

// AllLoginAuthMethods returns all supported authentication methods.
func AllLoginAuthMethods() []LoginAuthMethodType {
	return []LoginAuthMethodType{
		LoginAuthMethodSSH,
		LoginAuthMethodToken,
		LoginAuthMethodBrowser,
		LoginAuthMethodStatus,
	}
}

// AllLoginAuthMethodTypes is an alias for AllLoginAuthMethods.
func AllLoginAuthMethodTypes() []LoginAuthMethodType {
	return AllLoginAuthMethods()
}

// LoginInteractiveChoiceType represents user menu selections for interactive login.
type LoginInteractiveChoiceType string

const (
	LoginInteractiveChoiceSSH     LoginInteractiveChoiceType = "1"
	LoginInteractiveChoiceToken   LoginInteractiveChoiceType = "2"
	LoginInteractiveChoiceBrowser LoginInteractiveChoiceType = "3"
	LoginInteractiveChoiceStatus  LoginInteractiveChoiceType = "4"
)

// IsValid reports whether the interactive choice is recognized.
func (c LoginInteractiveChoiceType) IsValid() bool {
	switch c {
	case LoginInteractiveChoiceSSH, LoginInteractiveChoiceToken, LoginInteractiveChoiceBrowser, LoginInteractiveChoiceStatus:
		return true
	default:
		return false
	}
}

// String returns the string representation of the menu choice.
func (c LoginInteractiveChoiceType) String() string {
	return string(c)
}

// IsSSH reports whether the choice is SSH setup.
func (c LoginInteractiveChoiceType) IsSSH() bool {
	return c == LoginInteractiveChoiceSSH
}

// IsToken reports whether the choice is Personal Access Token.
func (c LoginInteractiveChoiceType) IsToken() bool {
	return c == LoginInteractiveChoiceToken
}

// IsBrowser reports whether the choice is Browser Device Flow.
func (c LoginInteractiveChoiceType) IsBrowser() bool {
	return c == LoginInteractiveChoiceBrowser
}

// IsStatus reports whether the choice is Authentication Status.
func (c LoginInteractiveChoiceType) IsStatus() bool {
	return c == LoginInteractiveChoiceStatus
}

// ParseLoginInteractiveChoice parses an interactive menu choice input.
func ParseLoginInteractiveChoice(s string) (LoginInteractiveChoiceType, error) {
	val := LoginInteractiveChoiceType(strings.TrimSpace(s))
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid login interactive choice: %q", s)
}

// ParseLoginInteractiveChoiceType is an alias for ParseLoginInteractiveChoice.
func ParseLoginInteractiveChoiceType(s string) (LoginInteractiveChoiceType, error) {
	return ParseLoginInteractiveChoice(s)
}

// AllLoginInteractiveChoices returns all interactive menu choices.
func AllLoginInteractiveChoices() []LoginInteractiveChoiceType {
	return []LoginInteractiveChoiceType{
		LoginInteractiveChoiceSSH,
		LoginInteractiveChoiceToken,
		LoginInteractiveChoiceBrowser,
		LoginInteractiveChoiceStatus,
	}
}

// AllLoginInteractiveChoiceTypes is an alias for AllLoginInteractiveChoices.
func AllLoginInteractiveChoiceTypes() []LoginInteractiveChoiceType {
	return AllLoginInteractiveChoices()
}

// GitAuthProviderType represents remote git hosting providers.
type GitAuthProviderType string

const (
	GitAuthProviderGitHub GitAuthProviderType = "github"
	GitAuthProviderGitLab GitAuthProviderType = "gitlab"
)

// IsValid reports whether the git provider is recognized.
func (p GitAuthProviderType) IsValid() bool {
	switch p {
	case GitAuthProviderGitHub, GitAuthProviderGitLab:
		return true
	default:
		return false
	}
}

// String returns the string representation of the provider.
func (p GitAuthProviderType) String() string {
	return string(p)
}

// IsGitHub reports whether the provider is GitHub.
func (p GitAuthProviderType) IsGitHub() bool {
	return p == GitAuthProviderGitHub
}

// IsGitLab reports whether the provider is GitLab.
func (p GitAuthProviderType) IsGitLab() bool {
	return p == GitAuthProviderGitLab
}

// ParseGitAuthProvider parses a git auth provider name.
func ParseGitAuthProvider(s string) (GitAuthProviderType, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	val := GitAuthProviderType(normalized)
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid git auth provider: %q", s)
}

// ParseGitAuthProviderType is an alias for ParseGitAuthProvider.
func ParseGitAuthProviderType(s string) (GitAuthProviderType, error) {
	return ParseGitAuthProvider(s)
}

// AllGitAuthProviders returns all supported git providers.
func AllGitAuthProviders() []GitAuthProviderType {
	return []GitAuthProviderType{
		GitAuthProviderGitHub,
		GitAuthProviderGitLab,
	}
}

// AllGitAuthProviderTypes is an alias for AllGitAuthProviders.
func AllGitAuthProviderTypes() []GitAuthProviderType {
	return AllGitAuthProviders()
}

// ClipboardToolType represents system clipboard utilities.
type ClipboardToolType string

const (
	ClipboardToolClip   ClipboardToolType = "clip"
	ClipboardToolPbcopy ClipboardToolType = "pbcopy"
	ClipboardToolWlCopy ClipboardToolType = "wl-copy"
	ClipboardToolXclip  ClipboardToolType = "xclip"
	ClipboardToolXsel   ClipboardToolType = "xsel"
)

// IsValid reports whether the clipboard utility is recognized.
func (c ClipboardToolType) IsValid() bool {
	switch c {
	case ClipboardToolClip, ClipboardToolPbcopy, ClipboardToolWlCopy, ClipboardToolXclip, ClipboardToolXsel:
		return true
	default:
		return false
	}
}

// String returns the string representation of the clipboard utility.
func (c ClipboardToolType) String() string {
	return string(c)
}

// IsClip reports whether the tool is Windows clip.
func (c ClipboardToolType) IsClip() bool {
	return c == ClipboardToolClip
}

// IsPbcopy reports whether the tool is macOS pbcopy.
func (c ClipboardToolType) IsPbcopy() bool {
	return c == ClipboardToolPbcopy
}

// IsWlCopy reports whether the tool is Wayland wl-copy.
func (c ClipboardToolType) IsWlCopy() bool {
	return c == ClipboardToolWlCopy
}

// IsXclip reports whether the tool is X11 xclip.
func (c ClipboardToolType) IsXclip() bool {
	return c == ClipboardToolXclip
}

// IsXsel reports whether the tool is X11 xsel.
func (c ClipboardToolType) IsXsel() bool {
	return c == ClipboardToolXsel
}

// ParseClipboardTool parses a clipboard tool name.
func ParseClipboardTool(s string) (ClipboardToolType, error) {
	normalized := strings.ToLower(strings.TrimSpace(s))
	val := ClipboardToolType(normalized)
	if val.IsValid() {
		return val, nil
	}

	return "", fmt.Errorf("invalid clipboard tool: %q", s)
}

// ParseClipboardToolType is an alias for ParseClipboardTool.
func ParseClipboardToolType(s string) (ClipboardToolType, error) {
	return ParseClipboardTool(s)
}

// AllClipboardTools returns all recognized clipboard utilities.
func AllClipboardTools() []ClipboardToolType {
	return []ClipboardToolType{
		ClipboardToolClip,
		ClipboardToolPbcopy,
		ClipboardToolWlCopy,
		ClipboardToolXclip,
		ClipboardToolXsel,
	}
}

// AllClipboardToolTypes is an alias for AllClipboardTools.
func AllClipboardToolTypes() []ClipboardToolType {
	return AllClipboardTools()
}
