package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

var pyDefRe = regexp.MustCompile(`^def\s+([a-z_][A-Za-z0-9_]*)\s*\(`)

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
