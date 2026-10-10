package cmd

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdwinutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"strings"
)

func RunCleanDevTopLevel(args []string) error {
	return cmdos.RunOSDevClean(args)
}

func RunDevToolTopLevel(args []string) error {
	return runDevTopLevel(args)
}

func RunWinUtilTopLevel(args []string) error {
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

func RunCleanTopLevel(args []string) error {
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
	if !isKnownCleanOrClearFlag(sub) {
		return handleUnknownCleanOrClearTarget("clean", sub)
	}
	return cmdos.RunOSDevClean(args)
}

func runCleanTerminalTopLevel(args []string) error {
	return cmdos.RunTerminalCleanCLI(args)
}

func RunClearTopLevel(args []string) error {
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
	if !isKnownCleanOrClearFlag(sub) {
		return handleUnknownCleanOrClearTarget("clear", sub)
	}
	return cmdos.RunTerminalCleanCLI(args)
}

func RunTerminalTopLevel(args []string) error {
	return cmdos.RunTerminalCleanCLI(args)
}

func isKnownCleanOrClearFlag(sub string) bool {
	if strings.HasPrefix(sub, "-") || sub == "/?" || sub == "/y" {
		return true
	}
	return sub == "help" || sub == "yes" || sub == "dry-run" || sub == "cache" || sub == "caches"
}

func handleUnknownCleanOrClearTarget(verb, sub string) error {
	fullCmd := verb + " " + sub
	sugg := []string{
		"gitmap clear devtools",
		"gitmap clear dev-tools-cache",
		"gitmap clear terminal",
		"gitmap clear-terminal",
		"gitmap agy clean-cache",
	}
	msg := fmt.Sprintf("unknown %s target '%s'", verb, sub)
	store.LogFailedCommand(fullCmd, strings.Join(os.Args[1:], " "), verb, "E1001", msg, sugg)

	fmt.Printf("\n  %s✗ Unknown %s target: '%s'%s\n\n", constants.ColorRed, verb, sub, constants.ColorReset)
	fmt.Printf("  %s💡 Available Clear & Clean Suggestions:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("    • gitmap clear devtools         - Clean Go, npm, pnpm, Cargo, pip, Maven & Gradle caches")
	fmt.Println("    • gitmap clear dev-tools-cache  - Alias for developer tools cache cleanup")
	fmt.Println("    • gitmap clear terminal         - Clear PowerShell, Bash & Zsh history and reseed suggestions")
	fmt.Println("    • gitmap clear-terminal         - Direct alias to clear terminal history & suggestions")
	fmt.Println("    • gitmap agy clean-cache        - Purge Antigravity IDE cache & logs")
	fmt.Println("    • gitmap failed-commands count  - Review total failed command count")
	fmt.Println("    • gitmap failed-commands clear  - Clear recorded failed commands log")
	fmt.Println("    • gitmap clean-artifacts        - Clean workspace temporary files and build artifacts")
	fmt.Println()
	return apperror.NewValidationError(msg)
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
