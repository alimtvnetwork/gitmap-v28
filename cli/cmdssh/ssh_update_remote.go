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

func parseUpdateOptions(args []string) SSHFleetUpdateOptions {
	opts := SSHFleetUpdateOptions{
		Target: "",
		Pkg:    "gitmap",
	}
	var clean []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isTargetFlag(arg) && i+1 < len(args) {
			opts.Target = args[i+1]
			i++
			continue
		}
		if isExceptFlag(arg) {
			var exceptTokens []string
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				exceptTokens = append(exceptTokens, args[i+1])
				i++
			}
			opts.Except = strings.Join(exceptTokens, ",")
			continue
		}
		if strings.HasPrefix(arg, "--node=") {
			opts.Target = strings.TrimPrefix(arg, "--node=")
			continue
		}
		if strings.HasPrefix(arg, "--target=") {
			opts.Target = strings.TrimPrefix(arg, "--target=")
			continue
		}
		if strings.HasPrefix(arg, "--remote=") {
			opts.Target = strings.TrimPrefix(arg, "--remote=")
			continue
		}
		if strings.HasPrefix(arg, "--except=") {
			opts.Except = strings.TrimPrefix(arg, "--except=")
			continue
		}
		if strings.HasPrefix(arg, "--excep=") {
			opts.Except = strings.TrimPrefix(arg, "--excep=")
			continue
		}
		if strings.HasPrefix(arg, "--exclude=") {
			opts.Except = strings.TrimPrefix(arg, "--exclude=")
			continue
		}
		if arg == "--dry-run" {
			opts.IsDryRun = true
			continue
		}
		if arg == "-f" || arg == "--force" {
			opts.IsForce = true
			continue
		}
		clean = append(clean, arg)
	}

	return populateTargetAndPackage(opts, clean)
}

func isTargetFlag(arg string) bool {
	return arg == "--node" || arg == "--remote" || arg == "-n" || arg == "-r" || arg == "-t" || arg == "--target"
}

func isExceptFlag(arg string) bool {
	return arg == "--except" || arg == "--excep" || arg == "--exclude" || arg == "-e"
}

func populateTargetAndPackage(opts SSHFleetUpdateOptions, clean []string) SSHFleetUpdateOptions {
	for _, token := range clean {
		if isAllTarget(token) {
			opts.Target = "all-nodes"
			continue
		}
		if isUpdateAction(token) {
			continue
		}
		if isKnownPackage(token) {
			opts.Pkg = token
			continue
		}
		if opts.Target == "" {
			opts.Target = token
		}
	}
	if opts.Target == "" {
		opts.Target = "all-nodes"
	}
	return opts
}

func isUpdateAction(token string) bool {
	low := strings.ToLower(token)
	return low == "update" || low == "up" || low == "update-all" || low == "updateall"
}

func isKnownPackage(token string) bool {
	low := strings.ToLower(token)
	return low == "gitmap" || low == "agm" || low == "ag-manager" || low == "antigravity-manager"
}

func runSSHUpdateCLI(args []string) error {
	conns, err := fetchAllSSHConnections()
	if err != nil {
		return apperror.WrapSimple(err, "fetchAllSSHConnections")
	}

	opts := parseUpdateOptions(args)
	executeFleetUpdateWithOptions(conns, opts)
	return nil
}

func executeFleetUpdateWithOptions(conns []db.SSHConnection, opts SSHFleetUpdateOptions) {
	fmt.Printf("\n%s Updating '%s' across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, opts.Pkg, opts.Target, constants.ColorReset)

	fleetOpts := FleetParallelOptions{
		Target:   opts.Target,
		Except:   opts.Except,
		TaskName: fmt.Sprintf("Update %s", opts.Pkg),
		IsDryRun: opts.IsDryRun,
	}

	RunParallelFleetExecution(conns, fleetOpts, func(c db.SSHConnection) (string, error) {
		return executeSingleSSHNodeUpdate(c, opts.Pkg, opts.IsDryRun)
	})
}

func executeSingleSSHNodeUpdate(c db.SSHConnection, pkg string, isDryRun bool) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if isDryRun {
		return fmt.Sprintf("[DRY-RUN] Would update %s on %s", pkg, c.Alias), nil
	}
	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		return "", fmt.Errorf("unable to connect to %s", header)
	}
	defer client.Close()

	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		osType = probed
	}

	cmd := resolveRemoteUpdateCommand(osType, pkg)
	out, err := crypto.RunCommand(client, cmd, resolveRemoteShell(osType))
	reportRemoteExecution(header, "Updated", out, err)
	return out, err
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
