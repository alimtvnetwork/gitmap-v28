package cmd

// Flag parser for `gitmap clone-from`. Split from clonefrom.go so
// that file stays under the 200-line per-file cap
// (mem://style/code-constraints, item 3).
//
// Behavioral notes worth preserving when editing this file:
//
//   - --emit-schema short-circuits BEFORE the <file> requirement so
//     `gitmap clone-from --emit-schema=report` works without a
//     dummy input file. The dispatcher then routes to
//     runCloneFromEmitSchema and exits without touching the rest of
//     the clone surface.
//   - Missing <file> => exit 2 (CLI-usage class), invalid checkout
//     mode => exit 2 (same bucket), invalid max-concurrency => exit 2.
//     Operation-time failures stay on exit 1 (see clonefrom.go).

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// parseCloneFromFlags wires flags + extracts the positional file
// argument. Exits 2 on caller errors (missing arg, invalid flag
// values) so shell scripts can distinguish them from operation
// failures (exit 1).
func parseCloneFromFlags(args []string) cloneFromFlags {
	var cfg cloneFromFlags
	fs := flag.NewFlagSet("clone-from", flag.ExitOnError)
	maxConcFlag := bindCloneFromFlags(fs, &cfg)
	reordered := reorderFlagsBeforeArgs(args)
	fs.Parse(reordered)
	if cfg.emitSchema != "" {
		return cfg
	}
	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, constants.MsgCloneFromMissingArg)
		os.Exit(2)
	}
	validateCheckoutFlag(cfg.checkout)
	resolvedConc, ok := cloneconcurrency.Resolve(*maxConcFlag)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrCloneMaxConcurrencyInvalid, *maxConcFlag)
		os.Exit(2)
	}
	cfg.maxConcurrency = resolvedConc
	cfg.file = fs.Arg(0)

	return cfg
}

// bindCloneFromFlags registers every flag on fs against cfg and
// returns the max-concurrency pointer (the one flag that can't be
// inlined into the struct because cloneconcurrency.Resolve must
// post-process it). Split out so parseCloneFromFlags stays under
// the 15-line function budget.
func bindCloneFromFlags(fs *flag.FlagSet, cfg *cloneFromFlags) *int {
	fs.BoolVar(&cfg.execute, constants.FlagCloneFromExecute, false,
		constants.FlagDescCloneFromExecute)
	fs.BoolVar(&cfg.quiet, constants.FlagCloneFromQuiet, false,
		constants.FlagDescCloneFromQuiet)
	fs.BoolVar(&cfg.noReport, constants.FlagCloneFromNoReport, false,
		constants.FlagDescCloneFromNoReport)
	fs.StringVar(&cfg.output, constants.FlagCloneFromOutput, "",
		constants.FlagDescCloneFromOutput)
	fs.BoolVar(&cfg.verifyCmdFaithful, constants.FlagCloneVerifyCmdFaithful,
		false, constants.FlagDescCloneVerifyCmdFaithful)
	fs.BoolVar(&cfg.verifyCmdFaithfulExitOnMismatch,
		constants.FlagCloneVerifyCmdFaithfulExitOnMismatch, false,
		constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch)
	fs.BoolVar(&cfg.printCloneArgv, constants.FlagClonePrintArgv,
		false, constants.FlagDescClonePrintArgv)
	fs.StringVar(&cfg.checkout, constants.FlagCloneFromCheckout, "",
		constants.FlagDescCloneFromCheckout)
	fs.StringVar(&cfg.emitSchema, constants.FlagCloneFromEmitSchema, "",
		constants.FlagDescCloneFromEmitSchema)
	fs.BoolVar(&cfg.noVSCodeSync, constants.FlagNoVSCodeSync, false,
		constants.FlagDescNoVSCodeSync)

	return fs.Int(constants.CloneFlagMaxConcurrency,
		constants.CloneDefaultMaxConcurrency, constants.FlagDescCloneMaxConcurrency)
}
