// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"strings"
)

// ParseURLDetails breaks down any Git or GitHub URL into canonical components.
func ParseURLDetails(rawURL string) (string, string, bool, string, error) {
	trimmed := strings.TrimSpace(rawURL)
	if trimmed == "" {
		return "", "", false, "", fmt.Errorf("empty URL")
	}

	if isShorthandRepo(trimmed) {
		cleanRepo := strings.TrimSuffix(trimmed, ".git")
		return strings.ToLower(cleanRepo), "https://github.com/" + cleanRepo + ".git", false, "", nil
	}

	if strings.HasPrefix(trimmed, "git@") {
		parts := strings.Split(trimmed, ":")
		if len(parts) == 2 {
			repoPath := strings.TrimSuffix(parts[1], ".git")
			repoPath = strings.Trim(repoPath, "/")
			return strings.ToLower(repoPath), trimmed, false, "", nil
		}
	}

	u := trimmed
	for _, prefix := range []string{"https://", "http://", "ssh://"} {
		if strings.HasPrefix(strings.ToLower(u), prefix) {
			u = u[len(prefix):]
			break
		}
	}

	slashIdx := strings.Index(u, "/")
	if slashIdx == -1 {
		if isShorthandRepo(trimmed) {
			cleanRepo := strings.TrimSuffix(trimmed, ".git")
			return strings.ToLower(cleanRepo), "https://github.com/" + cleanRepo + ".git", false, "", nil
		}
		return "", "", false, "", fmt.Errorf("invalid repository URL: %s", rawURL)
	}

	pathOnly := strings.Trim(u[slashIdx+1:], "/")
	segments := strings.Split(pathOnly, "/")
	if len(segments) < 2 {
		return "", "", false, "", fmt.Errorf("invalid repository path in URL: %s", rawURL)
	}

	owner := segments[0]
	repo := strings.TrimSuffix(segments[1], ".git")
	canonicalRepo := strings.ToLower(owner + "/" + repo)
	cloneURL := fmt.Sprintf("https://github.com/%s/%s.git", owner, repo)

	var isFileURL bool
	var filePath string
	if len(segments) > 2 {
		kind := strings.ToLower(segments[2])
		if kind == "blob" || kind == "tree" || kind == "raw" {
			isFileURL = true
			if len(segments) > 4 {
				filePath = strings.Join(segments[4:], "/")
			} else if len(segments) == 4 {
				filePath = segments[3]
			}
		}
	}

	return canonicalRepo, cloneURL, isFileURL, filePath, nil
}

func isFullURL(s string) bool {
	lower := strings.ToLower(strings.TrimSpace(s))
	return strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "git@") ||
		strings.HasPrefix(lower, "ssh://")
}

func isShorthandRepo(s string) bool {
	trimmed := strings.TrimSpace(s)
	if strings.Contains(trimmed, " ") || isFullURL(trimmed) {
		return false
	}
	parts := strings.Split(trimmed, "/")
	return len(parts) == 2 && parts[0] != "" && parts[1] != ""
}

func inferTitleFromURL(rawURL string) string {
	canonicalRepo, _, _, _, err := ParseURLDetails(rawURL)
	if err == nil && canonicalRepo != "" {
		return canonicalRepo
	}
	trimmed := strings.TrimRight(strings.TrimSuffix(rawURL, ".git"), "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "repository"
}
