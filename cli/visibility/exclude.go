// Package visibility — exclude.go: parser for the per-index exclusion
// string the user types at the bulk-visibility confirm prompt.
//
// Grammar (per 02-spec/01-app/116 §4):
//
//	""        → no exclusions (equivalent to "none")
//	"none"    → no exclusions
//	"all"     → exclude every index (caller should abort the run)
//	"1,3-5"   → exclude indices 1, 3, 4, 5
//	"7"       → exclude index 7
//
// Rules:
//   - 1-based indices, ascending ranges only ("5-3" → error).
//   - Any index outside [1, totalCount] → error with offending value.
//   - Non-numeric token → error citing the token.
//   - Duplicates are silently coalesced.
//
// Returns a sorted, deduped slice of 1-based int indices.
package visibility

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ExcludeAll is the sentinel returned when the user types "all".
// Callers compare with this to short-circuit the apply loop.
const ExcludeAll = "all"

// ExcludeNone is the no-op token.
const ExcludeNone = "none"

// ExclusionTokenParams captures parameters for parsing an exclusion token.
type ExclusionTokenParams struct {
	Token      string
	TokenIndex int
	TotalCount int
	Out        map[int]bool
}

// ParseExclusionList compiles the raw input. totalCount is the size of
// the currently-displayed matched list (1-based bounds enforced against
// it). Returns ([]int{}, false, nil) for "none"/"".
// Returns (nil, true, nil) for "all".
func ParseExclusionList(raw string, totalCount int) ([]int, bool, *apperror.AppError) {
	trimmed := strings.ToLower(strings.TrimSpace(raw))
	if len(trimmed) == 0 || trimmed == ExcludeNone {
		return []int{}, false, nil
	}

	if trimmed == ExcludeAll {
		return nil, true, nil
	}

	out := make(map[int]bool)
	for i, tok := range strings.Split(trimmed, ",") {
		if err := absorbExclusionToken(ExclusionTokenParams{
			Token:      strings.TrimSpace(tok),
			TokenIndex: i + 1,
			TotalCount: totalCount,
			Out:        out,
		}); err != nil {
			return nil, false, err
		}
	}

	return sortedKeys(out), false, nil
}

// absorbExclusionToken parses one comma-separated token ("3" or "3-5")
// and writes every covered index into `out`. Returns Code Red errors
// on any malformed token.
func absorbExclusionToken(params ExclusionTokenParams) *apperror.AppError {
	if len(params.Token) == 0 {
		return apperror.NewWithDetails(
			"parse-exclusion",
			"E1001",
			fmt.Sprintf("Error: empty exclusion token at position %d (operation: parse-exclusion, reason: blank between commas)", params.TokenIndex),
			"visibility",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	if strings.Contains(params.Token, "-") {
		return absorbExclusionRange(params)
	}

	n, err := strconv.Atoi(params.Token)
	if err != nil {
		return apperror.NewWithDetails(
			"parse-exclusion",
			"E1002",
			fmt.Sprintf("Error: non-numeric exclusion token %q at position %d (operation: parse-exclusion, reason: %s)", params.Token, params.TokenIndex, err.Error()),
			"visibility",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	if err := checkExclusionBounds(n, params.TotalCount, params.Token); err != nil {
		return err
	}

	params.Out[n] = true

	return nil
}

// absorbExclusionRange handles "A-B" tokens with ascending-only rule.
func absorbExclusionRange(params ExclusionTokenParams) *apperror.AppError {
	parts := strings.SplitN(params.Token, "-", 2)
	lo, errLo := strconv.Atoi(strings.TrimSpace(parts[0]))
	hi, errHi := strconv.Atoi(strings.TrimSpace(parts[1]))
	if errLo != nil || errHi != nil {
		return apperror.NewWithDetails(
			"parse-exclusion",
			"E1003",
			fmt.Sprintf("Error: malformed range %q at position %d (operation: parse-exclusion, reason: range bounds must be integers)", params.Token, params.TokenIndex),
			"visibility",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	if hi < lo {
		return apperror.NewWithDetails(
			"parse-exclusion",
			"E1004",
			fmt.Sprintf("Error: descending range %q at position %d (operation: parse-exclusion, reason: hi < lo)", params.Token, params.TokenIndex),
			"visibility",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	if err := checkExclusionBounds(lo, params.TotalCount, params.Token); err != nil {
		return err
	}

	if err := checkExclusionBounds(hi, params.TotalCount, params.Token); err != nil {
		return err
	}

	for n := lo; n <= hi; n++ {
		params.Out[n] = true
	}

	return nil
}

// checkExclusionBounds returns a Code Red error when n is outside
// [1, totalCount].
func checkExclusionBounds(n, totalCount int, tok string) *apperror.AppError {
	if n >= 1 && n <= totalCount {
		return nil
	}

	return apperror.NewWithDetails(
		"parse-exclusion",
		"E1005",
		fmt.Sprintf("Error: exclusion index %d out of range in token %q (operation: parse-exclusion, reason: valid range is 1..%d)", n, tok, totalCount),
		"visibility",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)
}

// sortedKeys returns the map keys ascending.
func sortedKeys(m map[int]bool) []int {
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	return keys
}
