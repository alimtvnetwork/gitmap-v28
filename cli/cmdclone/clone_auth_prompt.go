// Package cmdclone — clone_auth_prompt.go presents interactive terminal authentication options.
package cmdclone

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PromptTerminalAuth displays authentication choices with repository context.
func PromptTerminalAuth(repoName, repoURL string) (AuthPromptResult, error) {
	printAuthPromptBanner(repoName, repoURL)
	choice := readAuthChoice()

	switch choice {
	case "2":
		return handleBrowserAuthFlow(repoName)
	case "3":
		return AuthPromptResult{Method: AuthMethodSkip}, ErrAuthSkipped
	default:
		return handleTokenInputFlow(repoName)
	}
}

func printAuthPromptBanner(repoName, repoURL string) {
	fmt.Printf("\n%s================================================================================%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Printf("%s[!] Authentication Required for Repository:%s %s\n", constants.ColorYellow, constants.ColorReset, repoName)
	fmt.Printf("    Remote URL: %s\n", repoURL)
	fmt.Printf("--------------------------------------------------------------------------------\n")
	fmt.Printf("SSH access was not detected for this repository.\n")
	fmt.Printf("Please choose an authentication option:\n")
	fmt.Printf("  [1] Enter Personal Access Token (PAT)\n")
	fmt.Printf("  [2] Browser Login (OAuth / Device Flow)\n")
	fmt.Printf("  [3] Skip this repository\n")
	fmt.Printf("%s================================================================================%s\n", constants.ColorYellow, constants.ColorReset)
	fmt.Print("Select option [1-3] (default: 1): ")
}

func readAuthChoice() string {
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')

	return strings.TrimSpace(choice)
}
