// Package cmd — visibilitybulkprompt.go: interactive renderer +
// confirm/exclude loop for the bulk wildcard visibility commands.
//
// Renderer is pure (string in / string out) so it can be golden-tested.
// Prompt I/O is split into a thin wrapper that takes io.Reader / Writer
// for the same reason. -Y short-circuits BOTH prompts upstream in
// visibilityallbulk.go (plan step 13).
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §4.
package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/visibility"
)

// renderMatchedTable returns the human-readable table emitted before
// the confirm prompt. Format:
//
//	Matched N of TOTAL repos under OWNER:
//	   1  repo-name             (pattern)
//	   2  another-repo          (pattern)
//
// Right-aligns the 1-based index to the width of N so columns line up
// for 1..9999 without rewrap.
func renderMatchedTable(owner string, total int, matches []visibility.MatchedRepo) string {
	var b strings.Builder
	fmt.Fprintf(&b, constants.MsgBulkMatchedHeaderFmt, len(matches), total, owner)
	width := digitWidth(len(matches))
	nameWidth := longestRepoName(matches)
	for i, m := range matches {
		fmt.Fprintf(&b, "  %*d  %-*s  (%s)\n", width, i+1, nameWidth, m.RepoName, m.MatchedPattern)
	}

	return b.String()
}

// digitWidth returns the number of decimal digits in n (min 1).
func digitWidth(n int) int {
	if n < 10 {
		return 1
	}
	w := 0
	for n > 0 {
		n /= 10
		w++
	}

	return w
}

// longestRepoName returns the max RepoName length, with a sane floor.
func longestRepoName(matches []visibility.MatchedRepo) int {
	w := 1
	for _, m := range matches {
		if len(m.RepoName) > w {
			w = len(m.RepoName)
		}
	}

	return w
}

// promptConfirmOrExclude runs the y/N/exclude loop until the user
// either confirms a subset or aborts. Returns the final subset (after
// any exclusions applied) and a boolean indicating whether the run
// should proceed. Re-prompts on a parse error rather than exiting so a
// typo does not nuke the whole run.
func promptConfirmOrExclude(in io.Reader, out io.Writer, matches []visibility.MatchedRepo) ([]visibility.MatchedRepo, bool) {
	reader := bufio.NewReader(in)
	current := matches
	for {
		fmt.Fprintf(out, constants.MsgBulkConfirmFmt, len(current))
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprint(out, constants.ErrBulkPromptEOF)

			return nil, false
		}
		next, keepGoing, done := handlePromptLine(line, current, out)
		if done {
			return next, keepGoing
		}
		current = next
	}
}

// handlePromptLine returns (currentSet, proceed, done). When done is
// false the outer loop re-prompts (used for the parse-error retry +
// the exclusion-applied re-display path).
func handlePromptLine(raw string, current []visibility.MatchedRepo, out io.Writer) ([]visibility.MatchedRepo, bool, bool) {
	tok := strings.ToLower(strings.TrimSpace(raw))
	if tok == "y" || tok == "yes" {
		return current, true, true
	}
	if len(tok) == 0 || tok == "n" || tok == "no" {
		return nil, false, true
	}

	excluded, isAll, err := visibility.ParseExclusionList(tok, len(current))
	if err != nil {
		fmt.Fprintf(out, constants.ErrBulkExclusionFmt, err)

		return current, false, false
	}
	if isAll {
		return nil, false, true
	}
	next := applyExclusions(current, excluded)
	fmt.Fprintf(out, constants.MsgBulkExcludedFmt, len(excluded), len(next))

	return next, false, false
}

// applyExclusions returns a new slice with the 1-based indices removed.
// `excluded` MUST be sorted ascending (ParseExclusionList guarantees it).
func applyExclusions(matches []visibility.MatchedRepo, excluded []int) []visibility.MatchedRepo {
	skip := make(map[int]bool, len(excluded))
	for _, n := range excluded {
		skip[n] = true
	}
	out := make([]visibility.MatchedRepo, 0, len(matches)-len(excluded))
	for i, m := range matches {
		if skip[i+1] {
			continue
		}
		out = append(out, m)
	}

	return out
}
