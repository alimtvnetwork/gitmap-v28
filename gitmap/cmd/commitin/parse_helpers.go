package commitin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// sprintf is a thin wrapper kept so parse_types.go avoids a direct
// `fmt` dependency and the file stays under the size budget.
func sprintf(format string, args ...any) string {
	return fmt.Sprintf(format, args...)
}

// splitInputs implements spec §2.2 `INPUT (SEP INPUT)*`. SEP is one
// or more spaces, a comma, or a comma followed by spaces. Quoting is
// optional and may wrap each token OR the whole list — we strip a
// matched pair of surrounding double quotes once.
//
// Behavior is the union of two real-world invocations:
//
//	commit-in <s> a b c            → ["a","b","c"]
//	commit-in <s> a,b,c            → ["a","b","c"]
//	commit-in <s> "a, b, c"        → ["a","b","c"]
//	commit-in <s> "a"  "b"   "c"   → ["a","b","c"]
//	commit-in <s> "a, b" c         → ["a","b","c"]
func splitInputs(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, raw := range tokens {
		trimmed := stripOuterQuotes(strings.TrimSpace(raw))
		for _, part := range strings.Split(trimmed, constants.CommitInCsvSep) {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			out = append(out, stripOuterQuotes(part))
		}
	}
	return out
}

// stripOuterQuotes removes a single pair of matched outer double quotes.
// Mismatched or absent quotes pass through unchanged.
func stripOuterQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

// classifyKeyword recognizes the §2.4 special inputs. Returns
// (keyword, tail, isKeyword, err). `tail` is N for "-N" forms (N >= 1).
// "all" returns tail = 0.
func classifyKeyword(token string) (string, int, bool, *ParseError) {
	if token == constants.CommitInInputKeywordAll {
		return token, 0, true, nil
	}
	if !strings.HasPrefix(token, constants.CommitInInputKeywordTailDash) {
		return "", 0, false, nil
	}
	digits := token[len(constants.CommitInInputKeywordTailDash):]
	if digits == "" {
		return "", 0, true, &ParseError{
			ExitCode: constants.CommitInExitBadArgs,
			Message:  fmt.Sprintf(constants.CommitInErrInputKeyword, token),
		}
	}
	n, err := strconv.Atoi(digits)
	if err != nil || n < 1 {
		return "", 0, true, &ParseError{
			ExitCode: constants.CommitInExitBadArgs,
			Message:  fmt.Sprintf(constants.CommitInErrInputKeyword, token),
		}
	}
	return token, n, true, nil
}

// splitCSV splits a CSV value, trimming whitespace and dropping empty
// fragments. Used uniformly by every CSV-shaped flag.
func splitCSV(value string) []string {
	if value == "" {
		return nil
	}
	parts := strings.Split(value, constants.CommitInCsvSep)
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
