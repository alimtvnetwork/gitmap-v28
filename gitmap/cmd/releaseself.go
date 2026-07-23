package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// runReleaseSelf handles the 'release-self' command.
func runReleaseSelf(args []string) {
	checkHelp("release-self", args)
	requireOnline()
	version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, bin, listTargets, noCommit, yes := parseReleaseFlags(args)
	_ = verbose

	if listTargets {
		printListTargets(targets)

		return
	}

	validateReleaseFlags(version, bump, commit, branch)
	executeSelfRelease(version, assets, commit, branch, bump, notes, targets, zipGroups, zipItems, bundleName, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes)
}

// executeSelfRelease builds options and runs the self-release workflow.
func executeSelfRelease(version, assets, commit, branch, bump, notes, targets string, zipGroups, zipItems []string, bundleName string, draft, dryRun, verbose, compress, checksums, bin, noCommit, yes bool) {
	opts := release.Options{
		Version: version, Assets: assets,
		Commit: commit, Branch: branch,
		Bump: bump, Notes: notes, Targets: targets,
		ZipGroups:  zipGroups,
		ZipItems:   zipItems,
		BundleName: bundleName,
		IsDraft:    draft, DryRun: dryRun,
		Verbose:   verbose,
		Compress:  compress,
		Checksums: checksums,
		Bin:       bin,
		NoCommit:  noCommit,
		Yes:       yes,
	}

	err := release.ExecuteSelf(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	persistReleaseToDB()
}
