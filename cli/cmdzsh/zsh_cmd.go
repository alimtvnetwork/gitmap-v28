// Package cmdzsh provides the main CLI entrypoint for ZSH suite operations.
package cmdzsh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunZsh executes zsh commands and routes to subcommands.
func RunZsh(args []string) error {
	if len(args) == 0 {
		return dispatchStatus([]string{})
	}

	return result.AsError(routeZshSubcommand(args[0], args[1:]))
}

func buildZshRouter() map[string]func([]string) error {
	return map[string]func([]string) error{
		"install": dispatchInstall,
		"theme":   dispatchTheme,
		"themes":  dispatchThemes,
		"clean":   dispatchClean,
		"switch":  dispatchSwitch,
		"profile": dispatchProfile,
		"status":  dispatchStatus,
	}
}

func routeZshSubcommand(subcmd string, args []string) result.ErrorWrapper {
	fn, isFound := buildZshRouter()[subcmd]
	if isFound {
		return result.MatchWrapper(fn(args))
	}

	return routeFallback(subcmd, args)
}

func routeFallback(subcmd string, args []string) result.ErrorWrapper {
	isHelp := subcmd == "help" || subcmd == "-h" || subcmd == "--help"
	if isHelp {
		return result.MatchWrapper(printZshHelp())
	}

	return result.MatchWrapper(handleUnknownSubcommand(subcmd, args))
}

func handleUnknownSubcommand(subcmd string, args []string) error {
	hasDash := strings.HasPrefix(subcmd, "-")
	if hasDash {
		return printZshHelp()
	}

	return dispatchInstall(append([]string{subcmd}, args...))
}

func printZshHelp() error {
	printZshUsage()
	printZshCommands()

	return nil
}

func printZshUsage() {
	fmt.Println("Usage: gitmap zsh <command> [flags]")
	fmt.Println()
	fmt.Println("Manage ZSH, Oh-My-Zsh, themes, plugins, and user profiles.")
	fmt.Println()
}

func printZshCommands() {
	fmt.Println("Commands:")
	fmt.Println("  install   Install zsh, oh-my-zsh, plugins, and set theme")
	fmt.Println("  theme     Change ZSH theme in ~/.zshrc")
	fmt.Println("  themes    List 41 supported Oh-My-Zsh themes")
	fmt.Println("  clean     Backup and remove Oh-My-Zsh, with optional reinstall")
	fmt.Println("  switch    Switch default login shell to ZSH")
	fmt.Println("  profile   Provision workspace directories (scripts, gitlab, github, .ssh)")
	fmt.Println("  status    Inspect ZSH version, Oh-My-Zsh presence, and current theme")
	fmt.Println("  help      Show help for zsh commands")
}
