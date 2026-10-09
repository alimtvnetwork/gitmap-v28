// Package cmdlogin — login_status.go shows login state and removes credentials.
package cmdlogin

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// runLoginStatus prints the current GitHub login state with the token masked.
func runLoginStatus() error {
	token, source, err := secrets.Resolve()
	if err != nil || len(token) == 0 {
		fmt.Printf("%s✖ Not logged in to GitHub.%s\n", constants.ColorYellow, constants.ColorReset)
		fmt.Println("  Run `gitmap login` to log in.")
		return nil
	}

	fmt.Printf("%s✔ Logged in to GitHub.%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  • Token:  %s%s%s\n", constants.ColorCyan, maskLoginToken(token), constants.ColorReset)
	fmt.Printf("  • Source: %s\n", source)
	return nil
}

// runLogout removes the stored GitHub credential.
func runLogout(args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
		fmt.Println("  gitmap logout    # Remove the stored GitHub credential")
		return nil
	}

	token, _, err := secrets.Resolve()
	hasStoredToken := err == nil && len(token) > 0
	if !hasStoredToken {
		fmt.Printf("%sNo GitHub credential is currently stored.%s\n", constants.ColorYellow, constants.ColorReset)
		return nil
	}

	if err := removeStoredTokens(); err != nil {
		return err
	}

	fmt.Printf("%s✔ Logged out: GitHub credential removed.%s\n", constants.ColorGreen, constants.ColorReset)
	return nil
}

func removeStoredTokens() error {
	unset := exec.Command("git", "config", "--global", "--unset", "github.token")
	_ = unset.Run()

	credApprove := exec.Command("git", "credential", "reject")
	credApprove.Stdin = strings.NewReader("protocol=https\nhost=github.com\n\n")
	_ = credApprove.Run()

	return nil
}

func maskLoginToken(raw string) string {
	if len(raw) <= 8 {
		return "********"
	}
	return fmt.Sprintf("%s****%s", raw[:4], raw[len(raw)-4:])
}
