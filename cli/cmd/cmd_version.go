// Package cmd — version_tags.go handles querying and rendering GitHub release tags for GitMap.
package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunGitMapVersionTagsLS() error {
	fmt.Printf("Fetching available GitMap release tags from GitHub (%s/%s)...\n", constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug)
	tags, err := cmdagy.FetchGitHubReleaseTags(constants.UpdateRepoOwner, constants.UpdateCurrentRepoSlug, 30)
	if err != nil {
		fmt.Printf("Note: could not query GitHub API directly (%v); showing current running version.\n", err)
	}
	cmdagy.RenderReleaseTagsTable("GitMap CLI (gitmap-v28)", "v"+constants.Version, tags)
	fmt.Println("To install/switch to an older or specific GitMap release:")
	fmt.Println("  gitmap update <version>         (e.g. gitmap update v6.324.0)")
	fmt.Println("  gitmap update --version <ver>   (e.g. gitmap update --version v6.325.0)")
	return nil
}

// RunAGMVersionTagsLS queries and renders all GitHub release tags for Antigravity Manager.
