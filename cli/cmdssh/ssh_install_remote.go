package cmdssh

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runSSHInstallCLI(args []string) error {
	conns, err := loadSSHConnectionsForTarget("all")
	if err != nil {
		return err
	}

	pkg, target := parseInstallTargetAndPackage(conns, args)
	executeFleetInstall(conns, target, pkg)
	return nil
}

func executeFleetInstall(conns []db.SSHConnection, target, pkg string) {
	fmt.Printf("\n%s Installing / Updating '%s' across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, pkg, target, constants.ColorReset)

	for _, c := range filterConnectionsByTarget(conns, target) {
		installOrUpdateNodeWithPackage(c, pkg)
	}

	fmt.Printf("\nSSH Install '%s' complete.\n\n", pkg)
}

func parseInstallTargetAndPackage(conns []db.SSHConnection, args []string) (string, string) {
	if len(args) == 0 {
		return "gitmap", "all"
	}
	if len(args) == 1 {
		return parseSingleArgInstall(conns, args[0])
	}
	if args[0] == "gitmap" {
		return "gitmap", args[1]
	}

	return args[0], args[1]
}

func parseSingleArgInstall(conns []db.SSHConnection, arg string) (string, string) {
	if arg == "gitmap" {
		return "gitmap", "all"
	}
	if isKnownTarget(conns, arg) {
		return "gitmap", arg
	}

	return arg, "all"
}

func installOrUpdateNodeWithPackage(c db.SSHConnection, pkg string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !isNodeAvailable(c.IPAddress, header) {
		return
	}

	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	dispatchNodePackageAction(client, header, c.OS, pkg)
}

func dispatchNodePackageAction(client *ssh.Client, header, osType, pkg string) {
	if !isGitmapPresent(client) {
		installFreshGitmap(client, header, osType)
	}

	if pkg == "gitmap" {
		updateExistingGitmap(client, header, osType)
		return
	}

	installRemotePackage(client, header, osType, pkg)
}

func isGitmapPresent(client *ssh.Client) bool {
	out, err := crypto.RunCommand(client, "gitmap --version", "")

	return err == nil && strings.Contains(out, "gitmap")
}
