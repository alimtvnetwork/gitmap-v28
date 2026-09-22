// Package cmd — lowercasefix_scan.go handles repository file scanning for uppercase files.
package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func findRenameCandidates(root string, opts LowerCaseFixOptions) ([]RenamePair, int, error) {
	var pairs []RenamePair
	scannedCount := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return skipOrContinueDir(info.Name())
		}
		scannedCount++
		rel, _ := filepath.Rel(root, path)
		if pair, isMatch := checkFileCandidate(path, rel, info.Name(), opts); isMatch {
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

func checkFileCandidate(path, rel, base string, opts LowerCaseFixOptions) (RenamePair, bool) {
	if !hasUppercaseChars(base) {
		return RenamePair{}, false
	}
	if opts.IsReadmeOnly && !isRootReadme(rel, base) {
		return RenamePair{}, false
	}
	if !opts.IsReadmeOnly && !isMatchingAnyPattern(rel, base, opts.Patterns) {
		return RenamePair{}, false
	}

	lowBase := strings.ToLower(base)
	dir := filepath.Dir(path)
	return RenamePair{
		OldPath: path,
		NewPath: filepath.Join(dir, lowBase),
		OldBase: base,
		NewBase: lowBase,
		RelPath: rel,
	}, true
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

func isMatchingAnyPattern(rel, base string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}
	for _, p := range patterns {
		if isPatternMatch(rel, base, p) {
			return true
		}
	}
	return false
}

func isPatternMatch(rel, base, pattern string) bool {
	p := strings.ToLower(pattern)
	if p == "*" || p == "*.*" {
		return true
	}
	lowBase := strings.ToLower(base)
	if match, _ := filepath.Match(p, lowBase); match {
		return true
	}
	lowRel := filepath.ToSlash(strings.ToLower(rel))
	if match, _ := filepath.Match(p, lowRel); match {
		return true
	}
	if strings.HasPrefix(p, "*") && !strings.Contains(p[1:], "*") {
		return strings.HasSuffix(lowBase, p[1:])
	}
	return false
}
