// Package cmdprompt — prompt_child_repos.go discovers repositories inside a workspace folder.
package cmdprompt

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func DiscoverPromptChildRepos(rootDir string) PromptTargetSliceResult {
	repos, err := fsutil.DiscoverTopLevelGitRepos(rootDir)
	if err != nil {
		return result.NewFailureSlice[string](err)
	}

	return result.OkSlice(repos)
}
