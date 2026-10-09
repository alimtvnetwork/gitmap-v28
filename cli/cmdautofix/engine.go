// Package cmdautofix implements `gitmap autofix` — a native Go port of the
// 03-ai-scripts autofixers (04, 05, 07, 08, 10, 26, 27, 31) behind one
// parallel file-walking engine.
package cmdautofix

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// binaryProbeSize mirrors 02-shared-engine.py BINARY_PROBE_CHUNK_SIZE.
const binaryProbeSize = 8192

// Options carries the CLI-level knobs for an autofix run.
type Options struct {
	Apply      bool
	Workers    int
	Exts       []string // normalized, lowercase, leading dot; empty = category defaults
	Categories []string // empty = all registered categories
	URIPattern string   // effective repo-URI pattern for the paths category
	AsJSON     bool
	Root       string // scan root; set by Run for Exec categories (gofmt)
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
// is report-only: it never writes, even under --apply.
type Category struct {
	Name string
	Desc string
	Exts []string // nil = all text files (subject to --ext and the binary probe)
	// Check returns violations without touching disk.
	Check func(relPath string, src []byte, opts Options) []Violation
	// Fix returns the rewritten bytes plus violations found during the fix.
	// Returning bytes equal to src is a no-op (never written).
	Fix func(relPath string, src []byte, opts Options) ([]byte, []Violation)
	// Exec means the category shells out to a file on disk (gofmt): the
	// engine flushes byte-level fixes first, then runs Check+Fix against
	// the written file instead of the in-memory buffer.
	Exec bool
}

// FileResult aggregates one worker's outcome for a single file.
type FileResult struct {
	Path            string
	Violations      []Violation
	Fixed           bool
	FixedCategories []string
}

// RunResult is the aggregated, deterministically ordered outcome of a run.
type RunResult struct {
	Root       string
	FilesSeen  int
	Violations []Violation
	FixedFiles []string
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

// walkFiles emits candidate file paths (absolute) on fileCh, skipping .git dirs.
func walkFiles(root string, fileCh chan<- string) {
	defer close(fileCh)
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		fileCh <- path
		return nil
	})
}

// processFile runs every selected category over one file.
func processFile(absPath string, root string, cats []Category, opts Options) FileResult {
	rel := relSlash(root, absPath)
	result := FileResult{Path: rel}

	src, err := os.ReadFile(absPath)
	if err != nil {
		return result
	}
	if isBinaryProbe(src) {
		return result
	}

	buf := src
	fixedCats := []string{}
	byteCats := []Category{}
	execCats := []Category{}
	for _, cat := range cats {
		if !matchesScope(rel, cat, opts) {
			continue
		}
		if cat.Exec {
			execCats = append(execCats, cat)
			continue
		}
		byteCats = append(byteCats, cat)
	}

	for _, cat := range byteCats {
		result.Violations = append(result.Violations, cat.Check(rel, buf, opts)...)
		if !opts.Apply || cat.Fix == nil {
			continue
		}
		fixed, fixVios := cat.Fix(rel, buf, opts)
		result.Violations = append(result.Violations, fixVios...)
		if !bytes.Equal(fixed, buf) {
			buf = fixed
			fixedCats = append(fixedCats, cat.Name)
		}
	}

	if opts.Apply && !bytes.Equal(buf, src) {
		if writePreserveMode(absPath, buf) {
			result.Fixed = true
			result.FixedCategories = fixedCats
		}
	}

	for _, cat := range execCats {
		result.Violations = append(result.Violations, cat.Check(rel, buf, opts)...)
		if !opts.Apply || cat.Fix == nil {
			continue
		}
		fixed, fixVios := cat.Fix(rel, buf, opts)
		result.Violations = append(result.Violations, fixVios...)
		if !bytes.Equal(fixed, buf) {
			buf = fixed
			result.Fixed = true
			result.FixedCategories = append(result.FixedCategories, cat.Name)
		}
	}

	return result
}

