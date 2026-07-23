package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runTempRelease handles the "temp-release" command and routes subcommands.
func runTempRelease(args []string) {
	checkHelp("temp-release", args)

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrTRUsage)
		os.Exit(1)
	}

	sub := args[0]
	dispatchTempRelease(sub, args)
}

// dispatchTempRelease routes to list, remove, or create.
func dispatchTempRelease(sub string, args []string) {
	if sub == constants.SubCmdTRList {
		runTempReleaseList(args[1:])

		return
	}
	if sub == constants.SubCmdTRRemove {
		runTempReleaseRemove(args[1:])

		return
	}

	// Default: treat as create (first arg is count).
	runTempReleaseCreate(args)
}

// parseTempReleaseCreateFlags parses flags for the create subcommand.
func parseTempReleaseCreateFlags(args []string) (count int, pattern string, start int, dryRun, verbose bool) {
	fs := flag.NewFlagSet(constants.CmdTempRelease, flag.ExitOnError)
	startFlag := fs.Int("start", 0, constants.FlagDescTRStart)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescTRDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescTRVerbose)

	// Register -s as shorthand for --start.
	fs.IntVar(startFlag, "s", 0, constants.FlagDescTRStart)

	fs.Parse(args)

	if fs.NArg() < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrTRUsage)
		os.Exit(1)
	}

	count = parseCount(fs.Arg(0))
	pattern = fs.Arg(1)

	return count, pattern, *startFlag, *dryRunFlag, *verboseFlag
}

// parseCount converts the count argument to an integer with validation.
func parseCount(s string) int {
	var n int

	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 || n > constants.TempReleaseMaxCount {
		fmt.Fprintf(os.Stderr, constants.ErrTRInvalidCount+"\n", constants.TempReleaseMaxCount)
		os.Exit(1)
	}

	return n
}
