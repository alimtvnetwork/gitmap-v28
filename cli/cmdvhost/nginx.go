package cmdvhost

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runNginx handles CLI routing for the 'nginx' and 'ngx' top-level commands.
func RunNginx(args []string) error {
	checkHelp(constants.CmdNginx, args)
	if len(args) == 0 {
		return runNginxStatus()
	}

	sub := strings.ToLower(args[0])

	return dispatchNginxSubcommand(sub, args[1:])
}

func dispatchNginxVHostManage(sub string, rest []string) (bool, error) {
	switch sub {
	case "create", "add", "new":
		return true, RunVHostCreate(rest)
	case "rm", "remove", "delete", "del":
		return true, RunVHostRm(rest)
	case "enable", "en":
		return true, RunVHostEnable(rest)
	case "disable", "dis":
		return true, RunVHostDisable(rest)
	default:
		return false, nil
	}
}

func dispatchNginxVHostSubcommand(sub string, rest []string) (bool, error) {
	switch sub {
	case "vhost", "vh":
		return true, RunVHost(rest)
	case "list", "ls":
		return true, runNginxList(rest)
	case "ini", "inishowcase", "showcase":
		return true, runNginxIni(rest)
	default:
		return dispatchNginxVHostManage(sub, rest)
	}
}

func dispatchNginxOpsSubcommand(sub string, rest []string) (bool, error) {
	switch sub {
	case "test", "t", "check":
		return true, RunVHostTest(rest)
	case "reload", "r", "restart":
		return true, RunVHostReload(rest)
	case "install", "in":
		return true, cmdinstall.RunInstall(append([]string{constants.ToolNginx}, rest...))
	default:
		return false, nil
	}
}

func dispatchNginxSubcommand(sub string, rest []string) error {
	if isHandled, err := dispatchNginxVHostSubcommand(sub, rest); isHandled {
		return err
	}

	if sub == "status" || sub == "st" {
		return runNginxStatus()
	}

	if isHandled, err := dispatchNginxOpsSubcommand(sub, rest); isHandled {
		return err
	}

	return unknownNginxSubcommandError(sub)
}

func unknownNginxSubcommandError(sub string) error {
	msg := fmt.Sprintf("unknown nginx subcommand %q; see 'gitmap nginx --help'", sub)

	return apperror.NewWithDetails(
		"cmd.runNginx",
		"E4001",
		msg,
		"cmd.nginx",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		nil,
	)
}

func printNginxNotFound() error {
	fmt.Println("  Nginx Binary: not found in PATH")
	fmt.Println("  Install with: gitmap install nginx")

	return nil
}

func runNginxStatus() error {
	fmt.Println("▶ gitmap nginx status")
	path, err := exec.LookPath("nginx")
	if err != nil {
		return printNginxNotFound()
	}

	printNginxDetails(path)

	return nil
}

func printNginxDetails(path string) {
	fmt.Printf("  Nginx Binary: %s\n", path)
	version := cmdinstall.GetInstalledVersion("nginx")
	if version != "" {
		fmt.Printf("  Version: %s\n", version)
	}

	testOut, testErr := executeNginxTest()
	isTestOk := testErr == nil
	fmt.Printf("  Config Test: valid=%t (%s)\n", isTestOk, strings.TrimSpace(testOut))
}

func executeNginxTest() (string, error) {
	cmd := exec.Command("nginx", "-t")
	out, err := cmd.CombinedOutput()

	return string(out), err
}
