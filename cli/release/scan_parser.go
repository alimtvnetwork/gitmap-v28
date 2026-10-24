package release

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
)

var versionRegex = lazyregex.New(`(?i)^(?:(?:chore\(release\)|bump(?:\s+version)?(?:\s+to)?|release(?:\([^)]+\))?)[\s:]+(?:(?:bump|version|to)[\s:]+)*|\s*)(v?\d+\.\d+\.\d+(?:-[a-zA-Z0-9.-]+)?)`)

func ParseVersionFromCommit(commitMessage string) (string, bool) {
	matches := versionRegex.FindStringSubmatch(commitMessage)
	if len(matches) < 2 {
		return "", false
	}

	isFound := true

	return matches[1], isFound
}
