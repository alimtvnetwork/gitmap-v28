package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var tsArrowRe = lazyregex.New(lazyregex.FuncIntelTs)

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
