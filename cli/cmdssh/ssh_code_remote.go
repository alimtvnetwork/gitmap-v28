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

func runRemoteCodeBinary(c db.SSHConnection, header string, codeArgs []string) {
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		appErr := apperror.NewExecutionError("ssh connection authentication failed for " + header)
		fmt.Printf("  %s %sAuth Error:%s %v\n", header, constants.ColorRed, constants.ColorReset, appErr)
		return
	}
	defer client.Close()

	executeRemoteCodeSession(client, header, c.OS, codeArgs)
}

func executeRemoteCodeSession(client *ssh.Client, header, osType string, codeArgs []string) {
	cmdStr := "code " + strings.Join(codeArgs, " ")
	out, err := crypto.RunCommand(client, cmdStr, resolveRemoteShell(osType))
	if err != nil {
		appErr := apperror.WrapSimple(err, "runRemoteCodeBinary")
		fmt.Printf("  %s %sVS Code Error:%s %v\n%s\n", header, constants.ColorRed, constants.ColorReset, appErr, out)
		return
	}

	fmt.Printf("%s\n%s\n", header, strings.TrimSpace(out))
}
