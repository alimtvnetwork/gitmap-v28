package message

import (
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// stripRules drops every line matching any rule or bot marker, then collapses
// consecutive blank lines and trims trailing whitespace per §6.1 step 1.
func stripRules(msg string, rules []profile.MessageRule) string {
	lines := strings.Split(msg, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if !lineMatches(line, rules) && !isUnwantedBotLine(line) {
			kept = append(kept, line)
		}
	}

	return collapseBlankLines(kept)
}

func isUnwantedBotLine(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	if strings.HasPrefix(lower, "x-lovable") ||
		strings.Contains(lower, "x-lovable-edit-id") ||
		strings.Contains(lower, "lovable-edit-id") ||
		strings.HasPrefix(lower, "co-authored-by: gpt-engineer") ||
		strings.HasPrefix(lower, "signed-off-by: gpt-engineer") {
		return true
	}

	return strings.HasPrefix(lower, "co-authored-by:") && strings.Contains(lower, "noreply.github.com")
}

func lineMatches(line string, rules []profile.MessageRule) bool {
	for _, r := range rules {
		if matchOne(line, r) {
			return true
		}
	}

	return false
}

func matchOne(line string, r profile.MessageRule) bool {
	trimmed := strings.TrimSpace(line)
	kind := strings.ToLower(strings.TrimSpace(r.Kind))
	switch kind {
	case strings.ToLower(constants.CommitInMessageRuleKindStartsWith), "starts_with", "startswith":
		return strings.HasPrefix(line, r.Value) || strings.HasPrefix(trimmed, r.Value) ||
			strings.HasPrefix(strings.ToLower(trimmed), strings.ToLower(r.Value))
	case strings.ToLower(constants.CommitInMessageRuleKindEndsWith), "ends_with", "endswith":
		return strings.HasSuffix(line, r.Value) || strings.HasSuffix(trimmed, r.Value) ||
			strings.HasSuffix(strings.ToLower(trimmed), strings.ToLower(r.Value))
	case strings.ToLower(constants.CommitInMessageRuleKindContains), "contains":
		return strings.Contains(line, r.Value) || strings.Contains(strings.ToLower(line), strings.ToLower(r.Value))
	case "regex", "regexp":
		matched, err := regexp.MatchString(r.Value, line)
		return err == nil && matched
	}

	return false
}

func collapseBlankLines(lines []string) string {
	out := make([]string, 0, len(lines))
	prevBlank := false
	for _, l := range lines {
		isBlank := strings.TrimSpace(l) == ""
		if isBlank && prevBlank {
			continue
		}

		out = append(out, l)
		prevBlank = isBlank
	}

	return strings.TrimRight(strings.Join(out, "\n"), " \t\n")
}
