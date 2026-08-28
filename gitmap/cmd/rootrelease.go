package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// dispatchRelease routes release-related commands.
func dispatchRelease(command string) (bool, error) {
	return runDispatchTable(command, releaseDispatchEntries())
}

// releaseDispatchEntries returns the routing table for release commands.
func releaseDispatchEntries() []dispatchEntry {
	return []dispatchEntry{
		{[]string{constants.CmdRelease, constants.CmdReleaseShort}, func() error { return runRelease(argsTail()) }},
		{[]string{constants.CmdReleasePull, constants.CmdReleasePullAlias, constants.CmdReleasePullAlias2, constants.CmdReleasePullAlias3, constants.CmdReleasePullAlias4}, func() error { return runReleasePull(argsTail()) }},
		{[]string{constants.CmdPullReleaseCD, constants.CmdPullReleaseCDAlias}, func() error { return runPullReleaseCD(argsTail()) }},
		{
			[]string{constants.CmdReleaseSelf, constants.CmdReleaseSelfAlias, constants.CmdReleaseSelfAlias2},
			func() error { return runReleaseSelf(argsTail()) },
		},
		{[]string{constants.CmdReleaseBranch, constants.CmdReleaseBranchAlias}, func() error { return runReleaseBranch(argsTail()) }},
		{[]string{constants.CmdReleasePending, constants.CmdReleasePendingAlias}, func() error { return runReleasePending(argsTail()) }},
		{[]string{"release-scan-commits", "rsc-commits"}, func() error { return runReleaseScanCommits(argsTail()) }},
		{[]string{constants.CmdReleaseUndo, constants.CmdReleaseUndoAlias}, func() error { return runReleaseUndo(argsTail()) }},
		{[]string{constants.CmdChangelog, constants.CmdChangelogAlias}, func() error { return runChangelog(argsTail()) }},
		{[]string{constants.CmdChangelogMD}, func() error { return runChangelog([]string{constants.FlagOpenValue}) }},
		{[]string{constants.CmdClearReleaseJSON, constants.CmdClearReleaseJSONAlias}, func() error { return runClearReleaseJSON(argsTail()) }},
		{[]string{constants.CmdChangelogGen, constants.CmdChangelogGenAlias}, func() error { return runChangelogGen(argsTail()) }},
		{[]string{constants.CmdReleaseAlias, constants.CmdReleaseAliasShort}, func() error { return runReleaseAlias(argsTail(), false) }},
		{[]string{constants.CmdReleaseAliasPull, constants.CmdReleaseAliasPullShort}, func() error { return runReleaseAlias(argsTail(), true) }},
	}
}
