package cmdnodes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func isNodesTasksToken(tok string) bool {
	low := strings.ToLower(tok)
	return low == "tasks" || low == "task" || low == "task-queue" || low == "tasks-queue" || low == "queue" || low == "jobs" || low == "delegation"
}

// RunNodesTasksCLI executes fleet tasks queue inspection across local and remote nodes.
func RunNodesTasksCLI(args []string) error {
	for len(args) > 0 && (strings.EqualFold(args[0], "nodes") || isNodesTasksToken(args[0])) {
		args = args[1:]
	}
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

	if opts.isHeavyOnly {
		payload = filterPayloadHeavyOnly(payload)
	}

	if opts.isJSON {
		return printFleetTasksJSON(payload)
	}

	return renderFleetTasksTable(os.Stdout, payload)
}

func filterPayloadHeavyOnly(payload FleetTasksPayload) FleetTasksPayload {
	var filteredReports []NodeTaskReport
	var totalActive, totalPending, totalHeavy int
	for _, r := range payload.NodeReports {
		r.ActiveTasks = filterHeavyTasks(r.ActiveTasks)
		r.PendingTasks = filterHeavyTasks(r.PendingTasks)
		r.ActiveTaskCount = len(r.ActiveTasks)
		r.PendingTaskCount = len(r.PendingTasks)
		totalActive += r.ActiveTaskCount
		totalPending += r.PendingTaskCount
		totalHeavy += r.HeavyTaskCount
		filteredReports = append(filteredReports, r)
	}
	payload.NodeReports = filteredReports
	payload.TotalActive = totalActive
	payload.TotalPending = totalPending
	payload.TotalHeavyDuty = totalHeavy
	return payload
}

func filterHeavyTasks(tasks []TaskItemSummary) []TaskItemSummary {
	var res []TaskItemSummary
	for _, t := range tasks {
		if t.IsHeavyDuty {
			res = append(res, t)
		}
	}
	return res
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

	conns, errConns := cmdssh.FetchAllSSHConnections()
	connMap := make(map[string]db.SSHConnection)
	if errConns == nil {
		connMap = indexSSHConnections(conns)
	}

	var mu sync.Mutex
	var wg sync.WaitGroup
	sem := make(chan struct{}, 8)

	for _, n := range unifiedNodes {
		isLocalNode := strings.EqualFold(n.Host, "127.0.0.1") || strings.EqualFold(n.Host, constants.Localhost) || strings.EqualFold(n.Alias, "local")
		if isLocalNode {
			continue
		}
		if opts.isOpenOnly && !n.IsOnline {
			continue
		}
		if !isTargetMatching(n.Alias, n.Host, opts.targetFilter) {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(node UnifiedFleetNode) {
			defer func() {
				<-sem
				wg.Done()
			}()

			conn := resolveNodeSSHConnection(node, connMap)
			rep := queryRemoteNodeTasksAndHardware(ctx, node, conn, opts.limit)
			if opts.isOpenOnly && rep.Status == "UNREACHABLE" {
				return
			}
			mu.Lock()
			reports = append(reports, rep)
			mu.Unlock()
		}(n)
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
		NodeAlias:       alias,
		NodeHost:        host,
		NodePort:        0,
		IsLocal:         true,
		Status:          "ONLINE",
		RolePlaceholder: "?",
		ActiveTasks:     make([]TaskItemSummary, 0),
		PendingTasks:    make([]TaskItemSummary, 0),
	}

	hw, errHW := cmdos.GetHardwareMetrics()
	if errHW == nil {
		report.CPUPercent = hw.CPU.LoadPercent
		report.RAMPercent = hw.RAM.UsedPercent
		report.TotalRAMMB = hw.RAM.TotalMB
		report.UsedRAMMB = hw.RAM.UsedMB
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

	report.ActiveTasks = limitTaskSummaries(report.ActiveTasks, limit)
	report.PendingTasks = limitTaskSummaries(report.PendingTasks, limit)

	return report
}

func limitTaskSummaries(tasks []TaskItemSummary, limit int) []TaskItemSummary {
	hasLimit := limit > 0
	hasExcess := len(tasks) > limit
	if hasLimit && hasExcess {
		return tasks[:limit]
	}
	return tasks
}
