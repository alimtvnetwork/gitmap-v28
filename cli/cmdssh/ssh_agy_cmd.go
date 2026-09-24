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
	target, except := "", ""
	var clean []string
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isExceptFlag(arg) && i+1 < len(args) {
			except = args[i+1]
			i++
			continue
		}
		if isTargetFlag(arg) && i+1 < len(args) {
			target = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(arg, "--except=") {
			except = strings.TrimPrefix(arg, "--except=")
			continue
		}
		clean = append(clean, arg)
	}
	return resolveAgyTargetAndCommand(target, except, clean)
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
	conns, err := loadTargetNodes(target)
	if err != nil {
		return err
	}
	opts := FleetParallelOptions{
		Target:   target,
		Except:   except,
		TaskName: "AGY " + strings.Join(agyArgs, " "),
	}
	RunParallelFleetExecution(conns, opts, func(c db.SSHConnection) (string, error) {
		return executeAgyNodeWorker(c, agyArgs)
	})
	return nil
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
