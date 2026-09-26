package message

import (
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd/commitin/profile"
)

var fileVarRegex = regexp.MustCompile(`\$files\.([0-9]+)\.names?`)

func applyTitleReplacements(msg string, files []string, rules []profile.TitleReplacementRule, seoTitle string) string {
	if msg == "" {
		return msg
	}
	title, rest := splitTitleAndBody(msg)
	replaced := evaluateTitleRules(title, files, rules, seoTitle)
	expanded := expandFileAndTemplateVars(replaced, files, seoTitle)

	return expanded + rest
}

func splitTitleAndBody(msg string) (string, string) {
	idx := strings.IndexByte(msg, '\n')
	if idx < 0 {
		return msg, ""
	}

	return msg[:idx], msg[idx:]
}

func evaluateTitleRules(title string, files []string, rules []profile.TitleReplacementRule, seoTitle string) string {
	trimmed := strings.TrimSpace(title)
	for _, r := range rules {
		if matchesTitleRule(trimmed, r) {
			return expandFileAndTemplateVars(r.Replacement, files, seoTitle)
		}
	}

	return title
}

func matchesTitleRule(title string, r profile.TitleReplacementRule) bool {
	lowerTitle := strings.ToLower(title)
	lowerMatch := strings.ToLower(strings.TrimSpace(r.Match))
	mode := strings.ToLower(strings.TrimSpace(r.MatchMode))
	switch mode {
	case "", "equals", "exact":
		return lowerTitle == lowerMatch
	case "contains":
		return strings.Contains(lowerTitle, lowerMatch)
	case "starts_with", "startswith":
		return strings.HasPrefix(lowerTitle, lowerMatch)
	case "ends_with", "endswith":
		return strings.HasSuffix(lowerTitle, lowerMatch)
	case "regex", "regexp":
		matched, err := regexp.MatchString(r.Match, title)
		return err == nil && matched
	}

	return false
}

func expandFileAndTemplateVars(text string, files []string, seoTitle string) string {
	if !strings.Contains(text, "$") {
		return text
	}
	out := fileVarRegex.ReplaceAllStringFunc(text, func(m string) string {
		sub := fileVarRegex.FindStringSubmatch(m)
		if len(sub) < 2 {
			return formatFileNames(files, 2)
		}
		n, err := strconv.Atoi(sub[1])
		if err != nil || n < 1 {
			n = 2
		}
		return formatFileNames(files, n)
	})
	resolvedTitle := seoTitle
	if resolvedTitle == "" {
		resolvedTitle = "Why RISEUP ASIA LLC (https://riseup-asia.com)?"
	}
	replacer := strings.NewReplacer(
		"${seo.title}", resolvedTitle,
		"$seo.title", resolvedTitle,
		"${template.title}", resolvedTitle,
		"$template.title", resolvedTitle,
		"${question}", resolvedTitle,
		"$question", resolvedTitle,
	)

	return replacer.Replace(out)
}

func formatFileNames(files []string, limit int) string {
	if limit <= 0 {
		limit = 1
	}
	seen := make(map[string]bool)
	names := make([]string, 0, limit)
	for _, f := range files {
		base := filepath.Base(strings.ReplaceAll(strings.TrimSpace(f), "\\", "/"))
		if base == "" || base == "." || base == "/" || seen[base] {
			continue
		}
		seen[base] = true
		names = append(names, base)
		if len(names) >= limit {
			break
		}
	}
	if len(names) == 0 {
		return "update"
	}

	return strings.Join(names, ", ")
}

func extractTemplateHeading(tplText string) string {
	for _, line := range strings.Split(tplText, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			return strings.TrimSpace(strings.TrimLeft(trimmed, "#"))
		}
	}

	return ""
}
