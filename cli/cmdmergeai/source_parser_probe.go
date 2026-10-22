// Package cmdmergeai implements multi-repository amalgamation, collision sequencing, and AI merge handoff.
package cmdmergeai

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

func isExistingFile(p string) bool {
	if p == "" {
		return false
	}
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}

// ResolveSourceFilePath locates an input text file directly or inside repo-cache.
func ResolveSourceFilePath(inputPath string) (string, error) {
	trimmed := strings.TrimSpace(inputPath)
	if trimmed == "" {
		return "", apperror.NewValidationError("empty source file path")
	}

	candidates := []string{trimmed}
	lower := strings.ToLower(trimmed)
	if !strings.HasSuffix(lower, ".txt") && !strings.HasSuffix(lower, ".manifest") {
		candidates = append(candidates, trimmed+".txt")
	}

	for _, cand := range candidates {
		resolved, ok := probeSingleCandidate(cand)
		if ok {
			return resolved, nil
		}
	}

	return trimmed, apperror.NewNotFoundError("source file not found: " + trimmed)
}

func probeSingleCandidate(cand string) (string, bool) {
	if isExistingFile(cand) {
		return cand, true
	}
	resolvedCwd, okCwd := probeCwdCandidate(cand)
	if okCwd {
		return resolvedCwd, true
	}
	resolvedRc, okRc := probeRepoCacheCandidate(cand)
	if okRc {
		return resolvedRc, true
	}
	resolvedWb, okWb := probeWorkBaseCandidate(cand)
	if okWb {
		return resolvedWb, true
	}
	return "", false
}

func probeCwdCandidate(cand string) (string, bool) {
	cwd, err := os.Getwd()
	if err != nil || cwd == "" {
		return "", false
	}
	cwdCandidate := filepath.Join(cwd, cand)
	if isExistingFile(cwdCandidate) {
		return cwdCandidate, true
	}
	cwdRcCandidate := filepath.Join(cwd, "repo-cache", filepath.Base(cand))
	if isExistingFile(cwdRcCandidate) {
		return cwdRcCandidate, true
	}
	return "", false
}

func probeRepoCacheCandidate(cand string) (string, bool) {
	repoCacheRoot := cmdclone.ResolveRepoCacheRoot()
	if repoCacheRoot == "" {
		return "", false
	}
	rcCandidate := filepath.Join(repoCacheRoot, cand)
	if isExistingFile(rcCandidate) {
		return rcCandidate, true
	}
	rcBaseCandidate := filepath.Join(repoCacheRoot, filepath.Base(cand))
	if isExistingFile(rcBaseCandidate) {
		return rcBaseCandidate, true
	}
	return "", false
}

func probeWorkBaseCandidate(cand string) (string, bool) {
	workBase := resolveWorkBaseDir()
	if workBase == "" {
		return "", false
	}
	wbCandidate := filepath.Join(workBase, cand)
	if isExistingFile(wbCandidate) {
		return wbCandidate, true
	}
	wbBaseCandidate := filepath.Join(workBase, "repo-cache", filepath.Base(cand))
	if isExistingFile(wbBaseCandidate) {
		return wbBaseCandidate, true
	}
	return "", false
}
