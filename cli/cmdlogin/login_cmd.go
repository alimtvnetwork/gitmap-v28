// Package cmdlogin — login_cmd.go parses gitmap login/logout flags and dispatches flows.
package cmdlogin

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdlogin/enums"
	appconstants "github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type loginOptions struct {
	tokenValue         string
	hasToken           bool
	isBrowser          bool
	isStatus           bool
	isHelp             bool
	isVerify           bool
	hasAIAuto          bool
	hasRestrictionTime bool
	hasRevertOriginal  bool
}

func runLogin(args []string) error {
	if len(args) == 0 {
		return runInteractiveLogin()
	}

	isHandled, err := dispatchSubcommand(args)
	if isHandled {
		return err
	}

	return parseAndDispatchLogin(args)
}

func parseAndDispatchLogin(args []string) error {
	opts, rest, err := parseLoginArgs(args)
	if err != nil {
		return err
	}

	if len(rest) > 0 {
		return apperror.NewValidationError(fmt.Sprintf("unexpected argument %q (see `gitmap login --help`)", rest[0]))
	}

	return dispatchLoginOptions(opts)
}

func dispatchSubcommand(args []string) (bool, error) {
	method, err := enums.ParseLoginAuthMethod(args[0])
	if err != nil {
		return false, nil
	}

	return dispatchMethodAction(method, args[1:])
}

func dispatchMethodAction(method enums.LoginAuthMethodType, subArgs []string) (bool, error) {
	switch method {
	case enums.LoginAuthMethodSSH:
		return true, RunLoginSSH(subArgs)
	case enums.LoginAuthMethodBrowser:
		return true, runBrowserLogin()
	case enums.LoginAuthMethodToken:
		return true, runTokenSubcommand(subArgs)
	case enums.LoginAuthMethodStatus:
		return true, runLoginStatus()
	default:
		return false, nil
	}
}

func runTokenSubcommand(args []string) error {
	token, isVerify := parseTokenSubcommandArgs(args)
	if len(token) > 0 {
		return loginWithToken(token, isVerify)
	}

	rawToken, err := readSecureToken()
	if err != nil {
		return err
	}

	return loginWithToken(rawToken, isVerify)
}

func parseTokenSubcommandArgs(args []string) (string, bool) {
	isVerify := true
	token := ""
	for _, a := range args {
		if a == constants.FlagNoVerifyArg {
			isVerify = false
			continue
		}

		if !strings.HasPrefix(a, "-") && len(token) == 0 {
			token = a
		}
	}

	return token, isVerify
}

func dispatchLoginOptions(opts loginOptions) error {
	if opts.isHelp {
		printLoginUsage()
		return nil
	}
	if isHandled, err := dispatchAIAutoOptions(opts); isHandled {
		return err
	}

	return dispatchStandardLoginOptions(opts)
}

func dispatchAIAutoOptions(opts loginOptions) (bool, error) {
	if opts.hasRevertOriginal {
		return true, runRevertToOriginal()
	}
	if opts.hasAIAuto {
		return true, runUseAIAuto(opts.hasRestrictionTime)
	}
	if opts.hasRestrictionTime {
		return true, applyRestrictionTimeDefaults()
	}

	return false, nil
}

func dispatchStandardLoginOptions(opts loginOptions) error {
	if opts.isStatus {
		return runLoginStatus()
	}

	if opts.isBrowser {
		return runBrowserLogin()
	}

	if opts.hasToken {
		return loginWithToken(opts.tokenValue, opts.isVerify)
	}

	return runInteractiveLogin()
}

func parseLoginArgs(args []string) (loginOptions, []string, error) {
	opts := loginOptions{isVerify: true}
	rest := []string{}
	for i := 0; i < len(args); i++ {
		step, isParsed, err := applyLoginArgFlag(args, i, &opts)
		if err != nil {
			return opts, rest, err
		}

		if isParsed {
			i += step
			continue
		}

		rest = append(rest, args[i])
	}

	return opts, rest, nil
}

