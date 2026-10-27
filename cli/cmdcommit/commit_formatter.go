package cmdcommit

import (
	"strings"
	"unicode"

	commitconstants "github.com/alimtvnetwork/gitmap-v28/cli/cmdcommit/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
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
	hyphenRegex = lazyregex.New(commitconstants.PatternHyphen)
	// bracketedModuleRegex matches "(<module>): <summary>" or "(<module>) - <summary>"
	bracketedModuleRegex = lazyregex.New(commitconstants.PatternBracketedModule)
	// conventionalRegex matches standard conventional commits: "<type>(<module>): <summary>"
	conventionalRegex = lazyregex.New(commitconstants.PatternConventional)
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
	case "cpb", "cbf", "bug", "fix", "commit-push-bug":
		return "Bug: "
	case "cfr", "refactor", "commit-push-refactor":
		return "Refactor: "
	case "cpc", "chore", "commit-push-chore":
		return "Chore: "
	case "cpr", "release", "commit-push-release":
		return "Release: "
	}

	return ""
}

func resolveConventionalType(verb string) string {
	mapping := map[string]string{
		"cpf": "feat", "feature": "feat", "commit-push-feature": "feat", "feat": "feat",
		"cpb": "fix", "cbf": "fix", "bug": "fix", "commit-push-bug": "fix", "fix": "fix",
		"cpc": "chore", "chore": "chore", "commit-push-chore": "chore",
		"cfr": "refactor", "refactor": "refactor", "commit-push-refactor": "refactor",
		"cpr": "release", "release": "release", "commit-push-release": "release",
		"test": "test", "docs": "docs", "perf": "perf",
	}
	if val, hasVal := mapping[strings.ToLower(verb)]; hasVal {
		return val
	}

	return strings.ToLower(verb)
}

func isKnownCommitTag(tag string) bool {
	switch tag {
	case "feature", "feat", "release", "rel", "bug", "fix", "refactor", "refact", "chore", "docs", "spec", "test", "perf", "commit", "commits", "gitcommit":
		return true
	default:
		return false
	}
}

func stripSinglePrefix(s string) (string, bool) {
	colonIdx := strings.Index(s, ":")
	if colonIdx == -1 {
		return s, false
	}

	tag := strings.ToLower(strings.TrimSpace(s[:colonIdx]))
	if !isKnownCommitTag(tag) {
		return s, false
	}

	return strings.TrimSpace(s[colonIdx+1:]), true
}

func stripRedundantCommitPrefixes(s string) string {
	curr := s
	for {
		stripped, hasPrefix := stripSinglePrefix(curr)
		if !hasPrefix {
			return curr
		}

		curr = stripped
	}
}

func trimDuplicatePrefix(verb, s string) string {
	prefix := resolveVerbPrefix(verb)
	if prefix != "" && strings.HasPrefix(strings.ToLower(s), strings.ToLower(prefix)) {
		return strings.TrimSpace(s[len(prefix):])
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
	stripped := stripRedundantCommitPrefixes(raw)
	normalized := trimDuplicatePrefix(verb, stripped)

	return cleanseBracketArtifacts(normalized)
}

func isCommitTypeToken(mod string) bool {
	switch strings.ToUpper(strings.TrimSpace(mod)) {
	case "FEATURE", "FEAT", "RELEASE", "REL", "VERSION", "BUG", "CHORE", "REFACTOR", "REFACT", "SPEC", "COMMIT", "COMMITS", "GITCOMMIT":
		return true
	default:
		return false
	}
}

func isAllUpperShoutCase(s string) bool {
	hasLetters := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetters = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}

	return hasLetters
}

func normalizeModuleCasing(mod string) string {
	trimmed := strings.TrimSpace(mod)
	if isAllUpperShoutCase(trimmed) {
		return strings.ToLower(trimmed)
	}

	return trimmed
}

func isRedundantModule(prefix, mod string) bool {
	upperMod := strings.ToUpper(strings.Trim(mod, "[] "))
	if isCommitTypeToken(upperMod) {
		return true
	}

	prefixType := strings.ToUpper(strings.TrimSuffix(strings.TrimSpace(prefix), ":"))
	if prefixType != "" && prefixType == upperMod {
		return true
	}

	return false
}

func splitNestedHyphen(item string) (string, string, bool) {
	idx := strings.Index(item, " - ")
	if idx == -1 {
		return "", item, false
	}

	nestedMod := strings.TrimSpace(item[:idx])
	nestedItem := strings.TrimSpace(item[idx+3:])
	if nestedMod == "" || nestedItem == "" {
		return "", item, false
	}

	return nestedMod, nestedItem, true
}

func resolveHyphenModuleAndItem(prefix, mod, item string) (string, string, bool) {
	candidateMod := strings.Trim(mod, "[] ")
	upperMod := strings.ToUpper(candidateMod)
	cleanItem := strings.TrimSpace(item)

	for isRedundantModule(prefix, upperMod) {
		nestedMod, nestedItem, hasNested := splitNestedHyphen(cleanItem)
		if !hasNested {
			return "", cleanItem, false
		}

		candidateMod = strings.Trim(nestedMod, "[] ")
		upperMod = strings.ToUpper(candidateMod)
		cleanItem = strings.TrimSpace(nestedItem)
	}

	return normalizeModuleCasing(candidateMod), cleanItem, true
}

func makeHyphenRecord(raw, prefix, mod, item string, hasMod bool) ParsedCommitMessage {
	msg := prefix + mod + " - " + item
	if !hasMod {
		msg = prefix + item
	}

	return ParsedCommitMessage{
		RawInput:           raw,
		FormattedMessage:   msg,
		CommitType:         strings.TrimSuffix(prefix, ": "),
		Module:             mod,
		ItemDescription:    item,
		HasModule:          hasMod,
		HasHyphenSeparator: true,
	}
}

func buildHyphenRecord(raw, prefix, mod, item string) ParsedCommitMessage {
	finalMod, cleanItem, hasMod := resolveHyphenModuleAndItem(prefix, mod, item)

	return makeHyphenRecord(raw, prefix, finalMod, cleanItem, hasMod)
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
	cType := resolveConventionalType(ctype)
	cMod := normalizeModuleCasing(strings.Trim(mod, "() "))
	cleanItem := strings.TrimSpace(item)
	hasMod := cMod != "" && !isCommitTypeToken(cMod)
	formatted := cType + ": " + cleanItem
	if hasMod {
		formatted = cType + "(" + cMod + "): " + cleanItem
	} else {
		cMod = ""
	}

	return ParsedCommitMessage{
		RawInput:           raw,
		FormattedMessage:   formatted,
		CommitType:         cType,
		Module:             cMod,
		ItemDescription:    cleanItem,
		HasModule:          hasMod,
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
