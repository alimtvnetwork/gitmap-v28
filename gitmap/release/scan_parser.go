package release

import (
	"regexp"
)

var versionRegex = regexp.MustCompile(`(?i)(?:bump version to|release[ :]*)\s*(v?\d+\.\d+\.\d+(?:-[a-zA-Z0-9.-]+)?)`)

func ParseVersionFromCommit(commitMessage string) (string, bool) {
	matches := versionRegex.FindStringSubmatch(commitMessage)
	if len(matches) < 2 {
		return "", false
	}
	isFound := true
	return matches[1], isFound
}
