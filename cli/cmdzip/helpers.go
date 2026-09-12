// Package cmdzip provides archive compression, extraction, and zip-group management commands.
package cmdzip

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkHelp(command string, args []string) {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

type hintEntry struct {
	command     string
	description string
}

func printHints(hints []hintEntry) {
	fmt.Fprint(os.Stderr, constants.MsgHintHeader)
	for _, h := range hints {
		fmt.Fprintf(os.Stderr, constants.MsgHintRowFmt, h.command, h.description)
	}
}

func zipGroupCreateHints() []hintEntry {
	return []hintEntry{
		{constants.HintZGAdd, constants.HintZGAddDesc},
		{constants.HintZGRelease, constants.HintZGReleaseDesc},
	}
}

func zipGroupListHints() []hintEntry {
	return []hintEntry{
		{constants.HintZGCreate, constants.HintZGCreateDesc},
		{constants.HintZGShow, constants.HintZGShowDesc},
	}
}

func zipGroupShowHints() []hintEntry {
	return []hintEntry{
		{constants.HintZGAdd, constants.HintZGAddDesc},
		{constants.HintZGRelease, constants.HintZGReleaseDesc},
		{constants.HintZGDelete, constants.HintZGDeleteDesc},
	}
}

var knownValueFlags = map[string]bool{
	"--format": true, "--out": true, "-o": true,
	"-p": true, "--prefix": true,
}

func reorderFlagsBeforeArgs(args []string) []string {
	var flags []string
	var positional []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
			continue
		}

		flags = append(flags, arg)
		if knownValueFlags[arg] && i+1 < len(args) {
			flags = append(flags, args[i+1])
			i++
		}
	}

	return append(flags, positional...)
}
