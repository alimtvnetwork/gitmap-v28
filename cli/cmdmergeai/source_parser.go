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

	candidates := []string{trimmed}
	lower := strings.ToLower(trimmed)
	if !strings.HasSuffix(lower, ".txt") && !strings.HasSuffix(lower, ".manifest") {
		candidates = append(candidates, trimmed+".txt")
	}

	for _, cand := range candidates {
		if resolved, ok := probeSingleCandidate(cand); ok {
			return resolved, nil
		}
	}

	return trimmed, fmt.Errorf("source file not found: %s", trimmed)
}

func probeSingleCandidate(cand string) (string, bool) {
	if info, err := os.Stat(cand); err == nil && !info.IsDir() {
		return cand, true
	}

	if cwd, err := os.Getwd(); err == nil && cwd != "" {
		cwdCandidate := filepath.Join(cwd, cand)
		if info, err := os.Stat(cwdCandidate); err == nil && !info.IsDir() {
			return cwdCandidate, true
		}

		cwdRcCandidate := filepath.Join(cwd, "repo-cache", filepath.Base(cand))
		if info, err := os.Stat(cwdRcCandidate); err == nil && !info.IsDir() {
			return cwdRcCandidate, true
		}
	}

	repoCacheRoot := cmdclone.ResolveRepoCacheRoot()
	if repoCacheRoot != "" {
		rcCandidate := filepath.Join(repoCacheRoot, cand)
		if info, err := os.Stat(rcCandidate); err == nil && !info.IsDir() {
			return rcCandidate, true
		}

		rcBaseCandidate := filepath.Join(repoCacheRoot, filepath.Base(cand))
		if info, err := os.Stat(rcBaseCandidate); err == nil && !info.IsDir() {
			return rcBaseCandidate, true
		}
	}

	workBase := resolveWorkBaseDir()
	if workBase != "" {
		wbCandidate := filepath.Join(workBase, cand)
		if info, err := os.Stat(wbCandidate); err == nil && !info.IsDir() {
			return wbCandidate, true
		}

		wbBaseCandidate := filepath.Join(workBase, "repo-cache", filepath.Base(cand))
		if info, err := os.Stat(wbBaseCandidate); err == nil && !info.IsDir() {
			return wbBaseCandidate, true
		}
	}

	return "", false
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
		if strings.HasPrefix(lower, "# dest") || strings.HasPrefix(lower, "# target") ||
			strings.HasPrefix(lower, "dest:") || strings.HasPrefix(lower, "target:") ||
			strings.HasPrefix(lower, "destination:") || strings.HasPrefix(lower, "# destination:") {
			parts := strings.SplitN(trimmed, ":", 2)
			if len(parts) == 2 && destHint == "" {
				destHint = strings.TrimSpace(parts[1])
			}
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			heading := strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
			if heading != "" {
				cleanLines = append(cleanLines, heading)
			}
			continue
		}

		cleanLines = append(cleanLines, trimmed)
	}

	return destHint, cleanLines
}

func isGenericLabel(s string) bool {
	clean := strings.ToLower(strings.TrimRight(strings.TrimSpace(s), ":"))
	switch clean {
	case "title", "repo", "repository", "url", "source", "name", "link":
		return true
	default:
		return false
	}
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
			title := curr
			if isGenericLabel(title) {
				title = inferTitleFromURL(lines[i+1])
			}
			items = append(items, RawSourceItem{
				Title:   title,
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
	fullRepos := make(map[string]bool)
	for _, item := range rawItems {
		canonicalRepo, _, isFileURL, _, err := ParseURLDetails(item.URL)
		if err == nil && !isFileURL {
			fullRepos[canonicalRepo] = true
		}
	}

	seenRepos := make(map[string]bool)
	var resolved []ResolvedSourceItem
	var skipped []SkippedSourceItem

	for _, item := range rawItems {
		canonicalRepo, cloneURL, isFileURL, filePath, err := ParseURLDetails(item.URL)
		if err != nil {
			continue
		}

		if isFileURL {
			if fullRepos[canonicalRepo] || seenRepos[canonicalRepo] {
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
