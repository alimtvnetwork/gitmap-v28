package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// printHints prints contextual helper hints to stderr.
func printHints(hints []hintEntry) {
	fmt.Fprint(os.Stderr, constants.MsgHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, constants.MsgHintRowFmt, h.command, h.description)
	}
}

// hintEntry holds a command example and its description.
type hintEntry struct {
	command     string
	description string
}

// projectReposHints returns hints shown after go-repos, node-repos, etc.
func projectReposHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupAdd, constants.HintGroupAddDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
		{constants.HintPullGroup, constants.HintPullGroupDesc},
	}
}

// listHints returns hints shown after gitmap list.
func listHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintLsType, constants.HintLsTypeDesc},
		{constants.HintCDRepo, constants.HintCDRepoDesc},
	}
}

// listGroupsHints returns hints shown after gitmap ls groups.
func listGroupsHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
	}
}

// activeGroupHints returns hints shown after gitmap g (active group).
func activeGroupHints() []hintEntry {
	return []hintEntry{
		{constants.HintGPull, constants.HintGPullDesc},
		{constants.HintGStatus, constants.HintGStatusDesc},
		{constants.HintGClear, constants.HintGClearDesc},
	}
}

// groupListHints returns hints shown after gitmap group list.
func groupListHints() []hintEntry {
	return []hintEntry{
		{constants.HintGroupCreate, constants.HintGroupCreateDesc},
		{constants.HintGroupShow, constants.HintGroupShowDesc},
		{constants.HintGroupDelete, constants.HintGroupDeleteDesc},
	}
}

// aliasListHints returns hints shown after gitmap a list.
func aliasListHints() []hintEntry {
	return []hintEntry{
		{constants.HintAliasSet, constants.HintAliasSetDesc},
		{constants.HintAliasSuggest, constants.HintAliasSuggestDesc},
		{constants.HintAliasUse, constants.HintAliasUseDesc},
	}
}

// aliasSetHints returns hints shown after gitmap a set.
func aliasSetHints() []hintEntry {
	return []hintEntry{
		{constants.HintAliasList, constants.HintAliasListDesc},
		{constants.HintAliasUse, constants.HintAliasUseDesc},
	}
}

// aliasSuggestHints returns hints shown after gitmap a suggest.
func aliasSuggestHints() []hintEntry {
	return []hintEntry{
		{constants.HintAliasList, constants.HintAliasListDesc},
		{constants.HintAliasUse, constants.HintAliasUseDesc},
	}
}
