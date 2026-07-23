package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// copyFileContent copies file content from source to destination.
func copyFileContent(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := make([]byte, constants.TaskCopyBufferSize)
	_, err = io.CopyBuffer(destFile, srcFile, buf)

	return err
}

// loadGitignorePatterns reads .gitignore patterns from a directory.
func loadGitignorePatterns(dir string) []string {
	path := filepath.Join(dir, ".gitignore")
	data, err := os.ReadFile(path)

	if err != nil {
		return nil
	}

	return parseGitignoreLines(string(data))
}

// parseGitignoreLines extracts patterns from gitignore content.
func parseGitignoreLines(content string) []string {
	var patterns []string
	lines := strings.Split(content, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isGitignoreComment(trimmed) {
			continue
		}

		patterns = append(patterns, trimmed)
	}

	return patterns
}

// isGitignoreComment returns true for empty lines and comments.
func isGitignoreComment(line string) bool {
	if line == "" {
		return true
	}

	return strings.HasPrefix(line, "#")
}

// isIgnored checks if a path matches any gitignore pattern.
func isIgnored(relPath string, isDir bool, patterns []string) bool {
	if len(patterns) == 0 {
		return false
	}

	var wg sync.WaitGroup
	result := make(chan bool, len(patterns))

	for _, pattern := range patterns {
		wg.Add(1)

		go func(p string) {
			defer wg.Done()

			if matchesPattern(relPath, isDir, p) {
				result <- true
			}
		}(pattern)
	}

	go func() {
		wg.Wait()
		close(result)
	}()

	for matched := range result {
		if matched {
			return true
		}
	}

	return false
}

// matchesPattern checks if a path matches a single gitignore pattern.
func matchesPattern(relPath string, isDir bool, pattern string) bool {
	isDirPattern := strings.HasSuffix(pattern, "/")
	if isDirPattern && isDir {
		cleanPattern := strings.TrimSuffix(pattern, "/")

		return matchGitignoreGlob(relPath, cleanPattern)
	}

	if isDirPattern {
		return false
	}

	return matchGitignoreGlob(relPath, pattern)
}

// matchGitignoreGlob performs glob matching against path components.
func matchGitignoreGlob(relPath, pattern string) bool {
	matched, err := filepath.Match(pattern, filepath.Base(relPath))
	if err != nil {
		return false
	}

	return matched
}
