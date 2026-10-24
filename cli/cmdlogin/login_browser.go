// Package cmdlogin — login_browser.go handles browser-based GitHub login.
package cmdlogin

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/enums"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/term"
)

// runBrowserLogin logs in via the browser: gh CLI device flow when available,
// otherwise a guided PAT creation page plus secure paste.
func runBrowserLogin() error {
	if _, err := exec.LookPath(constants.BinGh); err == nil {
		return runGhBrowserLogin()
	}

	return runGuidedBrowserLogin()
}

func runGhBrowserLogin() error {
	if err := executeGhAuth(); err != nil {
		return err
	}

	return resolveAndStoreGhToken()
}

func executeGhAuth() error {
	fmt.Println("Launching browser login via GitHub CLI...")
	cmd := exec.Command(constants.BinGh, constants.GhCmdAuth, constants.GhCmdLogin, constants.GhHostFlag, constants.GhHostDomain, constants.GhWebFlag)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "browser login via gh failed")
	}

	return nil
}

func resolveAndStoreGhToken() error {
	token, source, err := secrets.Resolve()
	if err != nil || len(token) == 0 {
		return apperror.NewExecutionError("browser login did not produce a usable token (see `gitmap login --help`)")
	}

	if err := storeGitHubToken(token); err != nil {
		return err
	}

	fmt.Printf("%s✔ Logged in to GitHub via browser.%s\n", appconstants.ColorGreen, appconstants.ColorReset)
	fmt.Printf("  • Source: %s\n", source)

	return nil
}

func runGuidedBrowserLogin() error {
	printGuidedBrowserNotice(constants.GithubNewTokenURL)
	openBrowser(constants.GithubNewTokenURL)

	token, err := readSecureToken()
	if err != nil {
		return err
	}

	return loginWithToken(token, false)
}

func printGuidedBrowserNotice(tokenURL string) {
	fmt.Println("Opening GitHub token creation page in your browser...")
	fmt.Printf("If the browser did not open, visit:\n  %s%s%s\n", appconstants.ColorCyan, tokenURL, appconstants.ColorReset)
	fmt.Println("Create a token with the 'repo' scope, then paste it below.")
}

// readSecureToken reads a token without echoing it to the terminal.
func readSecureToken() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return "", apperror.NewValidationError("cannot read token securely: stdin is not a terminal (use `gitmap login --token <PAT>`)")
	}

	return promptTerminalPassword()
}

func promptTerminalPassword() (string, error) {
	fmt.Print("Paste your personal access token: ")
	raw, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", apperror.WrapSimple(err, "failed to read token")
	}

	token := strings.TrimSpace(string(raw))
	if len(token) == 0 {
		return "", apperror.NewValidationError("token cannot be empty")
	}

	return token, nil
}

func openBrowser(url string) {
	cmd := resolveBrowserCommand(url)
	_ = cmd.Start()
}

func resolveBrowserCommand(url string) *exec.Cmd {
	switch runtime.GOOS {
	case constants.PlatformDarwin:
		return exec.Command(constants.BrowserOpenDarwin, url)
	case constants.PlatformWindows:
		return exec.Command(constants.BrowserOpenWindows, constants.BrowserWindowsDllArg, url)
	default:
		return exec.Command(constants.BrowserOpenLinux, url)
	}
}

// runInteractiveLogin lets the user pick a login method from the terminal.
func runInteractiveLogin() error {
	if !term.IsTerminal(int(os.Stdin.Fd())) {
		return apperror.NewValidationError("no login method given and stdin is not a terminal (use `gitmap login --token <PAT>` or `gitmap login --web`)")
	}

	printLoginMethodMenu()
	choice := readMenuChoice()

	return handleInteractiveChoice(choice)
}

func handleInteractiveChoice(choice string) error {
	typedChoice, err := enums.ParseLoginInteractiveChoice(choice)
	if err != nil {
		return apperror.NewValidationError(fmt.Sprintf("unknown choice %q: pick 1, 2, 3, or 4", choice))
	}

	switch typedChoice {
	case enums.LoginInteractiveChoiceSSH:
		return RunLoginSSH(nil)
	case enums.LoginInteractiveChoiceToken:
		return promptAndLoginToken()
	case enums.LoginInteractiveChoiceBrowser:
		return runBrowserLogin()
	case enums.LoginInteractiveChoiceStatus:
		return runLoginStatus()
	default:
		return apperror.NewValidationError(fmt.Sprintf("unknown choice %q: pick 1, 2, 3, or 4", choice))
	}
}

func promptAndLoginToken() error {
	token, err := readSecureToken()
	if err != nil {
		return err
	}

	return loginWithToken(token, false)
}

func printLoginMethodMenu() {
	fmt.Println(appconstants.ColorCyan + "GitMap Authentication:" + appconstants.ColorReset)
	fmt.Println("  [1] SSH Key Setup & Copy (Recommended)")
	fmt.Println("  [2] Personal Access Token (PAT)")
	fmt.Println("  [3] Browser Device Flow")
	fmt.Println("  [4] Check Authentication Status")
}

func readMenuChoice() string {
	fmt.Print("Choose [1-4]: ")
	reader := bufio.NewReader(os.Stdin)
	choice, _ := reader.ReadString('\n')

	return strings.TrimSpace(choice)
}
