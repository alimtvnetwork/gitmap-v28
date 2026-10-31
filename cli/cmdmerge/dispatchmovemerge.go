package cmdmerge

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/movemerge"
)

// dispatchMoveMerge routes mv / move / merge-both / merge-left /
// merge-right (and their short aliases mb / ml / mrg), plus
// `merge write` / `mr` (open a GitHub pull request).
//
// NOTE (program 273, spec 05): `mr` is the merge-write alias. It used to be
// merge-right's short alias; merge-right keeps `mrg` so no alias is lost.
//
// Spec: 02-spec/01-app/97-move-and-merge.md
func DispatchMoveMerge(command string) (bool, error) {
	if command == constants.CmdMv || command == constants.CmdMove {
		RunMove(os.Args[2:])

		return true, nil
	}

	// Program 275: `mr list|view|merge|checkout` route to the mr family
	// before the merge-write branch below (bare `mr` still opens a PR).
	if handled, err := dispatchMrFamily(command, os.Args[2:]); handled {
		return true, err
	}

	if isMergeWriteInvocation(command, os.Args) {
		if appErr := runMergeWrite(stripMergeWriteSubcommand(os.Args[2:])); appErr != nil {
			return true, appErr
		}

		return true, nil
	}

	if spec, ok := mergeSpecFor(command); ok {
		runMerge(spec, os.Args[2:])

		return true, nil
	}

	if isMergeCommand(command) {
		return dispatchMergeCommand(command, os.Args[2:])
	}

	return false, nil
}

func isMergeCommand(command string) bool {
	return command == constants.CmdMerge || command == "merge-in" || command == "mergein"
}

func dispatchMergeCommand(command string, args []string) (bool, error) {
	if command == "merge-in" || command == "mergein" {
		if isHelpRequested(args) {
			RenderMergeHelp()

			return true, nil
		}

		if CommitInRunner != nil {
			return true, CommitInRunner(args)
		}

		return true, nil
	}

	if len(args) == 0 || isHelpRequested(args) {
		RenderMergeHelp()

		return true, nil
	}

	return dispatchMergeSubcommand(args)
}

func dispatchMergeSubcommand(args []string) (bool, error) {
	sub := strings.ToLower(args[0])
	switch sub {
	case "right":
		spec, _ := mergeSpecFor(constants.CmdMergeRight)
		return true, runMerge(spec, args[1:])
	case "left":
		spec, _ := mergeSpecFor(constants.CmdMergeLeft)
		return true, runMerge(spec, args[1:])
	case "both":
		spec, _ := mergeSpecFor(constants.CmdMergeBoth)
		return true, runMerge(spec, args[1:])
	case "in":
		if CommitInRunner != nil {
			return true, CommitInRunner(args[1:])
		}
		return true, nil
	default:
		RenderMergeHelp()

		return true, nil
	}
}

// mergeSpecFor maps a command/alias to its mergeSpec; ok=false otherwise.
func mergeSpecFor(command string) (mergeSpec, bool) {
	switch command {
	case constants.CmdMergeBoth, constants.CmdMergeBothA:
		return mergeSpec{constants.CmdMergeBoth, constants.LogPrefixMergeBoth,
			constants.CommitMsgMergeBoth, movemerge.DirBoth}, true
	case constants.CmdMergeLeft, constants.CmdMergeLeftA:
		return mergeSpec{constants.CmdMergeLeft, constants.LogPrefixMergeLeft,
			constants.CommitMsgMergeLeft, movemerge.DirLeftOnly}, true
	case constants.CmdMergeRight, constants.CmdMergeRgtA:
		return mergeSpec{constants.CmdMergeRight, constants.LogPrefixMergeRight,
			constants.CommitMsgMergeRight, movemerge.DirRightOnly}, true
	}

	return mergeSpec{}, false
}
