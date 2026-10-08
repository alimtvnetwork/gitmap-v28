package cmdsync

import (
	"path/filepath"
	"regexp"
	"strings"
)

var (
	spec21Regex   = regexp.MustCompile(`^02-spec/2[1-5]-`)
	bumpRegex     = regexp.MustCompile(`(?i)^bump|^version\.json$`)
	nonOwnedRegex = regexp.MustCompile(`(?i)(?:^|[\\/])(oh-my-zsh|ohmyzsh|zsh|omis|oh-my-posh|dotfiles|homebrew|\.cargo|\.rustup|\.nvm|\.asdf|\.pyenv)(?:[\\/]|$)`)
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
	return strings.HasPrefix(norm, "06-archive/") || strings.HasPrefix(norm, "06-old-prompts/")
}

// IsSecret checks if a file name or path matches credentials or private secrets.
func IsSecret(filename string) bool {
	base := strings.ToLower(filepath.Base(filename))
	if strings.HasPrefix(base, ".env") || strings.HasSuffix(base, ".pem") ||
		strings.HasSuffix(base, ".key") || strings.Contains(base, "id_rsa") ||
		strings.Contains(base, "id_ed25519") || base == "credentials.json" {
		return true
	}
	return false
}

// IsNonOwnedRepo checks if a repo name or path corresponds to non-owned third-party configurations.
func IsNonOwnedRepo(nameOrPath string) bool {
	norm := filepath.ToSlash(nameOrPath)
	return nonOwnedRegex.MatchString(norm)
}
