package folder

import (
	"path/filepath"
	"strings"
)

// FilterConfig holds patterns and predicates for scanning and filtering files.
type FilterConfig struct {
	ExceptGlobs []string
	Extensions  []string
	MaxDepth    int
	OnlyText    bool
	OnlyBinary  bool
}

// ParseExceptGlobs parses comma-separated exclusion strings into a slice of glob patterns.
func ParseExceptGlobs(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var globs []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			globs = append(globs, strings.ReplaceAll(trimmed, "\\", "/"))
		}
	}
	return globs
}

// IsPathExcluded tests whether a path or directory matches exclusion patterns.
func (fc *FilterConfig) IsPathExcluded(relPath string, isDir bool) bool {
	norm := strings.ReplaceAll(relPath, "\\", "/")
	base := filepath.Base(norm)

	if isDefaultIgnored(base) {
		return true
	}

	for _, g := range fc.ExceptGlobs {
		if matchGlob(g, norm, base, isDir) {
			return true
		}
	}
	return false
}

func isDefaultIgnored(base string) bool {
	return base == ".git"
}

func matchGlob(pattern, normPath, base string, isDir bool) bool {
	cleanPattern := strings.TrimSuffix(pattern, "/*")
	cleanPattern = strings.TrimSuffix(cleanPattern, "/**")
	cleanPattern = strings.TrimSuffix(cleanPattern, "/")

	if strings.HasPrefix(pattern, "*.") {
		return strings.HasSuffix(base, pattern[1:])
	}

	if matched, _ := filepath.Match(pattern, base); matched {
		return true
	}
	if matched, _ := filepath.Match(pattern, normPath); matched {
		return true
	}

	isPathMatch := (strings.Contains(pattern, "/") || strings.Contains(pattern, "**")) &&
		(strings.HasPrefix(normPath, cleanPattern+"/") || normPath == cleanPattern)
	if isPathMatch {
		return true
	}

	return false
}

// IsMetaAllowed validates whether extracted file metadata meets filter criteria.
func (fc *FilterConfig) IsMetaAllowed(meta *FileMeta) bool {
	if fc.OnlyText && meta.IsBinary {
		return false
	}
	if fc.OnlyBinary && !meta.IsBinary {
		return false
	}
	if len(fc.Extensions) > 0 {
		return hasMatchingExtension(meta.Extension, fc.Extensions)
	}
	return true
}

func hasMatchingExtension(ext string, allowed []string) bool {
	cleanExt := strings.TrimPrefix(strings.ToLower(ext), ".")
	for _, a := range allowed {
		if strings.TrimPrefix(strings.ToLower(a), ".") == cleanExt {
			return true
		}
	}
	return false
}
