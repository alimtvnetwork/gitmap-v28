package cmdnodes

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// NodesRemoteOptions configures fleet-wide remote command execution.
type NodesRemoteOptions struct {
	FilterOpts NodeFilterOptions
	Command    string
	IsJSON     bool
}

// NodeRemoteExecResult represents outcome of executing a command on a remote node.
type NodeRemoteExecResult struct {
	NodeAlias    string `json:"nodeAlias"`
	Host         string `json:"host"`
	OS           string `json:"os"`
	Command      string `json:"command"`
	IsSuccess    bool   `json:"isSuccess"`
	ExitCode     int    `json:"exitCode"`
	Stdout       string `json:"stdout,omitempty"`
	Stderr       string `json:"stderr,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	DurationMs   int64  `json:"durationMs"`
}

// RunNodesRemoteCLI dispatches a command concurrently across matching fleet SSH nodes.
func RunNodesRemoteCLI(args []string) error {
	for len(args) > 0 && isNodesRunToken(args[0]) {
		args = args[1:]
	}
	opts := parseNodesRemoteOptions(args)
	hasCommand := len(opts.Command) > 0
	if !hasCommand {
		return apperror.NewSimple("usage: gitmap nodes run [target] \"<command>\" [--json] [--except <list>]", "E5001")
	}

	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return apperror.WrapSimple(err, "fetch ssh connections:")
	}

	targets := FilterFleetNodes(conns, opts.FilterOpts)
	if len(targets) == 0 {
		fmt.Println("No matching fleet nodes found.")
		return nil
	}

	results := dispatchNodesRemoteConcurrent(targets, opts.Command)
	if opts.IsJSON {
		return printRemoteResultsJSON(results)
	}

	renderRemoteResultsTable(results)
	return nil
}

func isNodesRunToken(arg string) bool {
	low := strings.ToLower(arg)
	return low == "run" || low == "exec" || low == "r"
}

func parseNodesRemoteOptions(args []string) NodesRemoteOptions {
	opts := NodesRemoteOptions{
		FilterOpts: ParseNodeFilterOptions(args),
	}
	var positionals []string
	for i := 0; i < len(args); i++ {
		step := parseSingleRemoteFlag(&opts, args, i)
		if step > 0 {
			i += (step - 1)
			continue
		}
		if !strings.HasPrefix(args[i], "-") {
			positionals = append(positionals, args[i])
		}
	}
	hasFlagTarget := hasExplicitTargetFlag(args)
	assignRemotePositionals(&opts, positionals, hasFlagTarget)
	return opts
}

func hasExplicitTargetFlag(args []string) bool {
	for i, a := range args {
		if strings.HasPrefix(a, "--target=") || strings.HasPrefix(a, "-t=") {
			return true
		}
		if slices.Contains([]string{"--target", "-t"}, a) && i+1 < len(args) {
			return true
		}
	}
	return false
}

func parseSingleRemoteFlag(opts *NodesRemoteOptions, args []string, i int) int {
	arg := args[i]
	if slices.Contains([]string{"--json", "-j"}, arg) {
		opts.IsJSON = true
		return 1
	}
	return 0
}

func assignRemotePositionals(opts *NodesRemoteOptions, positionals []string, hasFlagTarget bool) {
	if len(positionals) == 0 {
		return
	}
	if hasFlagTarget {
		opts.Command = strings.Join(positionals, " ")
		return
	}
	if len(positionals) == 1 {
		opts.Command = positionals[0]
		opts.FilterOpts.Target = ""
		return
	}
	opts.FilterOpts.Target = positionals[0]
	opts.Command = strings.Join(positionals[1:], " ")
}

func resolveRemoteExecShell(osType string) string {
	isWin := strings.EqualFold(osType, "windows")
	return resolveNodeShell(isWin)
}

func dispatchNodesRemoteConcurrent(targets []db.SSHConnection, command string) []NodeRemoteExecResult {
	results := make([]NodeRemoteExecResult, len(targets))
	var wg sync.WaitGroup
	for i, conn := range targets {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = executeSingleRemoteNode(c, command)
		}(i, conn)
	}
	wg.Wait()
	return results
}

func executeSingleRemoteNode(conn db.SSHConnection, command string) NodeRemoteExecResult {
	shell := resolveRemoteExecShell(conn.OS)
	start := time.Now()
	out, errExec := currentSSHExecutor.Execute(conn, command, shell)
	duration := time.Since(start).Milliseconds()

	res := NodeRemoteExecResult{
		NodeAlias:  conn.Alias,
		Host:       conn.IPAddress,
		OS:         conn.OS,
		Command:    command,
		DurationMs: duration,
	}
	populateRemoteExecutionOutcome(&res, out, errExec)
	return res
}

func populateRemoteExecutionOutcome(res *NodeRemoteExecResult, out string, err error) {
	if err != nil {
		res.ErrorMessage = err.Error()
		res.Stderr = out
		res.ExitCode = 1
		res.IsSuccess = false
		return
	}
	res.Stdout = out
	res.ExitCode = 0
	res.IsSuccess = true
}

func printRemoteResultsJSON(results []NodeRemoteExecResult) error {
	envelope := map[string]any{
		"timestamp":  time.Now().Format(time.RFC3339),
		"totalNodes": len(results),
		"results":    results,
	}
	data, err := json.MarshalIndent(envelope, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal json results:")
	}
	fmt.Println(string(data))
	return nil
}

func renderRemoteResultsTable(results []NodeRemoteExecResult) {
	fmt.Printf("\n%-16s %-16s %-8s %-8s %s\n", "NODE", "HOST", "STATUS", "LATENCY", "OUTPUT")
	fmt.Println(strings.Repeat("-", 70))
	for _, r := range results {
		status := "FAIL"
		output := strings.TrimSpace(r.ErrorMessage)
		if r.IsSuccess {
			status = "OK"
			output = strings.TrimSpace(r.Stdout)
		}
		if len(output) > 30 {
			output = output[:27] + "..."
		}
		fmt.Printf("%-16s %-16s %-8s %-8dms %s\n", r.NodeAlias, r.Host, status, r.DurationMs, output)
	}
	fmt.Println()
}
