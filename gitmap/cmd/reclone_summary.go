package cmd

// Pre-execute summary for `gitmap reclone --execute`.
//
// Printed (stderr) right BEFORE the safety prompt so the user sees:
//
//   - source manifest + format
//   - effective --mode / --on-exists / --cwd
//   - row totals: how many will clone fresh vs. land on an existing dir
//   - a tree-style preview of the destination folder layout
//
// The dry-run renderer (Render in clonenow/render.go) already shows
// the per-row plan; this summary is a higher-altitude view designed
// for the moment between "I typed --execute" and "git starts running"
// — it answers "what's about to happen to my disk?" in one screen.

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printRecloneExecuteSummary emits the totals + folder-tree preview.
// No-op when the user passed --no-summary so terse CI logs stay
// terse. All output goes to stderr to keep stdout reserved for the
// machine-readable per-row results that follow.
func printRecloneExecuteSummary(plan clonenow.Plan, cfg cloneNowFlags) {
	if cfg.noSummary {

		return
	}
	resolvedCwd := cfg.cwd
	if resolvedCwd == "" {
		resolvedCwd = "."
	}
	existing := collectExistingDests(plan, cfg.cwd)
	printSummaryHeader(plan, cfg, resolvedCwd, len(existing))
	printSummaryTree(plan)
}

// printSummaryHeader writes the metadata block (source/mode/cwd) and
// the row-count line. Split out so the function stays under the
// project's 15-line guideline.
func printSummaryHeader(plan clonenow.Plan, cfg cloneNowFlags, cwd string, existing int) {
	fmt.Fprintf(os.Stderr, constants.MsgCloneNowSummaryHeaderFmt,
		plan.Source, plan.Format, plan.Mode, cfg.onExists, cwd)
	total := len(plan.Rows)
	fmt.Fprintf(os.Stderr, constants.MsgCloneNowSummaryCountsFmt,
		total, total-existing, existing)
}

// printSummaryTree renders the destination-folder layout as an
// indented, sorted tree — capped at CloneNowSummaryTreeLimit lines
// to keep big round-trips scannable.
func printSummaryTree(plan clonenow.Plan) {
	paths := collectSortedDestPaths(plan)
	if len(paths) == 0 {

		return
	}
	fmt.Fprint(os.Stderr, constants.MsgCloneNowSummaryTreeTitle)
	limit := constants.CloneNowSummaryTreeLimit
	shown := len(paths)
	if shown > limit {
		shown = limit
	}
	for _, line := range buildTreeLines(paths[:shown]) {
		fmt.Fprintf(os.Stderr, constants.MsgCloneNowSummaryTreeLineFmt, line)
	}
	if len(paths) > shown {
		fmt.Fprintf(os.Stderr, constants.MsgCloneNowSummaryTreeTruncFmt,
			len(paths)-shown)
	}
}

// collectSortedDestPaths returns the row RelativePaths normalized to
// forward slashes and sorted lexicographically so the tree renders
// deterministically across OSes (Windows scan output uses backslash
// separators in places).
func collectSortedDestPaths(plan clonenow.Plan) []string {
	out := make([]string, 0, len(plan.Rows))
	for _, r := range plan.Rows {
		if r.RelativePath == "" {

			continue
		}
		out = append(out, filepath.ToSlash(r.RelativePath))
	}
	sort.Strings(out)

	return out
}

// buildTreeLines converts a sorted slice of slash-separated paths
// into indented tree lines. Each path is rendered relative to its
// shared prefix with the previous path so common parents collapse
// into a single visual subtree:
//
//	apps/
//	  web/
//	    frontend
//	    backend
//	  cli
//
// Pure-string transform — no filesystem access — so the preview
// stays accurate even when destinations don't exist yet.
func buildTreeLines(paths []string) []string {
	lines := make([]string, 0, len(paths))
	prev := []string{}
	for _, p := range paths {
		segs := strings.Split(p, "/")
		shared := sharedPrefixLen(prev, segs[:len(segs)-1])
		for i := shared; i < len(segs)-1; i++ {
			lines = append(lines, indent(i)+segs[i]+"/")
		}
		lines = append(lines, indent(len(segs)-1)+segs[len(segs)-1])
		prev = segs[:len(segs)-1]
	}

	return lines
}

// sharedPrefixLen returns the count of leading equal entries between
// two slices. Used to collapse common parent directories in the
// tree renderer.
func sharedPrefixLen(a, b []string) int {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {

			return i
		}
	}

	return n
}

// indent returns the visual indent for a tree depth. Two spaces per
// level keeps the preview compact even on 80-column terminals.
func indent(depth int) string {
	return strings.Repeat("  ", depth)
}
