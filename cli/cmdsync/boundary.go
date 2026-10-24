package cmdsync

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var (
	spec21Regex   = lazyregex.New(`^02-spec/2[1-5]-`)
	bumpRegex     = lazyregex.New(`(?i)^bump|^version\.json$`)
	nonOwnedRegex = lazyregex.New(`(?i)(?:^|[\\/])(oh-my-zsh|ohmyzsh|zsh|omis|oh-my-posh|dotfiles|homebrew|\.cargo|\.rustup|\.nvm|\.asdf|\.pyenv)(?:[\\/]|$)`)
)

// IsSpec21OrHigher checks if a path belongs to private application specs (02-spec/21-* to 25-*).
func IsSpec21OrHigher(relPath string) bool {
	norm := filepath.ToSlash(relPath)

	return spec21Regex.MatchString(norm)
}

// IsBumpScript checks if a filename is a protected version bump script.
func IsBumpScript(filename string) bool {
	base := filepath.Base(filename)

	return bumpRegex.MatchString(base)
}

// IsMemoryOrPlans checks if a path belongs to target-repo operational memory or plans.
func IsMemoryOrPlans(relPath string) bool {
	norm := filepath.ToSlash(relPath)

	return strings.HasPrefix(norm, ".ai-memory/memory/") ||
		strings.HasPrefix(norm, ".ai-memory/plans/") ||
		strings.HasPrefix(norm, ".ai-memory/temp-agents/") ||
		strings.HasPrefix(norm, ".ai-memory/cicd-issues/") ||
		strings.HasPrefix(norm, ".ai-memory/ambiguous-questions/")
}

// IsArchive checks if a path belongs to excluded archive folders.
func IsArchive(relPath string) bool {
	norm := filepath.ToSlash(relPath)

	return strings.HasPrefix(norm, constants.DirArchive+"/") || strings.HasPrefix(norm, constants.DirOldPrompts+"/")
}

// IsSecret checks if a file name or path matches credentials or private secrets.
func IsSecret(filename string) bool {
	base := strings.ToLower(filepath.Base(filename))

	return strings.HasPrefix(base, ".env") || strings.HasSuffix(base, ".pem") ||
		strings.HasSuffix(base, ".key") || strings.Contains(base, "id_rsa") ||
		strings.Contains(base, "id_ed25519") || base == constants.FileCredentialsJSON
}

// IsNonOwnedRepo checks if a repo name or path corresponds to non-owned third-party configurations.
func IsNonOwnedRepo(nameOrPath string) bool {
	norm := filepath.ToSlash(nameOrPath)

	return nonOwnedRegex.MatchString(norm)
}
