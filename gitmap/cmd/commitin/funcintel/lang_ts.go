package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var tsArrowRe = regexp.MustCompile(`^(?:export\s+)?const\s+([A-Za-z_$][A-Za-z0-9_$]*)\s*=\s*(?:async\s*)?\(`)

type tsDetector struct{}

func (tsDetector) Detect(prev, next string) []string {
	return addedNames(
		extractByRegex(prev, matchTs),
		extractByRegex(next, matchTs),
	)
}

func matchTs(line string) (string, bool) {
	if name, ok := matchJsFunc(line); ok {
		return name, true
	}
	m := tsArrowRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}
	return m[1], true
}

func init() { register(constants.CommitInLanguageTypeScript, tsDetector{}) }
