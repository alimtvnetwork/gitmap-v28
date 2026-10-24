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
	// bracketedModuleRegex matches "(<module>): <summary>" or "(<module>) - <summary>"
	bracketedModuleRegex = regexp.MustCompile(`^\(([A-Za-z0-9_\-\./]+)\)\s*[:-]\s*(.+)$`)
	// conventionalRegex matches standard conventional commits: "<type>(<module>): <summary>"
	conventionalRegex = regexp.MustCompile(`^([A-Za-z]+)\(([A-Za-z0-9_\-\./]+)\):\s*(.+)$`)
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

func resolveUppercaseType(verb string) string {
	mapping := map[string]string{
		"cpf": "FEAT", "feature": "FEAT", "commit-push-feature": "FEAT", "feat": "FEAT",
		"cpb": "FIX", "cbf": "FIX", "bug": "FIX", "commit-push-bug": "FIX", "fix": "FIX",
		"cpc": "CHORE", "chore": "CHORE", "commit-push-chore": "CHORE",
		"cfr": "REFACTOR", "refactor": "REFACTOR", "commit-push-refactor": "REFACTOR",
		"cpr": "RELEASE", "release": "RELEASE", "commit-push-release": "RELEASE",
		"test": "TEST", "docs": "DOCS", "perf": "PERF",
	}
	if val, hasVal := mapping[strings.ToLower(verb)]; hasVal {
		return val
	}

	return strings.ToUpper(verb)
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

func buildConventionalRecord(raw, ctype, mod, item string) ParsedCommitMessage {
	uType := resolveUppercaseType(ctype)
	uMod := strings.ToUpper(strings.Trim(mod, "() "))
	cleanItem := strings.TrimSpace(item)
	formatted := uType + "(" + uMod + "): " + cleanItem

	return ParsedCommitMessage{
		RawInput:           raw,
		FormattedMessage:   formatted,
		CommitType:         uType,
		Module:             uMod,
		ItemDescription:    cleanItem,
		HasModule:          uMod != "",
		HasHyphenSeparator: false,
		IsConventional:     true,
	}
}

func tryParseBracketedModuleSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := bracketedModuleRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}

	record := buildConventionalRecord(s, verb, matches[1], matches[2])

	return record, true
}

func tryParseConventionalSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := conventionalRegex.FindStringSubmatch(s)
	if len(matches) < 4 || !isKnownConventionalType(strings.ToLower(matches[1])) {
		return ParsedCommitMessage{}, false
	}

	record := buildConventionalRecord(s, matches[1], matches[2], matches[3])

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

	if parsed, ok := tryParseBracketedModuleSyntax(verb, cleansed); ok {
		return parsed
	}

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
