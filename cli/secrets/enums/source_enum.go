package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// SourceType identifies which mechanism produced a token.
type SourceType string

const (
	SourceEnvGitHubToken  SourceType = "GITHUB_TOKEN env var"
	SourceEnvGHToken      SourceType = "GH_TOKEN env var"
	SourceWinRegistryUser SourceType = "Windows User Registry"
	SourceWinRegistrySys  SourceType = "Windows System Registry"
	SourceGitCredential   SourceType = "Git Credential Manager"
	SourceGhCLI           SourceType = "GitHub CLI (gh auth token)"
	SourceGhHostsFile     SourceType = "GitHub CLI config (hosts.yml)"
	SourceGitConfig       SourceType = "git config (github.token)"
	SourceNone            SourceType = ""
)

// IsValid reports whether the value matches a known variant.
func (s SourceType) IsValid() bool {
	switch s {
	case SourceEnvGitHubToken, SourceEnvGHToken, SourceWinRegistryUser, SourceWinRegistrySys, SourceGitCredential, SourceGhCLI, SourceGhHostsFile, SourceGitConfig, SourceNone:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (s SourceType) String() string {
	return string(s)
}

// ParseSourceType parses a raw string into a validated SourceType.
func ParseSourceType(s string) (SourceType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllSourceTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "sourceenvgithubtoken":
		return SourceEnvGitHubToken, nil
	case "envgithubtoken":
		return SourceEnvGitHubToken, nil
	case "sourceenvghtoken":
		return SourceEnvGHToken, nil
	case "envghtoken":
		return SourceEnvGHToken, nil
	case "sourcewinregistryuser":
		return SourceWinRegistryUser, nil
	case "winregistryuser":
		return SourceWinRegistryUser, nil
	case "sourcewinregistrysys":
		return SourceWinRegistrySys, nil
	case "winregistrysys":
		return SourceWinRegistrySys, nil
	case "sourcegitcredential":
		return SourceGitCredential, nil
	case "gitcredential":
		return SourceGitCredential, nil
	case "sourceghcli":
		return SourceGhCLI, nil
	case "ghcli":
		return SourceGhCLI, nil
	case "sourceghhostsfile":
		return SourceGhHostsFile, nil
	case "ghhostsfile":
		return SourceGhHostsFile, nil
	case "sourcegitconfig":
		return SourceGitConfig, nil
	case "gitconfig":
		return SourceGitConfig, nil
	case "sourcenone":
		return SourceNone, nil
	case "none":
		return SourceNone, nil
	}
	return "", apperror.NewValidation("ParseSourceType", "E1000", fmt.Sprintf("invalid sourcetype: %q", s)).WithContext("value", s)
}

// AllSourceTypes returns all valid variants.
func AllSourceTypes() []SourceType {
	return []SourceType{
		SourceEnvGitHubToken,
		SourceEnvGHToken,
		SourceWinRegistryUser,
		SourceWinRegistrySys,
		SourceGitCredential,
		SourceGhCLI,
		SourceGhHostsFile,
		SourceGitConfig,
		SourceNone,
	}
}

// IsEnvGitHubToken reports whether the value is SourceEnvGitHubToken.
func (s SourceType) IsEnvGitHubToken() bool {
	return s == SourceEnvGitHubToken
}

// IsEnvGHToken reports whether the value is SourceEnvGHToken.
func (s SourceType) IsEnvGHToken() bool {
	return s == SourceEnvGHToken
}

// IsWinRegistryUser reports whether the value is SourceWinRegistryUser.
func (s SourceType) IsWinRegistryUser() bool {
	return s == SourceWinRegistryUser
}

// IsWinRegistrySys reports whether the value is SourceWinRegistrySys.
func (s SourceType) IsWinRegistrySys() bool {
	return s == SourceWinRegistrySys
}

// IsGitCredential reports whether the value is SourceGitCredential.
func (s SourceType) IsGitCredential() bool {
	return s == SourceGitCredential
}

// IsGhCLI reports whether the value is SourceGhCLI.
func (s SourceType) IsGhCLI() bool {
	return s == SourceGhCLI
}

// IsGhHostsFile reports whether the value is SourceGhHostsFile.
func (s SourceType) IsGhHostsFile() bool {
	return s == SourceGhHostsFile
}

// IsGitConfig reports whether the value is SourceGitConfig.
func (s SourceType) IsGitConfig() bool {
	return s == SourceGitConfig
}

// IsNone reports whether the value is SourceNone.
func (s SourceType) IsNone() bool {
	return s == SourceNone
}

// HasSource reports whether the value is non-empty and valid.
func (s SourceType) HasSource() bool {
	return s.IsValid()
}
