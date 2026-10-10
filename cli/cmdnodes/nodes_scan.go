// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// NodeScanResult stores telemetry for a scan operation on an individual node.
type NodeScanResult struct {
	NodeAlias string        `json:"node_alias"`
	Host      string        `json:"host"`
	OS        string        `json:"os"`
	Status    string        `json:"status"`
	Latency   time.Duration `json:"latency"`
	LatencyMs int64         `json:"latency_ms"`
	Output    string        `json:"output,omitempty"`
	Error     string        `json:"error,omitempty"`
}

// RunNodesScan connects and runs gitmap scan across filtered fleet nodes via SSH.
func RunNodesScan(args []string) error {
	return executeFleetScanDispatch("scan", args)
}

// RunNodesRescan connects and runs gitmap rescan across filtered fleet nodes via SSH.
func RunNodesRescan(args []string) error {
	return executeFleetScanDispatch("rescan", args)
}

func executeFleetScanDispatch(subCmd string, args []string) error {
	if isScanHelp(args) {
		return printNodesScanHelp(subCmd)
	}

	opts := ParseNodeFilterOptions(args)
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return apperror.WrapSimple(err, "fetch fleet ssh connections")
	}

	if len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found", "E9070")
	}

	targets := FilterFleetNodes(conns, opts)
	if len(targets) == 0 {
		return apperror.NewSimple("no candidate fleet nodes found matching filter criteria", "E9071")
	}

	passArgs := extractScanPassArgs(args, opts.Target)
	results := runFleetScanParallel(targets, subCmd, passArgs)

	if isJSONRequested(args) {
		return emitNodeScanJSON(results)
	}

	renderNodeScanTable(results, subCmd)
	printNodeScanSummary(results, subCmd)
	return nil
}

func runFleetScanParallel(targets []db.SSHConnection, subCmd string, passArgs []string) []NodeScanResult {
	results := make([]NodeScanResult, len(targets))
	var wg sync.WaitGroup

	for i, c := range targets {
		wg.Add(1)
		go func(idx int, conn db.SSHConnection) {
			defer wg.Done()
			results[idx] = executeSingleNodeScan(conn, subCmd, passArgs)
		}(i, c)
	}

	wg.Wait()
	return results
}

func executeSingleNodeScan(conn db.SSHConnection, subCmd string, passArgs []string) NodeScanResult {
	start := time.Now()
	res := NodeScanResult{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OS:        conn.OS,
	}

	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return handleNodeScanConnectError(res, errConnect, start)
	}

	defer client.Close()

	out, errExec := executeRemoteScanCommand(client, conn, subCmd, passArgs)
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	res.Output = strings.TrimSpace(out)

	if errExec != nil {
		res.Status = "FAILED"
		res.Error = errExec.Error()
		return res
	}

	res.Status = "SUCCESS"
	return res
}

func handleNodeScanConnectError(res NodeScanResult, err error, start time.Time) NodeScanResult {
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	res.Error = err.Error()

	if isOfflineError(err.Error()) {
		res.Status = "OFFLINE"
		return res
	}

	res.Status = "FAILED"
	return res
}

func executeRemoteScanCommand(client *ssh.Client, conn db.SSHConnection, subCmd string, passArgs []string) (string, error) {
	isWin := isWindowsNode(conn)
	cmd, shell := resolveRemoteScanCommand(subCmd, passArgs, isWin)
	out, err := secrets.RunCommand(client, cmd, shell)

	if err != nil && !isWin && isBashMissingError(out, err) {
		out, err = secrets.RunCommand(client, cmd, "sh")
	}

	return out, err
}

func resolveRemoteScanCommand(subCmd string, passArgs []string, isWin bool) (string, string) {
	fullCmd := "gitmap " + subCmd
	if len(passArgs) > 0 {
		fullCmd += " " + strings.Join(passArgs, " ")
	}

	if isWin {
		return `powershell.exe -NoProfile -Command "` + fullCmd + `"`, "ps"
	}

	return fullCmd, "bash"
}

