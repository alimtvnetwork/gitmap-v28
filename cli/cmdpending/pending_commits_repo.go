package cmdpending

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// resolveProperRepoName ensures opaque hash names or UUIDs resolve to their canonical repo name.
func resolveProperRepoName(dir, currentName string) string {
	cleanName := strings.TrimSpace(currentName)
	if !looksLikeHashOrOpaqueID(cleanName) && cleanName != "" && cleanName != "." {
		return cleanName
	}

	remoteURL := queryRepoRemoteURL(dir)
	repoName := extractRepoNameFromRemoteURL(remoteURL)
	if remoteURL != "" && repoName != "" {
		return repoName
	}

	dbRepoName := lookupRepoNameInStore(dir, remoteURL)
	if dbRepoName != "" {
		return dbRepoName
	}

	base := filepath.Base(dir)
	if !looksLikeHashOrOpaqueID(base) && base != "" && base != "." {
		return base
	}

	return cleanName
}

func looksLikeHashOrOpaqueID(s string) bool {
	clean := strings.TrimSpace(s)
	if len(clean) >= 20 && isAllHexRunes(clean) {
		return true
	}

	return len(clean) == 36 && strings.Count(clean, "-") == 4
}

func isAllHexRunes(s string) bool {
	for _, r := range s {
		if !isHexRune(r) {
			return false
		}
	}
	return true
}

func isHexRune(r rune) bool {
	return (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
}

func queryRepoRemoteURL(dir string) string {
	out, err := currentPendingCommitsGitExecutor(dir, "config", "--get", "remote.origin.url")
	clean := strings.TrimSpace(out)
	if err == nil && clean != "" && !isGitErrorOutput(clean) {
		return clean
	}

	outRemote, errRemote := currentPendingCommitsGitExecutor(dir, "remote", "get-url", "origin")
	cleanRemote := strings.TrimSpace(outRemote)
	if errRemote == nil && cleanRemote != "" && !isGitErrorOutput(cleanRemote) {
		return cleanRemote
	}

	return ""
}

func extractRepoNameFromRemoteURL(url string) string {
	clean := strings.TrimRight(strings.TrimSpace(url), "/\\")
	clean = strings.TrimSuffix(clean, ".git")
	clean = strings.TrimRight(clean, "/\\")
	if idx := strings.LastIndexAny(clean, "/:"); idx != -1 {
		return clean[idx+1:]
	}
	return clean
}

func lookupRepoNameInStore(dir, remoteURL string) string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	repos, errList := db.ListRepos()
	if errList != nil {
		return ""
	}

	for _, r := range repos {
		if stringsEqualAbs(r.AbsolutePath, dir) {
			return r.RepoName
		}
		if remoteURL != "" && (strings.EqualFold(r.HTTPSUrl, remoteURL) || strings.EqualFold(r.SSHUrl, remoteURL)) {
			return r.RepoName
		}
	}

	return ""
}

func stringsEqualAbs(a, b string) bool {
	return strings.EqualFold(filepath.Clean(a), filepath.Clean(b))
}

func safeQueryRepoBranch(dir, fallback string) string {
	branch, err := currentPendingCommitsGitExecutor(dir, "rev-parse", "--abbrev-ref", "HEAD")
	branch = strings.TrimSpace(branch)
	if err == nil && branch != "" && !isGitErrorOutput(branch) {
		return branch
	}

	symBranch, errSym := currentPendingCommitsGitExecutor(dir, "symbolic-ref", "--short", "HEAD")
	symBranch = strings.TrimSpace(symBranch)
	if errSym == nil && symBranch != "" && !isGitErrorOutput(symBranch) {
		return symBranch
	}

	if fallback != "" && !isGitErrorOutput(fallback) {
		return fallback
	}

	return "main"
}

func isGitErrorOutput(s string) bool {
	lower := strings.ToLower(s)
	return strings.Contains(lower, "fatal:") || strings.Contains(lower, "error:") || strings.HasPrefix(lower, "fatal")
}
