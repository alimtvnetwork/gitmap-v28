// Package cmdrelease implements release management commands.
package cmdrelease

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os/exec"
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
	allPins := make([]string, 0, len(persistentPins)+len(opts.Pins))
	allPins = append(allPins, persistentPins...)
	allPins = append(allPins, opts.Pins...)

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
	} else if !confirmBatchDeletion(len(toDelete), lo.String(), hi.String()) {
		return apperror.NewSimple("release remove-between: operation canceled (to bypass confirmation prompt, pass --confirm or -y)", "E9001")
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
	for _, p := range strings.Split(mergedPin, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			pins = append(pins, trimmed)
		}
	}

	isConfirmed := false
	confirmSource := ""
	switch {
	case *confirm:
		isConfirmed = true
		confirmSource = "--confirm"
	case *yes:
		isConfirmed = true
		confirmSource = "-y"
	case *yesLong:
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
