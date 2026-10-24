package cmdplans

import (
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdplans/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var (
	idSlugRegex     = lazyregex.New(`^(\d+)-(.*)\.md$`)
	statusLineRegex = lazyregex.New(`(?i)(?:##\s*|>\s*\*\*)?Status:\s*\*?\*?\s*([a-zA-Z\s_-]+)`)
	headingRegex    = lazyregex.New(`^#\s+(.*)`)
)

func parseIDAndSlug(fileName string) (int, string) {
	matches := idSlugRegex.FindStringSubmatch(fileName)
	if len(matches) >= 3 {
		id, _ := strconv.Atoi(matches[1])
		return id, matches[2]
	}

	slug := strings.TrimSuffix(fileName, constants.ExtMarkdown)

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

func matchSecondaryStatus(upper string) enums.PlanStatusType {
	if strings.Contains(upper, "ARCHIV") {
		return enums.PlanStatusArchived
	}
	if strings.Contains(upper, "FAIL") {
		return enums.PlanStatusFailed
	}

	return enums.PlanStatusPending
}

func parseStatusToken(token string) enums.PlanStatusType {
	upper := strings.ToUpper(strings.TrimSpace(token))
	if strings.Contains(upper, "COMPLETE") || strings.Contains(upper, "DONE") {
		return enums.PlanStatusCompleted
	}
	if strings.Contains(upper, "PROGRESS") || strings.Contains(upper, "ACTIVE") {
		return enums.PlanStatusInProgress
	}

	return matchSecondaryStatus(upper)
}

func extractStatusFromLine(line string) (enums.PlanStatusType, bool) {
	trimmed := strings.TrimSpace(line)
	matches := statusLineRegex.FindStringSubmatch(trimmed)
	if len(matches) >= 2 {
		return parseStatusToken(matches[1]), true
	}

	return enums.PlanStatusPending, false
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
