package cmdworkdir

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
)

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
	}
}

func autoRegisterFirstWorkDir(absDir string, quiet bool) bool {
	return cmdscan.AutoRegisterFirstWorkDir(absDir, quiet)
}

func isWorkDirKeyword(name string) bool {
	return name == "work" || name == "workdir" || name == "wd" || name == "default"
}
