package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var rustFnRe = lazyregex.New(lazyregex.FuncIntelRust)

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
