package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var goFuncRe = lazyregex.New(lazyregex.FuncIntelGo)

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
