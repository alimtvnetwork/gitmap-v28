// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

// ParseSourceFile reads and intelligently parses a manifest text file.
func ParseSourceFile(filePath string) (string, []ResolvedSourceItem, []SkippedSourceItem, error) {
	resolvedPath, errResolve := ResolveSourceFilePath(filePath)
	if errResolve != nil {
		return "", nil, nil, errResolve
	}

	data, errRead := os.ReadFile(resolvedPath)
	if errRead != nil {
		return "", nil, nil, apperror.WrapSimple(errRead, fmt.Sprintf("failed to read source file %s", resolvedPath))
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

func isDestHintPrefix(lower string) bool {
	return strings.HasPrefix(lower, "# dest") || strings.HasPrefix(lower, "# target") ||
		strings.HasPrefix(lower, "dest:") || strings.HasPrefix(lower, "target:") ||
		strings.HasPrefix(lower, "destination:") || strings.HasPrefix(lower, "# destination:")
}

func extractDestValue(trimmed, currentDest string) string {
	if currentDest != "" {
		return currentDest
	}
	parts := strings.SplitN(trimmed, ":", 2)
	if len(parts) == 2 {
		return strings.TrimSpace(parts[1])
	}
	return currentDest
}

func appendHeadingLine(lines []string, raw string) []string {
	heading := strings.TrimSpace(strings.TrimLeft(raw, "#"))
	if heading != "" {
		return append(lines, heading)
	}
	return lines
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
		if isDestHintPrefix(lower) {
			destHint = extractDestValue(trimmed, destHint)
			continue
		}

		if strings.HasPrefix(trimmed, "#") {
			cleanLines = appendHeadingLine(cleanLines, trimmed)
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

func resolveHeaderTitle(curr, nextURL string) string {
	if isGenericLabel(curr) {
		return inferTitleFromURL(nextURL)
	}
	return curr
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
			title := resolveHeaderTitle(curr, lines[i+1])
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

func resolveItemTitle(title, url string) string {
	if title == "" {
		return inferTitleFromURL(url)
	}
	return title
}

func processFileURLItem(
	item RawSourceItem,
	canonicalRepo, cloneURL, filePath string,
	fullRepos, seenRepos map[string]bool,
	resolved []ResolvedSourceItem,
	skipped []SkippedSourceItem,
) ([]ResolvedSourceItem, []SkippedSourceItem) {
	if fullRepos[canonicalRepo] || seenRepos[canonicalRepo] {
		skipped = append(skipped, SkippedSourceItem{
			RawInput: item.URL,
			Title:    item.Title,
			Repo:     canonicalRepo,
			Reason:   fmt.Sprintf("file URL from repository '%s' which is already included", canonicalRepo),
			IsFile:   true,
			FilePath: filePath,
		})
		return resolved, skipped
	}

	seenRepos[canonicalRepo] = true
	title := resolveItemTitle(item.Title, item.URL)
	resolved = append(resolved, ResolvedSourceItem{
		Title:         title,
		RepoURL:       cloneURL,
		CanonicalRepo: canonicalRepo,
		OriginalInput: item.URL,
		IsFromFile:    true,
	})
	return resolved, skipped
}

func processRepoURLItem(
	item RawSourceItem,
	canonicalRepo, cloneURL string,
	seenRepos map[string]bool,
	resolved []ResolvedSourceItem,
	skipped []SkippedSourceItem,
) ([]ResolvedSourceItem, []SkippedSourceItem) {
	if seenRepos[canonicalRepo] {
		skipped = append(skipped, SkippedSourceItem{
			RawInput: item.URL,
			Title:    item.Title,
			Repo:     canonicalRepo,
			Reason:   fmt.Sprintf("duplicate repository '%s'", canonicalRepo),
			IsFile:   false,
		})
		return resolved, skipped
	}

	seenRepos[canonicalRepo] = true
	title := resolveItemTitle(item.Title, item.URL)
	resolved = append(resolved, ResolvedSourceItem{
		Title:         title,
		RepoURL:       cloneURL,
		CanonicalRepo: canonicalRepo,
		OriginalInput: item.URL,
		IsFromFile:    false,
	})
	return resolved, skipped
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
			resolved, skipped = processFileURLItem(item, canonicalRepo, cloneURL, filePath, fullRepos, seenRepos, resolved, skipped)
			continue
		}

		resolved, skipped = processRepoURLItem(item, canonicalRepo, cloneURL, seenRepos, resolved, skipped)
	}

	return resolved, skipped
}
