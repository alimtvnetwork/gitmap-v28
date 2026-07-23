package cmd

import (
	"flag"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/movemerge"
)

// movemergeFlagSet attaches every shared flag for mv/merge-* to fs.
type movemergeFlagSet struct {
	yes, accept                                   bool
	prefL, prefR, prefNewer, prefSkip             bool
	noPush, noCommit, forceFold, pullFold, dryRun bool
	initNew, includeVCS, includeNM                bool
}

// bindFlags wires every flag onto the provided FlagSet.
func (m *movemergeFlagSet) bindFlags(fs *flag.FlagSet) {
	fs.BoolVar(&m.yes, constants.FlagMMYes, false, "bypass conflict prompt")
	fs.BoolVar(&m.yes, constants.FlagMMYesShort, false, "bypass conflict prompt (short)")
	fs.BoolVar(&m.accept, constants.FlagMMAccept, false, "bypass conflict prompt (alias)")
	fs.BoolVar(&m.accept, constants.FlagMMAcceptShrt, false, "bypass conflict prompt (short alias)")
	fs.BoolVar(&m.prefL, constants.FlagMMPreferL, false, "LEFT always wins on conflict")
	fs.BoolVar(&m.prefR, constants.FlagMMPreferR, false, "RIGHT always wins on conflict")
	fs.BoolVar(&m.prefNewer, constants.FlagMMPreferNew, false, "newer mtime wins on conflict")
	fs.BoolVar(&m.prefSkip, constants.FlagMMPreferSkip, false, "skip every conflict")
	fs.BoolVar(&m.noPush, constants.FlagMMNoPush, false, "skip git push on URL endpoints")
	fs.BoolVar(&m.noCommit, constants.FlagMMNoCommit, false, "skip commit + push on URL endpoints")
	fs.BoolVar(&m.forceFold, constants.FlagMMForceFold, false, "replace folder whose origin doesn't match URL")
	fs.BoolVar(&m.pullFold, constants.FlagMMPullFold, false, "force git pull --ff-only on a folder endpoint")
	fs.BoolVar(&m.dryRun, constants.FlagMMDryRun, false, "print every action; perform none")
	fs.BoolVar(&m.initNew, constants.FlagMMInit, false, "git init RIGHT when freshly created (mv only)")
	fs.BoolVar(&m.includeVCS, constants.FlagMMIncludeVCS, false, "include .git/ in copy/diff")
	fs.BoolVar(&m.includeNM, constants.FlagMMIncludeNM, false, "include node_modules/ in copy/diff")
}

// toOptions converts parsed flags into the movemerge.Options struct.
func (m *movemergeFlagSet) toOptions(cmd, prefix, msgFmt string) movemerge.Options {
	return movemerge.Options{
		Yes:             m.yes || m.accept,
		Prefer:          pickPolicy(m),
		NoPush:          m.noPush,
		NoCommit:        m.noCommit,
		ForceFolder:     m.forceFold,
		PullFolder:      m.pullFold,
		InitNewRight:    m.initNew,
		DryRun:          m.dryRun,
		IncludeVCS:      m.includeVCS,
		IncludeNodeMods: m.includeNM,
		CommandName:     cmd,
		LogPrefix:       prefix,
		CommitMsgFmt:    msgFmt,
	}
}

// pickPolicy returns the first non-default --prefer-* flag.
func pickPolicy(m *movemergeFlagSet) movemerge.PreferPolicy {
	switch {
	case m.prefL:
		return movemerge.PreferLeft
	case m.prefR:
		return movemerge.PreferRight
	case m.prefNewer:
		return movemerge.PreferNewer
	case m.prefSkip:
		return movemerge.PreferSkip
	}

	return movemerge.PreferNone
}
