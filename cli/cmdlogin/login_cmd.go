// Package cmdlogin — login_cmd.go parses gitmap login/logout flags and dispatches flows.
package cmdlogin

import (
	"fmt"
	"strings"

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
	opts, rest, err := parseLoginArgs(args)
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return fmt.Errorf("unexpected argument %q (see `gitmap login --help`)", rest[0])
	}
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
		if arg == "--token" {
			if i+1 >= len(args) {
				return opts, rest, fmt.Errorf("missing value for --token (see `gitmap login --help`)")
			}
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
	fmt.Println("  gitmap login                  # Interactive: pick token paste or browser login")
	fmt.Println("  gitmap login --token <PAT>    # Log in with a personal access token (non-interactive)")
	fmt.Println("  gitmap login --web            # Log in via the browser (GitHub CLI device flow when available)")
	fmt.Println("  gitmap login --status         # Show current login state")
	fmt.Println("  gitmap logout                 # Remove the stored GitHub credential")
	fmt.Println("")
	fmt.Println(constants.ColorCyan + "Notes:" + constants.ColorReset)
	fmt.Println("  --token validates the token against api.github.com before storing it.")
	fmt.Println("  The token is stored in git global config (github.token) so clone/pull/push reuse it.")
	fmt.Println("  Use --no-verify with --token to store without validation (offline use).")
}
