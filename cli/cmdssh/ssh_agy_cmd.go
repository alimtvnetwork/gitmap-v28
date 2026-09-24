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
	if len(args) == 0 {
		showSSHAgyUsage()
		return nil
	}
	target, except, agyArgs := parseAgyFleetArgs(args)
	return executeAgyOnFleet(target, except, agyArgs)
}

func showSSHAgyUsage() {
	fmt.Printf("\n%s Usage:%s gitmap agy ssh [target] [flags] <agy-command...>\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Flags: -e, --except <nodes>   Exclude machines by alias or IP")
	fmt.Println("  Examples:")
	fmt.Println("    gitmap agy ssh status")
	fmt.Println("    gitmap agy ssh -e worker-2 prompt ls")
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
