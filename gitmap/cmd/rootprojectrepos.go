package cmd

import (
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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
	if command == constants.CmdReactRepos || command == constants.CmdReactReposAlias {
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
