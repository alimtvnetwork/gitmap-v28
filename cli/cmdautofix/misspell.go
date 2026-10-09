package cmdautofix

import (
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

// spellingMap ports the 30-word British→American dictionary from
// 27-misspell-auditor.py verbatim. The script's "delegates to the misspell
// binary" docstring is deliberately NOT ported — no shell-out here.
var spellingMap = map[string]string{
	"behaviour": "behavior", "behaviours": "behaviors",
	"colour": "color", "colours": "colors",
	"initialise": "initialize", "initialised": "initialized",
	"initialising": "initializing", "initialisation": "initialization",
	"customise": "customize", "customised": "customized",
	"customising": "customizing", "customisation": "customization",
	"synchronise": "synchronize", "synchronised": "synchronized",
	"synchronising": "synchronizing",
	"optimise": "optimize", "optimised": "optimized",
	"optimising": "optimizing", "optimisation": "optimization",
	"prioritise": "prioritize", "prioritised": "prioritized",
	"prioritising": "prioritizing",
	"serialise": "serialize", "serialised": "serialized",
	"serialising": "serializing", "serialisation": "serialization",
	"normalise": "normalize", "normalised": "normalized",
	"normalising": "normalizing", "normalisation": "normalization",
	"cancelling": "canceling", "cancelled": "canceled",
}

var misspellRe *regexp.Regexp

func init() {
	words := make([]string, 0, len(spellingMap))
	for brit := range spellingMap {
		words = append(words, brit)
	}
	// \b anchors make alternation order irrelevant (behaviour never matches
	// inside behaviours: no word boundary between r and s).
	misspellRe = regexp.MustCompile(`(?i)\b(` + strings.Join(words, "|") + `)\b`)
}

// isMisspellSkipped mirrors the script's `if "misspell" in tf.name: continue`.
func isMisspellSkipped(relPath string) bool {
	return strings.Contains(strings.ToLower(filepath.Base(relPath)), "misspell")
}

// preserveCase applies the D9 case-preserving replacement: ALL-UPPER→upper,
// Title→title, lower→dictionary form, anything else→dictionary (lowercase).
func preserveCase(word, replacement string) string {
	if isAllUpper(word) {
		return strings.ToUpper(replacement)
	}
	if isTitleWord(word) {
		return titleWord(replacement)
	}
	return replacement
}

func isAllUpper(s string) bool {
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				return false
			}
		}
	}
	return hasLetter
}

func isTitleWord(s string) bool {
	if s == "" {
		return false
	}
	runes := []rune(s)
	if !unicode.IsUpper(runes[0]) {
		return false
	}
	for _, r := range runes[1:] {
		if unicode.IsLetter(r) && !unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func titleWord(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	for i := 1; i < len(runes); i++ {
		runes[i] = unicode.ToLower(runes[i])
	}
	return string(runes)
}

// misspellCheck lists each hit (line, word) — script 27 port.
func misspellCheck(relPath string, src []byte, opts Options) []Violation {
	if isMisspellSkipped(relPath) {
		return nil
	}
	content := string(src)
	violations := []Violation{}
	for _, loc := range misspellRe.FindAllStringIndex(content, -1) {
		word := content[loc[0]:loc[1]]
		us := spellingMap[strings.ToLower(word)]
		violations = append(violations, Violation{
			Path:     relPath,
			Category: "misspell",
			Line:     lineAt(content, loc[0]),
			Detail:   "British spelling '" + word + "' → use American '" + preserveCase(word, us) + "'",
		})
	}
	return violations
}

// misspellFix rewrites with case-preserving replacement. The script's
// lowercasing bug (re.sub with a fixed lowercase replacement under
// re.IGNORECASE) is NOT replicated — case is preserved per D9.
func misspellFix(relPath string, src []byte, opts Options) ([]byte, []Violation) {
	if isMisspellSkipped(relPath) {
		return src, nil
	}
	fixed := misspellRe.ReplaceAllStringFunc(string(src), func(word string) string {
		return preserveCase(word, spellingMap[strings.ToLower(word)])
	})
	if fixed == string(src) {
		return src, nil
	}
	// Check already recorded these hits; Fix reports only new findings.
	return []byte(fixed), nil
}
