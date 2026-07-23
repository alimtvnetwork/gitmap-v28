package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// historyOpts is the parsed flag bundle for both history-* commands.
type historyOpts struct {
	yes         bool
	noPush      bool
	dryRun      bool
	keepSandbox bool
	quiet       bool
	message     string
	modeLabel   string
	pathCount   int
}

// parseHistoryArgs splits args into flags and positional paths,
// validates, and exits 4 on any invalid combination.
func parseHistoryArgs(args []string) (historyOpts, []string) {
	fs := flag.NewFlagSet(constants.CmdHistoryPurge, flag.ContinueOnError)
	raw := registerHistoryFlags(fs)
	flagsOnly, positional := splitHistoryFlagsAndArgs(args)
	if err := fs.Parse(flagsOnly); err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrBadArgs, err.Error())
		os.Exit(constants.HistoryExitBadArgs)
	}
	opts := assembleHistoryOpts(raw)
	paths := parseHistoryPaths(positional)
	validateHistoryOpts(opts, paths)
	return opts, paths
}

// rawHistoryFlags holds the pointers from FlagSet.* registrations.
type rawHistoryFlags struct {
	yes, yesShort, noPush, dryRun, keep, quiet, quietShort *bool
	message                                                *string
}

// registerHistoryFlags wires every flag onto fs and returns pointers.
func registerHistoryFlags(fs *flag.FlagSet) rawHistoryFlags {
	return rawHistoryFlags{
		yes:        fs.Bool(constants.HistoryFlagYes, false, constants.HistoryDescYes),
		yesShort:   fs.Bool(constants.HistoryFlagYesShort, false, constants.HistoryDescYes),
		noPush:     fs.Bool(constants.HistoryFlagNoPush, false, constants.HistoryDescNoPush),
		dryRun:     fs.Bool(constants.HistoryFlagDryRun, false, constants.HistoryDescDryRun),
		keep:       fs.Bool(constants.HistoryFlagKeepSandbox, false, constants.HistoryDescKeepSandbox),
		quiet:      fs.Bool(constants.HistoryFlagQuiet, false, constants.HistoryDescQuiet),
		quietShort: fs.Bool(constants.HistoryFlagQuietShort, false, constants.HistoryDescQuiet),
		message:    fs.String(constants.HistoryFlagMessage, "", constants.HistoryDescMessage),
	}
}

// assembleHistoryOpts collapses short/long aliases into one struct.
func assembleHistoryOpts(r rawHistoryFlags) historyOpts {
	return historyOpts{
		yes:         *r.yes || *r.yesShort,
		noPush:      *r.noPush,
		dryRun:      *r.dryRun,
		keepSandbox: *r.keep,
		quiet:       *r.quiet || *r.quietShort,
		message:     *r.message,
	}
}

// validateHistoryOpts enforces spec §5 / §6: at least one path,
// --yes and --no-push mutually exclusive.
func validateHistoryOpts(opts historyOpts, paths []string) {
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, constants.HistoryErrBadArgs, constants.HistoryErrNoPaths)
		os.Exit(constants.HistoryExitBadArgs)
	}
	if opts.yes && opts.noPush {
		fmt.Fprintf(os.Stderr, constants.HistoryErrBadArgs, constants.HistoryErrConflictFlags)
		os.Exit(constants.HistoryExitBadArgs)
	}
}

// splitHistoryFlagsAndArgs separates `-x`/`--x[=v]` tokens from
// positional path tokens so flag.Parse can see only flags.
func splitHistoryFlagsAndArgs(args []string) ([]string, []string) {
	flags, positional := make([]string, 0, len(args)), make([]string, 0, len(args))
	i := 0
	for i < len(args) {
		tok := args[i]
		if len(tok) > 1 && tok[0] == '-' {
			flags = append(flags, tok)
			if isHistoryStringFlag(tok) && i+1 < len(args) {
				flags = append(flags, args[i+1])
				i += 2
				continue
			}
		} else {
			positional = append(positional, tok)
		}
		i++
	}
	return flags, positional
}

// isHistoryStringFlag returns true for flags that consume a value.
// Only --message takes a value; bool flags don't.
func isHistoryStringFlag(tok string) bool {
	return tok == "--"+constants.HistoryFlagMessage ||
		tok == "-"+constants.HistoryFlagMessage
}
