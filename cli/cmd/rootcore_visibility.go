package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvisibility"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func coreVisibilityActionEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdMakePublic}, func() error { return cmdvisibility.RunMakePublic(argsTail()) }},
		{[]string{constants.CmdMakePrivate}, func() error { return cmdvisibility.RunMakePrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublic, constants.CmdMAPUB}, func() error { return cmdvisibility.RunMakeAllPublic(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivate, constants.CmdMAPRI}, func() error { return cmdvisibility.RunMakeAllPrivate(argsTail()) }},
		{[]string{constants.CmdMakeAllPublicExceptLatest, constants.CmdMAPUBXL}, func() error { return cmdvisibility.RunMakeAllPublicExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeAllPrivateExceptLatest, constants.CmdMAPRIXL}, func() error { return cmdvisibility.RunMakeAllPrivateExceptLatest(argsTail()) }},
		{[]string{constants.CmdMakeLastPublic, constants.CmdMLPUB}, func() error { return cmdvisibility.RunMakeLastPublic(argsTail()) }},
		{[]string{constants.CmdMakeLastPrivate, constants.CmdMLPRI}, func() error { return cmdvisibility.RunMakeLastPrivate(argsTail()) }},
	}
}

func coreVisibilityHistoryEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdVisibilityUndo, constants.CmdVisibilityUndoAlias}, func() error { return cmdvisibility.RunVisibilityUndo(argsTail()) }},
		{[]string{constants.CmdVisibilityRedo, constants.CmdVisibilityRedoAlias}, func() error { return cmdvisibility.RunVisibilityRedo(argsTail()) }},
		{[]string{constants.CmdVisibilityHistory, constants.CmdVisibilityHistoryAlias}, func() error { return cmdvisibility.RunVisibilityHistory(argsTail()) }},
	}
}
