package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHAgyCLI(args []string) error {
	if len(args) == 0 {
		showSSHAgyUsage()
		return nil
	}

	target, agyArgs := parseRemoteTargetAndArgs(args)
	conns, err := loadSSHConnectionsForTarget(target)
	if err != nil {
		return err
	}

	if len(conns) == 0 {
		return apperror.NewNotFoundError("no target SSH machines found for target: " + target)
	}

	for _, c := range conns {
		runAgyOnNode(c, agyArgs)
	}

	return nil
}

func showSSHAgyUsage() {
	fmt.Printf("\n%s Usage:%s gitmap ssh agy <target> <agy-command...>\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ssh agy w1 open /projects/my-app")
	fmt.Println("    gitmap ssh agy all status")
	fmt.Println("    gitmap ssh agy 192.168.1.8 ls\n")
}

func parseRemoteTargetAndArgs(args []string) (string, []string) {
	if len(args) == 1 {
		return "all", args
	}

	return args[0], args[1:]
}

func runAgyOnNode(c db.SSHConnection, agyArgs []string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)

	isOnline, reason := CheckConnLiveness(nil, c.IPAddress, 22, 0)
	if !isOnline {
		appErr := apperror.NewExecutionError(fmt.Sprintf("node %s is unreachable: %s", header, reason))
		fmt.Printf("  %s %sOffline:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		appErr := apperror.NewExecutionError("ssh connection authentication failed for " + header)
		fmt.Printf("  %s %sAuth Error:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return
	}
	defer client.Close()

	cmdStr := "agy " + strings.Join(agyArgs, " ")
	shell := "bash"
	if isWindowsOS(c.OS) {
		shell = "ps"
	}

	out, err := crypto.RunCommand(client, cmdStr, shell)
	if err != nil {
		appErr := apperror.WrapSimple(err, "runAgyOnNode")
		fmt.Printf("  %s %sAGY Error:%s %v\n%s\n", header, constants.ColorRed, constants.ColorReset, appErr, out)
		return
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))
}
