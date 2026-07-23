package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var rustFnRe = regexp.MustCompile(`^(?:pub(?:\([^)]*\))?\s+)?fn\s+([a-z_][A-Za-z0-9_]*)\s*[<(]`)

type rustDetector struct{}

func (rustDetector) Detect(prev, next string) []string {
	return addedNames(extractByRegex(prev, matchRust), extractByRegex(next, matchRust))
}

func matchRust(line string) (string, bool) {
	m := rustFnRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func init() { register(constants.CommitInLanguageRust, rustDetector{}) }
