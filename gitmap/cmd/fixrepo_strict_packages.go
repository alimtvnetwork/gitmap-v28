package cmd

// fixrepo_strict_packages.go — pure derivation of `go test` package
// patterns from the absolute paths of modified .go files. Split from
// the orchestration file (fixrepo_strict.go) so the mapping logic is
// trivially unit-testable without spawning `go test`.
//
// Why we feed `go test` directory patterns (./pkg/...) instead of
// individual .go files: `go test` operates per-package, not per-file,
// and refuses to compile a single file in isolation when its package
// has siblings (which is the common case). Mapping to the parent dir
// of every modified .go file gives us exactly the set of packages
// whose semantics could have shifted under the rewriter.

import (
	"path/filepath"
	"sort"
	"strings"
)

// derivePackagesFromGoFiles maps absolute .go file paths under repoRoot
// to a sorted, deduplicated slice of `go test` package patterns of the
// form `./<rel-dir>` (or `.` for files at the repo root). Files that
// fall outside repoRoot — which would happen only via a symlink that
// the scan layer should already have rejected — are silently skipped
// rather than emitting a malformed pattern, because passing absolute
// paths to `go test` would force GOPATH/module-mode resolution that
// the caller cannot predict.
//
// The `./...` recursive form is intentionally NOT used here: a file
// modified in package A should not trigger tests in unrelated
// sub-package A/B. The user opted into --strict to catch desyncs in
// touched packages, not to run the entire repo's test suite.
func derivePackagesFromGoFiles(repoRoot string, goFiles []string) []string {
	if len(goFiles) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(goFiles))
	for _, full := range goFiles {
		pattern, ok := goFileToPackagePattern(repoRoot, full)
		if !ok {
			continue
		}
		seen[pattern] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)

	return out
}

// goFileToPackagePattern converts one absolute .go path to its
// `./<rel-dir>` package pattern. Returns (pattern, true) on success
// or ("", false) when the file is not under repoRoot.
//
// filepath.Rel + ToSlash normalize Windows path separators so the
// emitted pattern is identical across platforms — important because
// `go test` on Windows accepts both forward and backslash paths but
// downstream log readers (CI, tests, humans) expect forward slashes.
func goFileToPackagePattern(repoRoot, fullPath string) (string, bool) {
	rel, err := filepath.Rel(repoRoot, fullPath)
	if err != nil {
		return "", false
	}
	rel = filepath.ToSlash(rel)
	if strings.HasPrefix(rel, "../") || rel == ".." {
		return "", false
	}
	dir := filepath.ToSlash(filepath.Dir(rel))
	if dir == "." || dir == "" {
		return ".", true
	}

	return "./" + dir, true
}
