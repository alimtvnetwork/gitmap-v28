package cmdautofix

import (
	"bytes"
	"os"
	"sort"
	"time"
)

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

// flaggedFiles returns the sorted unique relative paths in a violation list.
func flaggedFiles(violations []Violation) []string {
	seen := map[string]struct{}{}
	for _, v := range violations {
		seen[v.Path] = struct{}{}
	}
	out := make([]string, 0, len(seen))
	for p := range seen {
		out = append(out, p)
	}
	sort.Strings(out)
	return out
}

// Apply rewrites every flagged file: it runs each selected category's Fix
// (nil Fix = report-only, never written), writes only when bytes differ,
// and runs Exec categories (gofmt) against the flushed file.
func Apply(result *ScanResult, opts Options) (*ApplyResult, error) {
	start := time.Now()
	if opts.Root == "" {
		opts.Root = result.Root
	}
	if opts.Root == "" {
		opts.Root = "."
	}
	cats, err := selectedCategories(opts)
	if err != nil {
		return nil, err
	}

	byteCats := []Category{}
	execCats := []Category{}
	for _, cat := range cats {
		if cat.Fix == nil {
			continue
		}
		if cat.Exec {
			execCats = append(execCats, cat)
			continue
		}
		byteCats = append(byteCats, cat)
	}

	out := &ApplyResult{FixCounts: map[string]int{}}
	for _, rel := range flaggedFiles(result.Violations) {
		abs := absFromRel(opts.Root, rel)
		src, err := os.ReadFile(abs)
		if err != nil {
			continue
		}
		if isBinaryProbe(src) {
			continue
		}

		buf := src
		touched := map[string]bool{}
		for _, cat := range byteCats {
			if !matchesScope(rel, cat, opts) {
				continue
			}
			fixed, fixVios := cat.Fix(rel, buf, opts)
			out.Skipped = append(out.Skipped, fixVios...)
			if !bytes.Equal(fixed, buf) {
				buf = fixed
				touched[cat.Name] = true
			}
		}
		if !bytes.Equal(buf, src) && writePreserveMode(abs, buf) {
			out.FilesModified++
			out.ModifiedFiles = append(out.ModifiedFiles, rel)
		}

		for _, cat := range execCats {
			if !matchesScope(rel, cat, opts) {
				continue
			}
			if len(cat.Check(rel, buf, opts)) == 0 {
				continue
			}
			fixed, fixVios := cat.Fix(rel, buf, opts)
			out.Skipped = append(out.Skipped, fixVios...)
			if !bytes.Equal(fixed, buf) {
				buf = fixed
				touched[cat.Name] = true
				out.FilesModified++
				out.ModifiedFiles = append(out.ModifiedFiles, rel)
			}
		}
		for name := range touched {
			out.FixCounts[name]++
		}
	}

	sort.Strings(out.ModifiedFiles)
	out.Skipped = dedupeViolations(out.Skipped)
	out.Elapsed = time.Since(start)
	return out, nil
}

// categoryRegistry is the ordered category set (spec §2). A nil Exts means
// "all text files" (the null-byte probe is the real gate); --ext overrides
// every category's scope.
