package gitignoreagm

import (
	"strings"
)

// DefaultGitmapIgnoreEntries lists the canonical paths automatically ensured by GitMap.
var DefaultGitmapIgnoreEntries = []string{
	".gitmap/backup/",
	"antigravity-resume_task.json",
	".antigravity_resume_task.json",
	"antigravity_resume_task.json",
	".antigravity-resume_task.json",
}

// DeduplicateAndSanitizeGitignore cleans duplicate pattern lines and redundant comment blocks.
func DeduplicateAndSanitizeGitignore(content string, extraDefaults ...string) (string, bool) {
	lines := strings.Split(content, "\n")
	seenPatterns := make(map[string]bool, len(lines))
	seenComments := make(map[string]bool, len(lines))
	cleaned := make([]string, 0, len(lines))
	isModified := false

	for _, rawLine := range lines {
		trimmed := strings.TrimSpace(rawLine)
		if isRetainedLine(trimmed, &seenPatterns, &seenComments) {
			cleaned = append(cleaned, rawLine)
		} else {
			isModified = true
		}
	}

	appended, hasAppended := appendMissingDefaults(cleaned, seenPatterns, extraDefaults)
	if hasAppended {
		isModified = true
		cleaned = appended
	}

	result := normalizeGitignoreOutput(cleaned)
	return result, isModified || result != content
}

func isRetainedLine(trimmed string, seenPatterns, seenComments *map[string]bool) bool {
	if trimmed == "" {
		return true
	}
	if strings.HasPrefix(trimmed, "#") {
		return checkAndMarkComment(trimmed, seenComments)
	}
	return checkAndMarkPattern(trimmed, seenPatterns)
}

func checkAndMarkComment(comment string, seenComments *map[string]bool) bool {
	norm := strings.ToLower(strings.TrimSpace(comment))
	if (*seenComments)[norm] {
		return false
	}
	(*seenComments)[norm] = true
	return true
}

func checkAndMarkPattern(pattern string, seenPatterns *map[string]bool) bool {
	norm := strings.TrimPrefix(pattern, "/")
	if (*seenPatterns)[norm] {
		return false
	}
	(*seenPatterns)[norm] = true
	return true
}

func appendMissingDefaults(lines []string, seen map[string]bool, extras []string) ([]string, bool) {
	missing := collectMissingDefaults(seen, DefaultGitmapIgnoreEntries)
	missing = append(missing, collectMissingDefaults(seen, extras)...)
	if len(missing) == 0 {
		return lines, false
	}

	result := append(lines, "", "# GitMap & Task Persistence")
	for _, entry := range missing {
		result = append(result, entry)
	}
	return result, true
}

func collectMissingDefaults(seen map[string]bool, entries []string) []string {
	var missing []string
	for _, entry := range entries {
		norm := strings.TrimPrefix(entry, "/")
		if !seen[norm] {
			missing = append(missing, entry)
			seen[norm] = true
		}
	}
	return missing
}

func normalizeGitignoreOutput(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	text := strings.Join(lines, "\n")
	trimmed := strings.TrimRight(text, "\r\n")
	if trimmed == "" {
		return ""
	}
	return trimmed + "\n"
}
