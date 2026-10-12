// Package cmdrelease implements release management commands.
package cmdrelease

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// RemoveBetweenOptions configures a range-based release tag removal.
type RemoveBetweenOptions struct {
	Start         string
	End           string
	KeepRemote    bool
	DryRun        bool
	Confirm       bool
	ConfirmSource string
	Pins          []string
}

// RunReleaseRemoveBetween executes `gitmap release remove-between` (alias `rrb`).
func RunReleaseRemoveBetween(args []string) error {
	checkHelp(constants.CmdReleaseRemoveBetween, args)

	opts, err := parseRemoveBetweenFlags(args)
	if err != nil {
		return err
	}

	lo, hi, err := resolveReleaseBounds(opts.Start, opts.End)
	if err != nil {
		return err
	}

	workDir := "."
	persistentPins, _ := LoadPinnedTags(workDir)
	allPins := append(persistentPins, opts.Pins...)

	allTags, err := fetchRepositoryTags()
	if err != nil {
		return err
	}

	toDelete, toKeep := filterTagsBetweenBounds(allTags, lo, hi, allPins)

	// 1. Always evaluate and print the automatic dry-run preview first
	printAutoDryRunEvaluation(lo, hi, toDelete, toKeep, allPins)

	if len(toDelete) == 0 {
		fmt.Printf("\n\x1b[32m✔ No intermediate tags found to remove between %s and %s.\x1b[0m\n", lo.String(), hi.String())
		return nil
	}

	// 2. Explicit --dry-run exits after the preview without prompt or modifications
	if opts.DryRun {
		fmt.Println("\n\x1b[33mℹ Dry-run preview complete (--dry-run). No tags or releases were removed.\x1b[0m")
		return nil
	}

	// 3. Prompt before removal; ONLY bypassed by --confirm or -y / --yes
	if opts.Confirm {
		fmt.Printf("\n\x1b[1;33m⚡ Confirmation prompt bypassed via %s. Proceeding with removal...\x1b[0m\n", opts.ConfirmSource)
	} else {
		if !confirmBatchDeletion(len(toDelete), lo.String(), hi.String()) {
			return apperror.NewSimple("release remove-between: operation cancelled (to bypass confirmation prompt, pass --confirm or -y)", "E9001")
		}
	}

	steps := applyBatchTagRemoval(toDelete, opts.KeepRemote)
	optimizeRepositoryPostPrune(opts.KeepRemote)

	fmt.Printf("\n\x1b[1;32m✅ release remove-between complete: %d tags removed (%s), %d kept.\x1b[0m\n",
		len(toDelete), strings.Join(steps, ", "), len(toKeep))
	fmt.Printf("\x1b[2mRetained milestone anchors: %s and %s\x1b[0m\n", lo.String(), hi.String())

	return nil
}

func parseRemoveBetweenFlags(args []string) (RemoveBetweenOptions, error) {
	fs := flag.NewFlagSet(constants.CmdReleaseRemoveBetween, flag.ExitOnError)
	keepRemote := fs.Bool("keep-remote", false, "Do not delete the remote tags from origin")
	dryRun := fs.Bool("dry-run", false, "Preview which tags would be removed without executing")
	confirm := fs.Bool("confirm", false, "Bypass interactive confirmation prompt and proceed with removal")
	yes := fs.Bool("y", false, "Bypass interactive confirmation prompt and proceed with removal")
	yesLong := fs.Bool("yes", false, "Bypass interactive confirmation prompt and proceed with removal")
	pinStr := fs.String("pin", "", "Comma-separated list of additional tags to preserve")
	keepStr := fs.String("keep", "", "Alias for --pin")

	_ = fs.Parse(reorderFlagsBeforeArgs(args))

	var rawTokens []string
	for _, a := range fs.Args() {
		parts := strings.Split(a, ",")
		for _, p := range parts {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				rawTokens = append(rawTokens, trimmed)
			}
		}
	}

	if len(rawTokens) < 2 {
		return RemoveBetweenOptions{}, apperror.NewSimple(
			"release remove-between: requires both start and end version bounds (e.g. `gitmap rrb 4.162, 4.190`)",
			"E9002",
		)
	}

	var pins []string
	mergedPin := *pinStr
	if mergedPin == "" {
		mergedPin = *keepStr
	}
	if mergedPin != "" {
		for _, p := range strings.Split(mergedPin, ",") {
			trimmed := strings.TrimSpace(p)
			if trimmed != "" {
				pins = append(pins, trimmed)
			}
		}
	}

	isConfirmed := false
	confirmSource := ""
	if *confirm {
		isConfirmed = true
		confirmSource = "--confirm"
	} else if *yes {
		isConfirmed = true
		confirmSource = "-y"
	} else if *yesLong {
		isConfirmed = true
		confirmSource = "--yes"
	}

	return RemoveBetweenOptions{
		Start:         rawTokens[0],
		End:           rawTokens[1],
		KeepRemote:    *keepRemote,
		DryRun:        *dryRun,
		Confirm:       isConfirmed,
		ConfirmSource: confirmSource,
		Pins:          pins,
	}, nil
}

func resolveReleaseBounds(sRaw, eRaw string) (release.Version, release.Version, error) {
	sVer, err := release.Parse(sRaw)
	if err != nil {
		return release.Version{}, release.Version{}, apperror.WrapSimple(err, "E9003")
	}

	eVer, err := release.Parse(eRaw)
	if err != nil {
		return release.Version{}, release.Version{}, apperror.WrapSimple(err, "E9004")
	}

	if sVer.GreaterThan(eVer) {
		fmt.Printf("\x1b[33mℹ Inverted bounds provided (%s > %s). Normalizing to [%s .. %s].\x1b[0m\n",
			sVer.String(), eVer.String(), eVer.String(), sVer.String())
		return eVer, sVer, nil
	}

	return sVer, eVer, nil
}

