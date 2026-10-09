// Package cmdlogin — login_browser.go handles browser-based GitHub login.
package cmdlogin

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/term"
)

// runBrowserLogin logs in via the browser: gh CLI device flow when available,
// otherwise a guided PAT creation page plus secure paste.
func runBrowserLogin() error {
	if _, err := exec.LookPath("gh"); err == nil {
		return runGhBrowserLogin()
	}
	return runGuidedBrowserLogin()
}

func runGhBrowserLogin() error {
	fmt.Println("Launching browser login via GitHub CLI...")
	cmd := exec.Command("gh", "auth", "login", "-h", "github.com", "-w")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("browser login via gh failed: %w", err)
	}

	token, source, err := secrets.Resolve()
	if err != nil || len(token) == 0 {
		return fmt.Errorf("browser login did not produce a usable token (see `gitmap login --help`)")
	}

	if err := storeGitHubToken(token); err != nil {
		return err
	}

	fmt.Printf("%s✔ Logged in to GitHub via browser.%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("  • Source: %s\n", source)
	return nil
}

func runGuidedBrowserLogin() error {
	tokenURL := "https://github.com/settings/tokens/new"
	fmt.Println("Opening GitHub token creation page in your browser...")
	openBrowser(tokenURL)
	fmt.Printf("If the browser did not open, visit:\n  %s%s%s\n", constants.ColorCyan, tokenURL, constants.ColorReset)
	fmt.Println("Create a token with the 'repo' scope, then paste it below.")

	token, err := readSecureToken()
	if err != nil {
		return err
	}
	return loginWithToken(token, false)
}

// readSecureToken reads a token without echoing it to the terminal.
func readSecureToken() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", fmt.Errorf("cannot read token securely: stdin is not a terminal (use `gitmap login --token <PAT>`)")
	}
	fmt.Print("Paste your personal access token: ")
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", fmt.Errorf("failed to read token: %w", err)
	}
	token := strings.TrimSpace(string(raw))
	if len(token) == 0 {
		return "", fmt.Errorf("token cannot be empty")
	}
	return token, nil
}

func openBrowser(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	_ = cmd.Start()
}

// runInteractiveLogin lets the user pick a login method from the terminal.
func runInteractiveLogin() error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return fmt.Errorf("no login method given and stdin is not a terminal (use `gitmap login --token <PAT>` or `gitmap login --web`)")
	}
	printLoginMethodMenu()
	choice := readMenuChoice()
	switch choice {
	case "1":
		token, err := readSecureToken()
		if err != nil {
			return err
		}
		return loginWithToken(token, false)
	case "2":
		return runBrowserLogin()
	default:
		return fmt.Errorf("unknown choice %q: pick 1 or 2", choice)
	}
}

func printLoginMethodMenu() {
	fmt.Println(constants.ColorCyan + "Log in to GitHub:" + constants.ColorReset)
	fmt.Println("  1) Paste a personal access token")
	fmt.Println("  2) Log in via the browser")
}

func readMenuChoice() string {
	fmt.Print("Choose [1/2]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')
	return strings.TrimSpace(choice)
}
