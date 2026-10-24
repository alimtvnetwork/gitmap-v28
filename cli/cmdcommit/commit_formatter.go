package cmdcommit

import (
	"regexp"
	"strings"
)

// ParsedCommitMessage contains the decomposed components of a semantic commit message.
type ParsedCommitMessage struct {
	RawInput         string
	FormattedMessage string
	CommitType       string
	Module           string
	ItemDescription  string
	HasModule        bool
	HasBracketSyntax bool
	IsConventional   bool
}

var (
	bracketRegex      = regexp.MustCompile(`^\[([A-Za-z0-9_\-\./]+)\]:\s*(.+)$`)
	conventionalRegex = regexp.MustCompile(`^([a-z]+)(?:\((?:\[([A-Za-z0-9_\-\./]+)\]|([A-Za-z0-9_\-\./]+))\))?:\s*(.+)$`)
	legacyHyphenRegex = regexp.MustCompile(`^([A-Za-z0-9_\-\./]+)\s+-\s+(.+)$`)
)

func stripWrappingQuotes(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) >= 2 && ((trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"') ||
		(trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'')) {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	return trimmed
}

func resolveVerbPrefix(verb string) string {
	switch strings.ToLower(verb) {
	case "cpf", "feature", "commit-push-feature":
		return "Feature: "
	case "cpb", "cbf", "bug", "commit-push-bug":
		return "Bug: "
	case "cpc", "chore", "commit-push-chore":
		return "Chore: "
	case "cpr", "release", "commit-push-release":
		return "Release: "
	default:
		return ""
	}
}

func trimDuplicatePrefix(verb, s string) string {
	prefix := resolveVerbPrefix(verb)
	if prefix != "" && strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix)) {
		return strings.TrimSpace(s[len(prefix):])
	}

	if prefix == "Bug: " && strings.HasPrefix(strings.ToLower(s), "fix: ") {
		return strings.TrimSpace(s[5:])
	}

	return s
}

func buildBracketRecord(raw, prefix, mod, item string) ParsedCommitMessage {
	return ParsedCommitMessage{
		RawInput:         raw,
		FormattedMessage: prefix + "[" + mod + "]: " + item,
		CommitType:       strings.TrimSuffix(prefix, ": "),
		Module:           mod,
		ItemDescription:  item,
		HasModule:        true,
		HasBracketSyntax: true,
	}
}

func tryParseBracketSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := bracketRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}

	record := buildBracketRecord(s, resolveVerbPrefix(verb), matches[1], strings.TrimSpace(matches[2]))

	return record, true
}

func isKnownConventionalType(t string) bool {
	switch t {
	case "feat", "fix", "chore", "release", "refactor", "test", "docs", "perf":
		return true
	default:
		return false
	}
}

func extractConventionalModule(matches []string) (string, bool) {
	if matches[2] != "" {
		return matches[2], true
	}

	return matches[3], false
}

func buildConventionalRecord(s, ctype, mod, item string, hasBracket bool) ParsedCommitMessage {
	return ParsedCommitMessage{
		RawInput:         s,
		FormattedMessage: s,
		CommitType:       ctype,
		Module:           mod,
		ItemDescription:  item,
		HasModule:        mod != "",
		HasBracketSyntax: hasBracket,
		IsConventional:   true,
	}
}

func tryParseConventionalSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := conventionalRegex.FindStringSubmatch(s)
	if len(matches) < 5 || !isKnownConventionalType(matches[1]) {
		return ParsedCommitMessage{}, false
	}

	mod, hasBracket := extractConventionalModule(matches)
	record := buildConventionalRecord(s, matches[1], mod, strings.TrimSpace(matches[4]), hasBracket)

	return record, true
}

func tryParseLegacyHyphenSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := legacyHyphenRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}

	record := buildBracketRecord(s, resolveVerbPrefix(verb), matches[1], strings.TrimSpace(matches[2]))

	return record, true
}

func buildFallbackCommitMessage(verb, s string) ParsedCommitMessage {
	prefix := resolveVerbPrefix(verb)
	formatted := prefix + s

	return ParsedCommitMessage{
		RawInput:         s,
		FormattedMessage: formatted,
		CommitType:       strings.TrimSuffix(prefix, ": "),
		ItemDescription:  s,
	}
}

// ParseAndFormatCommitMessage processes raw arguments according to the target semantic verb.
func ParseAndFormatCommitMessage(verb string, rawArgs []string) ParsedCommitMessage {
	cleanMessage := strings.TrimSpace(strings.Join(rawArgs, " "))
	stripped := stripWrappingQuotes(cleanMessage)
	normalized := trimDuplicatePrefix(verb, stripped)

	if parsed, ok := tryParseBracketSyntax(verb, normalized); ok {
		return parsed
	}

	if parsed, ok := tryParseConventionalSyntax(verb, normalized); ok {
		return parsed
	}

	if parsed, ok := tryParseLegacyHyphenSyntax(verb, normalized); ok {
		return parsed
	}

	return buildFallbackCommitMessage(verb, normalized)
}

// FormatCommitMessage formats raw commit arguments into a standardized commit message string.
func FormatCommitMessage(verb string, rawArgs []string) string {
	parsed := ParseAndFormatCommitMessage(verb, rawArgs)

	return parsed.FormattedMessage
}
