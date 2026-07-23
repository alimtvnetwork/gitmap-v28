package message

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd/commitin/profile"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// stripRules drops every line matching any rule, then collapses
// consecutive blank lines and trims trailing whitespace per §6.1 step 1.
func stripRules(msg string, rules []profile.MessageRule) string {
	if len(rules) == 0 {
		return strings.TrimRight(msg, " \t\n")
	}
	lines := strings.Split(msg, "\n")
	kept := lines[:0]
	for _, line := range lines {
		if !lineMatches(line, rules) {
			kept = append(kept, line)
		}
	}
	return collapseBlankLines(kept)
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
	switch r.Kind {
	case constants.CommitInMessageRuleKindStartsWith:
		return strings.HasPrefix(line, r.Value)
	case constants.CommitInMessageRuleKindEndsWith:
		return strings.HasSuffix(line, r.Value)
	case constants.CommitInMessageRuleKindContains:
		return strings.Contains(line, r.Value)
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
