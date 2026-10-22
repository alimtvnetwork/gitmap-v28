package cmdplans

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	idSlugRegex     = regexp.MustCompile(`^(\d+)-(.*)\.md$`)
	statusLineRegex = regexp.MustCompile(`(?i)(?:##\s*|>\s*\*\*)?Status:\s*\*?\*?\s*([a-zA-Z\s_-]+)`)
	headingRegex    = regexp.MustCompile(`^#\s+(.*)`)
)

func parseIDAndSlug(fileName string) (int, string) {
	matches := idSlugRegex.FindStringSubmatch(fileName)
	if len(matches) >= 3 {
		id, _ := strconv.Atoi(matches[1])
		return id, matches[2]
	}

	slug := strings.TrimSuffix(fileName, ".md")

	return 0, slug
}

func extractTitleFromLine(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	matches := headingRegex.FindStringSubmatch(trimmed)
	if len(matches) >= 2 {
		return strings.TrimSpace(matches[1]), true
	}

	return "", false
}

func parseStatusToken(token string) PlanStatusType {
	upper := strings.ToUpper(strings.TrimSpace(token))
	if strings.Contains(upper, "COMPLETE") || strings.Contains(upper, "DONE") {
		return StatusCompleted
	}

	if strings.Contains(upper, "PROGRESS") || strings.Contains(upper, "ACTIVE") {
		return StatusInProgress
	}

	return StatusPending
}

func extractStatusFromLine(line string) (PlanStatusType, bool) {
	trimmed := strings.TrimSpace(line)
	matches := statusLineRegex.FindStringSubmatch(trimmed)
	if len(matches) >= 2 {
		return parseStatusToken(matches[1]), true
	}

	return StatusPending, false
}

func formatFallbackTitle(slug string) string {
	words := strings.Split(slug, "-")
	for i, w := range words {
		if len(w) > 0 {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}

	return strings.Join(words, " ")
}
