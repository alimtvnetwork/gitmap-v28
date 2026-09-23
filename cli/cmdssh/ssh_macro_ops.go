package cmdssh

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func runSSHMacroAddCLI(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printSSHMacroAddHelp()
		return nil
	}

	err := cmdmacro.RunMacroCmd(append([]string{"add"}, args...))
	if err != nil {
		return err
	}

	macroName := args[0]
	fmt.Printf("  %s[tip]%s Macro %q created locally.\n", constants.ColorCyan, constants.ColorReset, macroName)
	fmt.Printf("  To deploy to remote nodes, run: %sgitmap ssh macro deploy %s%s\n\n",
		constants.ColorGreen, macroName, constants.ColorReset)
	return nil
}

func printSSHMacroAddHelp() {
	fmt.Println("Create a new macro locally for SSH fleet orchestration.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro add <name> <steps...>")
	fmt.Println("\nExample:")
	fmt.Println("  gitmap ssh macro add ping-nodes \"gitmap status\" \"gitmap --version\"")
}

func runSSHMacroRmCLI(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("Usage: gitmap ssh macro rm <name>")
		return nil
	}
	return cmdmacro.RunMacroCmd(append([]string{"delete"}, args...))
}

func runSSHMacroLsCLI(args []string) error {
	return cmdmacro.RunMacroCmd(append([]string{"list"}, args...))
}

func runSSHMacroEditCLI(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		fmt.Println("Usage: gitmap ssh macro edit <name>")
		return nil
	}
	return cmdmacro.RunMacroCmd(append([]string{"edit"}, args...))
}

func extractMacroRemoteTarget(args []string) (string, bool) {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if isRemoteFlag(a) && i+1 < len(args) {
			return args[i+1], true
		}
		if strings.HasPrefix(a, "--target=") {
			return strings.TrimPrefix(a, "--target="), true
		}
		if strings.HasPrefix(a, "--remote=") {
			return strings.TrimPrefix(a, "--remote="), true
		}
		if a == "--deploy" || a == "--remote" {
			return "all-nodes", true
		}
	}
	return "", false
}

func isRemoteFlag(a string) bool {
	return a == "-t" || a == "--target" || a == "-r" || a == "--remote"
}

func runSSHMacroRunCLI(args []string) error {
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		printSSHMacroRunHelp()
		return nil
	}

	target, hasRemote := extractMacroRemoteTarget(args)
	if !hasRemote {
		return cmdmacro.RunMacroCmd(append([]string{"run"}, args...))
	}

	macroName := extractFirstNonFlag(args)
	if macroName == "" {
		return apperror.NewValidationError("macro name required to run")
	}

	return runRemoteMacroWorkflow(macroName, target, args)
}

func extractFirstNonFlag(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}

func printSSHMacroRunHelp() {
	fmt.Println("Execute a recorded macro locally or across remote SSH fleet nodes.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro run <name> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: local execution)")
	fmt.Println("      --deploy            Deploy macro JSON to remote nodes before execution")
	fmt.Println("      --dry-run           Preview macro execution without executing")
}

func runRemoteMacroWorkflow(macroName, target string, args []string) error {
	deployErr := runSSHMacroDeployCLI([]string{macroName, "--target", target, "--force"})
	if deployErr != nil {
		fmt.Fprintf(os.Stderr, "Notice: macro deploy warning: %v\n", deployErr)
	}

	execCmd := fmt.Sprintf("gitmap macro run %s", macroName)
	execArgs := []string{"--target", target, execCmd}
	return runSSHExec(execArgs)
}

func runSSHMacroDeployCLI(args []string) error {
	macroName := extractFirstNonFlag(args)
	if macroName == "" || macroName == "help" {
		printSSHMacroDeployHelp()
		return nil
	}

	localMacro, err := macro.LoadMacro(macroName)
	if err != nil || localMacro == nil {
		return apperror.NewValidationError(fmt.Sprintf(
			"macro %q not found locally. Create it first with: gitmap ssh macro add %s <steps...>",
			macroName, macroName,
		))
	}

	opts := parseMacroDeployOptions(args)
	return executeMacroFleetDeploy(*localMacro, opts)
}

func printSSHMacroDeployHelp() {
	fmt.Println("Deploy a verified local macro to remote SSH nodes via JSON reimport.")
	fmt.Println("\nUsage:")
	fmt.Println("  gitmap ssh macro deploy <name> [flags]")
	fmt.Println("\nFlags:")
	fmt.Println("  -t, --target string     Target machine alias or IP (default: all-nodes)")
	fmt.Println("      --except string     Exclude machines by alias, IP, or ID")
	fmt.Println("  -f, --force             Overwrite existing macro on remote nodes")
}

func parseMacroDeployOptions(args []string) SSHMacroSyncOptions {
	opts := SSHMacroSyncOptions{
		Target:  "all-nodes",
		IsForce: true,
	}

	for i := 0; i < len(args); i++ {
		a := args[i]
		if isRemoteFlag(a) && i+1 < len(args) {
			opts.Target = args[i+1]
			i++
			continue
		}
		if isExceptFlag(a) && i+1 < len(args) {
			opts.Except = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(a, "--except=") {
			opts.Except = strings.TrimPrefix(a, "--except=")
			continue
		}
		if a == "-f" || a == "--force" {
			opts.IsForce = true
		}
	}
	return opts
}

func executeMacroFleetDeploy(m macro.Macro, opts SSHMacroSyncOptions) error {
	conns, err := fetchAllSSHConnections()
	if err != nil {
		return apperror.WrapSimple(err, "fetchAllSSHConnections")
	}

	filtered := filterConnectionsByTarget(conns, opts.Target)
	if opts.Except != "" {
		filtered = filterSSHConns(filtered, opts.Except)
	}

	if len(filtered) == 0 {
		fmt.Printf("No matching target machines found to deploy macro %q.\n", m.Name)
		return nil
	}

	fmt.Printf("\n%s Deploying macro %q (%d step(s)) across SSH fleet (%s):%s\n\n",
		constants.ColorCyan, m.Name, len(m.Steps), opts.Target, constants.ColorReset)

	successCount := 0
	for _, c := range filtered {
		if deployMacroToSingleNode(c, m, opts.IsForce) {
			successCount++
		}
	}

	fmt.Printf("\nMacro %q deployed to %d/%d nodes.\n\n", m.Name, successCount, len(filtered))
	return nil
}

func deployMacroToSingleNode(c db.SSHConnection, m macro.Macro, isForce bool) bool {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		return false
	}
	defer client.Close()

	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		osType = probed
	}

	shell := resolveRemoteShell(osType)
	isOk := syncSingleMacroToClient(client, c, m, shell, isForce)
	if !isOk {
		fmt.Printf("  %s %s✖ Failed to deploy macro %q%s\n",
			header, constants.ColorRed, m.Name, constants.ColorReset)
		return false
	}

	fmt.Printf("  %s %s✔ Deployed and reimported macro %q successfully%s\n",
		header, constants.ColorGreen, m.Name, constants.ColorReset)
	return true
}
