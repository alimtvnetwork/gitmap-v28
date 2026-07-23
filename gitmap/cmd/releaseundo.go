package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// runReleaseUndo handles `gitmap release-undo [vX.Y.Z]`.
//
// Reverses a prior `gitmap release` by:
//  1. Deleting the local annotated tag (`git tag -d vX.Y.Z`).
//  2. Deleting the remote tag (`git push origin :refs/tags/vX.Y.Z`)
//     unless --keep-remote is set.
//  3. Removing the local sidecar `.gitmap/release/vX.Y.Z.json`.
//
// When no version is supplied the newest `.gitmap/release/v*.json` is used.
// The summary line is intentionally copy-friendly so the user can paste it
// into a task-completion report.
func runReleaseUndo(args []string) {
	checkHelp(constants.CmdReleaseUndo, args)

	version, keepRemote, dryRun, yes := parseReleaseUndoFlags(args)

	if version == "" {
		latest, ok := latestReleaseJSONVersion()
		if !ok {
			fmt.Fprintln(os.Stderr, "release-undo: no .gitmap/release/v*.json found and no version argument given")
			os.Exit(1)
		}
		version = latest
	}

	v, err := release.Parse(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReleaseInvalidVersion, version)
		os.Exit(1)
	}
	tag := "v" + v.String()
	jsonPath := filepath.Join(constants.DefaultReleaseDir, v.String()+constants.ExtJSON)

	fmt.Printf("\x1b[1;36m▶ release-undo\x1b[0m: %s\n", tag)
	fmt.Printf("  local tag      : git tag -d %s\n", tag)
	if !keepRemote {
		fmt.Printf("  remote tag     : git push origin :refs/tags/%s\n", tag)
	}
	fmt.Printf("  release json   : %s\n", jsonPath)

	if dryRun {
		fmt.Println("\x1b[33m(dry-run; no changes applied)\x1b[0m")
		return
	}
	if !yes && !confirmUndoRelease(tag) {
		fmt.Fprintln(os.Stderr, "release-undo: aborted")
		os.Exit(1)
	}

	steps := applyReleaseUndo(tag, jsonPath, keepRemote)
	summary := fmt.Sprintf("✅ release-undo complete — %s removed (%s)", tag, strings.Join(steps, ", "))
	fmt.Println("\x1b[1;32m" + summary + "\x1b[0m")
	fmt.Println("\x1b[2m(share this line in your task report)\x1b[0m")
}

// parseReleaseUndoFlags parses CLI flags for release-undo.
func parseReleaseUndoFlags(args []string) (version string, keepRemote, dryRun, yes bool) {
	fs := flag.NewFlagSet(constants.CmdReleaseUndo, flag.ExitOnError)
	fs.BoolVar(&keepRemote, "keep-remote", false, "Do not delete the remote tag")
	fs.BoolVar(&dryRun, "dry-run", false, "Preview without applying")
	fs.BoolVar(&yes, "yes", false, "Skip confirmation")
	fs.BoolVar(&yes, "y", false, "Skip confirmation")
	_ = fs.Parse(reorderFlagsBeforeArgs(args))
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}
	return version, keepRemote, dryRun, yes
}

// confirmUndoRelease prompts the user to confirm tag deletion.
func confirmUndoRelease(tag string) bool {
	fmt.Printf("Delete %s locally and on origin? [y/N]: ", tag)
	var reply string
	_, _ = fmt.Scanln(&reply)
	reply = strings.ToLower(strings.TrimSpace(reply))
	return reply == "y" || reply == "yes"
}

// applyReleaseUndo executes the three undo steps and returns which ones succeeded.
func applyReleaseUndo(tag, jsonPath string, keepRemote bool) []string {
	var done []string
	if err := runGitQuiet("tag", "-d", tag); err == nil {
		done = append(done, "local tag")
	} else {
		fmt.Fprintf(os.Stderr, "  ⚠ local tag delete failed: %v\n", err)
	}
	if !keepRemote {
		if err := runGitQuiet("push", "origin", ":refs/tags/"+tag); err == nil {
			done = append(done, "remote tag")
		} else {
			fmt.Fprintf(os.Stderr, "  ⚠ remote tag delete failed: %v\n", err)
		}
	}
	if err := os.Remove(jsonPath); err == nil {
		done = append(done, "release json")
	} else if !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "  ⚠ json delete failed: %v\n", err)
	}
	return done
}

// runGitQuiet runs `git <args>` suppressing stdout; stderr is preserved.
func runGitQuiet(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// latestReleaseJSONVersion returns the highest semver under .gitmap/release/.
func latestReleaseJSONVersion() (string, bool) {
	entries, err := os.ReadDir(constants.DefaultReleaseDir)
	if err != nil {
		return "", false
	}
	var versions []release.Version
	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, constants.ExtJSON) {
			continue
		}
		raw := strings.TrimSuffix(name, constants.ExtJSON)
		v, err := release.Parse(raw)
		if err != nil {
			continue
		}
		versions = append(versions, v)
	}
	if len(versions) == 0 {
		return "", false
	}
	sort.Slice(versions, func(i, j int) bool {
		a, b := versions[i], versions[j]
		if a.Major != b.Major {
			return a.Major < b.Major
		}
		if a.Minor != b.Minor {
			return a.Minor < b.Minor
		}
		return a.Patch < b.Patch
	})
	return versions[len(versions)-1].String(), true
}
