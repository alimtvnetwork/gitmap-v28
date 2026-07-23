package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var phpFnRe = regexp.MustCompile(`^(?:(?:public|protected|private)\s+)?function\s+([A-Za-z_][A-Za-z0-9_]*)\s*\(`)

type phpDetector struct{}

func (phpDetector) Detect(prev, next string) []string {
	return addedNames(extractByRegex(prev, matchPhp), extractByRegex(next, matchPhp))
}

func matchPhp(line string) (string, bool) {
	m := phpFnRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func init() { register(constants.CommitInLanguagePhp, phpDetector{}) }
