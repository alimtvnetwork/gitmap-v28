package cmdfoldertree

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// SearchMode defines the pattern matching strategy for tree searching.
type SearchMode string

const (
	ModeWildcard   SearchMode = "wildcard"
	ModeStartsWith SearchMode = "startsWith"
	ModeContains   SearchMode = "contains"
	ModeEndsWith   SearchMode = "endsWith"
	ModeGrep       SearchMode = "grep"
)

// TreeSearchFilter configures query parameters and search options.
type TreeSearchFilter struct {
	Pattern       string
	SearchMode    SearchMode
	DirsOnly      bool
	FilesOnly     bool
	CaseSensitive bool
	Limit         int
}

// MatchTreeItem tests whether an individual TreeItem satisfies the search mode and pattern.
func MatchTreeItem(item TreeItem, mode SearchMode, pattern string) bool {
	matched, _ := matchTreeItemInternal(item, mode, pattern, false)
	return matched
}

// matchTreeItemInternal evaluates the match and returns compile errors if in grep mode.
func matchTreeItemInternal(item TreeItem, mode SearchMode, pattern string, caseSensitive bool) (bool, error) {
	cleanPat := strings.TrimSpace(pattern)
	if cleanPat == "" {
		return true, nil
	}

	normRel := item.RelPath
	normName := item.FileName
	normPattern := cleanPat

	if !caseSensitive {
		normRel = strings.ToLower(normRel)
		normName = strings.ToLower(normName)
		normPattern = strings.ToLower(normPattern)
	}

	switch mode {
	case ModeWildcard:
		return matchWildcard(item, normPattern, caseSensitive), nil

	case ModeStartsWith:
		// Strip trailing '*' if user typed e.g. "cli*"
		prefix := strings.TrimRight(normPattern, "*")
		prefix = strings.TrimRight(prefix, "/\\")

		// If pattern contains slash (folder path specified), match on RelPath
		if strings.Contains(normPattern, "/") || strings.Contains(normPattern, "\\") {
			normPrefix := normalizeSlash(prefix)
			if strings.HasPrefix(normRel, normPrefix+"/") || normRel == normPrefix {
				return true, nil
			}
			return strings.HasPrefix(normRel, normPrefix), nil
		}

		// Otherwise match prefix of RelPath or FileName
		if strings.HasPrefix(normRel, prefix+"/") || strings.HasPrefix(normRel, prefix) {
			return true, nil
		}
		return strings.HasPrefix(normName, prefix), nil

	case ModeContains:
		// Substring check across RelPath or FileName
		return strings.Contains(normRel, normPattern) || strings.Contains(normName, normPattern), nil

	case ModeEndsWith:
		// Strip leading '*' if user typed e.g. "*.go"
		suffix := strings.TrimLeft(normPattern, "*")
		if strings.HasSuffix(normRel, suffix) || strings.HasSuffix(normName, suffix) {
			return true, nil
		}
		// Match against extension (with or without dot)
		ext := strings.ToLower(item.Extension)
		if strings.EqualFold(ext, suffix) {
			return true, nil
		}
		if !strings.HasPrefix(suffix, ".") && strings.EqualFold(ext, "."+suffix) {
			return true, nil
		}
		return false, nil

	case ModeGrep:
		regexPattern := cleanPat
		if !caseSensitive && !strings.HasPrefix(regexPattern, "(?i)") {
			regexPattern = "(?i)" + regexPattern
		}
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return false, fmt.Errorf("invalid regex pattern %q: %w", pattern, err)
		}
		return re.MatchString(item.RelPath) || re.MatchString(item.FileName), nil

	default:
		// Default to wildcard
		return matchWildcard(item, normPattern, caseSensitive), nil
	}
}

// matchWildcard evaluates wildcard glob patterns (* and ?) against item fields.
func matchWildcard(item TreeItem, normPattern string, caseSensitive bool) bool {
	normRel := item.RelPath
	normName := item.FileName
	if !caseSensitive {
		normRel = strings.ToLower(normRel)
		normName = strings.ToLower(normName)
		normPattern = strings.ToLower(normPattern)
	}

	rePattern := wildcardToRegex(normPattern)
	re, err := regexp.Compile(rePattern)
	if err != nil {
		return strings.Contains(normRel, normPattern)
	}

	if strings.Contains(normPattern, "/") {
		return re.MatchString(normRel)
	}
	return re.MatchString(normName) || re.MatchString(normRel)
}

// wildcardToRegex converts a glob pattern (*, ?) to an anchored regular expression.
func wildcardToRegex(pattern string) string {
	var sb strings.Builder
	sb.WriteString("^")
	for i := 0; i < len(pattern); i++ {
		c := pattern[i]
		switch c {
		case '*':
			sb.WriteString(".*")
		case '?':
			sb.WriteString(".")
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			sb.WriteByte('\\')
			sb.WriteByte(c)
		default:
			sb.WriteByte(c)
		}
	}
	sb.WriteString("$")
	return sb.String()
}

// FilterTreeItems applies the search filter to a list of TreeItem entries.
func FilterTreeItems(items []TreeItem, filter TreeSearchFilter) ([]TreeItem, error) {
	// Pre-validate regex if grep mode
	var grepRE *regexp.Regexp
	if filter.SearchMode == ModeGrep && strings.TrimSpace(filter.Pattern) != "" {
		pat := strings.TrimSpace(filter.Pattern)
		if !filter.CaseSensitive && !strings.HasPrefix(pat, "(?i)") {
			pat = "(?i)" + pat
		}
		re, err := regexp.Compile(pat)
		if err != nil {
			return nil, fmt.Errorf("invalid regex %q: %w", filter.Pattern, err)
		}
		grepRE = re
	}

	matched := make([]TreeItem, 0, len(items))

	for _, item := range items {
		// Type filter
		if filter.DirsOnly && !item.IsDir {
			continue
		}
		if filter.FilesOnly && item.IsDir {
			continue
		}

		if strings.TrimSpace(filter.Pattern) == "" {
			matched = append(matched, item)
			if filter.Limit > 0 && len(matched) >= filter.Limit {
				break
			}
			continue
		}

		var isMatch bool
		if grepRE != nil {
			isMatch = grepRE.MatchString(item.RelPath) || grepRE.MatchString(item.FileName)
		} else {
			isMatch = MatchTreeItem(item, filter.SearchMode, filter.Pattern)
		}

		if isMatch {
			matched = append(matched, item)
			if filter.Limit > 0 && len(matched) >= filter.Limit {
				break
			}
		}
	}

	return matched, nil
}

// FormatTreeItemsText renders matched items in clean ANSI/text format.
func FormatTreeItemsText(items []TreeItem, queryDesc string) string {
	if len(items) == 0 {
		return fmt.Sprintf("No matches found for %s.\n", queryDesc)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d match(es) for %s:\n\n", len(items), queryDesc))

	for _, it := range items {
		if it.IsDir {
			sb.WriteString(fmt.Sprintf("  [DIR]  %s/\n", it.RelPath))
		} else {
			sizeStr := ""
			if it.SizeBytes > 0 {
				sizeStr = fmt.Sprintf(" (%s)", formatItemBytes(it.SizeBytes))
			}
			sb.WriteString(fmt.Sprintf("  [FILE] %s%s\n", it.RelPath, sizeStr))
		}
	}

	return sb.String()
}

// FormatTreeItemsJSON renders items as a pretty-printed JSON array.
func FormatTreeItemsJSON(items []TreeItem) (string, error) {
	b, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func formatItemBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
