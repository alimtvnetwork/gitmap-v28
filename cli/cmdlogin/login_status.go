// Package cmdlogin — login_status.go shows login state and removes credentials.
package cmdlogin

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// runLoginStatus prints the current GitHub login state with the token masked.
func runLoginStatus() error {
	token, source, err := secrets.Resolve()
	if err != nil || len(token) == 0 {
		fmt.Printf("%s✖ Not logged in to GitHub.%s\n", appconstants.ColorYellow, appconstants.ColorReset)
		fmt.Println("  Run `gitmap login` to log in.")
		return nil
	}

	fmt.Printf("%s✔ Logged in to GitHub.%s\n", appconstants.ColorGreen, appconstants.ColorReset)
	fmt.Printf("  • Token:  %s%s%s\n", appconstants.ColorCyan, maskLoginToken(token), appconstants.ColorReset)
	fmt.Printf("  • Source: %s\n", source)

	return nil
}

// runLogout removes the stored GitHub credential.
func runLogout(args []string) error {
	if isLogoutHelp(args) {
		printLogoutUsage()
		return nil
	}

	token, _, err := secrets.Resolve()
	hasStoredToken := err == nil && len(token) > 0
	if !hasStoredToken {
		fmt.Printf("%sNo GitHub credential is currently stored.%s\n", appconstants.ColorYellow, appconstants.ColorReset)
		return nil
	}

	if err := removeStoredTokens(); err != nil {
		return err
	}

	fmt.Printf("%s✔ Logged out: GitHub credential removed.%s\n", appconstants.ColorGreen, appconstants.ColorReset)

	return nil
}

func isLogoutHelp(args []string) bool {
	return len(args) > 0 && (args[0] == constants.FlagHelpShortArg || args[0] == constants.FlagHelpArg)
}

func printLogoutUsage() {
	fmt.Println(appconstants.ColorCyan + "Usage:" + appconstants.ColorReset)
	fmt.Println("  gitmap logout    # Remove the stored GitHub credential")
}

func removeStoredTokens() error {
	unset := exec.Command(constants.GitBin, constants.GitConfig, constants.GitConfigGlobalFlag, constants.GitConfigUnsetFlag, constants.GitConfigGithubToken)
	_ = unset.Run()

	credApprove := exec.Command(constants.GitBin, constants.GitCredential, constants.GitCredentialReject)
	credApprove.Stdin = strings.NewReader(constants.GitCredentialRejectPayload)
	_ = credApprove.Run()

	return nil
}

func maskLoginToken(raw string) string {
	if len(raw) <= 8 {
		return constants.TokenMaskFull
	}

	return fmt.Sprintf("%s****%s", raw[:4], raw[len(raw)-4:])
}
