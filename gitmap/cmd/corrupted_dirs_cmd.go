package cmd

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func parseCleanCorruptedFlags(args []string) (CleanOptions, bool) {
	fs := flag.NewFlagSet("clean-corrupted", flag.ContinueOnError)
	var dryRun, force, jsonOut bool
	fs.BoolVar(&dryRun, "dry-run", false, "Preview changes without deleting")
	fs.BoolVar(&force, "force", false, "Bypass confirmation prompt")
	fs.BoolVar(&jsonOut, "json", false, "Output results in JSON format")
	_ = fs.Parse(args)

	return CleanOptions{IsDryRun: dryRun, IsForce: force}, jsonOut
}

func printCleanResultHuman(res CleanResult, isDryRun bool) {
	if len(res.DetectedDirs) == 0 {
		fmt.Printf("  %s✓%s No corrupted installation directories found.\n", constants.ColorGreen, constants.ColorReset)

		return
	}

	prefix := "Removed"
	if isDryRun {
		prefix = "[dry-run] Would remove"
	}

	for _, dir := range res.RemovedDirs {
		fmt.Printf("  %s✓%s %s: %s\n", constants.ColorGreen, constants.ColorReset, prefix, dir)
	}

	for _, rec := range res.RecoveredFiles {
		fmt.Printf("  %s✓%s Recovered asset: %s\n", constants.ColorGreen, constants.ColorReset, rec)
	}

	if res.EscapedFrom != "" {
		fmt.Printf("  %s!%s Escaped working directory to: %s\n", constants.ColorYellow, constants.ColorReset, res.EscapedTo)
	}
}

func runCleanCorrupted(args []string) error {
	opts, isJson := parseCleanCorruptedFlags(args)
	res, err := CleanCorruptedDirs(opts)
	if err != nil {
		return apperror.Wrap(err, "failed to clean corrupted directories", nil)
	}

	if isJson {
		data, _ := json.MarshalIndent(res, "", "  ")
		fmt.Println(string(data))

		return nil
	}

	printCleanResultHuman(res, opts.IsDryRun)

	return nil
}
