package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var pyDefRe = lazyregex.New(lazyregex.FuncIntelPython)

type pyDetector struct{}

func (pyDetector) Detect(prev, next string) []string {
	return addedNames(extractByRegex(prev, matchPy), extractByRegex(next, matchPy))
}

func matchPy(line string) (string, bool) {
	m := pyDefRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}

	return m[1], true
}

func init() { register(constants.CommitInLanguagePython, pyDetector{}) }
