// Package cmd — prompt_child_repos.go discovers repositories inside a workspace folder.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/fsutil"
)

func DiscoverPromptChildRepos(rootDir string) ([]string, error) {
	return fsutil.DiscoverTopLevelGitRepos(rootDir)
}