func applyLoginArgFlag(args []string, i int, opts *loginOptions) (int, bool, error) {
	arg := args[i]
	if isSimpleFlag(arg, opts) {
		return 0, true, nil
	}

	return parseTokenFlag(args, i, opts)
}

func isSimpleFlag(arg string, opts *loginOptions) bool {
	if isAIFlag(arg, opts) {
		return true
	}

	return isStandardFlag(arg, opts)
}

func isAIFlag(arg string, opts *loginOptions) bool {
	switch arg {
	case "--use-ai-auto", "--ai-auto":
		opts.hasAIAuto = true
		return true
	case "--restriction-time", "--time-restriction":
		opts.hasRestrictionTime = true
		return true
	case "--revert-to-original", "--revert-original":
		opts.hasRevertOriginal = true
		return true
	default:
		return false
	}
}

func isStandardFlag(arg string, opts *loginOptions) bool {
	switch arg {
	case constants.FlagHelpShortArg, constants.FlagHelpArg:
		opts.isHelp = true
		return true
	case constants.FlagStatusArg:
		opts.isStatus = true
		return true
	case constants.FlagWebArg, constants.FlagBrowserArg:
		opts.isBrowser = true
		return true
	case constants.FlagNoVerifyArg:
		opts.isVerify = false
		return true
	default:
		return false
	}
}

func parseTokenFlag(args []string, i int, opts *loginOptions) (int, bool, error) {
	arg := args[i]
	if strings.HasPrefix(arg, constants.FlagTokenPrefix) {
		opts.tokenValue = strings.TrimSpace(strings.TrimPrefix(arg, constants.FlagTokenPrefix))
		opts.hasToken = true

		return 0, true, nil
	}

	if arg != constants.FlagTokenArg {
		return 0, false, nil
	}

	if i+1 >= len(args) {
		return 0, false, apperror.NewValidationError("missing value for --token (see `gitmap login --help`)")
	}

	opts.tokenValue = strings.TrimSpace(args[i+1])
	opts.hasToken = true

	return 1, true, nil
}

func printLoginUsage() {
	printUsageExamples()
	printUsageNotes()
}

func printUsageExamples() {
	fmt.Println(appconstants.ColorCyan + "Usage:" + appconstants.ColorReset)
	fmt.Println("  gitmap login                  # Interactive: choose SSH, PAT, Browser, or Status")
	fmt.Println("  gitmap login ssh [alias] [-y] # Generate or reuse SSH key and copy public key to clipboard")
	fmt.Println("  gitmap login token [PAT]      # Log in with a personal access token")
	fmt.Println("  gitmap login browser-flow     # Log in via the browser (GitHub CLI device flow)")
	fmt.Println("  gitmap login --status         # Show current login state")
	fmt.Println("  gitmap login --use-ai-auto [--restriction-time] # Switch to AI bot and apply wait defaults")
	fmt.Println("  gitmap login --revert-to-original # Revert Git identity back to original user")
	fmt.Println("  gitmap logout                 # Remove the stored GitHub credential")
	fmt.Println("")
}

func printUsageNotes() {
	fmt.Println(appconstants.ColorCyan + "Flags:" + appconstants.ColorReset)
	fmt.Println("  --use-ai-auto        Switch identity to AI automation bot and backup original user")
	fmt.Println("  --restriction-time   Auto-apply time restriction & dynamic waiting defaults (20s)")
	fmt.Println("  --revert-to-original Revert Git identity back to backed-up original user info")
	fmt.Println("")
	fmt.Println(appconstants.ColorCyan + "Notes:" + appconstants.ColorReset)
	fmt.Println("  --token validates the token against api.github.com before storing it.")
	fmt.Println("  The token is stored in git global config (github.token) so clone/pull/push reuse it.")
	fmt.Println("  Use --no-verify with --token to store without validation (offline use).")
	fmt.Println("  --use-ai-auto sets user: gitmap automation bot - runs autonmously <info@rgitmap.com>.")
}
