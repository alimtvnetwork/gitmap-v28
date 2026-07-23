package cmd

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/committransfer"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/movemerge"
)

// commitTransferSpec describes one of the three commit-transfer commands.
type commitTransferSpec struct {
	Name      string // e.g. constants.CmdCommitLeft
	LogPrefix string // e.g. constants.LogPrefixCommitLeft
}

// runCommitTransfer is the single entry point for commit-left,
// commit-right, and commit-both.
//
// Phase 1 (v3.76.0): commit-right.
// Phase 2 (v3.102.0): commit-left wired through committransfer.RunLeft.
// Phase 3 (v3.102.0): commit-both wired through committransfer.RunBoth.
func runCommitTransfer(spec commitTransferSpec, args []string) {
	checkHelp(spec.Name, args)
	executeCommitTransfer(spec, args)
}

// executeCommitTransfer parses flags, resolves endpoints, and dispatches
// to the directional runner that matches spec.Name. Kept separate from
// runCommitTransfer so the help-check stays a one-liner and the dispatch
// table reads top-to-bottom without nested branching.
func executeCommitTransfer(spec commitTransferSpec, args []string) {
	opts, positional := parseCommitTransferArgs(spec, args)
	if len(positional) != 2 {
		fmt.Fprintf(os.Stderr, constants.ErrCTArgCountFmt, spec.Name, len(positional))
		fmt.Fprintf(os.Stderr, constants.MsgCTUsageFmt, spec.Name, spec.Name)
		os.Exit(1)
	}
	if opts.Interleave && spec.Name != constants.CmdCommitBoth {
		fmt.Fprintf(os.Stderr,
			"%s --interleave is only valid for commit-both (got %s)\n",
			opts.LogPrefix, spec.Name)
		os.Exit(2)
	}
	left, right, resolveErr := resolveCommitEndpoints(positional[0], positional[1], opts)
	if resolveErr != nil {
		fmt.Fprintf(os.Stderr, "%s endpoint resolve failed: %v\n", opts.LogPrefix, resolveErr)
		os.Exit(1)
	}
	opts.Message.SourceDisplayName = pickSourceDisplayName(spec.Name, left, right)
	if err := dispatchDirection(spec.Name, left.WorkingDir, right.WorkingDir, opts); err != nil {
		fmt.Fprintf(os.Stderr, "%s replay failed: %v\n", opts.LogPrefix, err)
		os.Exit(1)
	}
}

// dispatchDirection routes to the right RunX function. LEFT/RIGHT
// positional ordering matches the spec — `commit-left LEFT RIGHT`
// writes commits onto LEFT (using RIGHT as source). For commit-both,
// --interleave switches to the author-date merged stream.
func dispatchDirection(name, leftDir, rightDir string, opts committransfer.Options) error {
	switch name {
	case constants.CmdCommitLeft:
		return committransfer.RunLeft(leftDir, rightDir, opts)
	case constants.CmdCommitBoth:
		if opts.Interleave {
			return committransfer.RunBothInterleaved(leftDir, rightDir, opts)
		}

		return committransfer.RunBoth(leftDir, rightDir, opts)
	default:
		// commit-right (and any future direction defaulting to L→R).
		return committransfer.RunRight(leftDir, rightDir, opts)
	}
}

// pickSourceDisplayName labels the provenance footer with the side
// that's being read from. commit-left reads from RIGHT; commit-right
// reads from LEFT; commit-both reads from both — we pick LEFT as the
// canonical label because the RunBoth implementation uses (L→R) first.
func pickSourceDisplayName(name string, left, right movemerge.Endpoint) string {
	if name == constants.CmdCommitLeft {
		return right.DisplayName
	}

	return left.DisplayName
}

// resolveCommitEndpoints reuses the merge-* endpoint resolver. LEFT is
// the source for commit-right; we mark it as the "left" side for the
// resolver's missing-folder semantics.
func resolveCommitEndpoints(leftRaw, rightRaw string, _ committransfer.Options,
) (movemerge.Endpoint, movemerge.Endpoint, error) {
	mmOpts := movemerge.Options{}
	left, err := movemerge.ResolveEndpoint(leftRaw, true, mmOpts)
	if err != nil {
		return left, movemerge.Endpoint{}, err
	}
	right, err := movemerge.ResolveEndpoint(rightRaw, false, mmOpts)

	return left, right, err
}

// parseCommitTransferArgs builds the Options struct + positional args.
// One function per concern would be cleaner, but the flag.FlagSet API
// keeps us under the per-function line cap as long as helpers extract
// the message-policy block.
func parseCommitTransferArgs(spec commitTransferSpec, args []string,
) (committransfer.Options, []string) {
	fs := flag.NewFlagSet(spec.Name, flag.ExitOnError)
	opts := committransfer.Options{
		CommandName: spec.Name, LogPrefix: spec.LogPrefix,
		IncludeMerges: true, // v6.0.0 default — merge commits preserved
		Message: committransfer.MessagePolicy{
			DropPatterns: committransfer.DefaultDropPatterns,
			Conventional: true, Provenance: true,
			CommandName: spec.Name,
		},
	}
	registerCommitTransferBools(fs, &opts)
	registerCommitTransferStrings(fs, &opts)
	fs.Parse(reorderFlagsBeforeArgs(args))

	return opts, fs.Args()
}

