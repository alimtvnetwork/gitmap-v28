package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var phpFnRe = lazyregex.New(lazyregex.FuncIntelPhp)

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
