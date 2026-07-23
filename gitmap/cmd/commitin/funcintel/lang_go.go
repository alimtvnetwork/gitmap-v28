package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

var goFuncRe = regexp.MustCompile(`^func\s+(?:\([^)]*\)\s+)?([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

type goDetector struct{}

func (goDetector) Detect(prev, next string) []string {
	return addedNames(
		extractByRegex(prev, matchGoFunc),
		extractByRegex(next, matchGoFunc),
	)
}

func matchGoFunc(line string) (string, bool) {
	m := goFuncRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func init() { register(constants.CommitInLanguageGo, goDetector{}) }
