package cmdignore

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// IgnoreGroup represents a named collection of gitignore patterns.
type IgnoreGroup struct {
	Name      string   `json:"name"`
	Patterns  []string `json:"patterns"`
	IsDefault bool     `json:"isDefault"`
}

// IgnoreScanIssue represents an issue detected in a repository's gitignore configuration.
type IgnoreScanIssue struct {
	RepoPath         string   `json:"repoPath"`
	RepoName         string   `json:"repoName"`
	HasDuplicate     bool     `json:"hasDuplicate"`
	DuplicateCount   int      `json:"duplicateCount"`
	HasResumeTask    bool     `json:"hasResumeTask"`
	TrackedResume    []string `json:"trackedResume,omitempty"`
	MissingGitmapDir bool     `json:"missingGitmapDir"`
}

// IsIgnored evaluates if a file is ignored without regex.
func (g *IgnoreGroup) IsIgnored(path string) bool {
	if isPathBinary(path) {
		return true
	}
	return matchesAnyPattern(path, g.Patterns)
}

func matchesAnyPattern(path string, patterns []string) bool {
	for _, p := range patterns {
		if strings.Contains(path, strings.Trim(p, "/")) {
			return true
		}
	}
	return false
}

func isPathBinary(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	return isContentBinary(f)
}

func isContentBinary(r io.Reader) bool {
	buf := make([]byte, 512)
	n, _ := r.Read(buf)
	return bytes.Contains(buf[:n], []byte{0})
}
