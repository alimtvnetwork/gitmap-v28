// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

// RawSourceItem represents an extracted title/URL pairing from manifest text.
type RawSourceItem struct {
	Title   string
	URL     string
	RawLine string
	LineNum int
}

// ResolvedSourceItem represents a validated unique source repository to merge.
type ResolvedSourceItem struct {
	Title         string
	RepoURL       string
	CanonicalRepo string
	OriginalInput string
	IsFromFile    bool
}

// SkippedSourceItem represents a source skipped by manifest intelligence.
type SkippedSourceItem struct {
	RawInput string
	Title    string
	Repo     string
	Reason   string
	IsFile   bool
	FilePath string
}

// ResolveSourceFilePath locates an input text file directly or inside repo-cache.
func ResolveSourceFilePath(inputPath string) (string, error) {
	trimmed := strings.TrimSpace(inputPath)
	if trimmed == "" {
		return "", fmt.Errorf("empty source file path")
	}

	if info, err := os.Stat(trimmed); err == nil && !info.IsDir() {
		return trimmed, nil
	}

	repoCacheRoot := cmdclone.ResolveRepoCacheRoot()
	if repoCacheRoot != "" {
		rcCandidate := filepath.Join(repoCacheRoot, trimmed)
		if info, err := os.Stat(rcCandidate); err == nil && !info.IsDir() {
			return rcCandidate, nil
		}

		baseName := filepath.Base(trimmed)
		rcBaseCandidate := filepath.Join(repoCacheRoot, baseName)
		if info, err := os.Stat(rcBaseCandidate); err == nil && !info.IsDir() {
			return rcBaseCandidate, nil
		}
	}

	return trimmed, fmt.Errorf("source file not found: %s", trimmed)
}

// ParseSourceFile reads and intelligently parses a manifest text file.
func ParseSourceFile(filePath string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	resolvedPath, errResolve := ResolveSourceFilePath(filePath)
	if errResolve != nil {
		return "", nil, nil, errResolve
	}

	data, errRead := os.ReadFile(resolvedPath)
	if errRead != nil {
		return "", nil, nil, fmt.Errorf("failed to read source file %s: %w", resolvedPath, errRead)
	}

	return ParseSourceManifestContent(string(data))
}

// ParseSourceManifestContent parses lines of text containing URLs, headers, and file paths.
func ParseSourceManifestContent(content string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	destHint, cleanLines := extractDestHintAndCleanLines(content)
	rawItems := pairHeadersAndURLs(cleanLines)

	resolved, skipped := filterAndDeduplicateSources(rawItems)
	return destHint, resolved, skipped, nil
}

func extractDestHintAndCleanLines(content string) (string, []string) {
	var destHint string
	var cleanLines []string

	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if trimmed == "" {
			continue
		}

		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "# dest") || strings.HasPrefix(lower, "# target") || strings.HasPrefix(lower, "dest:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 && destHint == "" {
				destHint = strings.TrimSpace(parts[1])
			}
			continue
		}

		if !strings.HasPrefix(trimmed, "#") {
			cleanLines = append(cleanLines, trimmed)
		}
	}

	return destHint, cleanLines
}

func pairHeadersAndURLs(lines []string) []RawSourceItem {
	var items []RawSourceItem
	i := 0
	for i < len(lines) {
		curr := lines[i]

		if isFullURL(curr) {
			items = append(items, RawSourceItem{
				Title:   inferTitleFromURL(curr),
				URL:     curr,
				RawLine: curr,
				LineNum: i + 1,
			})
			i++
			continue
		}

		if i+1 < len(lines) && isFullURL(lines[i+1]) {
			items = append(items, RawSourceItem{
				Title:   curr,
				URL:     lines[i+1],
				RawLine: curr + "\n" + lines[i+1],
				LineNum: i + 1,
			})
			i += 2
			continue
		}

		if isShorthandRepo(curr) {
			items = append(items, RawSourceItem{
				Title:   curr,
				URL:     "https://github.com/" + curr + ".git",
				RawLine: curr,
				LineNum: i + 1,
			})
			i++
			continue
		}

		i++
	}

	return items
}

func filterAndDeduplicateSources(rawItems []RawSourceItem) ([]ResolvedSourceItem, []SkippedSourceItem) {
	seenRepos := make(map[string]bool)
	var resolved []ResolvedSourceItem
	var skipped []SkippedSourceItem

	for _, item := range rawItems {
		canonicalRepo, cloneURL, isFileURL, filePath, err := ParseURLDetails(item.URL)
		if err != nil {
			continue
		}

		if isFileURL {
			if seenRepos[canonicalRepo] {
				skipped = append(skipped, SkippedSourceItem{
					RawInput: item.URL,
					Title:    item.Title,
					Repo:     canonicalRepo,
					Reason:   fmt.Sprintf("file URL from repository '%s' which is already included", canonicalRepo),
					IsFile:   true,
					FilePath: filePath,
				})
				continue
			}

			seenRepos[canonicalRepo] = true
			title := item.Title
			if title == "" {
				title = inferTitleFromURL(item.URL)
			}
			resolved = append(resolved, ResolvedSourceItem{
				Title:         title,
				RepoURL:       cloneURL,
				CanonicalRepo: canonicalRepo,
				OriginalInput: item.URL,
				IsFromFile:    true,
			})
			continue
		}

		if seenRepos[canonicalRepo] {
			skipped = append(skipped, SkippedSourceItem{
				RawInput: item.URL,
				Title:    item.Title,
				Repo:     canonicalRepo,
				Reason:   fmt.Sprintf("duplicate repository '%s'", canonicalRepo),
				IsFile:   false,
			})
			continue
		}

		seenRepos[canonicalRepo] = true
		title := item.Title
		if title == "" {
			title = inferTitleFromURL(item.URL)
		}
		resolved = append(resolved, ResolvedSourceItem{
			Title:         title,
			RepoURL:       cloneURL,
			CanonicalRepo: canonicalRepo,
			OriginalInput: item.URL,
			IsFromFile:    false,
		})
	}

	return resolved, skipped
}

// ParseURLDetails breaks down any Git or GitHub URL into components.
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
