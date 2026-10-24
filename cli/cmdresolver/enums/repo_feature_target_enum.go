package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// RepoFeatureTargetType enumerates classification modes for destination targets.
type RepoFeatureTargetType string

const (
	TargetTypeRemoteURL  RepoFeatureTargetType = "remote-url"
	TargetTypeLocalGit   RepoFeatureTargetType = "local-git"
	TargetTypeLocalNoGit RepoFeatureTargetType = "local-no-git"
	TargetTypeBareSlug   RepoFeatureTargetType = "bare-slug"
)

// IsValid reports whether the value matches a known variant.
func (r RepoFeatureTargetType) IsValid() bool {
	switch r {
	case TargetTypeRemoteURL, TargetTypeLocalGit, TargetTypeLocalNoGit, TargetTypeBareSlug:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (r RepoFeatureTargetType) String() string {
	return string(r)
}

// ParseRepoFeatureTargetType parses a raw string into a validated RepoFeatureTargetType.
func ParseRepoFeatureTargetType(s string) (RepoFeatureTargetType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllRepoFeatureTargetTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "targettyperemoteurl":
		return TargetTypeRemoteURL, nil
	case "remoteurl":
		return TargetTypeRemoteURL, nil
	case "targettypelocalgit":
		return TargetTypeLocalGit, nil
	case "localgit":
		return TargetTypeLocalGit, nil
	case "targettypelocalnogit":
		return TargetTypeLocalNoGit, nil
	case "localnogit":
		return TargetTypeLocalNoGit, nil
	case "targettypebareslug":
		return TargetTypeBareSlug, nil
	case "bareslug":
		return TargetTypeBareSlug, nil
	}
	return "", apperror.NewValidation("ParseRepoFeatureTargetType", "E1000", fmt.Sprintf("invalid repofeaturetargettype: %q", s)).WithContext("value", s)
}

// AllRepoFeatureTargetTypes returns all valid variants.
func AllRepoFeatureTargetTypes() []RepoFeatureTargetType {
	return []RepoFeatureTargetType{
		TargetTypeRemoteURL,
		TargetTypeLocalGit,
		TargetTypeLocalNoGit,
		TargetTypeBareSlug,
	}
}

// IsRemoteURL reports whether the value is TargetTypeRemoteURL.
func (r RepoFeatureTargetType) IsRemoteURL() bool {
	return r == TargetTypeRemoteURL
}

// IsLocalGit reports whether the value is TargetTypeLocalGit.
func (r RepoFeatureTargetType) IsLocalGit() bool {
	return r == TargetTypeLocalGit
}

// IsLocalNoGit reports whether the value is TargetTypeLocalNoGit.
func (r RepoFeatureTargetType) IsLocalNoGit() bool {
	return r == TargetTypeLocalNoGit
}

// IsBareSlug reports whether the value is TargetTypeBareSlug.
func (r RepoFeatureTargetType) IsBareSlug() bool {
	return r == TargetTypeBareSlug
}

// HasFeatureTarget reports whether the value is non-empty and valid.
func (r RepoFeatureTargetType) HasFeatureTarget() bool {
	return r.IsValid()
}
