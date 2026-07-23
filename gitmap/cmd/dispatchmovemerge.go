package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/movemerge"
)

// dispatchMoveMerge routes mv / move / merge-both / merge-left /
// merge-right (and their short aliases mb / ml / mr).
//
// Spec: spec/01-app/97-move-and-merge.md
func dispatchMoveMerge(command string) bool {
	if command == constants.CmdMv || command == constants.CmdMove {
		runMove(os.Args[2:])

		return true
	}
	if spec, ok := mergeSpecFor(command); ok {
		runMerge(spec, os.Args[2:])

		return true
	}

	return false
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
