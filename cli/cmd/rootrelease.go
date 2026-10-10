package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchangelog"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdrelease"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchRelease routes release-related commands.
func dispatchRelease(command string) (bool, error) {
	return runDispatchTable(command, releaseDispatchEntries())
}

// releaseDispatchEntries returns the routing table for release commands.
func releaseDispatchEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdRelease, constants.CmdReleaseShort}, func() error { return cmdrelease.RunRelease(argsTail()) }},
		{[]string{constants.CmdReleasePull, constants.CmdReleasePullAlias, constants.CmdReleasePullAlias2, constants.CmdReleasePullAlias3, constants.CmdReleasePullAlias4}, func() error { return cmdrelease.RunReleasePull(argsTail()) }},
		{[]string{constants.CmdPullReleaseCD, constants.CmdPullReleaseCDAlias}, func() error { return cmdpull.RunPullReleaseCD(argsTail()) }},
		{
			[]string{constants.CmdReleaseSelf, constants.CmdReleaseSelfAlias, constants.CmdReleaseSelfAlias2},
			func() error { return cmdrelease.RunReleaseSelf(argsTail()) },
		},
		{[]string{constants.CmdReleaseBranch, constants.CmdReleaseBranchAlias}, func() error { return cmdrelease.RunReleaseBranch(argsTail()) }},
		{[]string{constants.CmdReleasePending, constants.CmdReleasePendingAlias}, func() error { return cmdrelease.RunReleasePending(argsTail()) }},
		{[]string{"release-scan-commits", "rsc-commits"}, func() error { return cmdrelease.RunReleaseScanCommits(argsTail()) }},
		{[]string{constants.CmdReleaseUndo, constants.CmdReleaseUndoAlias}, func() error { return cmdrelease.RunReleaseUndo(argsTail()) }},
		{[]string{constants.CmdChangelog, constants.CmdChangelogAlias}, func() error { return cmdchangelog.RunChangelog(argsTail()) }},
		{[]string{constants.CmdChangelogMD}, func() error { return cmdchangelog.RunChangelog([]string{constants.FlagOpenValue}) }},
		{[]string{constants.CmdClearReleaseJSON, constants.CmdClearReleaseJSONAlias}, func() error { return cmdrelease.RunClearReleaseJSON(argsTail()) }},
		{[]string{constants.CmdChangelogGen, constants.CmdChangelogGenAlias}, func() error { return cmdchangelog.RunChangelogGen(argsTail()) }},
		{[]string{constants.CmdReleaseAlias, constants.CmdReleaseAliasShort}, func() error { return cmdrelease.RunReleaseAlias(argsTail(), false) }},
		{[]string{constants.CmdReleaseAliasPull, constants.CmdReleaseAliasPullShort}, func() error { return cmdrelease.RunReleaseAlias(argsTail(), true) }},
	}
}
