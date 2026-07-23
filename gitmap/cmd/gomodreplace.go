package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// readModulePath reads the module path from go.mod in the current directory.
func readModulePath() string {
	data, err := os.ReadFile(constants.GoModFile)
	if err != nil {
		fmt.Fprint(os.Stderr, constants.ErrGoModNoFile)
		os.Exit(1)
	}

	return parseModuleLine(string(data))
}

// parseModuleLine extracts the module path from go.mod content.
func parseModuleLine(content string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, constants.GoModModuleLine) {
			return strings.TrimSpace(strings.TrimPrefix(line, constants.GoModModuleLine))
		}
	}

	fmt.Fprint(os.Stderr, constants.ErrGoModNoModule)
	os.Exit(1)

	return ""
}

// replaceModulePath replaces all occurrences of oldPath with newPath across the repo.
func replaceModulePath(oldPath, newPath string, verbose bool, exts []string) int {
	replaceInGoMod(oldPath, newPath)
	files := findFilesWithPath(oldPath, exts)

	if len(files) == 0 {
		fmt.Print(constants.MsgGoModNoImports)

		return 0
	}

	for _, f := range files {
		replaceInFile(f, oldPath, newPath)
		if verbose {
			fmt.Printf(constants.MsgGoModVerboseFile, f)
		}
	}

	return len(files)
}

// replaceInGoMod replaces the module line in go.mod.
func replaceInGoMod(oldPath, newPath string) {
	data, err := os.ReadFile(constants.GoModFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModReadFailed, constants.GoModFile, err)
		os.Exit(1)
	}

	updated := strings.ReplaceAll(string(data), oldPath, newPath)
	writeFileContent(constants.GoModFile, updated)
}

// findFilesWithPath walks the repo and returns files containing oldPath.
// If exts is empty, all files are considered. Otherwise only files matching
// the given extensions (e.g. ".go", ".md") are checked.
// go.mod itself is excluded since it is handled separately.
func findFilesWithPath(oldPath string, exts []string) []string {
	var matches []string

	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && isGoModExcludedDir(info.Name()) {
			return filepath.SkipDir
		}
		if info.IsDir() {
			return nil
		}
		if path == constants.GoModFile {
			return nil
		}
		if matchesGoModExt(path, exts) && fileContains(path, oldPath) {
			matches = append(matches, path)
		}

		return nil
	})

	return matches
}

// matchesGoModExt returns true if the file matches the extension filter.
// An empty exts slice means all files match. Domain-prefixed to avoid
// colliding with the replace package's matchesExtFilter (which carries
// an additional case-sensitivity flag).
func matchesGoModExt(path string, exts []string) bool {
	if len(exts) == 0 {
		return true
	}

	ext := filepath.Ext(path)
	for _, e := range exts {
		if ext == e {
			return true
		}
	}

	return false
}

// parseExtFlag parses a comma-separated extension string like "*.go,*.md"
// into a slice of extensions like [".go", ".md"].
func parseExtFlag(raw string) []string {
	if len(raw) == 0 {
		return nil
	}

	parts := strings.Split(raw, ",")
	var exts []string

	for _, p := range parts {
		p = strings.TrimSpace(p)
		p = strings.TrimPrefix(p, "*")
		if len(p) > 0 {
			exts = append(exts, p)
		}
	}

	return exts
}

// isGoModExcludedDir checks if a directory should be skipped.
// Domain-prefixed to avoid colliding with the replace package's
// isExcludedDir (which uses ReplaceExcludedDirs).
func isGoModExcludedDir(name string) bool {
	for _, d := range constants.GoModExcludeDirs {
		if name == d {
			return true
		}
	}

	return false
}

// fileContains checks if a file contains the given substring.
func fileContains(path, substr string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	return strings.Contains(string(data), substr)
}

// replaceInFile replaces all occurrences of oldPath with newPath in a file.
func replaceInFile(path, oldPath, newPath string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModReadFailed, path, err)

		return
	}

	updated := strings.ReplaceAll(string(data), oldPath, newPath)
	writeFileContent(path, updated)
}

// writeFileContent writes content to a file preserving its permissions.
func writeFileContent(path, content string) {
	info, err := os.Stat(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModWriteFailed, path, err)
		os.Exit(1)
	}

	err = os.WriteFile(path, []byte(content), info.Mode())
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrGoModWriteFailed, path, err)
		os.Exit(1)
	}
}
