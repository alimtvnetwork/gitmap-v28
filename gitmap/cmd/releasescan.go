package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/scanner"
)

// scanReleaseTarget represents one repo discovered by the scan-dir release
// planner together with its current and proposed-next semver.
type scanReleaseTarget struct {
	AbsolutePath string
	RelativePath string
	Current      release.Version
	Next         release.Version
}

// tryRunReleaseScanDir attempts the multi-repo bare-release flow. Returns
// true when scan-dir mode handled the command (with or without releases),
// false when the caller should fall back to the existing self-release path.
func tryRunReleaseScanDir(yes bool) bool {
	cwd, err := os.Getwd()
	if err != nil {
		return false
	}

	repos, err := scanner.ScanDir(cwd, defaultScanExcludes())
	if err != nil || len(repos) == 0 {
		return false
	}

	targets := planScanReleaseTargets(repos)
	if len(targets) == 0 {
		return false
	}

	executeScanReleasePlan(targets, yes)

	return true
}

// defaultScanExcludes returns the dir-name set we never descend into.
func defaultScanExcludes() []string {
	return []string{"node_modules", ".gitmap", "vendor", "dist", "build"}
}

// planScanReleaseTargets keeps only repos that have a prior release on disk
// and computes the next minor version for each.
func planScanReleaseTargets(repos []scanner.RepoInfo) []scanReleaseTarget {
	targets := make([]scanReleaseTarget, 0, len(repos))
	for _, info := range repos {
		target, ok := planOneScanTarget(info)
		if !ok {
			continue
		}
		targets = append(targets, target)
	}

	return targets
}

// planOneScanTarget reads the per-repo latest.json and bumps minor.
func planOneScanTarget(info scanner.RepoInfo) (scanReleaseTarget, bool) {
	latestPath := filepath.Join(info.AbsolutePath, constants.DefaultReleaseDir, constants.DefaultLatestFile)
	if _, err := os.Stat(latestPath); err != nil {
		return scanReleaseTarget{}, false
	}

	current, next, ok := readRepoNextMinor(info.AbsolutePath)
	if !ok {
		return scanReleaseTarget{}, false
	}

	return scanReleaseTarget{
		AbsolutePath: info.AbsolutePath,
		RelativePath: info.RelativePath,
		Current:      current,
		Next:         next,
	}, true
}

// readRepoNextMinor temporarily switches into the repo so the release
// helpers (which read DefaultReleaseDir relative to cwd) see its manifest.
func readRepoNextMinor(repoDir string) (release.Version, release.Version, bool) {
	origDir, err := os.Getwd()
	if err != nil {
		return release.Version{}, release.Version{}, false
	}
	if err := os.Chdir(repoDir); err != nil {
		return release.Version{}, release.Version{}, false
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			fmt.Fprintf(os.Stderr, "failed to restore original directory: %v\n", err)
		}
	}()

	return peekNextMinorVersion()
}

// executeScanReleasePlan prints the summary, prompts once, then runs each
// release. Skips repos that fail; aggregates errors at the end.
func executeScanReleasePlan(targets []scanReleaseTarget, yes bool) {
	printScanReleaseSummary(targets)
	if !confirmScanReleasePlan(yes) {
		fmt.Print(constants.MsgReleaseScanAborted)

		return
	}

	failed := runScanReleaseTargets(targets)
	if failed > 0 {
		fmt.Fprintf(os.Stderr, constants.MsgReleaseScanPartial, failed, len(targets))

		return
	}

	fmt.Printf(constants.MsgReleaseScanDone, len(targets))
}

// printScanReleaseSummary prints the list of planned per-repo bumps.
func printScanReleaseSummary(targets []scanReleaseTarget) {
	fmt.Printf(constants.MsgReleaseScanHeader, len(targets))
	for _, t := range targets {
		fmt.Printf(constants.MsgReleaseScanRow, t.RelativePath, t.Current.String(), t.Next.String())
	}
}

// confirmScanReleasePlan prompts once for the whole batch.
func confirmScanReleasePlan(yes bool) bool {
	if yes {
		fmt.Print(constants.MsgReleaseScanYes)

		return true
	}
	fmt.Print(constants.MsgReleaseScanPrompt)

	return readYesNo()
}

// runScanReleaseTargets executes one minor-bump release per target.
// Returns the count of failures.
func runScanReleaseTargets(targets []scanReleaseTarget) int {
	failed := 0
	for _, t := range targets {
		if !runOneScanRelease(t) {
			failed++
		}
	}

	return failed
}

// runOneScanRelease releases a single repo using the existing workflow,
// reusing the persist-to-DB path so the new Release.RepoId FK is honored.
func runOneScanRelease(t scanReleaseTarget) bool {
	fmt.Printf(constants.MsgReleaseScanRunning, t.RelativePath, t.Next.String())
	origDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgReleaseScanFail, t.RelativePath, err)

		return false
	}
	if err := os.Chdir(t.AbsolutePath); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgReleaseScanFail, t.RelativePath, err)

		return false
	}
	defer func() {
		if err := os.Chdir(origDir); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgReleaseScanFail, t.RelativePath, err)
		}
	}()

	return invokeScanRelease(t)
}

// invokeScanRelease wraps release.Execute + DB persistence with consistent
// error handling so the caller stays under the func length limit.
func invokeScanRelease(t scanReleaseTarget) bool {
	opts := release.Options{Bump: constants.BumpMinor, Yes: true}
	if err := release.Execute(opts); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgReleaseScanFail, t.RelativePath, err)

		return false
	}
	persistReleaseToDB()

	return true
}
