package cmd

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
)

func init() {
	hookGitOpsDispatch()
}

func hookGitOpsDispatch() {
	if len(os.Args) < 2 {
		return
	}

	cmd := strings.ToLower(strings.TrimSpace(os.Args[1]))
	if handleRootGitOps(cmd) {
		return
	}

	handleGitSubcommandOps(cmd)
}

func handleRootGitOps(cmd string) bool {
	if cmd == "show" {
		executeShowAndExit(os.Args[2:])

		return true
	}

	if cmd == "log" {
		executeLogAndExit(os.Args[2:])

		return true
	}

	if isShowLogCommand(cmd) {
		executeShowLogAndExit(os.Args[2:])

		return true
	}

	return false
}

func isShowLogCommand(cmd string) bool {
	return cmd == "show-log" || cmd == "showlog" || cmd == "log-ui"
}

func handleGitSubcommandOps(cmd string) {
	if cmd != "git" || len(os.Args) <= 2 {
		return
	}

	sub := strings.ToLower(strings.TrimSpace(os.Args[2]))
	if sub == "show" {
		executeShowAndExit(os.Args[3:])
	}

	if sub == "log" {
		executeLogAndExit(os.Args[3:])
	}

	if isShowLogCommand(sub) {
		executeShowLogAndExit(os.Args[3:])
	}
}

func executeShowAndExit(args []string) {
	if err := cmdgit.RunShow(args); err != nil {
		cliexit.HandleError(err, 1)
	}

	cliexit.Exit(0)
}

func executeLogAndExit(args []string) {
	if err := cmdgit.RunLog(args); err != nil {
		cliexit.HandleError(err, 1)
	}

	cliexit.Exit(0)
}

func executeShowLogAndExit(args []string) {
	if err := cmdgit.RunShowLog(args); err != nil {
		cliexit.HandleError(err, 1)
	}

	cliexit.Exit(0)
}

// DispatchGitOps handles programmatic invocation of git inspection operations.
func DispatchGitOps(command string, args []string) (bool, error) {
	clean := strings.ToLower(strings.TrimSpace(command))

	switch clean {
	case "show":
		return true, cmdgit.RunShow(args)
	case "log":
		return true, cmdgit.RunLog(args)
	case "show-log", "showlog", "log-ui":
		return true, cmdgit.RunShowLog(args)
	default:
		return false, nil
	}
}
