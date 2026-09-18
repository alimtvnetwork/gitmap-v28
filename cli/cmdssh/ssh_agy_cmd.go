package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

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
	return executeAgyOnFleet(target, agyArgs)
}

func executeAgyOnFleet(target string, agyArgs []string) error {
	conns, err := loadTargetNodes(target)
	if err != nil {
		return err
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
	fmt.Println("    gitmap ssh agy 192.168.1.8 ls")
	fmt.Println()
}

func parseRemoteTargetAndArgs(args []string) (string, []string) {
	if len(args) == 1 {
		return "all", args
	}

	return args[0], args[1:]
}

func runAgyOnNode(c db.SSHConnection, agyArgs []string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := establishAgyClient(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	executeAgyRemoteCommand(client, header, c.OS, agyArgs)
}

func establishAgyClient(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return nil, false
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		appErr := apperror.NewExecutionError("ssh connection authentication failed for " + header)
		printAppErrorWithStack(header, "Auth Error", appErr)
		return nil, false
	}

	return client, true
}

func executeAgyRemoteCommand(client *ssh.Client, header, osType string, agyArgs []string) {
	cmdStr := "agy " + strings.Join(agyArgs, " ")
	out, err := crypto.RunCommand(client, cmdStr, resolveRemoteShell(osType))
	if err != nil {
		appErr := apperror.WrapSimple(err, "runAgyOnNode")
		printAppErrorWithStack(header, "AGY Error", appErr)
		return
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))
}
