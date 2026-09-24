package cmdtoken

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
)

func RunTokenCommand(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs || args[0] == "-h" || args[0] == "--help" {
		printTokenUsage()
		return nil
	}

	subcmd := strings.ToLower(args[0])
	subArgs := args[1:]

	switch subcmd {
	case "list", "status", "show":
		return runTokenList()
	case "add", "set":
		return runTokenAdd(subArgs)
	case "remove", "rm", "delete":
		return runTokenRemove()
	case "deploy", "sync":
		return cmdssh.RunTokenDeployCLI(subArgs)
	default:
		printTokenUsage()
		return nil
	}
}

func printTokenUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap token list                 # Display current active token status & source")
	fmt.Println("  gitmap token add <token>          # Set global git token in local configuration")
	fmt.Println("  gitmap token remove               # Remove globally configured token")
	fmt.Println("  gitmap token deploy [--nodes ...] # Distribute token to all SSH fleet nodes")
}

func runTokenList() error {
	token, source, err := ghtoken.Resolve()
	hasToken := err == nil && len(token) > 0
	if !hasToken {
		fmt.Printf("%s✖ No GitHub access token currently resolved.%s\n", constants.ColorYellow, constants.ColorReset)
		fmt.Println("  Set GH_TOKEN, login via `gh auth login`, or run `gitmap token add <token>`")
		return nil
	}

	masked := maskToken(token)
	fmt.Printf("%s✔ Active GitHub Access Token:%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  • Token:  %s%s%s\n", constants.ColorCyan, masked, constants.ColorReset)
	fmt.Printf("  • Source: %s\n", source)

	return nil
}

func runTokenAdd(args []string) error {
	hasToken := len(args) > 0
	if !hasToken {
		fmt.Printf("%sUsage: gitmap token add <token>%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	tok := strings.TrimSpace(args[0])
	cmd := exec.Command("git", "config", "--global", "github.token", tok)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to save git config token: %w", err)
	}

	fmt.Printf("%s✔ Successfully saved GitHub token to git global configuration.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func runTokenRemove() error {
	cmd := exec.Command("git", "config", "--global", "--unset", "github.token")
	_ = cmd.Run()
	fmt.Printf("%s✔ Removed GitHub token from git global configuration.%s\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func maskToken(raw string) string {
	if len(raw) <= 8 {
		return "********"
	}

	prefix := raw[:4]
	suffix := raw[len(raw)-4:]

	return fmt.Sprintf("%s****%s", prefix, suffix)
}
