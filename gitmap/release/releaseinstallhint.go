package release

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printInstallHint prints install one-liner commands if the current repo
// matches the gitmap source repository prefix.
func printInstallHint(v Version) {
	url := getRemoteURL()
	if ShouldPrintInstallHint(url) {
		fmt.Printf(constants.MsgInstallHintHeader, v.String())
		fmt.Print(constants.MsgInstallHintWindows)
		fmt.Print(constants.MsgInstallHintUnix)
	}
}

// ShouldPrintInstallHint returns true if the remote URL matches the
// gitmap source repository prefix.
func ShouldPrintInstallHint(remoteURL string) bool {
	if len(remoteURL) == 0 {
		return false
	}

	normalized := normalizeInstallHintRemoteURL(remoteURL)
	repoName := extractInstallHintRepoName(normalized)

	return isVersionedGitmapRepo(repoName)
}

func normalizeInstallHintRemoteURL(remoteURL string) string {
	normalized := strings.TrimSuffix(strings.ToLower(remoteURL), ".git")
	if idx := strings.Index(normalized, "@"); idx >= 0 {
		return strings.Replace(normalized[idx+1:], ":", "/", 1)
	}

	return normalized
}

func extractInstallHintRepoName(normalized string) string {
	idx := strings.Index(normalized, constants.GitmapRepoOwner)
	if idx < 0 {
		return ""
	}

	tail := normalized[idx+len(constants.GitmapRepoOwner):]
	parts := strings.SplitN(tail, "/", 2)

	return parts[0]
}

func isVersionedGitmapRepo(repoName string) bool {
	if !strings.HasPrefix(repoName, constants.GitmapRepoNamePrefix) {
		return false
	}

	return isNumericString(strings.TrimPrefix(repoName, constants.GitmapRepoNamePrefix))
}

func isNumericString(value string) bool {
	if len(value) == 0 {
		return false
	}

	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return false
		}
	}

	return true
}
