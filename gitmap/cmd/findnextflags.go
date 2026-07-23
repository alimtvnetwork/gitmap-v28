package cmd

// Strict flag parser for `gitmap find-next`.
//
// Replaces the previous best-effort loop that silently ignored
// unknown tokens, mistyped flags, and `--json=true`-style boolean
// misuse. The new contract is: any token the parser cannot
// confidently interpret produces a stderr error + exit 2 from the
// caller. This catches typos like `--jsno` or `--scan_folder` early
// instead of letting the user wonder why their filter "didn't work."

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// parseFindNextFlags walks args once, accepting:
//   - --json                 (boolean)
//   - --scan-folder <id>     (integer; space-separated)
//   - --scan-folder=<id>     (integer; equals-separated, for symmetry)
//
// Returns (scanFolderID, jsonOut, err). On err, the caller prints to
// stderr and exits with the usage code.
func parseFindNextFlags(args []string) (int64, bool, error) {
	var scanFolderID int64
	jsonOut := false

	for i := 0; i < len(args); i++ {
		tok := args[i]
		consumed, nextID, nextJSON, err := classifyFindNextToken(tok, args, i, scanFolderID, jsonOut)
		if err != nil {
			return 0, false, err
		}
		scanFolderID, jsonOut = nextID, nextJSON
		i += consumed
	}

	return scanFolderID, jsonOut, nil
}

// classifyFindNextToken inspects one token and returns how many
// EXTRA args it consumed past `i` (0 or 1), the updated state, and
// any error. Splitting the dispatch out keeps parseFindNextFlags
// under the 15-line guideline.
func classifyFindNextToken(tok string, args []string, i int,
	curID int64, curJSON bool) (int, int64, bool, error) {
	if eq := strings.IndexByte(tok, '='); eq > 0 {
		return classifyEqualsForm(tok[:eq], tok[eq+1:], curID, curJSON)
	}
	switch tok {
	case constants.FindNextFlagJSON:
		return 0, curID, true, nil
	case constants.FindNextFlagScanFolder:
		return classifyScanFolderSpaceForm(args, i, curJSON)
	}

	return 0, 0, false, unknownOrPositional(tok)
}

// classifyEqualsForm handles `--flag=value` tokens uniformly so the
// space-separated and equals-separated paths share validation rules.
func classifyEqualsForm(name, value string, curID int64, curJSON bool) (int, int64, bool, error) {
	switch name {
	case constants.FindNextFlagJSON:
		return 0, 0, false, fmt.Errorf(constants.ErrFindNextBoolTakesNoValueFmt,
			constants.FindNextFlagJSON, value)
	case constants.FindNextFlagScanFolder:
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return 0, 0, false, fmt.Errorf(constants.ErrFindNextBadIntFmt,
				constants.FindNextFlagScanFolder, value)
		}

		return 0, v, curJSON, nil
	}

	return 0, 0, false, unknownOrPositional(name)
}

// classifyScanFolderSpaceForm handles the `--scan-folder <id>` form.
// Returns 1 to signal that the value token at i+1 was consumed.
func classifyScanFolderSpaceForm(args []string, i int, curJSON bool) (int, int64, bool, error) {
	if i+1 >= len(args) || strings.HasPrefix(args[i+1], "--") {
		return 0, 0, false, fmt.Errorf(constants.ErrFindNextMissingValueFmt,
			constants.FindNextFlagScanFolder)
	}
	raw := args[i+1]
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, 0, false, fmt.Errorf(constants.ErrFindNextBadIntFmt,
			constants.FindNextFlagScanFolder, raw)
	}

	return 1, v, curJSON, nil
}

// unknownOrPositional formats the right error for a token the parser
// did not recognize. Tokens starting with `--` are treated as unknown
// flags (with a typo suggestion when one is close); anything else is
// reported as an unexpected positional argument.
func unknownOrPositional(tok string) error {
	if !strings.HasPrefix(tok, "--") {
		return fmt.Errorf(constants.ErrFindNextUnexpectedArgFmt, tok)
	}
	if guess, ok := suggestFindNextFlag(tok); ok {
		return fmt.Errorf(constants.ErrFindNextUnknownFlagSuggestFmt, tok, guess)
	}

	return fmt.Errorf(constants.ErrFindNextUnknownFlagFmt, tok)
}

// suggestFindNextFlag returns the closest known flag when `tok` is
// within edit distance 2 of one of them. Threshold 2 catches single
// transpositions ("--jsno" → "--json") and one missing/extra char
// ("--scanfolder" → "--scan-folder") without firing on totally
// unrelated tokens like "--quiet".
func suggestFindNextFlag(tok string) (string, bool) {
	best := ""
	bestDist := 3
	for _, known := range constants.FindNextKnownFlags {
		d := levenshtein(tok, known)
		if d < bestDist {
			best, bestDist = known, d
		}
	}
	if bestDist <= 2 {
		return best, true
	}

	return "", false
}

// levenshtein is a tiny edit-distance implementation used only by
// the suggestion engine. Worst-case input is two short flag names
// (~20 chars), so the O(n*m) cost is trivial and dependency-free.
func levenshtein(a, b string) int {
	ar, br := []rune(a), []rune(b)
	rows := make([][]int, len(ar)+1)
	for i := range rows {
		rows[i] = make([]int, len(br)+1)
		rows[i][0] = i
	}
	for j := range rows[0] {
		rows[0][j] = j
	}
	fillLevenshteinMatrix(rows, ar, br)

	return rows[len(ar)][len(br)]
}

// fillLevenshteinMatrix is split out so levenshtein itself stays
// under the 15-line function budget.
func fillLevenshteinMatrix(rows [][]int, ar, br []rune) {
	for i := 1; i <= len(ar); i++ {
		for j := 1; j <= len(br); j++ {
			cost := 1
			if ar[i-1] == br[j-1] {
				cost = 0
			}
			rows[i][j] = minInt3(
				rows[i-1][j]+1,
				rows[i][j-1]+1,
				rows[i-1][j-1]+cost,
			)
		}
	}
}

// minInt3 returns the smallest of three ints. Local helper so we
// don't pull in a generics-heavy utility for two callsites.
func minInt3(a, b, c int) int {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}

	return m
}
