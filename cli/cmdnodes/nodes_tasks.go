package cmdnodes

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// RunNodesTasksCLI executes fleet tasks queue inspection across local and remote nodes.
func RunNodesTasksCLI(args []string) error {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			return printNodesTasksHelp()
		}
	}
	opts := parseNodesTasksOptions(args)

	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
	defer cancel()

	payload, err := executeFleetTasksInterrogation(ctx, opts)
	if err != nil {
		return err
	}

	if opts.isJSON {
		return printFleetTasksJSON(payload)
	}

	return renderFleetTasksTable(os.Stdout, payload)
}

func executeFleetTasksInterrogation(ctx context.Context, opts nodesTasksOptions) (FleetTasksPayload, error) {
	var reports []NodeTaskReport

	shouldIncludeLocal := isTargetMatching("local", "127.0.0.1", opts.targetFilter)
	if shouldIncludeLocal {
		localReport := collectLocalNodeTasks(opts.limit)
		reports = append(reports, localReport)
	}

	isLocalFilterOnly := strings.EqualFold(opts.targetFilter, "local") || strings.EqualFold(opts.targetFilter, "127.0.0.1")
	if isLocalFilterOnly {
		return buildFleetTasksPayload(reports), nil
	}

	unifiedNodes, errCollect := collectUnifiedNodes(ctx, true)
	hasCollectErr := errCollect != nil
	if hasCollectErr {
		return buildFleetTasksPayload(reports), nil
	}

	conns, _ := cmdssh.FetchAllSSHConnections()
	connMap := make(map[string]db.SSHConnection, len(conns)*2)
	for _, c := range conns {
		if c.IPAddress != "" {
			connMap[strings.ToLower(c.IPAddress)] = c
		}
		if c.Alias != "" {
			connMap[strings.ToLower(c.Alias)] = c
		}
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for _, n := range unifiedNodes {
		isLocalNode := strings.EqualFold(n.Host, "127.0.0.1") || strings.EqualFold(n.Host, "localhost") || strings.EqualFold(n.Alias, "local")
		if isLocalNode {
			continue
		}
		if isTargetMatching(n.Alias, n.Host, opts.targetFilter) {
			wg.Add(1)
			sem <- struct{}{}
			go func(node UnifiedFleetNode) {
				defer func() {
					<-sem
					wg.Done()
				}()

				conn, _ := connMap[strings.ToLower(node.Host)]
				isMissingConn := len(conn.IPAddress) == 0
				if isMissingConn {
					conn, _ = connMap[strings.ToLower(node.Alias)]
					isMissingConn = len(conn.IPAddress) == 0
					if isMissingConn {
						conn = db.SSHConnection{
							Alias:     node.Alias,
							IPAddress: node.Host,
							Username:  node.User,
						}
					}
				}

				rep := queryRemoteNodeTasks(ctx, node, conn, opts.limit)
				mu.Lock()
				reports = append(reports, rep)
				mu.Unlock()
			}(n)
		}
	}

	wg.Wait()
	return buildFleetTasksPayload(reports), nil
}

func buildFleetTasksPayload(reports []NodeTaskReport) FleetTasksPayload {
	var totalActive, totalPending, totalHeavy int
	for _, r := range reports {
		totalActive += r.ActiveTaskCount
		totalPending += r.PendingTaskCount
		totalHeavy += r.HeavyTaskCount
	}
	return FleetTasksPayload{
		Timestamp:      time.Now().Format(time.RFC3339),
		TotalNodes:     len(reports),
		TotalActive:    totalActive,
		TotalPending:   totalPending,
		TotalHeavyDuty: totalHeavy,
		NodeReports:    reports,
	}
}

func collectLocalNodeTasks(limit int) NodeTaskReport {
	host := "127.0.0.1"
	alias := "local (me)"
	if h, err := os.Hostname(); err == nil && h != "" {
		alias = fmt.Sprintf("local (%s)", h)
	}

	report := NodeTaskReport{
		NodeAlias:    alias,
		NodeHost:     host,
		NodePort:     0,
		IsLocal:      true,
		Status:       "ONLINE",
		ActiveTasks:  make([]TaskItemSummary, 0),
		PendingTasks: make([]TaskItemSummary, 0),
	}

	dbPaths := discoverAgentTaskDBs()
	for _, dbPath := range dbPaths {
		queryTasksFromDB(dbPath, &report)
	}

	report.ActiveTaskCount = len(report.ActiveTasks)
	report.PendingTaskCount = len(report.PendingTasks)

	var heavyCount int
	for _, t := range report.ActiveTasks {
		if t.IsHeavyDuty {
			heavyCount++
		}
	}
	for _, t := range report.PendingTasks {
		if t.IsHeavyDuty {
			heavyCount++
		}
	}
	report.HeavyTaskCount = heavyCount

	if limit > 0 {
		if len(report.ActiveTasks) > limit {
			report.ActiveTasks = report.ActiveTasks[:limit]
		}
		if len(report.PendingTasks) > limit {
			report.PendingTasks = report.PendingTasks[:limit]
		}
	}

	return report
}

func queryRemoteNodeTasks(ctx context.Context, node UnifiedFleetNode, conn db.SSHConnection, limit int) NodeTaskReport {
	alias := node.Alias
	isDashAlias := alias == "" || alias == "-"
	if isDashAlias {
		alias = node.Host
	}
	port := node.Port
	isZeroPort := port <= 0
	if isZeroPort {
		port = 22
	}

	report := NodeTaskReport{
		NodeAlias:    alias,
		NodeHost:     node.Host,
		NodePort:     port,
		IsLocal:      false,
		Status:       "ONLINE",
		ActiveTasks:  make([]TaskItemSummary, 0),
		PendingTasks: make([]TaskItemSummary, 0),
	}

	remoteCmd := fmt.Sprintf("gitmap tasks --json --limit %d", limit)
	shell := "sh"
	isWin := isWindowsNode(conn) || strings.Contains(strings.ToLower(node.Role), "win")
	if isWin {
		shell = "powershell"
	}

	out, errExec := currentSSHExecutor.Execute(conn, remoteCmd, shell)
	hasExecErr := errExec != nil
	if hasExecErr {
		report.Status = "UNREACHABLE"
		report.ErrorMessage = errExec.Error()
		return report
	}

	var parsedReport NodeTaskReport
	errJSON := json.Unmarshal([]byte(strings.TrimSpace(out)), &parsedReport)
	hasValidJSON := errJSON == nil
	if hasValidJSON {
		parsedReport.NodeAlias = alias
		parsedReport.NodeHost = node.Host
		parsedReport.NodePort = port
		parsedReport.IsLocal = false
		return parsedReport
	}

	var parsedPayload FleetTasksPayload
	errPayload := json.Unmarshal([]byte(strings.TrimSpace(out)), &parsedPayload)
	hasPayload := errPayload == nil && len(parsedPayload.NodeReports) > 0
	if hasPayload {
		first := parsedPayload.NodeReports[0]
		first.NodeAlias = alias
		first.NodeHost = node.Host
		first.NodePort = port
		first.IsLocal = false
		return first
	}

	report.Status = "ACTIVE"
	return report
}


