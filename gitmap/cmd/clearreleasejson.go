package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/release"
)

// parseClearReleaseJSONFlags parses flags for the clear-release-json command.
func parseClearReleaseJSONFlags(args []string) (string, bool) {
	fs := flag.NewFlagSet("clear-release-json", flag.ExitOnError)

	var dryRun bool

	fs.BoolVar(&dryRun, "dry-run", false, "Preview which file would be removed without deleting it")
	_ = fs.Parse(args)

	var version string
	if fs.NArg() > 0 {
		version = fs.Arg(0)
	}

	return version, dryRun
}

// runClearReleaseJSON handles the "clear-release-json" subcommand.
func runClearReleaseJSON(args []string) {
	checkHelp("clear-release-json", args)

	version, dryRun := parseClearReleaseJSONFlags(args)

	if version == "" {
		fmt.Fprintln(os.Stderr, constants.ErrClearReleaseUsage)
		os.Exit(1)
	}

	v, err := release.Parse(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReleaseInvalidVersion, version)
		os.Exit(1)
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.ErrClearReleaseNotFound, v.String())
		os.Exit(1)
	}

	if dryRun {
		fmt.Printf(constants.MsgClearReleaseDryRun, path)
		return
	}

	err = os.Remove(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrClearReleaseFailed, path, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgClearReleaseDone, v.String())
}
