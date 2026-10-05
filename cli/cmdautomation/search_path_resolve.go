package cmdautomation

import (
	"os"
	"path/filepath"
	"strings"
)

var commonFileExtensions = []string{
	".tsx", ".ts", ".jsx", ".js", ".go", ".py", ".rs",
	".json", ".md", ".html", ".css", ".yaml", ".yml", ".sql",
}

// resolveSearchCandidateFiles finds matching target files or returns empty if directory walk is needed.
func resolveSearchCandidateFiles(target string) ([]string, bool) {
	clean := filepath.Clean(strings.TrimSpace(target))
	if clean == "" || clean == "." {
		return nil, false
	}
	if files, isFound := tryResolveExactFile(clean); isFound {
		return files, true
	}
	if files, found := tryResolveWithExtensions(clean); found {
		return files, true
	}
	if files, found := tryResolvePrefixMatches(clean); found {
		return files, true
	}
	return nil, false
}

func tryResolveExactFile(clean string) ([]string, bool) {
	fi, err := os.Stat(clean)
	if err != nil || fi.IsDir() {
		return nil, false
	}
	return []string{clean}, true
}

func tryResolveWithExtensions(basePath string) ([]string, bool) {
	for _, ext := range commonFileExtensions {
		candidate := basePath + ext
		if fi, err := os.Stat(candidate); err == nil && !fi.IsDir() {
			return []string{candidate}, true
		}
	}
	return nil, false
}

func tryResolvePrefixMatches(target string) ([]string, bool) {
	parent := filepath.Dir(target)
	base := filepath.Base(target)
	entries, err := os.ReadDir(parent)
	if err != nil {
		return nil, false
	}
	var matches []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(strings.ToLower(name), strings.ToLower(base)) {
			matches = append(matches, filepath.Join(parent, name))
		}
	}
	if len(matches) > 0 {
		return matches, true
	}
	return nil, false
}
