// Package cmdclone — clone_auth_browser.go handles browser-based authentication flow.
package cmdclone

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
)

func handleBrowserAuthFlow(repoName string) (AuthPromptResult, error) {
	fmt.Println("\nLaunching browser authentication...")
	runBrowserLoginCommand()

	tok, _, err := ghtoken.Resolve()
	if err != nil || tok == "" {
		fmt.Printf("Browser login was not completed. Please create a token at: https://github.com/settings/tokens/new\n")
		return handleTokenInputFlow(repoName)
	}

	isReused := askTokenReuseConfirmation()
	if isReused {
		SetGlobalAccessToken(tok)
	}

	return AuthPromptResult{Token: tok, IsReused: isReused, Method: AuthMethodBrowser}, nil
}

func runBrowserLoginCommand() {
	if _, err := exec.LookPath("gh"); err == nil {
		cmd := exec.Command("gh", "auth", "login", "-h", "github.com", "-w")
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}
}
