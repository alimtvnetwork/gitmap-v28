package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var jsFuncRe = lazyregex.New(lazyregex.FuncIntelJavascript)

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
