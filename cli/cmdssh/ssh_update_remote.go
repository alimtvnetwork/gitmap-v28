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
	if !isNodeAvailable(c.IPAddress, header) {
		return "", fmt.Errorf("offline: node %s is unreachable", header)
	}
	client, isConnected := connectSSHClient(c, header)
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
		return resolveRemoteAgmUpdateCommand(isWin)
	default:
		return resolveRemoteGitmapUpdateCommand(isWin)
	}
}

func resolveRemoteGitmapUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (gitmap version 2>$null | Out-String).Trim(); & { $env:GITMAP_UPDATING='1'; irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.ps1 | iex } *>$null; $curr = (gitmap version 2>$null | Out-String).Trim(); if ($curr) { $msg = if ($prev -and $prev -ne $curr) { 'Upgraded: ' + $prev + ' -> ' + $curr } else { 'Version: ' + $curr }; @{ success = $true; current_version = $curr; previous_version = $prev; details = $msg } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'Update failed to verify binary' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'PREV=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); GITMAP_UPDATING=1 curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/cli/scripts/install.sh | bash >/dev/null 2>&1; CURR=$(gitmap version 2>/dev/null | grep -oE \"v[0-9]+\\.[0-9]+\\.[0-9]+\" | head -n1); if [ -n \"$CURR\" ]; then if [ -n \"$PREV\" ] && [ \"$PREV\" != \"$CURR\" ]; then DET=\"Upgraded: $PREV -> $CURR\"; else DET=\"Version: $CURR\"; fi; printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"previous_version\\\":\\\"%s\\\",\\\"details\\\":\\\"%s\\\"}\" \"$CURR\" \"$PREV\" \"$DET\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"Update failed to verify binary\\\"}\"; fi'"
}

func resolveRemoteAgmUpdateCommand(isWin bool) string {
	if isWin {
		return "powershell -NoProfile -ExecutionPolicy Bypass -Command \"& { $ErrorActionPreference='SilentlyContinue'; $WarningPreference='SilentlyContinue'; $ProgressPreference='SilentlyContinue'; $prev = (agm version 2>$null | Select-Object -First 1); & { irm https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.ps1 | iex } *>$null; $curr = (agm version 2>$null | Select-Object -First 1); if ($curr) { @{ success = $true; current_version = $curr; details = ('Version: ' + $curr) } | ConvertTo-Json -Compress } else { @{ success = $false; details = 'AGM update completed' } | ConvertTo-Json -Compress } }\""
	}
	return "sh -c 'curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash >/dev/null 2>&1; CURR=$(agm version 2>/dev/null | head -n1); if [ -n \"$CURR\" ]; then printf \"{\\\"success\\\":true,\\\"current_version\\\":\\\"%s\\\",\\\"details\\\":\\\"Version: %s\\\"}\" \"$CURR\" \"$CURR\"; else printf \"{\\\"success\\\":false,\\\"details\\\":\\\"AGM update completed\\\"}\"; fi'"
}
