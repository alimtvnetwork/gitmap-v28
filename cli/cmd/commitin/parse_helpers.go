package commitin

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
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

			cleaned := stripOuterQuotes(part)
			out = append(out, expandRangeToken(cleaned)...)
		}
	}

	return out
}

func expandRangeToken(token string) []string {
	if items, ok := expandBraceRange(token); ok {
		return items
	}
	if items, ok := expandShorthandRange(token); ok {
		return items
	}
	return []string{token}
}

func expandBraceRange(token string) ([]string, bool) {
	s := strings.Index(token, "{")
	e := strings.Index(token, "}")
	if s < 0 || e <= s {
		return nil, false
	}
	parts := strings.Split(token[s+1:e], "..")
	if len(parts) != 2 {
		return nil, false
	}
	start, err1 := strconv.Atoi(parts[0])
	end, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || start > end {
		return nil, false
	}
	var res []string
	for i := start; i <= end; i++ {
		res = append(res, token[:s]+strconv.Itoa(i)+token[e+1:])
	}
	return res, true
}

func expandShorthandRange(token string) ([]string, bool) {
	if !strings.Contains(token, "..") || strings.Contains(token, "://") {
		return nil, false
	}
	parts := strings.Split(token, "..")
	if len(parts) != 2 {
		return nil, false
	}
	s1, start := extractDigits(strings.TrimSpace(parts[0]))
	_, end := extractDigits(strings.TrimSpace(parts[1]))
	if start <= 0 || end < start {
		return nil, false
	}
	var res []string
	for i := start; i <= end; i++ {
		prefix := s1
		if prefix == "" {
			prefix = "gitmap-v"
		}
		res = append(res, fmt.Sprintf("https://github.com/alimtvnetwork/%s%d", prefix, i))
	}
	return res, true
}

func extractDigits(s string) (string, int) {
	idx := strings.IndexAny(s, "0123456789")
	if idx < 0 {
		return s, 0
	}
	val, err := strconv.Atoi(s[idx:])
	if err != nil {
		return s, 0
	}
	return s[:idx], val
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
