package constants

// SpecialPathType represents short and long alias tokens for special system directories.
type SpecialPathType string

const (
	SpecialPathRepoCache      SpecialPathType = "repo-cache"
	SpecialPathRepoCacheShort SpecialPathType = "rc"
	SpecialPathRepoSecrets    SpecialPathType = "repo-secrets"
	SpecialPathRepoSecretsShort SpecialPathType = "rs"
)

// IsValid reports whether p is a recognized special path token.
func (p SpecialPathType) IsValid() bool {
	switch p {
	case SpecialPathRepoCache, SpecialPathRepoCacheShort, SpecialPathRepoSecrets, SpecialPathRepoSecretsShort:
		return true
	default:
		return false
	}
}

// IsRepoCache reports whether p represents the repo-cache directory.
func (p SpecialPathType) IsRepoCache() bool {
	return p == SpecialPathRepoCache || p == SpecialPathRepoCacheShort
}

// IsRepoSecrets reports whether p represents the repo-secrets directory.
func (p SpecialPathType) IsRepoSecrets() bool {
	return p == SpecialPathRepoSecrets || p == SpecialPathRepoSecretsShort
}

// String returns the string representation.
func (p SpecialPathType) String() string {
	return string(p)
}
