package movemerge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ClassifyEndpoint inspects raw and returns its kind and (for URLs)
// the canonical URL plus optional :branch suffix.
func ClassifyEndpoint(raw string) (kind EndpointKind, url, branch, display string) {
	display = strings.TrimRight(raw, "/\\")
	lower := strings.ToLower(display)
	for _, p := range []string{"https://", "http://", "ssh://", "git@"} {
		if strings.HasPrefix(lower, p) {
			url, branch = splitBranchSuffix(display)

			return EndpointURL, url, branch, display
		}
	}

	return EndpointFolder, "", "", display
}

// splitBranchSuffix peels an optional `:branch` off a URL. The
// `git@host:user/repo` form is preserved by re-attaching the first
// colon group when no `/` separates it.
func splitBranchSuffix(raw string) (url, branch string) {
	idx := strings.LastIndex(raw, ":")
	if idx == -1 {
		return raw, ""
	}
	tail := raw[idx+1:]
	if strings.Contains(tail, "/") || tail == "" {
		return raw, ""
	}
	// scp-like git@host:user/repo has no branch suffix.
	if strings.HasPrefix(strings.ToLower(raw), "git@") &&
		strings.Count(raw[:idx], ":") == 0 {
		return raw, ""
	}

	return raw[:idx], tail
}

// MapURLToFolder derives the candidate working folder for a URL endpoint.
// `https://github.com/owner/repo.git` -> `<cwd>/repo`.
func MapURLToFolder(cwd, url string) string {
	base := url
	if i := strings.LastIndex(base, "/"); i >= 0 {
		base = base[i+1:]
	}
	base = strings.TrimSuffix(base, ".git")

	return filepath.Join(cwd, base)
}

// FolderExists reports whether path exists and is a directory.
func FolderExists(path string) (bool, error) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("stat %s: %w", path, err)
	}

	return info.IsDir(), nil
}

// IsGitRepo reports whether path/.git exists.
func IsGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))

	return err == nil && (info.IsDir() || info.Mode().IsRegular())
}
