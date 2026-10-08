package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runReset handles the "reset" subcommand: deletes all database files from disk,
// recreates the schemas, and reseeds clean baseline data.
func runReset(args []string) error {
	checkHelp(constants.CmdReset, args)
	opts := cmddb.ParseResetOptionsExport(args)
	needsPrompt := !opts.IsDryRun && !opts.IsConfirm
	if needsPrompt && !cmddb.ConfirmOrSkip(fmt.Sprintf("\n  %s⚠ Warning:%s This will permanently delete all database files and rebuild them from scratch. [y/N]: ", constants.ColorYellow, constants.ColorReset), args) {
		printResetNoConfirm(opts.IsRescan)

		return nil
	}
	if needsPrompt {
		opts.IsConfirm = true
	}

	if err := cmddb.PerformComprehensiveResetExport(opts); err != nil {
		return err
	}

	if opts.IsRescan && !opts.IsDryRun {
		_ = cmdscan.RunRescan()
	}

	return nil
}

func printResetNoConfirm(isRescan bool) {
	cmd := "gitmap reset --confirm"
	if isRescan {
		cmd += " --rescan"
	}

	fmt.Println()
	fmt.Printf("  %s⚠ Warning:%s This will permanently delete the database files and rebuild them from scratch.\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("  Run with %s--confirm%s (or %s-y%s) to proceed:\n\n", constants.ColorCyan, constants.ColorReset, constants.ColorCyan, constants.ColorReset)
	fmt.Printf("    %s%s%s\n\n", constants.ColorGreen, cmd, constants.ColorReset)
}
