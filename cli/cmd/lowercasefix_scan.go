// Package cmd — lowercasefix_scan.go handles repository file scanning for uppercase files.
package cmd

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

func findRenameCandidates(root string, opts LowerCaseFixOptions) ([]RenamePair, int, error) {
	var pairs []RenamePair
	scannedCount := 0
	err := filepath.Walk(root, func(filePath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return skipOrContinueDir(info.Name())
		}
		scannedCount++
		rel, _ := filepath.Rel(root, filePath)
		if pair, isMatch := checkFileCandidate(filePath, rel, info.Name(), opts); isMatch {
			pairs = append(pairs, pair)
		}
		return nil
	})
	return pairs, scannedCount, err
}

func skipOrContinueDir(name string) error {
	if isIgnoredScanDir(name) {
		return filepath.SkipDir
	}
	return nil
}

func isIgnoredScanDir(name string) bool {
	low := strings.ToLower(name)
	return low == ".git" || low == "node_modules" || low == "vendor" ||
		low == ".gemini" || low == "dist" || low == "build" || low == "target"
}

func isIgnoredScanFile(base string) bool {
	return strings.HasSuffix(base, ".tmp-lcf")
}

func checkFileCandidate(pathStr, rel, base string, opts LowerCaseFixOptions) (RenamePair, bool) {
	if !hasUppercaseChars(base) || isIgnoredScanFile(base) {
		return RenamePair{}, false
	}
	if opts.IsReadmeOnly {
		return checkReadmeCandidate(pathStr, rel, base)
	}

	return checkPatternCandidate(pathStr, rel, base, opts.Patterns)
}

func checkReadmeCandidate(pathStr, rel, base string) (RenamePair, bool) {
	if !isRootReadme(rel, base) {
		return RenamePair{}, false
	}

	return buildRenamePair(pathStr, rel, base, "root README filter"), true
}

func checkPatternCandidate(pathStr, rel, base string, patterns []string) (RenamePair, bool) {
	matched, pattern := findMatchingPattern(rel, base, patterns)
	if !matched {
		return RenamePair{}, false
	}

	return buildRenamePair(pathStr, rel, base, "filter ["+pattern+"]"), true
}

func buildRenamePair(pathStr, rel, base, matchedBy string) RenamePair {
	lowBase := strings.ToLower(base)
	dir := filepath.Dir(pathStr)

	return RenamePair{
		OldPath:   pathStr,
		NewPath:   filepath.Join(dir, lowBase),
		OldBase:   base,
		NewBase:   lowBase,
		RelPath:   rel,
		MatchedBy: matchedBy,
	}
}

func isRootReadme(rel, base string) bool {
	dir := filepath.Dir(rel)
	if dir != "." && dir != "" {
		return false
	}
	low := strings.ToLower(base)
	return strings.HasPrefix(low, "readme")
}

func hasUppercaseChars(s string) bool {
	return strings.ToLower(s) != s
}

func findMatchingPattern(rel, base string, patterns []string) (bool, string) {
	if len(patterns) == 0 {
		return true, "*"
	}
	for _, p := range patterns {
		if isPatternMatch(rel, base, p) {
			return true, p
		}
	}

	return false, ""
}

func isPatternMatch(rel, base, pattern string) bool {
	p := strings.ToLower(filepath.ToSlash(strings.TrimSpace(pattern)))
	p = strings.TrimPrefix(p, "./")
	if isUniversalWildcard(p) {
		return true
	}
	lowBase := strings.ToLower(base)
	lowRel := strings.TrimPrefix(strings.ToLower(filepath.ToSlash(rel)), "./")

	if isDirectMatch(p, lowBase, lowRel) {
		return true
	}
	if isPrefixOrSuffixMatch(p, lowBase, lowRel) {
		return true
	}
	if isSubstringMatch(p, lowBase, lowRel) {
		return true
	}

	return isDirectoryOrGlobMatch(p, lowBase, lowRel)
}

func isUniversalWildcard(p string) bool {
	return p == "" || p == "*" || p == "*.*" || p == "." || p == "./" || p == ".\\"
}

func isDirectMatch(p, lowBase, lowRel string) bool {
	if !strings.Contains(p, "/") {
		return lowBase == p || isGlobMatch(p, lowBase)
	}

	return lowRel == p || isGlobMatch(p, lowRel)
}

func isPrefixOrSuffixMatch(p, lowBase, lowRel string) bool {
	if isSingleSuffixPattern(p) {
		return matchSuffixPattern(p, lowBase, lowRel)
	}
	if isSinglePrefixPattern(p) {
		return matchPrefixPattern(p, lowBase, lowRel)
	}

	return false
}

func matchSuffixPattern(p, lowBase, lowRel string) bool {
	suffix := p[1:]
	if !strings.Contains(p, "/") {
		return strings.HasSuffix(lowBase, suffix)
	}

	return strings.HasSuffix(lowRel, suffix)
}

func matchPrefixPattern(p, lowBase, lowRel string) bool {
	prefix := strings.TrimSuffix(p, "*")
	if !strings.Contains(p, "/") {
		return strings.HasPrefix(lowBase, prefix)
	}

	return strings.HasPrefix(lowRel, prefix)
}

func isSingleSuffixPattern(p string) bool {
	return strings.HasPrefix(p, "*") && !strings.Contains(p[1:], "*")
}

func isSinglePrefixPattern(p string) bool {
	return strings.HasSuffix(p, "*") && !strings.Contains(p[:len(p)-1], "*")
}

func isSubstringMatch(p, lowBase, lowRel string) bool {
	if !isEnclosedWildcard(p) {
		return false
	}
	token := strings.Trim(p, "*")
	if !strings.Contains(p, "/") {
		return strings.Contains(lowBase, token)
	}

	return strings.Contains(lowRel, token)
}

func isEnclosedWildcard(p string) bool {
	return strings.HasPrefix(p, "*") && strings.HasSuffix(p, "*") && len(p) > 2
}

func isDirectoryOrGlobMatch(p, lowBase, lowRel string) bool {
	if strings.HasPrefix(p, "**/") && isGlobMatch(strings.TrimPrefix(p, "**/"), lowBase) {
		return true
	}
	dirPrefix := strings.TrimSuffix(p, "/") + "/"

	return strings.HasPrefix(lowRel, dirPrefix)
}

func isGlobMatch(pattern, target string) bool {
	match, _ := path.Match(pattern, target)

	return match
}
