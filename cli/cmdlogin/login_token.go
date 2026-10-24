// Package cmdlogin — login_token.go validates a GitHub token and stores it.
package cmdlogin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type githubUser struct {
	Login string `json:"login"`
}

func resolveTokenUsername(token string, isVerify bool) (string, error) {
	if !isVerify {
		return "", nil
	}

	return validateGitHubToken(token)
}

// loginWithToken validates the token (unless skipped) and stores it for reuse.
func loginWithToken(token string, isVerify bool) error {
	trimmedToken := strings.TrimSpace(token)
	if len(trimmedToken) == 0 {
		return apperror.NewValidationError("token cannot be empty (see `gitmap login --help`)")
	}

	username, err := resolveTokenUsername(trimmedToken, isVerify)
	if err != nil {
		return err
	}

	if err := storeGitHubToken(trimmedToken); err != nil {
		return err
	}

	printLoginSuccess(username, isVerify)

	return nil
}

// validateGitHubToken calls api.github.com/user and returns the GitHub username.
func validateGitHubToken(token string) (string, error) {
	req, err := newGitHubUserRequest(token)
	if err != nil {
		return "", err
	}

	resp, err := executeGitHubUserRequest(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return parseGitHubUserResponse(resp)
}

func newGitHubUserRequest(token string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, constants.GithubUserAPI, nil)
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to build validation request")
	}

	req.Header.Set(constants.HeaderAuthorization, constants.HeaderAuthorizationBearerPrefix+token)
	req.Header.Set(constants.HeaderAccept, constants.HeaderGithubJSON)

	return req, nil
}

func executeGitHubUserRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{Timeout: constants.TokenValidationTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, apperror.WrapSimple(err, "token validation failed: cannot reach api.github.com")
	}

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, apperror.NewExecutionError(fmt.Sprintf("GitHub rejected the token (HTTP %d): check the token value and its scopes", resp.StatusCode))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, apperror.NewExecutionError(fmt.Sprintf("GitHub returned HTTP %d during token validation", resp.StatusCode))
	}

	return resp, nil
}

func parseGitHubUserResponse(resp *http.Response) (string, error) {
	user := githubUser{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return "", apperror.WrapSimple(err, "failed to parse GitHub validation response")
	}

	if len(strings.TrimSpace(user.Login)) == 0 {
		return "", apperror.NewNotFoundError("GitHub validation response did not contain a username")
	}

	return user.Login, nil
}

// storeGitHubToken persists the token in git global config for reuse by clone/pull/push.
func storeGitHubToken(token string) error {
	cmd := exec.Command(constants.GitBin, constants.GitConfig, constants.GitConfigGlobalFlag, constants.GitConfigGithubToken, token)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to store GitHub token")
	}

	return nil
}

func printLoginSuccess(username string, isVerify bool) {
	fmt.Printf("%s✔ Logged in to GitHub.%s\n", appconstants.ColorGreen, appconstants.ColorReset)
	if len(username) > 0 {
		fmt.Printf("  • Account: %s%s%s\n", appconstants.ColorCyan, username, appconstants.ColorReset)
	}

	if !isVerify {
		fmt.Printf("  %s! Token stored without validation (--no-verify).%s\n", appconstants.ColorYellow, appconstants.ColorReset)
	}

	fmt.Println("  Private-repo clone/pull/push will now reuse this credential.")
}
