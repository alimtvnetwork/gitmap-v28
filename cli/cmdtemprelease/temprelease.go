package cmdtemprelease

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// RunTempRelease handles the "temp-release" command and routes subcommands.
func RunTempRelease(args []string) error {
	checkHelp("temp-release", args)

	if len(args) == 0 {
		return apperror.NewWithDetails(
			"cmd.temprelease.run",
			"E1116",
			"missing required arguments for temp-release",
			"cmd.temprelease",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	return dispatchTempRelease(args[0], args)
}

// dispatchTempRelease routes to list, remove, or create.
func dispatchTempRelease(sub string, args []string) error {
	if sub == constants.SubCmdTRList {
		return runTempReleaseList(args[1:])
	}

	if sub == constants.SubCmdTRRemove {
		return runTempReleaseRemove(args[1:])
	}

	return runTempReleaseCreate(args)
}

// parseTempReleaseCreateFlags parses flags for the create subcommand.
func parseTempReleaseCreateFlags(args []string) (int, string, int, bool, bool) {
	fs, startFlag, dryRunFlag, verboseFlag := initTempReleaseCreateFlags()
	fs.Parse(args)
	validateCreateArgs(fs.NArg())

	count := parseCount(fs.Arg(0))
	pattern := fs.Arg(1)

	return count, pattern, *startFlag, *dryRunFlag, *verboseFlag
}

func initTempReleaseCreateFlags() (*flag.FlagSet, *int, *bool, *bool) {
	fs := flag.NewFlagSet(constants.CmdTempRelease, flag.ExitOnError)
	startFlag := fs.Int("start", 0, constants.FlagDescTRStart)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescTRDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescTRVerbose)
	fs.IntVar(startFlag, "s", 0, constants.FlagDescTRStart)

	return fs, startFlag, dryRunFlag, verboseFlag
}

func validateCreateArgs(narg int) {
	if narg < 2 {
		err := apperror.NewWithDetails(
			"cmd.temprelease.parseCreateFlags",
			"E1117",
			"insufficient arguments for temp-release create; usage: gitmap temp-release <count> <pattern>",
			"cmd.temprelease",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleValidationError(err)
	}
}

// parseCount converts the count argument to an integer with validation.
func parseCount(s string) int {
	var n int

	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 || n > constants.TempReleaseMaxCount {
		appErr := apperror.NewWithDetails(
			"cmd.temprelease.parseCount",
			"E1118",
			fmt.Sprintf("invalid count '%s', must be between 1 and %d", s, constants.TempReleaseMaxCount),
			"cmd.temprelease",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			map[string]any{"raw": s},
		)
		cliexit.HandleValidationError(appErr)
	}

	return n
}