func extractScanPassArgs(args []string, target string) []string {
	var pass []string

	for i := 0; i < len(args); i++ {
		step := isFilterFlagStep(args, i)
		if step > 0 {
			i += (step - 1)
			continue
		}

		a := args[i]
		if isScanIgnoredArg(a, target) {
			continue
		}

		pass = append(pass, a)
	}

	return pass
}

func isScanIgnoredArg(arg, target string) bool {
	if arg == "--json" || arg == "-j" || arg == "--help" || arg == "-h" || arg == "help" {
		return true
	}

	if target != "" && arg == target {
		return true
	}

	return false
}

func isFilterFlagStep(args []string, i int) int {
	a := args[i]
	if a == "--include-main" || a == "--open-only" {
		return 1
	}

	if isPrefixFilterFlag(a) {
		return 1
	}

	if isTwoTokenFilterFlag(a) && i+1 < len(args) {
		return 2
	}

	return 0
}

func isPrefixFilterFlag(a string) bool {
	if strings.HasPrefix(a, "--target=") || strings.HasPrefix(a, "-t=") {
		return true
	}

	if strings.HasPrefix(a, "--except=") || strings.HasPrefix(a, "-e=") || strings.HasPrefix(a, "--exclude=") {
		return true
	}

	return strings.HasPrefix(a, "--include=") || strings.HasPrefix(a, "--accept=")
}

func isTwoTokenFilterFlag(a string) bool {
	switch a {
	case "--target", "-t", "--except", "-e", "--exclude", "--include", "--accept":
		return true
	default:
		return false
	}
}

func isScanHelp(args []string) bool {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return true
		}
	}

	return false
}

func printNodesScanHelp(subCmd string) error {
	fmt.Printf("\n  gitmap nodes %s - Run gitmap %s across fleet nodes via SSH\n\n", subCmd, subCmd)
	fmt.Println("  Usage:")
	fmt.Printf("    gitmap nodes %s [flags] [target]\n\n", subCmd)
	fmt.Println("  Flags:")
	fmt.Println("    -t, --target <alias>     Filter to a single target node alias or IP")
	fmt.Println("    -e, --except <list>      Exclude specific nodes (comma-separated)")
	fmt.Println("        --include <list>     Whitelist specific nodes (comma-separated)")
	fmt.Println("        --include-main       Include the 'main' node (excluded by default)")
	fmt.Println("        --open-only          Only run on currently reachable/online nodes")
	fmt.Println("    -j, --json               Output structured JSON telemetry")
	fmt.Println()
	return nil
}

func emitNodeScanJSON(results []NodeScanResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(results); err != nil {
		return apperror.WrapSimple(err, "encode node scan json")
	}

	return nil
}

func formatScanStatusBadge(status string) string {
	switch status {
	case "SUCCESS":
		return constants.ColorGreen + "● SUCCESS" + constants.ColorReset
	case "OFFLINE":
		return constants.ColorDim + "○ OFFLINE" + constants.ColorReset
	default:
		return constants.ColorRed + "▲ " + status + constants.ColorReset
	}
}

func renderNodeScanTable(results []NodeScanResult, action string) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-14s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST", "OS", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 76))

	for _, r := range results {
		statusStr := formatScanStatusBadge(r.Status)
		osStr := r.OS
		if osStr == "" {
			osStr = "linux"
		}

		fmt.Printf("  %-18s %-16s %-10s %-14s %-10v\n",
			r.NodeAlias, r.Host, osStr, statusStr, r.Latency.Round(time.Millisecond))
	}

	fmt.Println()
}

func printNodeScanSummary(results []NodeScanResult, action string) {
	successCount := 0
	offlineCount := 0
	failedCount := 0

	for _, r := range results {
		if r.Status == "SUCCESS" {
			successCount++
			continue
		}

		if r.Status == "OFFLINE" {
			offlineCount++
			continue
		}

		failedCount++
	}

	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	if failedCount > 0 {
		prefix = constants.ColorYellow + "▲" + constants.ColorReset
	}

	fmt.Printf("  %s Completed %s across %d/%d node(s) (%d offline, %d failed)\n\n",
		prefix, action, successCount, len(results), offlineCount, failedCount)
}
