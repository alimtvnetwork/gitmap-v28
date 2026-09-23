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

// RunSSHUpdateCLI executes remote package updates across SSH targets.
func RunSSHUpdateCLI(args []string) error {
	return runSSHUpdateCLI(args)
}

func extractNodeAndPackageFlags(args []string) (string, []string) {
	target := ""
	var clean []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "--node" || arg == "--remote" || arg == "-n" || arg == "-r") && i+1 < len(args) {
			target = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(arg, "--node=") {
			target = strings.TrimPrefix(arg, "--node=")
			continue
		}
		if strings.HasPrefix(arg, "--remote=") {
			target = strings.TrimPrefix(arg, "--remote=")
			continue
		}
		clean = append(clean, arg)
	}
	return target, clean
}

func resolveTargetAndPkg(conns []db.SSHConnection, args []string) (string, string) {
	flagTarget, clean := extractNodeAndPackageFlags(args)
	if flagTarget != "" {
		pkg := "gitmap"
		if len(clean) > 0 {
			pkg = clean[0]
		}
		return pkg, flagTarget
	}
	return parseInstallTargetAndPackage(conns, args)
}

func runSSHUpdateCLI(args []string) error {
	conns, err := loadSSHConnectionsForTarget("all")
	if err != nil {
		return apperror.WrapSimple(err, "loadSSHConnectionsForTarget")
	}

	pkg, target := resolveTargetAndPkg(conns, args)
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

	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		osType = probed
	}

	executeRemoteUpdate(client, header, osType, pkg)
}

func connectSSHNode(c db.SSHConnection, header string) (*ssh.Client, bool) {
	if !isNodeAvailable(c.IPAddress, header) {
		return nil, false
	}
	return connectSSHClient(c, header)
}

func resolveRemoteUpdateCommand(osType, pkg string) string {
	isWin := isWindowsOS(osType)
	switch strings.ToLower(pkg) {
	case "agm", "ag-manager", "antigravity-manager":
		if isWin {
			return constants.AgManagerWindowsInstallCmd
		}
		return constants.AgManagerUnixInstallCmd
	default:
		if isWin {
			return "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.ps1 | iex"
		}
		return "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.sh | bash"
	}
}

func executeRemoteUpdate(client *ssh.Client, header, osType, pkg string) {
	cmd := resolveRemoteUpdateCommand(osType, pkg)
	out, err := crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Updated", out, err)
}
