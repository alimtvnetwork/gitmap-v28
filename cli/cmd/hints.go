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
