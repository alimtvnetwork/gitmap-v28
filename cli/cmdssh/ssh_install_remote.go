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

func runSSHInstallCLI(args []string) error {
	conns, err := loadSSHConnectionsForTarget("all")
	if err != nil {
		return apperror.WrapSimple(err, "loadSSHConnectionsForTarget")
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
	return parseMultiArgsInstall(conns, args)
}

func parseMultiArgsInstall(conns []db.SSHConnection, args []string) (string, string) {
	if len(args) == 2 {
		return parseTwoArgsInstall(conns, args[0], args[1])
	}
	return parseExcessArgsInstall(conns, args)
}

func parseExcessArgsInstall(conns []db.SSHConnection, args []string) (string, string) {
	for _, arg := range args {
		if isAllFlag(arg) {
			return args[0], "all"
		}
	}
	return args[0], args[1]
}

func parseTwoArgsInstall(conns []db.SSHConnection, first, second string) (string, string) {
	if isAllFlag(first) {
		return second, "all"
	}
	if isAllFlag(second) {
		return first, "all"
	}
	if first == "gitmap" {
		return "gitmap", second
	}
	if isKnownTarget(conns, first) {
		return second, first
	}
	return first, second
}

func isAllFlag(arg string) bool {
	return arg == "--all" || arg == "-a" || arg == "all"
}

func parseSingleArgInstall(conns []db.SSHConnection, arg string) (string, string) {
	if isAllFlag(arg) || arg == "gitmap" {
		return "gitmap", "all"
	}
	if isKnownTarget(conns, arg) {
		return "gitmap", arg
	}
	return arg, "all"
}

func installOrUpdateNodeWithPackage(c db.SSHConnection, pkg string) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		return
	}
	defer client.Close()

	dispatchNodePackageAction(client, header, c.OS, pkg)
}

func dispatchNodePackageAction(client *ssh.Client, header, osType, pkg string) {
	isMissing := !isGitmapPresent(client)
	if isMissing {
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
