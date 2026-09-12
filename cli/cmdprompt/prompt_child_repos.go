// Package cmdprompt — prompt_child_repos.go discovers repositories inside a workspace folder.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

func DiscoverPromptChildRepos(rootDir string) ([]string, error) {
	return fsutil.DiscoverTopLevelGitRepos(rootDir)
}
