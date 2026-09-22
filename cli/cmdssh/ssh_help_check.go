package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

func hasHelpToken(arg string) bool {
	isMatch := arg == "--help" || arg == "-h" || arg == "help"

	return isMatch
}

func hasHelpFlag(args []string) bool {
	for _, a := range args {
		isHelp := hasHelpToken(a)
		if isHelp {
			return true
		}
	}

	return false
}

func hasJoinToken(args []string) bool {
	for _, a := range args {
		isJoin := a == "join" || a == "sj"
		if isJoin {
			return true
		}
	}

	return false
}

func checkSSHHelp(args []string) bool {
	hasJoin := hasJoinToken(args)
	hasHelp := hasHelpFlag(args)
	hasJoinHelp := hasJoin && hasHelp

	if hasJoinHelp {
		helptext.Print("ssh-join")

		return true
	}

	if hasHelp {
		RenderSSHHelp()

		return true
	}

	return false
}