// registerCommitTransferBools wires every boolean flag from spec §8.
func registerCommitTransferBools(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.BoolVar(&opts.Yes, constants.FlagCTYes, false, constants.FlagDescCTYes)
	fs.BoolVar(&opts.Yes, "y", false, constants.FlagDescCTYes)
	fs.BoolVar(&opts.DryRun, constants.FlagCTDryRun, false, constants.FlagDescCTDryRun)
	fs.BoolVar(&opts.NoPush, constants.FlagCTNoPush, false, constants.FlagDescCTNoPush)
	fs.BoolVar(&opts.NoCommit, constants.FlagCTNoCommit, false, constants.FlagDescCTNoCommit)
	fs.BoolFunc(constants.FlagCTIncludeMerges, constants.FlagDescCTIncludeMerges,
		func(string) error { opts.IncludeMerges = true; return nil })
	fs.BoolFunc(constants.FlagCTNoIncludeMerges, constants.FlagDescCTNoIncludeMerges,
		func(string) error { opts.IncludeMerges = false; return nil })
	fs.BoolVar(&opts.Mirror, constants.FlagCTMirror, false, constants.FlagDescCTMirror)
	fs.BoolVar(&opts.ForceReplay, constants.FlagCTForceReplay, false, constants.FlagDescCTForceReplay)
	fs.BoolVar(&opts.Interleave, constants.FlagCTInterleave, false, constants.FlagDescCTInterleave)
	registerMessagePolicyToggles(fs, opts)
}

// registerMessagePolicyToggles wires the on/off pairs for §6 stages.
// Uses fs.BoolFunc (Go 1.21+) so negations don't consume a value.
// Order on the command line is the order of effect (last wins).
func registerMessagePolicyToggles(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.BoolFunc(constants.FlagCTConventional, constants.FlagDescCTConventional,
		func(string) error { opts.Message.Conventional = true; return nil })
	fs.BoolFunc(constants.FlagCTNoConventional, constants.FlagDescCTNoConventional,
		func(string) error { opts.Message.Conventional = false; return nil })
	fs.BoolFunc(constants.FlagCTProvenance, constants.FlagDescCTProvenance,
		func(string) error { opts.Message.Provenance = true; return nil })
	fs.BoolFunc(constants.FlagCTNoProvenance, constants.FlagDescCTNoProvenance,
		func(string) error { opts.Message.Provenance = false; return nil })
}

// registerCommitTransferStrings wires value-taking flags + repeatable
// regex patterns. --no-strip and --no-drop are BoolFunc (no value).
func registerCommitTransferStrings(fs *flag.FlagSet, opts *committransfer.Options) {
	fs.IntVar(&opts.Limit, constants.FlagCTLimit, 0, constants.FlagDescCTLimit)
	fs.StringVar(&opts.Since, constants.FlagCTSince, "", constants.FlagDescCTSince)
	fs.IntVar(&opts.MaxHistoryScan, constants.FlagCTMaxHistoryScan, 0, constants.FlagDescCTMaxHistoryScan)
	fs.Func(constants.FlagCTStrip, constants.FlagDescCTStrip, func(v string) error {
		opts.Message.StripPatterns = append(opts.Message.StripPatterns, v)

		return nil
	})
	fs.Func(constants.FlagCTDrop, constants.FlagDescCTDrop, func(v string) error {
		opts.Message.DropPatterns = append(opts.Message.DropPatterns, v)

		return nil
	})
	fs.BoolFunc(constants.FlagCTNoStrip, constants.FlagDescCTNoStrip, func(string) error {
		opts.Message.StripPatterns = nil

		return nil
	})
	fs.BoolFunc(constants.FlagCTNoDrop, constants.FlagDescCTNoDrop, func(string) error {
		opts.Message.DropPatterns = nil

		return nil
	})
}

// commitTransferSpecFor maps a command name or alias to its spec.
func commitTransferSpecFor(command string) (commitTransferSpec, bool) {
	switch command {
	case constants.CmdCommitLeft, constants.CmdCommitLeftA:
		return commitTransferSpec{
			Name: constants.CmdCommitLeft, LogPrefix: constants.LogPrefixCommitLeft,
		}, true
	case constants.CmdCommitRight, constants.CmdCommitRightA:
		return commitTransferSpec{
			Name: constants.CmdCommitRight, LogPrefix: constants.LogPrefixCommitRight,
		}, true
	case constants.CmdCommitBoth, constants.CmdCommitBothA:
		return commitTransferSpec{
			Name: constants.CmdCommitBoth, LogPrefix: constants.LogPrefixCommitBoth,
		}, true
	}

	return commitTransferSpec{}, false
}