// writePreserveMode writes data back with the file's original permission bits.
func writePreserveMode(absPath string, data []byte) bool {
	info, err := os.Stat(absPath)
	if err != nil {
		return false
	}
	if err := os.WriteFile(absPath, data, info.Mode().Perm()); err != nil {
		return false
	}
	return true
}

// Run executes the parallel autofix pipeline over root.
func Run(root string, opts Options) (*RunResult, error) {
	opts.Root = root
	cats, err := selectedCategories(opts)
	if err != nil {
		return nil, err
	}

	for _, cat := range cats {
		if cat.Exec {
			if terr := catToolCheck(cat.Name); terr != nil {
				return nil, terr
			}
		}
	}

	workers := opts.Workers
	if workers < 1 {
		workers = 1
	}

	fileCh := make(chan string, workers*4)
	resultCh := make(chan FileResult, workers*4)

	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for absPath := range fileCh {
				resultCh <- processFile(absPath, root, cats, opts)
			}
		}()
	}

	go walkFiles(root, fileCh)
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	out := &RunResult{Root: root}
	for res := range resultCh {
		out.FilesSeen++
		out.Violations = append(out.Violations, res.Violations...)
		if res.Fixed {
			out.FixedFiles = append(out.FixedFiles, res.Path)
		}
	}

	// Deterministic output regardless of worker scheduling.
	sort.Slice(out.Violations, func(i, j int) bool {
		if out.Violations[i].Path != out.Violations[j].Path {
			return out.Violations[i].Path < out.Violations[j].Path
		}
		if out.Violations[i].Category != out.Violations[j].Category {
			return out.Violations[i].Category < out.Violations[j].Category
		}
		return out.Violations[i].Line < out.Violations[j].Line
	})
	out.Violations = dedupeViolations(out.Violations)
	sort.Strings(out.FixedFiles)

	return out, nil
}

// dedupeViolations drops exact duplicate findings (same path, category,
// line, detail) so Check+Fix pairs never double-report.
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

// categoryRegistry is the ordered category set (spec §2). A nil Exts means
// "all text files" (the null-byte probe is the real gate); --ext overrides
// every category's scope.
var categoryRegistry = []Category{
	{
		Name:  "encoding",
		Desc:  "UTF-8 no-BOM + LF normalization",
		Check: encodingCheck,
		Fix:   encodingFix,
	},
	{
		Name:  "newlines",
		Desc:  "CRLF→LF, trim trailing whitespace, one final newline",
		Check: newlinesCheck,
		Fix:   newlinesFix,
	},
	{
		Name:  "naming",
		Desc:  "boolean-comparison style audit (report-only)",
		Exts:  codeExts,
		Check: namingCheck,
		// Fix is nil: report-only by design (D11).
	},
	{
		Name:  "paths",
		Desc:  "Windows absolute path sanitizer (docs)",
		Check: pathsCheck,
		Fix:   pathsFix,
	},
	{
		Name:  "gofmt",
		Desc:  "gofmt -w over .go files",
		Exts:  []string{".go"},
		Check: gofmtCheck,
		Fix:   gofmtFix,
		Exec:  true,
	},
	{
		Name:  "misspell",
		Desc:  "British→American spelling, case-preserving",
		Exts:  []string{".md", ".go", ".ts", ".tsx", ".js", ".py", ".json", ".sh", ".ps1", ".txt"},
		Check: misspellCheck,
		Fix:   misspellFix,
	},
	{
		Name:  "markdown",
		Desc:  "collapse 3+ blank lines in .md",
		Exts:  []string{".md"},
		Check: markdownCheck,
		Fix:   markdownFix,
	},
	{
		Name:  "guidelines",
		Desc:  "composite: newlines + naming",
		Check: guidelinesCheck,
		Fix:   guidelinesFix,
	},
}
