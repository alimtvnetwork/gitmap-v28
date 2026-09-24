package cmdssh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunSSHAgyCLI executes AGY commands across remote SSH nodes in parallel.
func RunSSHAgyCLI(args []string) error {
	return runSSHAgyCLI(args)
}

func runSSHAgyCLI(args []string) error {
	if len(args) == 0 || isHelpArgToken(args[0]) {
		showSSHAgyUsage()
		return nil
	}
	target, except, agyArgs := parseAgyFleetArgs(args)
	return executeAgyOnFleet(target, except, agyArgs)
}

func isHelpArgToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "-h" || low == "--help" || low == "help"
}

func showSSHAgyUsage() {
	fmt.Printf("\n%s╔════════════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s║   Bidirectional SSH & AGY Fleet Orchestrator (gitmap ssh agy / agy ssh)   ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("%s╚════════════════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Execute ANY Antigravity (AGY) or GitMap command across all remote SSH machines in parallel.")
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap ssh agy [flags] <agy-or-gitmap-command...>")
	fmt.Println("    gitmap agy ssh [flags] <agy-or-gitmap-command...>")
	fmt.Println()
	fmt.Println("  Fleet Filter Flags:")
	fmt.Println("    -e, --except, --excep, --accept <id,ip,alias>   Exclude machines by Worker ID (1, worker-1), IP, or Alias")
	fmt.Println("    -t, --target <node>                             Run only on a specific target node or group")
	fmt.Println("    -n, --dry-run                                   Simulate fleet dispatch without SSH execution")
	fmt.Println()
	fmt.Println("  Supported AGY & GitMap Commands Across Fleet:")
	fmt.Println("    • rerun [1|2|3|4|all|queue]   Restart IDE, re-inject active prompt + images + 5 queued items with prefix")
	fmt.Println("    • rop [N]                     Re-read, optimize, and repair N Antigravity projects with temp backup")
	fmt.Println("    • status / ping / ls / stats  Check Antigravity IDE status, health, and project tables on all nodes")
	fmt.Println("    • prompt ls / list-prompts    List all prompt templates and active conversation prompts")
	fmt.Println("    • cache-clear (cc) / ccko     Clean Antigravity runtime cache across all remote machines")
	fmt.Println("    • fix-pipeline (fp)           Trigger autonomous CI/CD pipeline repair across remote nodes")
	fmt.Println()
	fmt.Println("  Examples:")
	fmt.Println("    gitmap ssh agy status --except 1,worker-2")
	fmt.Println("    gitmap agy ssh rerun all --except 192.168.1.50,dev-laptop")
	fmt.Println("    gitmap ssh agy rerun queue --except 2")
	fmt.Println("    gitmap agy ssh rop 4 --except alpha-win")
	fmt.Println()
}

func parseAgyFleetArgs(args []string) (string, string, []string) {
	target := ""
	var exceptParts, clean []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		low := strings.ToLower(arg)
		if isExceptOrExcepFlag(low) && i+1 < len(args) {
			exceptParts = append(exceptParts, args[i+1])
			i++
			continue
		}
		if isTargetFlag(low) && i+1 < len(args) {
			target = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--excep=") || strings.HasPrefix(low, "--exclude=") {
			idx := strings.IndexByte(arg, '=')
			exceptParts = append(exceptParts, arg[idx+1:])
			continue
		}
		clean = append(clean, arg)
	}
	return resolveAgyTargetAndCommand(target, strings.Join(exceptParts, ","), clean)
}

func resolveAgyTargetAndCommand(target, except string, clean []string) (string, string, []string) {
	if len(clean) == 0 {
		return "all", except, []string{"status"}
	}
	if target != "" {
		return target, except, clean
	}

	return inferAgyTargetAndCommand(except, clean)
}

func inferAgyTargetAndCommand(except string, clean []string) (string, string, []string) {
	if isAllTarget(clean[0]) {
		return "all", except, clean[1:]
	}

	conns, err := fetchAllSSHConnections()
	if err == nil && len(filterConnectionsByTarget(conns, clean[0])) > 0 {
		return clean[0], except, clean[1:]
	}

	return "all", except, clean
}

func executeAgyOnFleet(target, except string, agyArgs []string) error {
	isDryRun := hasDryRunToken(agyArgs)
	conns, err := loadTargetNodes(target)
	if err != nil {
		conns, _ = fetchAllSSHConnections()
	}
	filtered := FilterSSHConnectionsByExcept(conns, except)
	if isDryRun || len(filtered) == 0 {
		fmt.Printf("✓ Dispatched 'gitmap agy %s' across %d SSH node(s) [target=%s, except=%q]\n",
			strings.Join(agyArgs, " "), len(filtered), target, except)
		return nil
	}
	opts := FleetParallelOptions{
		Target:   target,
		Except:   except,
		TaskName: "AGY " + strings.Join(agyArgs, " "),
	}
	RunParallelFleetExecution(filtered, opts, func(c db.SSHConnection) (string, error) {
		return executeAgyNodeWorker(c, agyArgs)
	})
	return nil
}

func hasDryRunToken(args []string) bool {
	for _, a := range args {
		if a == "--dry-run" || a == "-n" {
			return true
		}
	}
	return false
}

func executeAgyNodeWorker(c db.SSHConnection, agyArgs []string) (string, error) {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	if !checkRemoteNodeOnline(c.IPAddress, header) {
		return "", fmt.Errorf("node %s is offline", header)
	}
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return "", fmt.Errorf("connection/auth failed for %s", header)
	}
	defer client.Close()

	cmdStr := "gitmap agy " + strings.Join(agyArgs, " ")
	return crypto.RunCommand(client, cmdStr, "")
}
