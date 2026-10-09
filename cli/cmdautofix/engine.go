package cmdautofix

import (
	"bytes"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// binaryProbeSize mirrors 02-shared-engine.py BINARY_PROBE_CHUNK_SIZE.
const binaryProbeSize = 8192

// Options carries the CLI-level knobs for a fix run. There is no Apply bool:
// the flow is check-driven (Scan → summary → prompt → Apply).
type Options struct {
	Root       string
	Workers    int
	Exts       []string // normalized, lowercase, leading dot; empty = category scope
	Categories []string // empty = all registered categories
	URIPattern string   // effective repo-URI pattern for the paths category
	AsJSON     bool
	NoCache    bool
	Yes        bool
}

// Violation is one hygiene finding. Path is always relative to the scan
// root and slash-separated. Line is 1-based; 0 means file-level.
type Violation struct {
	Path     string `json:"path"`
	Category string `json:"category"`
	Line     int    `json:"line"`
	Detail   string `json:"detail"`
}

// Category pairs a Check with an optional Fix. A nil Fix means the category
// is report-only: it never writes, even when the user confirms the prompt.
type Category struct {
	Name string
	Desc string
	Exts []string // nil = all text files (subject to --ext and the binary probe)
	// Check returns violations without touching disk.
	Check func(relPath string, src []byte, opts Options) []Violation
	// Fix returns the rewritten bytes plus violations found during the fix.
	// Returning bytes equal to src is a no-op (never written).
	Fix func(relPath string, src []byte, opts Options) ([]byte, []Violation)
	// Exec means the category shells out against the file on disk (gofmt):
	// Apply flushes byte-level fixes first, then runs Check+Fix on the
	// written file instead of the in-memory buffer.
	Exec bool
}

// ScanResult is the aggregated, deterministically ordered outcome of Scan.
type ScanResult struct {
	Root           string
	Categories     []string
	FilesScanned   int
	FilesFromCache int
	Violations     []Violation
	Elapsed        time.Duration
}

// ApplyResult aggregates one Apply pass over a ScanResult.
type ApplyResult struct {
	FilesModified int
	ModifiedFiles []string
	FixCounts     map[string]int
	Skipped       []Violation // findings that could not be fixed (never lossy-written)
	Elapsed       time.Duration
}

// errToolMissing signals a tool-level failure (exit code 2).
var errToolMissing = errors.New("required tool missing")

// selectedCategories resolves opts.Categories against the registry.
func selectedCategories(opts Options) ([]Category, error) {
	if len(opts.Categories) == 0 {
		return categoryRegistry, nil
	}
	selected := make([]Category, 0, len(opts.Categories))
	for _, name := range opts.Categories {
		found := false
		for _, cat := range categoryRegistry {
			if cat.Name != name {
				continue
			}
			selected = append(selected, cat)
			found = true
			break
		}
		if !found {
			return nil, fmt.Errorf("unknown category %q (valid: %s)", name, strings.Join(categoryNames(), ", "))
		}
	}
	return selected, nil
}

// categoryNames lists registered category names in registry order.
func categoryNames() []string {
	names := make([]string, 0, len(categoryRegistry))
	for _, cat := range categoryRegistry {
		names = append(names, cat.Name)
	}
	return names
}

// matchesScope reports whether relPath falls in a category's extension scope.
func matchesScope(relPath string, cat Category, opts Options) bool {
	ext := strings.ToLower(filepath.Ext(relPath))
	if len(opts.Exts) > 0 {
		return extInList(ext, opts.Exts)
	}
	if len(cat.Exts) == 0 {
		return true
	}
	return extInList(ext, cat.Exts)
}

func extInList(ext string, list []string) bool {
	for _, e := range list {
		if e == ext {
			return true
		}
	}
	return false
}

// isBinaryProbe reports true when the first chunk of data holds a null byte.
func isBinaryProbe(data []byte) bool {
	limit := len(data)
	if limit > binaryProbeSize {
		limit = binaryProbeSize
	}
	return bytes.IndexByte(data[:limit], 0x00) >= 0
}

// relSlash converts an absolute path to a scan-root-relative slash path.
func relSlash(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return filepath.ToSlash(abs)
	}
	return filepath.ToSlash(rel)
}

// absFromRel maps a scan-root-relative slash path back to disk.
func absFromRel(root, rel string) string {
	return filepath.Join(root, filepath.FromSlash(rel))
}

func dedupeViolations(in []Violation) []Violation {
	seen := make(map[Violation]struct{}, len(in))
	out := make([]Violation, 0, len(in))
	for _, v := range in {
		if _, dup := seen[v]; dup {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
