package cmdclone

import (
	"regexp"
	"strings"
)

var (
	rxMarkdownLink = regexp.MustCompile(`\[.*?\]\((https?://[^\s\)]+)\)`)
	rxDirectURL    = regexp.MustCompile(`(https?://[^\s]+|git@[a-zA-Z0-9_.-]+:[^\s]+)`)
	rxOwnerRepo    = regexp.MustCompile(`^([a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+)`)
)

// ParseMultiCloneText parses raw text, markdown blocks, and URL lists into unique git URLs.
func ParseMultiCloneText(text string) []string {
	lines := strings.Split(text, "\n")
	seen := make(map[string]struct{})
	var results []string

	for _, rawLine := range lines {
		url := processMultiCloneLine(rawLine)
		isValid := len(url) > 0
		if !isValid {
			continue
		}

		key := normaliseURLKey(url)
		_, isDuplicate := seen[key]
		if isDuplicate {
			continue
		}

		seen[key] = struct{}{}
		results = append(results, url)
	}

	return results
}

func processMultiCloneLine(rawLine string) string {
	line := strings.TrimSpace(rawLine)
	isIgnored := len(line) == 0 || isMarkdownFenceLine(line)
	if isIgnored {
		return ""
	}

	return extractRepoURLCandidate(line)
}

func isMarkdownFenceLine(line string) bool {
	return strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~")
}

func extractRepoURLCandidate(line string) string {
	linkMatch := rxMarkdownLink.FindStringSubmatch(line)
	if len(linkMatch) > 1 {
		return cleanExtractedURL(linkMatch[1])
	}

	urlMatch := rxDirectURL.FindStringSubmatch(line)
	if len(urlMatch) > 1 {
		return cleanExtractedURL(urlMatch[1])
	}

	return extractOwnerRepoCandidate(line)
}

func extractOwnerRepoCandidate(line string) string {
	parts := strings.Fields(line)
	if len(parts) == 0 {
		return ""
	}

	token := strings.TrimSuffix(parts[0], ":")
	repoMatch := rxOwnerRepo.FindStringSubmatch(token)
	if len(repoMatch) > 1 {
		return "https://github.com/" + repoMatch[1] + ".git"
	}

	return ""
}

func cleanExtractedURL(raw string) string {
	cleaned := strings.TrimRight(raw, "/\\\"'`")
	cleaned = strings.TrimPrefix(cleaned, "<")
	cleaned = strings.TrimSuffix(cleaned, ">")

	return strings.TrimSpace(cleaned)
}
