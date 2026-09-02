package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runFindDuplicatesGit() error {
	mainDB, err := store.OpenDefault()
	if err != nil {
		fmt.Println("  " + constants.ColorYellow + "Gitmap database not available." + constants.ColorReset)
		return nil
	}
	defer mainDB.Close()

	repos, err := mainDB.ListRepos()
	if err != nil || len(repos) == 0 {
		fmt.Println("  " + constants.ColorDim + "No repositories tracked in Gitmap database." + constants.ColorReset)
		return nil
	}
	dupGroups := groupGitDuplicates(repos)
	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ Git: No duplicate tracked repositories found. Total active: %d%s\n\n",
			constants.ColorGreen, len(repos), constants.ColorReset)
		return nil
	}
	printGitDupFindings(dupGroups)
	printGitRemediations(dupGroups)
	return nil
}

func groupGitDuplicates(repos []model.ScanRecord) map[string][]model.ScanRecord {
	groups := make(map[string][]model.ScanRecord)
	for _, r := range repos {
		remote := r.DiscoveredURL
		if remote == "" {
			remote = r.HTTPSUrl
		}
		if remote == "" {
			remote = r.SSHUrl
		}
		key := strings.ToLower(strings.TrimSuffix(strings.TrimSpace(remote), ".git"))
		if key == "" {
			key = strings.ToLower(filepath.Clean(r.AbsolutePath))
		}
		groups[key] = append(groups[key], r)
	}
	dupGroups := make(map[string][]model.ScanRecord)
	for k, list := range groups {
		if len(list) > 1 {
			dupGroups[k] = list
		}
	}
	return dupGroups
}

func printGitDupFindings(dupGroups map[string][]model.ScanRecord) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Git Tracked Duplicate Repositories ──" + constants.ColorReset)
	totalDups := 0
	for _, list := range dupGroups {
		totalDups += len(list) - 1
	}
	fmt.Printf("  Found %s%d%s duplicate repository group(s) (%s%d%s duplicate records total):\n\n",
		constants.ColorWhite, len(dupGroups), constants.ColorReset,
		constants.ColorYellow, totalDups, constants.ColorReset)

	groupNum := 1
	for key, list := range dupGroups {
		fmt.Printf("  Group %d: Target: %s%s%s (%d entries)\n", groupNum, constants.ColorCyan, key, constants.ColorReset, len(list))
		fmt.Printf("    %-6s %-24s %s\n", "ID", "SLUG", "PATH")
		fmt.Println("    " + strings.Repeat("─", 74))
		for _, r := range list {
			fmt.Printf("    %-6d %-24s %s\n", r.ID, truncateStr(r.Slug, 23), truncateStr(r.AbsolutePath, 42))
		}
		fmt.Println()
		groupNum++
	}
}

func printGitRemediations(dupGroups map[string][]model.ScanRecord) {
	var sampleID int64
	var sampleSlug string
	for _, list := range dupGroups {
		if len(list) > 1 {
			sampleID = list[1].ID
			sampleSlug = list[1].Slug
			break
		}
	}
	fmt.Println("  " + constants.ColorCyan + "Remediation & Fix Commands for Git Repositories:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 74))
	fmt.Printf("  ● Fix Single (Untrack duplicate record from database without deleting folder):\n")
	fmt.Printf("    %sgitmap rm --db-only %d%s\n", constants.ColorGreen, sampleID, constants.ColorReset)
	fmt.Printf("    %sgitmap rm --db-only %s%s\n\n", constants.ColorGreen, sampleSlug, constants.ColorReset)
	fmt.Printf("  ● Fix Single (Delete duplicate clone directory & untrack completely):\n")
	fmt.Printf("    %sgitmap rm %s%s\n\n", constants.ColorGreen, sampleSlug, constants.ColorReset)
	fmt.Printf("  ● Fix All Together (Auto-clean repeated clone directories & sync state):\n")
	fmt.Printf("    %sgitmap clone --fix%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %sgitmap rescan%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %sgitmap reconcile%s\n\n", constants.ColorGreen, constants.ColorReset)
}
