package cmdssh

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHCodeCLI(args []string) error {
	if len(args) == 0 {
		showSSHCodeUsage()
		return nil
	}

	target, codeArgs := parseRemoteTargetAndArgs(args)
	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return err
	}

	if len(conns) == 0 {
		return apperror.NewNotFoundError("no target SSH machines found for target: " + target)
	}

	for _, c := range conns {
		runCodeOnNode(c, codeArgs)
	}

	return nil
}

func showSSHCodeUsage() {
	fmt.Printf("\n%s Usage:%s gitmap ssh code <target> [open] <remote-path>\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ssh code w1 /var/www/my-project")
	fmt.Println("    gitmap ssh code 192.168.1.8 open /home/user/workspace\n")
}

func runCodeOnNode(c db.SSHConnection, codeArgs []string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)

	isOnline, reason := CheckConnLiveness(nil, c.IPAddress, 22, 0)
	if !isOnline {
		appErr := apperror.NewExecutionError(fmt.Sprintf("node %s is unreachable: %s", header, reason))
		fmt.Printf("  %s %sOffline:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return
	}

	pathArg := "."
	if len(codeArgs) > 0 {
		pathArg = codeArgs[len(codeArgs)-1]
	}

	// Try launching local VS Code SSH remote if available
	remoteURI := fmt.Sprintf("ssh-remote+%s@%s", c.Username, c.IPAddress)
	cmd := exec.Command("code", "--remote", remoteURI, pathArg)
	if err := cmd.Start(); err == nil {
		fmt.Printf("  %s %sOpening remote VS Code window for %s:%s%s\n",
			header, constants.ColorGreen, remoteURI, pathArg, constants.ColorReset)
		return
	}

	runRemoteCodeBinary(c, header, codeArgs)
}

func runRemoteCodeBinary(c db.SSHConnection, header string, codeArgs []string) {
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		appErr := apperror.NewExecutionError("ssh connection authentication failed for " + header)
		fmt.Printf("  %s %sAuth Error:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return
	}
	defer client.Close()

	cmdStr := "code " + strings.Join(codeArgs, " ")
	shell := "bash"
	if isWindowsOS(c.OS) {
		shell = "ps"
	}

	out, err := crypto.RunCommand(client, cmdStr, shell)
	if err != nil {
		appErr := apperror.WrapSimple(err, "runRemoteCodeBinary")
		fmt.Printf("  %s %sVS Code Error:%s %v\n%s\n", header, constants.ColorRed, constants.ColorReset, appErr, out)
		return
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))
}
