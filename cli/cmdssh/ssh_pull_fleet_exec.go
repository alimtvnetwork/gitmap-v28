package cmdssh

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
)

// FleetNodeRemediationResult encapsulates the execution outcome of a command dispatched to a node.
type FleetNodeRemediationResult struct {
	NodeName  string
	IP        string
	IsLocal   bool
	Command   string
	Output    string
	IsSuccess bool
	ErrorMsg  string
}

// ExecuteCommandOnFleetNode runs a command on a local or remote node and returns structured output.
func ExecuteCommandOnFleetNode(c db.SSHConnection, isLocal bool, cmdStr string) FleetNodeRemediationResult {
	if isLocal {
		return executeCommandOnLocalVM(cmdStr)
	}

	return executeCommandOnRemoteNode(c, cmdStr)
}

func executeCommandOnLocalVM(cmdStr string) FleetNodeRemediationResult {
	bin := resolveGitmapExecutable()
	args := strings.Fields(cmdStr)
	if len(args) > 0 && args[0] == "gitmap" {
		args = args[1:]
	}
	cmd := exec.Command(bin, args...)
	out, err := cmd.CombinedOutput()
	outStr := strings.TrimSpace(string(out))
	if err != nil {
		return buildRemediationResult("Local VM", "127.0.0.1", true, cmdStr, outStr, false, err.Error())
	}

	return buildRemediationResult("Local VM", "127.0.0.1", true, cmdStr, outStr, true, "")
}

func executeCommandOnRemoteNode(c db.SSHConnection, cmdStr string) FleetNodeRemediationResult {
	header := fmt.Sprintf("[%s|%s]", c.Alias, c.IPAddress)
	client, isConnected := connectSSHClient(c, header)
	if !isConnected {
		return buildRemediationResult(c.Alias, c.IPAddress, false, cmdStr, "", false, "SSH connection failed")
	}
	defer client.Close()

	osType := resolveTargetNodeOSFast(client, c)
	shellType := resolveRemoteShell(osType)
	out, err := secrets.RunCommand(client, cmdStr, shellType)
	outStr := strings.TrimSpace(out)
	if err != nil {
		return buildRemediationResult(c.Alias, c.IPAddress, false, cmdStr, outStr, false, err.Error())
	}

	return buildRemediationResult(c.Alias, c.IPAddress, false, cmdStr, outStr, true, "")
}

func buildRemediationResult(name, ip string, isLocal bool, cmd, out string, isSuccess bool, errMsg string) FleetNodeRemediationResult {
	return FleetNodeRemediationResult{
		NodeName:  name,
		IP:        ip,
		IsLocal:   isLocal,
		Command:   cmd,
		Output:    out,
		IsSuccess: isSuccess,
		ErrorMsg:  errMsg,
	}
}

// DispatchCommandToAffectedNodes runs a command across all nodes identified in issues list.
func DispatchCommandToAffectedNodes(cmdStr string, issues []FleetNodeIssue, conns []db.SSHConnection) []FleetNodeRemediationResult {
	var results []FleetNodeRemediationResult
	for _, issue := range issues {
		conn, hasConn := resolveConnectionForNode(issue.NodeName, issue.IP, conns)
		if !issue.IsLocal && !hasConn {
			results = append(results, buildRemediationResult(issue.NodeName, issue.IP, false, cmdStr, "", false, "Node connection details not found"))
			continue
		}
		res := ExecuteCommandOnFleetNode(conn, issue.IsLocal, cmdStr)
		results = append(results, res)
	}

	return results
}

// ResolveConnectionForNode matches node name or IP against known SSH connections.
func resolveConnectionForNode(nodeName, ip string, conns []db.SSHConnection) (db.SSHConnection, bool) {
	for _, c := range conns {
		if strings.EqualFold(c.Alias, nodeName) || strings.EqualFold(c.IPAddress, ip) {
			return c, true
		}
	}

	return db.SSHConnection{}, false
}

// RenderNodeExecResults prints clear execution feedback for each node result.
func RenderNodeExecResults(results []FleetNodeRemediationResult, actionLabel string) {
	fmt.Printf("\n  %s🔧 Execution Summary for '%s':%s\n", constants.ColorCyan, actionLabel, constants.ColorReset)
	for _, r := range results {
		renderSingleRemediationResult(r)
	}
	fmt.Println()
}

func renderSingleRemediationResult(r FleetNodeRemediationResult) {
	status := resolveRemediationResultStatus(r)
	fmt.Printf("    • [%s] (%s): %s\n", r.NodeName, r.IP, status)
	if len(r.Output) > 0 {
		printTruncatedRemediationOutput(r.Output)
	}
}

func resolveRemediationResultStatus(r FleetNodeRemediationResult) string {
	if r.IsSuccess {
		return fmt.Sprintf("%s✔ Success%s", constants.ColorGreen, constants.ColorReset)
	}

	return fmt.Sprintf("%s✖ Failed: %s%s", constants.ColorRed, r.ErrorMsg, constants.ColorReset)
}

func printTruncatedRemediationOutput(out string) {
	lines := strings.Split(out, "\n")
	limit := 6
	if len(lines) < limit {
		limit = len(lines)
	}
	for i := 0; i < limit; i++ {
		fmt.Printf("        %s%s%s\n", constants.ColorDim, strings.TrimRight(lines[i], "\r"), constants.ColorReset)
	}
	if len(lines) > limit {
		fmt.Printf("        %s... (%d more lines)%s\n", constants.ColorDim, len(lines)-limit, constants.ColorReset)
	}
}

func resolveGitmapExecutable() string {
	self, err := os.Executable()
	if err == nil && self != "" {
		return self
	}

	return "gitmap"
}
