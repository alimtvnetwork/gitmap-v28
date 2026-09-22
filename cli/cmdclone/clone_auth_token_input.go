// Package cmdclone — clone_auth_token_input.go captures user token and reuse preference.
package cmdclone

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func readTokenFromTerminal(repoName string) string {
	fmt.Printf("\nEnter Personal Access Token for %s: ", repoName)
	reader := bufio.NewReader(os.Stdin)
	token, _ := reader.ReadString('\n')

	return strings.TrimSpace(token)
}

func askTokenReuseConfirmation() bool {
	fmt.Print("Reuse this access token for all repositories in this session? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	ans, _ := reader.ReadString('\n')
	ans = strings.ToLower(strings.TrimSpace(ans))

	return ans == "" || ans == "y" || ans == "yes"
}

func handleTokenInputFlow(repoName string) (AuthPromptResult, error) {
	token := readTokenFromTerminal(repoName)
	if token == "" {
		return AuthPromptResult{}, fmt.Errorf("token cannot be empty")
	}

	isReused := askTokenReuseConfirmation()
	if isReused {
		SetGlobalAccessToken(token)
	}

	return AuthPromptResult{Token: token, IsReused: isReused, Method: AuthMethodToken}, nil
}
