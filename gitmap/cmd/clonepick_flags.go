package cmd

// clonepick_flags.go: flag binding for `gitmap clone-pick`. Split
// out of clonepick.go so the dispatcher file stays under the
// 200-line cap (core rule: <200 lines/file). Validation that needs
// cross-flag knowledge still lives in clonepick.ParseArgs -- this
// file is pure flag plumbing.

import (
	"flag"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/clonepick"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// clonePickParsed bundles every output of parseClonePickFlags so a
// new audit/debug toggle can be added without churning the call
// site signature each time. Fields are exported because the struct
// itself stays unexported (cmd-package-internal).
type clonePickParsed struct {
	RawURL                          string
	RawPaths                        string
	Flags                           clonepick.Flags
	Output                          string
	VerifyCmdFaithful               bool
	VerifyCmdFaithfulExitOnMismatch bool
	PrintCloneArgv                  bool
	// NoVSCodeSync suppresses the post-clone update of the
	// alefragnani.project-manager projects.json file. Mirrors
	// `gitmap scan --no-vscode-sync`. Default false. See
	// spec/01-vscode-project-manager-sync/02-clone-sync.md.
	NoVSCodeSync bool
}

// parseClonePickFlags binds every clone-pick flag and extracts the
// two positional args. Validation that needs cross-flag knowledge
// happens in clonepick.ParseArgs so this stays focused on flag
// binding.
func parseClonePickFlags(args []string) clonePickParsed {
	defaults := clonepick.DefaultFlags()
	flags := defaults
	fs := flag.NewFlagSet("clone-pick", flag.ExitOnError)
	bindClonePickCoreFlags(fs, &flags, defaults)
	output, audit, noVSCodeSync := bindClonePickAuxFlags(fs)

	reordered := reorderFlagsBeforeArgs(args)
	fs.Parse(reordered)

	rawURL, rawPaths := requireClonePickPositional(fs, flags)

	return clonePickParsed{
		RawURL:                          rawURL,
		RawPaths:                        rawPaths,
		Flags:                           flags,
		Output:                          *output,
		VerifyCmdFaithful:               *audit.verify,
		VerifyCmdFaithfulExitOnMismatch: *audit.verifyExit,
		PrintCloneArgv:                  *audit.printArgv,
		NoVSCodeSync:                    *noVSCodeSync,
	}
}

// bindClonePickCoreFlags binds every user-facing clone-pick flag
// (everything documented in helptext/clone-pick.md). Audit/debug
// toggles live in bindClonePickAuxFlags so this stays under the
// function-length cap.
func bindClonePickCoreFlags(fs *flag.FlagSet, flags *clonepick.Flags, defaults clonepick.Flags) {
	fs.BoolVar(&flags.Ask, constants.FlagClonePickAsk, defaults.Ask,
		constants.FlagDescClonePickAsk)
	fs.StringVar(&flags.Name, constants.FlagClonePickName, defaults.Name,
		constants.FlagDescClonePickName)
	fs.StringVar(&flags.Mode, constants.FlagClonePickMode, defaults.Mode,
		constants.FlagDescClonePickMode)
	fs.StringVar(&flags.Branch, constants.FlagClonePickBranch, defaults.Branch,
		constants.FlagDescClonePickBranch)
	fs.IntVar(&flags.Depth, constants.FlagClonePickDepth, defaults.Depth,
		constants.FlagDescClonePickDepth)
	fs.BoolVar(&flags.Cone, constants.FlagClonePickCone, defaults.Cone,
		constants.FlagDescClonePickCone)
	fs.StringVar(&flags.Dest, constants.FlagClonePickDest, defaults.Dest,
		constants.FlagDescClonePickDest)
	fs.BoolVar(&flags.KeepGit, constants.FlagClonePickKeepGit, defaults.KeepGit,
		constants.FlagDescClonePickKeepGit)
	fs.BoolVar(&flags.DryRun, constants.FlagClonePickDryRun, defaults.DryRun,
		constants.FlagDescClonePickDryRun)
	fs.BoolVar(&flags.Quiet, constants.FlagClonePickQuiet, defaults.Quiet,
		constants.FlagDescClonePickQuiet)
	fs.BoolVar(&flags.Force, constants.FlagClonePickForce, defaults.Force,
		constants.FlagDescClonePickForce)
	fs.StringVar(&flags.Replay, constants.FlagClonePickReplay, defaults.Replay,
		constants.FlagDescClonePickReplay)
}

// clonePickAuditFlags bundles the three verify/debug toggles so
// bindClonePickAuxFlags can return them without exceeding the
// "max 3 returns is fine, 4 is noisy" rule of thumb.
type clonePickAuditFlags struct {
	verify     *bool
	verifyExit *bool
	printArgv  *bool
}

// bindClonePickAuxFlags binds output + audit toggles that aren't
// part of the user-facing clone-pick surface but are needed by
// CI / regression harnesses.
func bindClonePickAuxFlags(fs *flag.FlagSet) (*string, clonePickAuditFlags, *bool) {
	output := fs.String(constants.FlagCloneTermOutput, "",
		constants.FlagDescCloneTermOutput)
	audit := clonePickAuditFlags{
		verify: fs.Bool(constants.FlagCloneVerifyCmdFaithful, false,
			constants.FlagDescCloneVerifyCmdFaithful),
		verifyExit: fs.Bool(constants.FlagCloneVerifyCmdFaithfulExitOnMismatch,
			false, constants.FlagDescCloneVerifyCmdFaithfulExitOnMismatch),
		printArgv: fs.Bool(constants.FlagClonePrintArgv, false,
			constants.FlagDescClonePrintArgv),
	}
	noVSCodeSync := fs.Bool(constants.FlagNoVSCodeSync, false,
		constants.FlagDescNoVSCodeSync)

	return output, audit, noVSCodeSync
}

// requireClonePickPositional enforces the "URL required unless
// --replay" contract and pulls the positional URL + paths out of fs.
func requireClonePickPositional(fs *flag.FlagSet, flags clonepick.Flags) (string, string) {
	if fs.NArg() < 1 && len(flags.Replay) == 0 {
		fmt.Fprintln(os.Stderr, constants.MsgClonePickMissingURL)
		os.Exit(2)
	}
	rawURL := ""
	if fs.NArg() >= 1 {
		rawURL = fs.Arg(0)
	}
	rawPaths := ""
	if fs.NArg() >= 2 {
		rawPaths = fs.Arg(1)
	}

	return rawURL, rawPaths
}
