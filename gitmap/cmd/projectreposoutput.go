// Package cmd — projectreposoutput.go formats project query output.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// printProjectsTerminal prints projects in labeled terminal format.
func printProjectsTerminal(projects []model.DetectedProject) {
	for _, p := range projects {
		slug := deriveProjectSlug(p)
		fmt.Printf("  %-6s %s\n", p.ProjectType, p.ProjectName)
		fmt.Printf("         Repo:      %s\n", slug)
		fmt.Printf("         Path:      %s\n", p.AbsolutePath)
		fmt.Printf("         Indicator: %s\n", p.PrimaryIndicator)
		fmt.Printf("         → gitmap cd %s\n\n", slug)
	}
}

// deriveProjectSlug returns the best cd-friendly name for a project.
// It prefers RepoName (the Git folder name) and falls back to the
// last segment of AbsolutePath.
func deriveProjectSlug(p model.DetectedProject) string {
	if p.RepoName != "" {
		return p.RepoName
	}

	return filepath.Base(p.AbsolutePath)
}

// printProjectsSummary prints count summary after the list.
func printProjectsSummary(projects []model.DetectedProject) {
	fmt.Fprintf(os.Stderr, constants.MsgProjectListCount, len(projects))
}

// printProjectsJSON prints projects as formatted JSON via stablejson.
func printProjectsJSON(projects []model.DetectedProject) {
	if err := encodeProjectReposJSON(os.Stdout, projects); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode projects JSON: %v\n", err)
	}
}
