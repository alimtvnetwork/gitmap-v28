package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
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
func runClearReleaseJSON(args []string) error {
	checkHelp("clear-release-json", args)

	version, dryRun := parseClearReleaseJSONFlags(args)

	if version == "" {
		return apperror.NewSimple(constants.ErrClearReleaseUsage, "E9000")
	}

	v, err := release.Parse(version)
	if err != nil {
		return apperror.NewSimple(constants.ErrReleaseInvalidVersion, "E9000")
	}

	filename := v.String() + constants.ExtJSON
	path := filepath.Join(constants.DefaultReleaseDir, filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, constants.ErrClearReleaseNotFound, v.String())
		return apperror.NewSimple("fatal error", "E9000")
	}

	if dryRun {
		fmt.Printf(constants.MsgClearReleaseDryRun, path)
		return nil
	}

	err = os.Remove(path)
	if err != nil {
		return apperror.NewSimple(constants.ErrClearReleaseFailed, "E9000")
	}

	fmt.Printf(constants.MsgClearReleaseDone, v.String())
	return nil
}
