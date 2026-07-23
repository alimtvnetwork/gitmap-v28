package cmd

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// repoRoot returns the absolute path of the current git repo. Exits 1
// when not inside a repo, per the zero-swallow error policy.
func repoRoot() string {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceNotInRepo, err)
		os.Exit(1)
	}
	return strings.TrimSpace(string(out))
}

// walkRepoFiles returns every text file under root, honoring directory
// and prefix exclusions plus an optional extension allow-list. Pass
// nil/empty exts to disable extension filtering. caseInsensitive
// controls how the file's extension is compared against the list — see
// matchesExtFilter for the exact contract.
func walkRepoFiles(root string, exts []string, caseInsensitive bool) ([]string, error) {
	out := make([]string, 0, 1024)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		return visitReplaceEntry(root, path, d, exts, caseInsensitive, &out)
	})
	return out, err
}

// visitReplaceEntry implements one step of walkRepoFiles. Split out so
// the closure stays under the 15-line ceiling.
func visitReplaceEntry(
	root, path string, d fs.DirEntry,
	exts []string, caseInsensitive bool, out *[]string,
) error {
	if d.IsDir() {
		if isExcludedDir(d.Name()) || isExcludedPrefix(root, path) {
			return filepath.SkipDir
		}
		return nil
	}
	if !matchesExtFilter(path, exts, caseInsensitive) {
		return nil
	}
	if isBinaryFile(path) {
		return nil
	}
	*out = append(*out, path)
	return nil
}

// matchesExtFilter returns true when exts is empty (no filter) or when
// the file's extension matches any allow-list entry. When
// caseInsensitive is true the file extension is lowercased to match the
// pre-normalized list; otherwise the comparison is byte-exact and the
// original filename casing is respected.
func matchesExtFilter(path string, exts []string, caseInsensitive bool) bool {
	if len(exts) == 0 {
		return true
	}
	got := filepath.Ext(path)
	if got == "" {
		return false
	}
	if caseInsensitive {
		got = strings.ToLower(got)
	}
	for _, want := range exts {
		if got == want {
			return true
		}
	}
	return false
}

// isExcludedDir checks the directory's base name against the deny list.
func isExcludedDir(name string) bool {
	for _, ex := range constants.ReplaceExcludedDirs {
		if name == ex {
			return true
		}
	}
	return false
}

// isExcludedPrefix matches a path against the prefix deny list, using
// forward slashes for portability.
func isExcludedPrefix(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	rel = filepath.ToSlash(rel)
	for _, p := range constants.ReplaceExcludedPrefixes {
		if rel == p || strings.HasPrefix(rel, p+"/") {
			return true
		}
	}
	return false
}

// isBinaryFile sniffs the first ReplaceBinarySniffBytes bytes for a
// null byte — git's classic heuristic.
func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, constants.ReplaceBinarySniffBytes)
	n, _ := f.Read(buf)
	return bytes.IndexByte(buf[:n], 0) >= 0
}
