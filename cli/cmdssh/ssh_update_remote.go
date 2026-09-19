package cmdssh

import (
	"fmt"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHUpdateCLI(args []string) error {
	conns, err := loadSSHConnectionsForTarget("all")
	if err != nil {
		return apperror.WrapSimple(err, "loadSSHConnectionsForTarget")
	}

	pkg, target := parseInstallTargetAndPackage(conns, args)
	executeFleetUpdate(conns, target, pkg)
	return nil
}

func executeFleetUpdate(conns []db.SSHConnection, target, pkg string) {
	fmt.Printf("\n%s Updating '%s' across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, pkg, target, constants.ColorReset)

	for _, c := range filterConnectionsByTarget(conns, target) {
		updateSingleSSHNode(c, pkg)
	}

	fmt.Printf("\nSSH Fleet Update '%s' complete.\n\n", pkg)
}

func updateSingleSSHNode(c db.SSHConnection, pkg string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	executeRemoteUpdate(client, header, c.OS, pkg)
}

func connectSSHNode(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if !isNodeAvailable(c.IPAddress, header) {
		return nil, false
	}
	return connectSSHClient(c, header)
}

func executeRemoteUpdate(client *ssh.Client, header, osType, pkg string) {
	cmd := "gitmap update"
	if pkg != "gitmap" {
		cmd = "gitmap update " + pkg
	}

	out, err := crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Updated", out, err)
}
