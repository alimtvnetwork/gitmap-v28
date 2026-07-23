package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

var jsFuncRe = regexp.MustCompile(`^(?:export\s+)?(?:async\s+)?function\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*\(`)

type jsDetector struct{}

func (jsDetector) Detect(prev, next string) []string {
	return addedNames(
		extractByRegex(prev, matchJsFunc),
		extractByRegex(next, matchJsFunc),
	)
}

func matchJsFunc(line string) (string, bool) {
	m := jsFuncRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func init() { register(constants.CommitInLanguageJavaScript, jsDetector{}) }
