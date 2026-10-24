package funcintel

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var javaMethodRe = lazyregex.New(lazyregex.FuncIntelJava)

type javaDetector struct{}

func (javaDetector) Detect(prev, next string) []string {
	return addedNames(extractByRegex(prev, matchJava), extractByRegex(next, matchJava))
}

func matchJava(line string) (string, bool) {
	m := javaMethodRe.FindStringSubmatch(line)
	if m == nil {
		return "", false
	}

	return m[1], true
}

func init() {
	register(constants.CommitInLanguageJava, javaDetector{})
	register(constants.CommitInLanguageCSharp, javaDetector{})
}
