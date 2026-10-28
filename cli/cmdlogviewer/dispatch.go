package cmdlogviewer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdgit"
)

// FallbackDashboardUIFn can be set by cmdui to launch the Web UI dashboard on non-git container folders.
var FallbackDashboardUIFn func(page string, preferredPort int) error

func init() {
	cmdgit.RunLogUIFn = RunShowLog
}

// DispatchShowLog routes the command if it matches show-log variants.
func DispatchShowLog(command string, args []string) (bool, error) {
	if IsShowLogCommand(command) {
		return true, RunShowLog(args)
	}

	return false, nil
}

// IsShowLogCommand reports whether command is a show-log alias.
func IsShowLogCommand(command string) bool {
	clean := strings.ToLower(strings.TrimSpace(command))

	return clean == "show-log" || clean == "showlog" || clean == "log-ui"
}

// IsLogUIArgs checks if the argument slice begins with the "ui" subcommand of log.
func IsLogUIArgs(args []string) bool {
	if len(args) == 0 {
		return false
	}

	clean := strings.ToLower(strings.TrimSpace(args[0]))

	return clean == "ui"
}
