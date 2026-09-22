// Package cmd — lowercasefix_scan.go handles repository file scanning for uppercase files.
package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func findRenameCandidates(root string, patterns []string) ([]RenamePair, error) {
	var pairs []RenamePair
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info == nil {
			return nil
		}
		if info.IsDir() {
			return skipOrContinueDir(info.Name())
		}
		if pair, isMatch := checkFileCandidate(path, info.Name(), patterns); isMatch {
			pairs = append(pairs, pair)
		}
		return nil
	})
	return pairs, err
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

func checkFileCandidate(path, base string, patterns []string) (RenamePair, bool) {
	if !hasUppercaseChars(base) {
		return RenamePair{}, false
	}
	if !isMatchingAnyPattern(base, patterns) {
		return RenamePair{}, false
	}

	lowBase := strings.ToLower(base)
	dir := filepath.Dir(path)
	return RenamePair{
		OldPath: path,
		NewPath: filepath.Join(dir, lowBase),
		OldBase: base,
		NewBase: lowBase,
	}, true
}

func hasUppercaseChars(s string) bool {
	return strings.ToLower(s) != s
}

func isMatchingAnyPattern(base string, patterns []string) bool {
	if len(patterns) == 0 {
		return strings.EqualFold(base, "readme.md")
	}
	for _, p := range patterns {
		if match, _ := filepath.Match(strings.ToLower(p), strings.ToLower(base)); match {
			return true
		}
	}
	return false
}
