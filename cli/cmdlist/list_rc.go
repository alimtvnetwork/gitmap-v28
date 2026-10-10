package cmdlist

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
)

// RunListRC inspects discovered repo-cache manifests and renders an executive TreeView summary.
func RunListRC(args []string) error {
	repoCacheRoot := cmdclone.ResolveRepoCacheRoot()
	manifests, err := cmdclone.DiscoverRCManifests(repoCacheRoot)
	if err != nil {
		return fmt.Errorf("failed to discover repo-cache manifests: %w", err)
	}

	displayRoot := repoCacheRoot
	if !strings.HasSuffix(displayRoot, "/") && !strings.HasSuffix(displayRoot, "\\") {
		displayRoot += "/"
	}

	fmt.Printf("=== GitMap Repo-Cache Manifests (%s) ===\n\n", displayRoot)

	hasManifests := len(manifests) > 0
	if !hasManifests {
		fmt.Printf("No manifests found in repo-cache (%s).\n", repoCacheRoot)
		return nil
	}

	totalAllRepos := 0
	totalAllSSH := 0
	totalAllHTTPS := 0

	for _, m := range manifests {
		totalAllRepos += m.TotalRepos
		totalAllSSH += m.SSHCount
		totalAllHTTPS += m.HTTPSCount

		fmt.Printf("📦 [%d] %s\n", m.Index, m.RelativePath)
		fmt.Printf("    Total Repos: %d (%d Public HTTPS, %d SSH)\n", m.TotalRepos, m.HTTPSCount, m.SSHCount)

		previewLimit := 3
		if len(m.Entries) <= previewLimit {
			for idx, e := range m.Entries {
				connector := "├── "
				if idx == len(m.Entries)-1 {
					connector = "└── "
				}
				protoTag := formatProtocolTag(e.CloneUrl)
				branchInfo := formatBranchInfo(e.Branch)
				fmt.Printf("    %s%s %s/%s%s\n", connector, protoTag, e.Owner, e.RepoName, branchInfo)
			}
		} else {
			for idx := 0; idx < previewLimit; idx++ {
				e := m.Entries[idx]
				protoTag := formatProtocolTag(e.CloneUrl)
				branchInfo := formatBranchInfo(e.Branch)
				fmt.Printf("    ├── %s %s/%s%s\n", protoTag, e.Owner, e.RepoName, branchInfo)
			}
			remaining := len(m.Entries) - previewLimit
			fmt.Printf("    └── ... (%d more)\n", remaining)
		}

		fmt.Printf("    💡 Command to import: gitmap clone --rc %d (or: gitmap cfr --rc %d)\n\n", m.Index, m.Index)
	}

	manifestWord := "manifests"
	if len(manifests) == 1 {
		manifestWord = "manifest"
	}
	repoWord := "repositories"
	if totalAllRepos == 1 {
		repoWord = "repository"
	}

	fmt.Printf("Summary: %d %s found | %d total %s (%d Public HTTPS, %d SSH)\n",
		len(manifests), manifestWord, totalAllRepos, repoWord, totalAllHTTPS, totalAllSSH)

	return nil
}

func formatProtocolTag(url string) string {
	proto := cmdclone.ClassifyProtocol(url)
	if proto == "[SSH]" {
		return "[SSH]         "
	}
	return "[Public HTTPS]"
}

func formatBranchInfo(branch string) string {
	trimmed := strings.TrimSpace(branch)
	if len(trimmed) > 0 {
		return fmt.Sprintf(" (branch: %s)", trimmed)
	}
	return ""
}
