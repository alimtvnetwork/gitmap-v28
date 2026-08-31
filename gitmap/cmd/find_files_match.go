package cmd

import (
	"path/filepath"
	"strings"
)

func parseExtensionList(raw string) []string {
	parts := strings.Split(raw, ",")
	var exts []string
	for _, p := range parts {
		clean := strings.TrimPrefix(strings.TrimSpace(strings.ToLower(p)), ".")
		if clean != "" {
			exts = append(exts, clean)
		}
	}
	return exts
}

func isFileMatching(opts FindFilesOptions, baseName, relPath string) bool {
	if len(opts.Exts) > 0 && !hasAllowedExtension(baseName, opts.Exts) {
		return false
	}
	if opts.Query == "" {
		return true
	}
	return evaluateMatchKind(opts.Kind, opts.Query, baseName, relPath)
}

func hasAllowedExtension(filename string, allowedExts []string) bool {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), ".")
	for _, a := range allowedExts {
		if a == ext {
			return true
		}
	}
	return false
}

func evaluateMatchKind(kind MatchKind, query, baseName, relPath string) bool {
	queryLower := strings.ToLower(query)
	baseLower := strings.ToLower(baseName)
	relLower := strings.ToLower(relPath)

	switch kind {
	case MatchExact:
		return baseLower == queryLower || relLower == queryLower
	case MatchContains:
		return strings.Contains(baseLower, queryLower) || strings.Contains(relLower, queryLower)
	case MatchStartsWith:
		return strings.HasPrefix(baseLower, queryLower) || strings.HasPrefix(relLower, queryLower)
	case MatchEndsWith:
		return strings.HasSuffix(baseLower, queryLower) || strings.HasSuffix(relLower, queryLower)
	case MatchWildcard:
		return matchWildcardPattern(queryLower, baseLower, relLower)
	default:
		return strings.Contains(baseLower, queryLower)
	}
}

func matchWildcardPattern(query, baseName, relPath string) bool {
	if matched, _ := filepath.Match(query, baseName); matched {
		return true
	}
	if matched, _ := filepath.Match(query, relPath); matched {
		return true
	}
	hasLeading := strings.HasPrefix(query, "*")
	hasTrailing := strings.HasSuffix(query, "*")
	clean := strings.Trim(query, "*")

	if hasLeading && hasTrailing {
		return strings.Contains(baseName, clean) || strings.Contains(relPath, clean)
	}
	if hasLeading {
		return strings.HasSuffix(baseName, clean) || strings.HasSuffix(relPath, clean)
	}
	if hasTrailing {
		return strings.HasPrefix(baseName, clean) || strings.HasPrefix(relPath, clean)
	}
	return strings.Contains(baseName, query) || strings.Contains(relPath, query)
}
