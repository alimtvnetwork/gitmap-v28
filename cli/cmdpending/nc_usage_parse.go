// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"strings"
)

func parseSubcommandDocLine(line string) (sub, desc string, isSub bool) {
	m := usageSubcommandLine.FindStringSubmatch(line)
	if m == nil {
		return "", "", false
	}
	return m[1], strings.TrimSpace(m[3]), true
}

func isContinuationLine(line, pendingSub string) bool {
	if pendingSub == "" {
		return false
	}
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return false
	}
	// Continuation lines are indented deeper than the 2-space subcommand base.
	return strings.HasPrefix(line, "   ") || strings.HasPrefix(line, "\t")
}

// scanUsageCommandLine extracts `gitmap <parent> <sub>` compounds and
// "(alias: ...)" hints from one usage line.
func scanUsageCommandLine(line string, existing []minedCompound) []minedCompound {
	var out []minedCompound
	loc := usageCommandLine.FindStringSubmatchIndex(line)
	for loc != nil {
		parent := line[loc[2]:loc[3]]
		rest := strings.TrimSpace(line[loc[4]:loc[5]])
		sub := firstBareToken(rest)
		alias := ""
		if m := usageAliasHint.FindStringSubmatch(line); m != nil {
			alias = m[1]
		}
		if sub != "" && !compoundSeen(existing, out, parent, sub) {
			out = append(out, minedCompound{parent: parent, sub: sub, alias: alias})
		} else if sub == "" && alias != "" {
			// Alias-only line (e.g. "(alias: gitmap te all)") — attach to
			// the most recently discovered compound.
			attachAliasHint(existing, out, alias)
		}
		next := usageCommandLine.FindStringSubmatchIndex(line[loc[1]:])
		if next == nil {
			break
		}
		for i := range next {
			next[i] += loc[1]
		}
		loc = next
	}
	return out
}

// usageSubStopWords are English glue words that are never real subcommands;
// they show up when a placeholder (e.g. <target>) precedes them.
var usageSubStopWords = map[string]bool{
	"to": true, "or": true, "and": true, "a": true, "an": true,
	"the": true, "for": true, "with": true, "from": true, "on": true,
	"in": true, "of": true, "as": true, "is": true, "it": true,
	"be": true, "at": true, "by": true, "via": true, "vs": true,
	"gitmap": true,
}

// firstBareToken returns the first plain-word token, skipping placeholders
// like [flags], <repo>, --json, "quoted", and (paren) tokens. When the first
// bare token is an English glue word (placeholder preceded it), there is no
// real subcommand on this line.
func firstBareToken(rest string) string {
	for _, tok := range strings.Fields(rest) {
		if !isBareCommandToken(tok) {
			continue
		}
		if usageSubStopWords[tok] {
			return ""
		}
		return tok
	}
	return ""
}

func isBareCommandToken(tok string) bool {
	if tok == "" {
		return false
	}
	switch tok[0] {
	case '[', '<', '-', '"', '\'', '(', '{', '|':
		return false
	}
	for _, r := range tok {
		if r == '_' || r == '-' || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func compoundSeen(existing, out []minedCompound, parent, sub string) bool {
	for _, c := range existing {
		if c.parent == parent && c.sub == sub {
			return true
		}
	}
	for _, c := range out {
		if c.parent == parent && c.sub == sub {
			return true
		}
	}
	return false
}

func attachAliasHint(existing, out []minedCompound, alias string) {
	if len(out) > 0 {
		if out[len(out)-1].alias == "" {
			out[len(out)-1].alias = alias
		}
		return
	}
	if len(existing) > 0 && existing[len(existing)-1].alias == "" {
		existing[len(existing)-1].alias = alias
	}
}
