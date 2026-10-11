package cmdpending

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func runFleetSSHPendingCommits(opts PendingCommitsOptions) error {
	conns, errFetch := cmdssh.FetchAllSSHConnections()
	if errFetch != nil {
		return apperror.WrapSimple(errFetch, "fetch fleet ssh connections")
	}

	if len(conns) == 0 {
		return handleEmptyFleet(opts.IsJSON)
	}

	results := make([]NodePendingCommitsRecord, len(conns))
	var wg sync.WaitGroup

	for i, conn := range conns {
		wg.Add(1)
		go func(idx int, c db.SSHConnection) {
			defer wg.Done()
			results[idx] = querySingleNodePendingCommits(c, opts)
		}(i, conn)
	}

	wg.Wait()

	onlineCount := 0
	for _, r := range results {
		if r.IsOnline {
			onlineCount++
		}
	}

	fleetPayload := NodesPendingCommitsPayload{
		Timestamp:    time.Now().UTC(),
		TotalNodes:   len(conns),
		OnlineNodes:  onlineCount,
		FleetResults: results,
	}

	if opts.IsJSON {
		return emitPendingCommitsFleetJSON(fleetPayload)
	}

	renderPendingCommitsFleetTerminal(fleetPayload)
	return nil
}

func handleEmptyFleet(isJSON bool) error {
	if isJSON {
		return emitPendingCommitsFleetJSON(NodesPendingCommitsPayload{
			Timestamp:    time.Now().UTC(),
			TotalNodes:   0,
			OnlineNodes:  0,
			FleetResults: []NodePendingCommitsRecord{},
		})
	}
	fmt.Printf("  %sNo registered SSH cluster nodes found.%s\n", constants.ColorYellow, constants.ColorReset)
	return nil
}

func querySingleNodePendingCommits(conn db.SSHConnection, opts PendingCommitsOptions) NodePendingCommitsRecord {
	start := time.Now()
	rec := NodePendingCommitsRecord{
		NodeAlias: conn.Alias,
		Host:      conn.IPAddress,
		OSType:    conn.OS,
	}

	isWin := strings.EqualFold(conn.OS, "windows")
	cmdStr := "gitmap pending-commits --json"
	if opts.SortMode != "" && opts.SortMode != "priority" {
		cmdStr += " --sort=" + opts.SortMode
	}
	if opts.DetailMode == "all" {
		cmdStr += " --detail=all"
	}
	if opts.IsAll {
		cmdStr += " --all"
	}

	shell := "bash"
	if isWin {
		cmdStr = fmt.Sprintf(`powershell.exe -NoProfile -Command "%s"`, cmdStr)
		shell = "ps"
	}

	out, errExec := currentPendingCommitsSSHExecutor.Execute(conn, cmdStr, shell)
	rec.LatencyMs = time.Since(start).Milliseconds()

	if errExec != nil {
		return handleSSHExecFailure(conn, errExec, rec.LatencyMs)
	}

	rec.IsOnline = true
	rec.IsSuccess = true

	cleanJSON := extractJSONPayload(out)
	var nodePayload PendingCommitsPayload
	if errUnmarshal := json.Unmarshal([]byte(cleanJSON), &nodePayload); errUnmarshal == nil {
		rec.Payload = &nodePayload
	}

	return rec
}

func handleSSHExecFailure(conn db.SSHConnection, errExec error, latencyMs int64) NodePendingCommitsRecord {
	isOffline := isOfflineErrorString(errExec.Error())
	return NodePendingCommitsRecord{
		NodeAlias:    conn.Alias,
		Host:         conn.IPAddress,
		OSType:       conn.OS,
		IsOnline:     !isOffline,
		IsSuccess:    false,
		LatencyMs:    latencyMs,
		ErrorMessage: errExec.Error(),
	}
}

func isOfflineErrorString(errStr string) bool {
	low := strings.ToLower(errStr)
	return strings.Contains(low, "network unreachable") ||
		strings.Contains(low, "connection refused") ||
		strings.Contains(low, "timeout") ||
		strings.Contains(low, "i/o timeout") ||
		strings.Contains(low, "offline")
}

func extractJSONPayload(s string) string {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start >= 0 && end > start {
		return s[start : end+1]
	}
	return s
}

func emitPendingCommitsFleetJSON(payload NodesPendingCommitsPayload) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal fleet pending commits json")
	}
	fmt.Println(string(b))
	return nil
}

func renderPendingCommitsFleetTerminal(fleet NodesPendingCommitsPayload) {
	border := constants.ColorCyan
	reset := constants.ColorReset
	bold := constants.ColorBold

	fmt.Println()
	fmt.Printf("  %s┌──────────────────────────────────────────────────────────────────────────────┐%s\n", border, reset)
	fmt.Printf("  %s│%s%s               GITMAP FLEET PENDING COMMITS (SSH CLUSTER)                     %s%s│%s\n", border, reset, bold, reset, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s Total Nodes: %-6d │ Online Nodes: %-6d                                   %s│%s\n",
		border, reset, fleet.TotalNodes, fleet.OnlineNodes, border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)
	fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %7s %s│%s\n",
		border, reset,
		"NODE", "HOST", "STATUS", "DIRTY", "UNCOMMIT", "UNPUSHED", "LATENCY",
		border, reset)
	fmt.Printf("  %s├──────────────────────────────────────────────────────────────────────────────┤%s\n", border, reset)

	for _, node := range fleet.FleetResults {
		statusLabel, dirtyStr, uncommitStr, unpushedStr := formatFleetNodeRow(node, reset)

		fmt.Printf("  %s│%s %-16s %-15s %-7s %-7s %-8s %-10s %5dms %s│%s\n",
			border, reset,
			node.NodeAlias,
			node.Host,
			statusLabel,
			dirtyStr,
			uncommitStr,
			unpushedStr,
			node.LatencyMs,
			border, reset)
	}

	fmt.Printf("  %s└──────────────────────────────────────────────────────────────────────────────┘%s\n", border, reset)
}

func formatFleetNodeRow(node NodePendingCommitsRecord, reset string) (string, string, string, string) {
	if !node.IsOnline {
		return constants.ColorRed + "OFFLINE" + reset, "-", "-", "-"
	}

	if node.IsSuccess && node.Payload != nil {
		statusLabel := constants.ColorGreen + "ONLINE " + reset
		dirtyStr := fmt.Sprintf("%d", node.Payload.TotalDirtyRepos)
		uncommitStr := fmt.Sprintf("%d", node.Payload.TotalUncommittedFiles)
		unpushedStr := fmt.Sprintf("%d", node.Payload.TotalUnpushedCommits)
		return statusLabel, dirtyStr, uncommitStr, unpushedStr
	}

	return constants.ColorYellow + "ERROR  " + reset, "-", "-", "-"
}
