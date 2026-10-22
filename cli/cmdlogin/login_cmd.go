// Package cmdlogin — login_cmd.go parses gitmap login/logout flags and dispatches flows.
package cmdlogin

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type loginOptions struct {
	tokenValue string
	hasToken   bool
	useBrowser bool
	showStatus bool
	showHelp   bool
	skipVerify bool
}

func runLogin(args []string) error {
	if len(args) == 0 {
		return runInteractiveLogin()
	}

	handled, err := dispatchSubcommand(args)
	if handled {
		return err
	}

	opts, rest, parseErr := parseLoginArgs(args)
	if parseErr != nil {
		return parseErr
	}
	if len(rest) > 0 {
		return apperror.NewValidationError(fmt.Sprintf("unexpected argument %q (see `gitmap login --help`)", rest[0]))
	}

	return dispatchLoginOptions(opts)
}

func dispatchSubcommand(args []string) (bool, error) {
	sub := strings.ToLower(args[0])
	switch sub {
	case "ssh":
		return true, RunLoginSSH(args[1:])
	case "browser-flow", "browser":
		return true, runBrowserLogin()
	case "token", "pat":
		return true, runTokenSubcommand(args[1:])
	case "status":
		return true, runLoginStatus()
	default:
		return false, nil
	}
}

func runTokenSubcommand(args []string) error {
	skipVerify := false
	token := ""
	for _, a := range args {
		if a == "--no-verify" {
			skipVerify = true
			continue
		}
		if !strings.HasPrefix(a, "-") && len(token) == 0 {
			token = a
		}
	}

	if len(token) > 0 {
		return loginWithToken(token, skipVerify)
	}

	rawToken, err := readSecureToken()
	if err != nil {
		return err
	}
	return loginWithToken(rawToken, skipVerify)
}

func dispatchLoginOptions(opts loginOptions) error {
	if opts.showHelp {
		printLoginUsage()
		return nil
	}
	if opts.showStatus {
		return runLoginStatus()
	}
	if opts.useBrowser {
		return runBrowserLogin()
	}
	if opts.hasToken {
		return loginWithToken(opts.tokenValue, opts.skipVerify)
	}
	return runInteractiveLogin()
}

func parseLoginArgs(args []string) (loginOptions, []string, error) {
	opts := loginOptions{}
	rest := []string{}
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "-h" || arg == "--help" {
			opts.showHelp = true
			i++
			continue
		}
		if arg == "--status" {
			opts.showStatus = true
			i++
			continue
		}
		if arg == "--web" || arg == "--browser" {
			opts.useBrowser = true
			i++
			continue
		}
		if arg == "--no-verify" {
			opts.skipVerify = true
			i++
			continue
		}
		if arg == "--token" && i+1 >= len(args) {
			return opts, rest, apperror.NewValidationError("missing value for --token (see `gitmap login --help`)")
		}
		if arg == "--token" {
			opts.tokenValue = strings.TrimSpace(args[i+1])
			opts.hasToken = true
			i += 2
			continue
		}
		if strings.HasPrefix(arg, "--token=") {
			opts.tokenValue = strings.TrimSpace(strings.TrimPrefix(arg, "--token="))
			opts.hasToken = true
			i++
			continue
		}
		rest = append(rest, arg)
		i++
	}
	return opts, rest, nil
}

func printLoginUsage() {
	fmt.Println(constants.ColorCyan + "Usage:" + constants.ColorReset)
	fmt.Println("  gitmap login                  # Interactive: choose SSH, PAT, Browser, or Status")
	fmt.Println("  gitmap login ssh [alias] [-y] # Generate or reuse SSH key and copy public key to clipboard")
	fmt.Println("  gitmap login token [PAT]      # Log in with a personal access token")
	fmt.Println("  gitmap login browser-flow     # Log in via the browser (GitHub CLI device flow)")
	fmt.Println("  gitmap login --status         # Show current login state")
	fmt.Println("  gitmap logout                 # Remove the stored GitHub credential")
	fmt.Println("")
	fmt.Println(constants.ColorCyan + "Notes:" + constants.ColorReset)
	fmt.Println("  --token validates the token against api.github.com before storing it.")
	fmt.Println("  The token is stored in git global config (github.token) so clone/pull/push reuse it.")
	fmt.Println("  Use --no-verify with --token to store without validation (offline use).")
}
