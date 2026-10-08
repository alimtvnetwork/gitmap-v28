package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func hasGitDupFixFlag() bool {
	for _, arg := range os.Args {
		if strings.EqualFold(arg, "--fix") {
			return true
		}
	}
	return false
}

func runFindDuplicatesGit() error {
	mainDB, err := store.OpenDefault()
	if err != nil {
		fmt.Println("  " + constants.ColorYellow + "Gitmap database not available." + constants.ColorReset)

		return nil
	}

	defer mainDB.Close()

	if hasGitDupFixFlag() {
		return runFixDuplicatesGit(mainDB)
	}

	dupGroups, err := mainDB.FindDuplicateRepos()
	if err != nil {
		fmt.Printf("  %s✗ Error querying duplicate repositories: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return nil
	}

	if len(dupGroups) == 0 {
		fmt.Printf("  %s✓ Git: No duplicate tracked repositories found in SQLite database.%s\n\n",
			constants.ColorGreen, constants.ColorReset)

		return nil
	}

	printGitDupGroupFindings(dupGroups)
	printGitDupGroupRemediations(dupGroups)

	return nil
}

func printGitDupGroupFindings(dupGroups []store.DuplicateRepoGroup) {
	fmt.Println()
	fmt.Println("  " + constants.ColorMagenta + "── Git Tracked Duplicate Repositories ──" + constants.ColorReset)
	totalDups := 0
	for _, g := range dupGroups {
		totalDups += len(g.Duplicates)
	}

	fmt.Printf("  Found %s%d%s duplicate repository group(s) (%s%d%s duplicate records total):\n\n",
		constants.ColorWhite, len(dupGroups), constants.ColorReset,
		constants.ColorYellow, totalDups, constants.ColorReset)

	for i, g := range dupGroups {
		all := append([]model.ScanRecord{g.Keeper}, g.Duplicates...)
		fmt.Printf("  Group %d: Target: %s%s%s (%d entries, keeper: ID %d)\n", i+1, constants.ColorCyan, g.CleanKey, constants.ColorReset, g.Count, g.Keeper.ID)
		fmt.Printf("    %-6s %-24s %s\n", "ID", "SLUG", "PATH")
		fmt.Println("    " + strings.Repeat("─", 74))
		for _, r := range all {
			marker := "  "
			if r.ID == g.Keeper.ID {
				marker = "* "
			}
			fmt.Printf("  %s%-6d %-24s %s\n", marker, r.ID, cmddb.TruncateStr(r.Slug, 23), cmddb.TruncateStr(r.AbsolutePath, 42))
		}

		fmt.Println()
	}
}

func printGitDupGroupRemediations(dupGroups []store.DuplicateRepoGroup) {
	var sampleID int64
	var sampleSlug string
	for _, g := range dupGroups {
		if len(g.Duplicates) > 0 {
			sampleID = g.Duplicates[0].ID
			sampleSlug = g.Duplicates[0].Slug
			break
		}
	}

	fmt.Println("  " + constants.ColorCyan + "Remediation & Fix Commands for Git Repositories:" + constants.ColorReset)
	fmt.Println("  " + strings.Repeat("─", 74))
	fmt.Printf("  ● Automated Database Deduplication (Remap FKs & Prune Redundancies):\n")
	fmt.Printf("    %sgitmap find-duplicates git --fix%s\n\n", constants.ColorGreen, constants.ColorReset)
	if sampleID > 0 {
		fmt.Printf("  ● Fix Single (Untrack duplicate record from database without deleting folder):\n")
		fmt.Printf("    %sgitmap rm --db-only %d%s\n", constants.ColorGreen, sampleID, constants.ColorReset)
		fmt.Printf("    %sgitmap rm --db-only %s%s\n\n", constants.ColorGreen, sampleSlug, constants.ColorReset)
	}
	fmt.Printf("  ● Clean Filesystem Clones & Sync State:\n")
	fmt.Printf("    %sgitmap clone --fix%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %sgitmap rescan%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    %sgitmap reconcile%s\n\n", constants.ColorGreen, constants.ColorReset)
}

func runFixDuplicatesGit(mainDB *store.DB) error {
	summary, err := mainDB.DeduplicateRepos(false)
	if err != nil {
		fmt.Printf("  %s✗ Error deduplicating repositories: %v%s\n\n", constants.ColorRed, err, constants.ColorReset)

		return nil
	}

	if summary == nil || summary.RowsPurged == 0 {
		fmt.Printf("  %s✓ Git: No redundant repositories to clean.%s\n\n", constants.ColorGreen, constants.ColorReset)

		return nil
	}

	fmt.Printf("  %s✓ SQLite: Successfully deduplicated %d redundant repository record(s) across %d group(s).%s\n\n",
		constants.ColorGreen, summary.RowsPurged, summary.GroupsFound, constants.ColorReset)

	return nil
}
