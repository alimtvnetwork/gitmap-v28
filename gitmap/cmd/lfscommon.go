package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// lfsCommonPatterns is the curated default set of file extensions that
// gitmap lfs-common will register with Git LFS in the current repo.
//
// Order is preserved so the resulting .gitattributes diff is stable and
// predictable across runs and machines.
var lfsCommonPatterns = []string{
	"*.pptx",
	"*.ppt",
	"*.eps",
	"*.psd",
	"*.ttf",
	"*.wott",
	"*.svg",
	"*.ai",
	"*.jpg",
	"*.bmp",
	"*.png",
	"*.zip",
	"*.gz",
	"*.tar",
	"*.rar",
	"*.7z",
	"*.mp4",
	"*.aep",
}

// lfsCommonFlags holds parsed flags for the lfs-common command.
type lfsCommonFlags struct {
	dryRun bool
}

// runLFSCommon implements `gitmap lfs-common`. It ensures Git LFS is
// installed in the current repo and tracks each entry from
// lfsCommonPatterns via `git lfs track`, which writes (or merges) the
// canonical "<pattern> filter=lfs diff=lfs merge=lfs -text" lines into
// .gitattributes. Idempotent: re-running prints which patterns were
// added vs. already tracked.
func runLFSCommon(args []string) {
	checkHelp("lfs-common", args)

	flags := parseLFSCommonFlags(args)

	if !insideGitRepo() {
		fmt.Fprintln(os.Stderr, "  ✗ Not inside a Git repository.")
		fmt.Fprintln(os.Stderr, "    Run this command from the root of a repo (where .git/ lives).")
		os.Exit(1)
	}

	if !lfsAvailable() {
		fmt.Fprintln(os.Stderr, "  ✗ Git LFS is not installed or not on PATH.")
		fmt.Fprintln(os.Stderr, "    Install it from https://git-lfs.com and re-run.")
		os.Exit(1)
	}

	printLFSCommonBanner(flags.dryRun)

	if flags.dryRun {
		printLFSCommonDryRun()

		return
	}

	if err := runGitLFSInstall(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ git lfs install failed: %v\n", err)
	}

	added, existing, failed := trackLFSPatterns(lfsCommonPatterns)
	printLFSCommonSummary(added, existing, failed)
}

// parseLFSCommonFlags parses CLI flags for lfs-common.
func parseLFSCommonFlags(args []string) lfsCommonFlags {
	fs := flag.NewFlagSet(constants.CmdLFSCommon, flag.ExitOnError)
	dryRun := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Could not parse flags: %v\n", err)
		os.Exit(1)
	}

	return lfsCommonFlags{dryRun: *dryRun}
}

// insideGitRepo returns true when the current working directory is part
// of a Git working tree.
func insideGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(out)) == "true"
}

// lfsAvailable returns true when `git lfs` is installed and runnable.
func lfsAvailable() bool {
	cmd := exec.Command("git", "lfs", "version")

	return cmd.Run() == nil
}

// runGitLFSInstall runs `git lfs install` for the current repo. It is
// safe to call repeatedly — Git LFS treats it as idempotent.
func runGitLFSInstall() error {
	cmd := exec.Command("git", "lfs", "install", "--local")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// trackLFSPatterns runs `git lfs track <pattern>` for each entry and
// returns counts of newly-added, already-tracked, and failed patterns.
//
// We detect "already tracked" by inspecting `git lfs track` output and
// by comparing .gitattributes against the pattern before invoking track.
func trackLFSPatterns(patterns []string) (added, existing []string, failed []lfsTrackFailure) {
	preExisting := loadTrackedPatterns()

	for _, p := range patterns {
		if preExisting[p] {
			existing = append(existing, p)

			continue
		}

		if err := trackOnePattern(p); err != nil {
			failed = append(failed, lfsTrackFailure{Pattern: p, Err: err})

			continue
		}

		added = append(added, p)
	}

	return added, existing, failed
}

// lfsTrackFailure pairs a pattern with the error encountered tracking it.
type lfsTrackFailure struct {
	Pattern string
	Err     error
}

// trackOnePattern runs `git lfs track "<pattern>"` for a single pattern.
func trackOnePattern(pattern string) error {
	cmd := exec.Command("git", "lfs", "track", pattern)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out)))
	}

	return nil
}

// loadTrackedPatterns parses .gitattributes (if present) and returns a
// set of LFS-tracked patterns already on disk.
func loadTrackedPatterns() map[string]bool {
	tracked := map[string]bool{}

	root, err := gitTopLevel()
	if err != nil {
		return tracked
	}

	data, err := os.ReadFile(filepath.Join(root, ".gitattributes"))
	if err != nil {
		return tracked
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "filter=lfs") {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) > 0 {
			tracked[fields[0]] = true
		}
	}

	return tracked
}

// (gitTopLevel is provided by as.go; we reuse it here.)

// printLFSCommonBanner prints the header for the run.
func printLFSCommonBanner(dryRun bool) {
	fmt.Println()
	fmt.Printf("  %s■ gitmap lfs-common —%s tracking common binary types with Git LFS\n",
		constants.ColorCyan, constants.ColorReset)
	if dryRun {
		fmt.Printf("  %s[dry-run]%s no files will be modified\n",
			constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println()
}

// printLFSCommonDryRun lists the patterns that would be tracked.
func printLFSCommonDryRun() {
	preExisting := loadTrackedPatterns()
	for _, p := range lfsCommonPatterns {
		status := "would add"
		color := constants.ColorGreen
		if preExisting[p] {
			status = "already tracked"
			color = constants.ColorDim
		}
		fmt.Printf("  %s%-15s%s %s\n", color, status, constants.ColorReset, p)
	}
	fmt.Println()
}

// printLFSCommonSummary prints the per-pattern results plus a total line.
func printLFSCommonSummary(added, existing []string, failed []lfsTrackFailure) {
	for _, p := range added {
		fmt.Printf("  %s+ added%s          %s\n", constants.ColorGreen, constants.ColorReset, p)
	}
	for _, p := range existing {
		fmt.Printf("  %s· already tracked%s %s\n", constants.ColorDim, constants.ColorReset, p)
	}
	for _, f := range failed {
		fmt.Printf("  %s✗ failed%s         %s — %v\n",
			constants.ColorYellow, constants.ColorReset, f.Pattern, f.Err)
	}

	fmt.Println()
	fmt.Printf("  %sSummary:%s %d added, %d already tracked, %d failed (of %d total)\n",
		constants.ColorCyan, constants.ColorReset,
		len(added), len(existing), len(failed), len(lfsCommonPatterns))

	if len(added) > 0 {
		fmt.Println()
		fmt.Printf("  %sNext step:%s commit the updated .gitattributes:\n",
			constants.ColorYellow, constants.ColorReset)
		fmt.Println("    git add .gitattributes")
		fmt.Println("    git commit -m \"chore: track common binary types with Git LFS\"")
	}
	fmt.Println()
}
