// Package cmdlogin — login_token.go validates a GitHub token and stores it.
package cmdlogin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const githubUserAPI = "https://api.github.com/user"

type githubUser struct {
	Login string `json:"login"`
}

// loginWithToken validates the token (unless skipped) and stores it for reuse.
func loginWithToken(token string, skipVerify bool) error {
	if len(strings.TrimSpace(token)) == 0 {
		return fmt.Errorf("token cannot be empty (see `gitmap login --help`)")
	}
	token = strings.TrimSpace(token)

	username := ""
	if !skipVerify {
		name, err := validateGitHubToken(token)
		if err != nil {
			return err
		}
		username = name
	}

	if err := storeGitHubToken(token); err != nil {
		return err
	}

	printLoginSuccess(username, skipVerify)
	return nil
}

// validateGitHubToken calls api.github.com/user and returns the GitHub username.
func validateGitHubToken(token string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, githubUserAPI, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build validation request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("token validation failed: cannot reach api.github.com: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return "", fmt.Errorf("GitHub rejected the token (HTTP %d): check the token value and its scopes", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub returned HTTP %d during token validation", resp.StatusCode)
	}

	user := githubUser{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", fmt.Errorf("failed to parse GitHub validation response: %w", err)
	}
	if len(strings.TrimSpace(user.Login)) == 0 {
		return "", fmt.Errorf("GitHub validation response did not contain a username")
	}

	return user.Login, nil
}

// storeGitHubToken persists the token in git global config for reuse by clone/pull/push.
func storeGitHubToken(token string) error {
	cmd := exec.Command("git", "config", "--global", "github.token", token)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to store GitHub token: %w", err)
	}
	return nil
}

func printLoginSuccess(username string, skippedVerify bool) {
	fmt.Printf("%s✔ Logged in to GitHub.%s\n", constants.ColorGreen, constants.ColorReset)
	if len(username) > 0 {
		fmt.Printf("  • Account: %s%s%s\n", constants.ColorCyan, username, constants.ColorReset)
	}
	if skippedVerify {
		fmt.Printf("  %s! Token stored without validation (--no-verify).%s\n", constants.ColorYellow, constants.ColorReset)
	}
	fmt.Println("  Private-repo clone/pull/push will now reuse this credential.")
}
