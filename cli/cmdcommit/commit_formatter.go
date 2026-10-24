package cmdcommit

import (
	"regexp"
	"strings"
)

// ParsedCommitMessage contains the decomposed components of a semantic commit message.
type ParsedCommitMessage struct {
	RawInput           string
	FormattedMessage   string
	CommitType         string
	Module             string
	ItemDescription    string
	HasModule          bool
	HasHyphenSeparator bool
	IsConventional     bool
}

var (
	// hyphenRegex matches "<module> - <summary>" or "<module>: <summary>"
	hyphenRegex = regexp.MustCompile(`^([A-Za-z0-9_\-\./]+)\s*[:-]\s*(.+)$`)
	// conventionalRegex matches standard conventional commits without brackets: "<type>(<module>): <summary>"
	conventionalRegex = regexp.MustCompile(`^([a-z]+)\(([A-Za-z0-9_\-\./]+)\):\s*(.+)$`)
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
	case "cfr", "refactor", "commit-push-refactor":
		return "Refactor: "
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

	if prefix == "Refactor: " && strings.HasPrefix(strings.ToLower(s), "refactor: ") {
		return strings.TrimSpace(s[10:])
	}

	return s
}

func cleanseBracketArtifacts(s string) string {
	cleaned := strings.ReplaceAll(s, "([", "(")
	cleaned = strings.ReplaceAll(cleaned, "])", ")")
	if !strings.HasPrefix(cleaned, "[") {
		return cleaned
	}

	idx := strings.Index(cleaned, "]")
	if idx == -1 {
		return cleaned
	}

	return strings.TrimSpace(cleaned[1:idx]) + cleaned[idx+1:]
}

func sanitizeInput(verb string, rawArgs []string) string {
	raw := stripWrappingQuotes(strings.Join(rawArgs, " "))
	normalized := trimDuplicatePrefix(verb, raw)

	return cleanseBracketArtifacts(normalized)
}

func buildHyphenRecord(raw, prefix, mod, item string) ParsedCommitMessage {
	upperMod := strings.ToUpper(strings.Trim(mod, "[] "))
	cleanItem := strings.TrimSpace(item)

	return ParsedCommitMessage{
		RawInput:           raw,
		FormattedMessage:   prefix + upperMod + " - " + cleanItem,
		CommitType:         strings.TrimSuffix(prefix, ": "),
		Module:             upperMod,
		ItemDescription:    cleanItem,
		HasModule:          true,
		HasHyphenSeparator: true,
	}
}

func tryParseHyphenSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := hyphenRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}

	record := buildHyphenRecord(s, resolveVerbPrefix(verb), matches[1], matches[2])

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

func buildConventionalRecord(s, ctype, mod, item string) ParsedCommitMessage {
	return ParsedCommitMessage{
		RawInput:           s,
		FormattedMessage:   s,
		CommitType:         ctype,
		Module:             mod,
		ItemDescription:    item,
		HasModule:          mod != "",
		HasHyphenSeparator: false,
		IsConventional:     true,
	}
}

func tryParseConventionalSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := conventionalRegex.FindStringSubmatch(s)
	if len(matches) < 4 || !isKnownConventionalType(matches[1]) {
		return ParsedCommitMessage{}, false
	}

	record := buildConventionalRecord(s, matches[1], matches[2], strings.TrimSpace(matches[3]))

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
	cleansed := sanitizeInput(verb, rawArgs)

	if parsed, ok := tryParseConventionalSyntax(verb, cleansed); ok {
		return parsed
	}

	if parsed, ok := tryParseHyphenSyntax(verb, cleansed); ok {
		return parsed
	}

	return buildFallbackCommitMessage(verb, cleansed)
}

// FormatCommitMessage formats raw commit arguments into a standardized commit message string.
func FormatCommitMessage(verb string, rawArgs []string) string {
	parsed := ParseAndFormatCommitMessage(verb, rawArgs)

	return parsed.FormattedMessage
}
