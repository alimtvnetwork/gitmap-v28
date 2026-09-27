package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// dispatchProjectRepos routes project type query commands.
func dispatchProjectRepos(command string) (bool, error) {
	if command == constants.CmdGoRepos || command == constants.CmdGoReposAlias {
		runProjectRepos(constants.ProjectKeyGo, os.Args[2:])

		return true, nil
	}

	if command == constants.CmdNodeRepos || command == constants.CmdNodeReposAlias {
		runProjectRepos(constants.ProjectKeyNode, os.Args[2:])

		return true, nil
	}

	if isReactReposCommand(command, os.Args[2:]) {
		runProjectRepos(constants.ProjectKeyReact, os.Args[2:])

		return true, nil
	}

	if command == constants.CmdCppRepos || command == constants.CmdCppReposAlias {
		runProjectRepos(constants.ProjectKeyCpp, os.Args[2:])

		return true, nil
	}

	if command == constants.CmdCsharpRepos || command == constants.CmdCsharpAlias {
		runProjectRepos(constants.ProjectKeyCsharp, os.Args[2:])

		return true, nil
	}

	return false, nil
}

func isReactReposCommand(cmd string, args []string) bool {
	if cmd == constants.CmdReactRepos {
		return true
	}
	if cmd != constants.CmdReactReposAlias {
		return false
	}

	return isReactReposOnlyArgs(args)
}

func isReactReposOnlyArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}
	first := strings.ToLower(strings.TrimSpace(args[0]))

	return first == "--count"
}
