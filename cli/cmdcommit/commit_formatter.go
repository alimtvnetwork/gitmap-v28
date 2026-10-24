package cmdcommit

import (
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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
	bracketRegex      = regexp.MustCompile(`^(?:\[|\(\[|\()([A-Za-z0-9_\-\./]+)(?:\]|\)\]|\]\)|\)):\s*(.+)$`)
	conventionalRegex = regexp.MustCompile(`(?i)^([a-z]+)(?:\((?:\[([A-Za-z0-9_\-\./]+)\]|([A-Za-z0-9_\-\./]+))\))?:\s*(.+)$`)
	legacyHyphenRegex = regexp.MustCompile(`^([A-Za-z0-9_\-\./]+)\s+-\s+(.+)$`)
)

func stripWrappingQuotes(s string) string {
	trimmed := strings.TrimSpace(s)
	if len(trimmed) < 2 {
		return trimmed
	}
	isDouble := trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"'
	isSingle := trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\''
	if isDouble || isSingle {
		return strings.TrimSpace(trimmed[1 : len(trimmed)-1])
	}

	return trimmed
}

func resolveCommitType(verb string) string {
	switch strings.ToLower(verb) {
	case "cpf", "feature", "feat", "commit-push-feature":
		return string(constants.ConventionalFeat)
	case "cpb", "cbf", "bug", "fix", "commit-push-bug", "commit-push-fix":
		return string(constants.ConventionalFix)
	case "cfr", "refactor", "commit-push-refactor":
		return string(constants.ConventionalRefactor)
	case "cpc", "chore", "commit-push-chore":
		return string(constants.ConventionalChore)
	case "cpr", "release", "commit-push-release":
		return string(constants.ConventionalRelease)
	default:
		return resolveOtherCommitType(verb)
	}
}

func resolveOtherCommitType(verb string) string {
	switch strings.ToLower(verb) {
	case "docs":
		return "DOCS"
	case "test":
		return "TEST"
	case "perf":
		return "PERF"
	default:
		return ""
	}
}

func resolveLegacyPrefix(ctype string) string {
	switch ctype {
	case "FEAT":
		return "Feature: "
	case "FIX":
		return "Bug: "
	case "CHORE":
		return "Chore: "
	case "RELEASE":
		return "Release: "
	case "REFACTOR":
		return "Refactor: "
	default:
		return ""
	}
}

func trimDuplicatePrefix(verb, s string) string {
	ctype := resolveCommitType(verb)
	prefix := resolveLegacyPrefix(ctype)
	if prefix != "" && strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix)) {
		return strings.TrimSpace(s[len(prefix):])
	}
	if ctype != "" && strings.HasPrefix(strings.ToLower(s), strings.ToLower(ctype)+": ") {
		return strings.TrimSpace(s[len(ctype)+2:])
	}
	if ctype == "FIX" && strings.HasPrefix(strings.ToLower(s), "fix: ") {
		return strings.TrimSpace(s[5:])
	}

	return s
}

func cleanModule(mod string) string {
	trimmed := strings.Trim(mod, "[]()")

	return strings.ToUpper(strings.TrimSpace(trimmed))
}

func sanitizeDescription(desc string) string {
	s := strings.TrimSpace(desc)
	s = strings.ReplaceAll(s, "[", "(")
	s = strings.ReplaceAll(s, "]", ")")

	return s
}

func buildFormattedMessage(ctype, mod, item string) string {
	if mod != "" && ctype != "" {
		return ctype + "(" + mod + "): " + item
	}
	if ctype != "" {
		return ctype + ": " + item
	}
	if mod != "" {
		return "(" + mod + "): " + item
	}

	return item
}

func determineCommitType(verb, matched string) string {
	verbType := resolveCommitType(verb)
	if verbType != "" {
		return verbType
	}
	matchedType := resolveCommitType(matched)
	if matchedType != "" {
		return matchedType
	}

	return string(constants.ConventionalFeat)
}

func buildBracketRecord(raw, ctype, mod, item string) ParsedCommitMessage {
	cleanMod := cleanModule(mod)
	sanitized := sanitizeDescription(item)
	formatted := buildFormattedMessage(ctype, cleanMod, sanitized)

	return ParsedCommitMessage{
		RawInput:         raw,
		FormattedMessage: formatted,
		CommitType:       ctype,
		Module:           cleanMod,
		ItemDescription:  sanitized,
		HasModule:        cleanMod != "",
		HasBracketSyntax: true,
		IsConventional:   true,
	}
}

func tryParseBracketSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := bracketRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}
	ctype := determineCommitType(verb, "")
	record := buildBracketRecord(s, ctype, matches[1], strings.TrimSpace(matches[2]))

	return record, true
}

func isKnownConventionalType(t string) bool {
	switch strings.ToUpper(t) {
	case string(constants.ConventionalFeat), string(constants.ConventionalFix),
		string(constants.ConventionalChore), string(constants.ConventionalRelease),
		string(constants.ConventionalRefactor), "TEST", "DOCS", "PERF":
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

func buildConventionalRecord(raw, ctype, mod, item string, hasBracket bool) ParsedCommitMessage {
	cleanMod := cleanModule(mod)
	sanitized := sanitizeDescription(item)
	formatted := buildFormattedMessage(ctype, cleanMod, sanitized)

	return ParsedCommitMessage{
		RawInput:         raw,
		FormattedMessage: formatted,
		CommitType:       ctype,
		Module:           cleanMod,
		ItemDescription:  sanitized,
		HasModule:        cleanMod != "",
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
	ctype := determineCommitType(verb, matches[1])
	record := buildConventionalRecord(s, ctype, mod, strings.TrimSpace(matches[4]), hasBracket)

	return record, true
}

func tryParseLegacyHyphenSyntax(verb, s string) (ParsedCommitMessage, bool) {
	matches := legacyHyphenRegex.FindStringSubmatch(s)
	if len(matches) < 3 {
		return ParsedCommitMessage{}, false
	}
	ctype := determineCommitType(verb, "")
	record := buildBracketRecord(s, ctype, matches[1], strings.TrimSpace(matches[2]))

	return record, true
}

func buildFallbackCommitMessage(verb, s string) ParsedCommitMessage {
	ctype := resolveCommitType(verb)
	sanitized := sanitizeDescription(s)
	formatted := buildFormattedMessage(ctype, "", sanitized)

	return ParsedCommitMessage{
		RawInput:         s,
		FormattedMessage: formatted,
		CommitType:       ctype,
		ItemDescription:  sanitized,
	}
}

// ParseAndFormatCommitMessage processes raw arguments according to the target semantic verb.
func ParseAndFormatCommitMessage(verb string, rawArgs []string) ParsedCommitMessage {
	cleanMessage := strings.TrimSpace(strings.Join(rawArgs, " "))
	stripped := stripWrappingQuotes(cleanMessage)
	normalized := trimDuplicatePrefix(verb, stripped)

	if parsed, ok := tryParseConventionalSyntax(verb, normalized); ok {
		return parsed
	}
	if parsed, ok := tryParseBracketSyntax(verb, normalized); ok {
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
