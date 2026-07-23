package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/config"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/tui"
)

// parseInteractiveFlags parses flags for the interactive command.
func parseInteractiveFlags(args []string) int {
	fs := flag.NewFlagSet(constants.CmdInteractive, flag.ExitOnError)
	refreshFlag := fs.Int(constants.FlagRefresh, 0, constants.FlagDescRefresh)
	fs.Parse(args)

	return *refreshFlag
}

// runInteractive launches the full-screen TUI.
func runInteractive() {
	checkHelp("interactive", os.Args[2:])

	refresh := parseInteractiveFlags(os.Args[2:])

	cfg, cfgErr := config.LoadFromFile(constants.DefaultConfigPath)
	if cfgErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not load config: %v\n", cfgErr)
	}

	if refresh > 0 {
		cfg.DashboardRefresh = refresh
	}

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrTUIDBOpen, err)
		os.Exit(1)
	}
	defer db.Close()

	if err := tui.Run(db, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
