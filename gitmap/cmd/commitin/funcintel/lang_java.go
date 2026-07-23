package funcintel

import (
	"regexp"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

var javaMethodRe = regexp.MustCompile(`^\s*(?:public|protected|private)?\s*(?:static\s+)?[A-Za-z_][A-Za-z0-9_<>\[\]]*\s+([A-Za-z_][A-Za-z0-9_]*)\s*\([^)]*\)\s*\{`)

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