func fetchRepositoryTags() ([]string, error) {
	cmd := exec.Command("git", "tag", "-l")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, apperror.WrapSimple(err, "E9005")
	}

	var tags []string
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			tags = append(tags, line)
		}
	}

	return tags, nil
}

func filterTagsBetweenBounds(
	allTags []string,
	lo, hi release.Version,
	pinned []string,
) ([]string, []string) {
	var toDelete []string
	var kept []string

	for _, tag := range allTags {
		v, err := release.Parse(tag)
		if err != nil {
			kept = append(kept, tag)
			continue
		}

		if isVersionEqual(v, lo) || isVersionEqual(v, hi) {
			kept = append(kept, tag)
			continue
		}

		if IsTagPinned(tag, pinned) {
			kept = append(kept, tag)
			continue
		}

		if isTagStrictlyBetween(v, lo, hi) {
			toDelete = append(toDelete, tag)
		} else {
			kept = append(kept, tag)
		}
	}

	return toDelete, kept
}

func isVersionEqual(a, b release.Version) bool {
	return a.Major == b.Major && a.Minor == b.Minor && a.Patch == b.Patch && a.PreRelease == b.PreRelease
}

func isTagStrictlyBetween(v, lo, hi release.Version) bool {
	// Base SemVer comparison
	if v.GreaterThan(lo) && hi.GreaterThan(v) {
		return true
	}

	// Handle pre-release tags tied to lower bound (e.g. v4.162.0-tmp, v4.162.1..5)
	if v.Major == lo.Major && v.Minor == lo.Minor {
		if v.Patch == lo.Patch && v.IsPreRelease() {
			return true
		}
		if v.Patch > lo.Patch && (v.Major < hi.Major || v.Minor < hi.Minor || v.Patch < hi.Patch) {
			return true
		}
	}

	return false
}

func printAutoDryRunEvaluation(lo, hi release.Version, toDelete, toKeep, pinned []string) {
	fmt.Printf("\x1b[1;36m▶ release remove-between\x1b[0m: [%s .. %s]\n", lo.String(), hi.String())
	fmt.Printf("  • Milestone anchors preserved : \x1b[1;32m%s\x1b[0m and \x1b[1;32m%s\x1b[0m\n", lo.String(), hi.String())
	if len(pinned) > 0 {
		fmt.Printf("  • Pinned tags preserved       : \x1b[1;33m%s\x1b[0m\n", strings.Join(pinned, ", "))
	}
	fmt.Printf("  • Candidate tags to remove    : \x1b[1;31m%d\x1b[0m\n", len(toDelete))
	fmt.Printf("  • Total tags retained         : \x1b[1;32m%d\x1b[0m\n", len(toKeep))

	if len(toDelete) > 0 {
		fmt.Println("\n\x1b[1;33m📋 Candidate intermediate tag(s) queued for removal:\x1b[0m")
		for _, t := range toDelete {
			fmt.Printf("    \x1b[31m✖ %s\x1b[0m\n", t)
		}
	}
}

func confirmBatchDeletion(count int, lo, hi string) bool {
	fmt.Printf("\n\x1b[1;31m⚠ Permanently delete %d intermediate tag(s) between %s and %s locally and on origin? [y/N]: \x1b[0m", count, lo, hi)
	reader := bufio.NewReader(os.Stdin)
	reply, err := reader.ReadString('\n')
	if err != nil && len(reply) == 0 {
		return false
	}
	t := strings.TrimSpace(reply)
	return strings.EqualFold(t, "y") || strings.EqualFold(t, "yes")
}

func applyBatchTagRemoval(tags []string, keepRemote bool) []string {
	var done []string

	// 1. Delete GitHub Release objects via gh if available
	for _, t := range tags {
		_ = exec.Command("gh", "release", "delete", t, "-y", "--cleanup-tag=false").Run()
	}
	done = append(done, "gh releases")

	// 2. Delete remote tags in batch chunks
	if !keepRemote {
		batchSize := 40
		for i := 0; i < len(tags); i += batchSize {
			end := i + batchSize
			if end > len(tags) {
				end = len(tags)
			}
			chunk := tags[i:end]
			args := append([]string{"push", "origin", "--delete"}, chunk...)
			_ = exec.Command("git", args...).Run()
		}
		done = append(done, "remote origin tags")
	}

	// 3. Delete local tags
	args := append([]string{"tag", "-d"}, tags...)
	_ = exec.Command("git", args...).Run()
	done = append(done, "local git tags")

	// 4. Delete sidecar json metadata
	for _, t := range tags {
		v, err := release.Parse(t)
		if err == nil {
			jsonPath := filepath.Join(constants.DefaultReleaseDir, v.String()+constants.ExtJSON)
			_ = os.Remove(jsonPath)
		}
	}
	done = append(done, "metadata sidecars")

	return done
}

func optimizeRepositoryPostPrune(keepRemote bool) {
	fmt.Println("  • Pruning references & optimizing repository pack files...")
	if !keepRemote {
		_ = exec.Command("git", "fetch", "--prune", "--prune-tags", "origin").Run()
	}
	_ = exec.Command("git", "gc", "--prune=now").Run()
}
