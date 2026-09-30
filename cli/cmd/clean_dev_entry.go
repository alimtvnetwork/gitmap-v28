package cmd

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
)

func runCleanDevTopLevel(args []string) error {
	return cmdos.RunOSDevClean(args)
}

func runDevToolTopLevel(args []string) error {
	return cmdos.RunDevToolCLI(args)
}

func runWinUtilTopLevel(args []string) error {
	return cmdwinutil.RunWinUtil(args)
}

func runDevTopLevel(args []string) error {
	if len(args) == 0 {
		return cmdos.RunDevToolCLI(args)
	}
	sub := strings.ToLower(args[0])
	if isDevCleanActionVerb(sub) {
		return cmdos.RunOSDevClean(args[1:])
	}
	return cmdos.RunDevToolCLI(args)
}

func runCleanTopLevel(args []string) error {
	if len(args) == 0 {
		return cmdos.RunOSDevClean(args)
	}
	sub := strings.ToLower(args[0])
	if isTerminalCleanSubToken(sub) {
		return runCleanTerminalTopLevel(args[1:])
	}
	if isDevCleanSubToken(sub) {
		return cmdos.RunOSDevClean(args[1:])
	}
	return cmdos.RunOSDevClean(args)
}

func runCleanTerminalTopLevel(args []string) error {
	return cmdos.RunTerminalCleanCLI(args)
}

func runClearTopLevel(args []string) error {
	if len(args) == 0 {
		return cmdos.RunTerminalCleanCLI(args)
	}
	sub := strings.ToLower(args[0])
	if isDevCleanSubToken(sub) {
		return cmdos.RunOSDevClean(args[1:])
	}
	if isTerminalCleanSubToken(sub) {
		return cmdos.RunTerminalCleanCLI(args[1:])
	}
	return cmdos.RunTerminalCleanCLI(args)
}

func runTerminalTopLevel(args []string) error {
	return cmdos.RunTerminalCleanCLI(args)
}

func isTerminalCleanSubToken(sub string) bool {
	return sub == "terminal" || sub == "term" || sub == "console" ||
		sub == "shell-history" || sub == "history" ||
		sub == "terminal-history" || sub == "terminal-suggestions" ||
		sub == "terminals" || sub == "clear-terminal" || sub == "clean-terminal" ||
		sub == "clear" || sub == "terminal-clear"
}

func isDevCleanActionVerb(sub string) bool {
	return sub == "clean" || sub == "clear" || sub == "cleanup" || sub == "cache"
}

func isDevCleanSubToken(sub string) bool {
	return sub == "dev" || sub == "devs" || sub == "developer" ||
		sub == "devtool" || sub == "devtools" ||
		sub == "dev-tool" || sub == "dev-tools" ||
		sub == "devtool-cache" || sub == "devtools-cache" ||
		sub == "dev-tool-cache" || sub == "dev-tools-cache"
}
