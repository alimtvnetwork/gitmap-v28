package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runTempReleaseRemove handles "tr remove <version>|<v1> to <v2>|all".
func runTempReleaseRemove(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrTRRemoveUsage)
		os.Exit(1)
	}

	if args[0] == "all" {
		removeTempReleaseAll()

		return
	}

	if len(args) >= 3 && args[1] == "to" {
		removeTempReleaseRange(args[0], args[2])

		return
	}

	removeTempReleaseSingle(args[0])
}

// removeTempReleaseSingle removes one temp-release branch.
func removeTempReleaseSingle(version string) {
	branchName := resolveTRBranch(version)

	fmt.Printf(constants.MsgTRRemovePrompt, branchName)
	if !confirmAction() {
		return
	}

	removeBranches([]string{branchName})
	fmt.Printf(constants.MsgTRRemovedOne, branchName)
}

// removeTempReleaseRange removes branches from v1 to v2.
func removeTempReleaseRange(from, to string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ DB migration failed: %v\n", err)
	}

	targets := collectRangeTargets(db, from, to)

	if len(targets) == 0 {
		fmt.Print(constants.MsgTRNoneToRemove)

		return
	}

	fmt.Printf(constants.MsgTRRemoveRange, len(targets))
	printBranchList(targets)

	fmt.Print(constants.MsgTRRemoveConfirm)
	if !confirmAction() {
		return
	}

	removeBranches(targets)
	cleanupTRFromDB(db, targets)
	fmt.Printf(constants.MsgTRRemoved, len(targets))
}

// collectRangeTargets finds all branches between from and to (inclusive).
func collectRangeTargets(db *store.DB, from, to string) []string {
	releases, _ := db.ListTempReleases()
	fromBranch := resolveTRBranch(from)
	toBranch := resolveTRBranch(to)

	var targets []string

	inRange := false
	for _, r := range releases {
		if r.Branch == fromBranch {
			inRange = true
		}
		if inRange {
			targets = append(targets, r.Branch)
		}
		if r.Branch == toBranch {
			break
		}
	}

	return targets
}

// removeTempReleaseAll removes all temp-release branches.
func removeTempReleaseAll() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ DB migration failed: %v\n", err)
	}

	releases, listErr := db.ListTempReleases()
	if listErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not list temp releases: %v\n", listErr)
	}
	if len(releases) == 0 {
		fmt.Print(constants.MsgTRNoneToRemove)

		return
	}

	var branches []string
	for _, r := range releases {
		branches = append(branches, r.Branch)
	}

	fmt.Printf(constants.MsgTRRemoveAll, len(branches))
	printBranchList(branches)

	fmt.Print(constants.MsgTRRemoveConfirm)
	if !confirmAction() {
		return
	}

	removeBranches(branches)
	if err := db.DeleteAllTempReleases(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not delete all temp releases from DB: %v\n", err)
	}
	fmt.Printf(constants.MsgTRRemoved, len(branches))
}

// printBranchList prints each branch name on its own line.
func printBranchList(branches []string) {
	for _, b := range branches {
		fmt.Printf(constants.MsgTRRemoveBranch, b)
	}
}

// resolveTRBranch adds the temp-release/ prefix if not present.
func resolveTRBranch(version string) string {
	if strings.HasPrefix(version, constants.TempReleaseBranchPrefix) {
		return version
	}

	return constants.TempReleaseBranchPrefix + version
}

// removeBranches deletes branches locally and from remote.
func removeBranches(branches []string) {
	for _, b := range branches {
		if err := release.DeleteLocalBranch(b); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not delete local branch %s: %v\n", b, err)
		}
	}

	if len(branches) > 0 {
		if err := release.DeleteRemoteBranches(branches); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not delete remote branches: %v\n", err)
		}
	}
}

// cleanupTRFromDB removes temp-release records from the database.
func cleanupTRFromDB(db *store.DB, branches []string) {
	for _, b := range branches {
		if err := db.DeleteTempRelease(b); err != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not delete temp release %s from DB: %v\n", b, err)
		}
	}
}

// confirmAction reads a y/N prompt from stdin.
func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}
