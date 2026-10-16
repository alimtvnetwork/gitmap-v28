package cmdnodes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	_ "modernc.org/sqlite"
)

// TaskItemSummary represents an enqueued or running agent subtask.
type TaskItemSummary struct {
	TaskID       int64  `json:"taskId"`
	TaskSlug     string `json:"taskSlug"`
	TaskCode     string `json:"taskCode"`
	Title        string `json:"title"`
	AssignedRole string `json:"assignedRole"`
	Status       string `json:"status"`
	IsHeavyDuty  bool   `json:"isHeavyDuty"`
	CurrentStep  int    `json:"currentStep"`
	TotalSteps   int    `json:"totalSteps"`
	CreatedAt    string `json:"createdAt"`
}

// NodeTaskReport summarizes active tasks and queue depth for a single fleet node.
type NodeTaskReport struct {
	NodeAlias        string            `json:"nodeAlias"`
	NodeHost         string            `json:"nodeHost"`
	NodePort         int               `json:"nodePort"`
	IsLocal          bool              `json:"isLocal"`
	Status           string            `json:"status"`
	ActiveTaskCount  int               `json:"activeTaskCount"`
	PendingTaskCount int               `json:"pendingTaskCount"`
	HeavyTaskCount   int               `json:"heavyTaskCount"`
	ActiveTasks      []TaskItemSummary `json:"activeTasks"`
	PendingTasks     []TaskItemSummary `json:"pendingTasks"`
	ErrorMessage     string            `json:"errorMessage,omitempty"`
}

// FleetTasksPayload encapsulates the complete response across all cluster nodes.
type FleetTasksPayload struct {
	Timestamp      string           `json:"timestamp"`
	TotalNodes     int              `json:"totalNodes"`
	TotalActive    int              `json:"totalActive"`
	TotalPending   int              `json:"totalPending"`
	TotalHeavyDuty int              `json:"totalHeavyDuty"`
	NodeReports    []NodeTaskReport `json:"nodeReports"`
}

type nodesTasksOptions struct {
	isJSON       bool
	limit        int
	targetFilter string
}

func parseNodesTasksOptions(args []string) nodesTasksOptions {
	opts := nodesTasksOptions{
		limit: 10,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" || a == "-j" {
			opts.isJSON = true
		} else if (a == "--limit" || a == "-l") && i+1 < len(args) {
			i++
			if n, err := strconv.Atoi(args[i]); err == nil && n > 0 {
				opts.limit = n
			}
		} else if isTaskFilterCandidate(a) && opts.targetFilter == "" {
			opts.targetFilter = a
		}
	}
	return opts
}

func isTaskFilterCandidate(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	tok := strings.ToLower(arg)
	isSubcmd := tok == "tasks" || tok == "task" || tok == "task-queue" || tok == "tasks-queue" || tok == "queue" || tok == "jobs" || tok == "delegation"
	if isSubcmd {
		return false
	}
	return true
}

func detectHeavyDutyTask(totalStepsBudget int, title, taskCode string) bool {
	isBudgetHeavy := totalStepsBudget >= 100
	if isBudgetHeavy {
		return true
	}
	text := strings.ToLower(title + " " + taskCode)
	keywords := []string{"heavy", "deep-scan", "migration", "full-audit", "rebuild-all"}
	for _, kw := range keywords {
		hasKw := strings.Contains(text, kw)
		if hasKw {
			return true
		}
	}
	return false
}

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

func discoverAgentTaskDBs() []string {
	var paths []string
	root := filepath.Join(".ai-memory", "temp-agents")
	entries, err := os.ReadDir(root)
	hasErr := err != nil
	if hasErr {
		return paths
	}
	for _, e := range entries {
		isDir := e.IsDir()
		if isDir {
			dbFile := filepath.Join(root, e.Name(), "agent-task.db")
			if _, errStat := os.Stat(dbFile); errStat == nil {
				paths = append(paths, dbFile)
			}
		}
	}
	directDB := filepath.Join(root, "agent-task.db")
	if _, errStat := os.Stat(directDB); errStat == nil {
		paths = append(paths, directDB)
	}
	return paths
}

func queryTasksFromDB(dbPath string, report *NodeTaskReport) {
	dsn := fmt.Sprintf("file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", filepath.ToSlash(dbPath))
	dbConn, err := sql.Open("sqlite", dsn)
	hasErr := err != nil
	if hasErr {
		return
	}
	defer dbConn.Close()

	query := `
		SELECT s.SubtaskId, p.TaskSlug, s.TaskCode, s.Title, s.AssignedAgentRole, s.Status,
		       p.TotalStepsBudget, p.CurrentStep, s.CreatedAt
		FROM Subtask s
		JOIN ParentTask p ON s.ParentTaskId = p.ParentTaskId
		WHERE s.Status IN ('PENDING', 'CLAIMED', 'ACTIVE', 'IN_PROGRESS')
		ORDER BY s.SubtaskId ASC
	`
	rows, errQuery := dbConn.Query(query)
	hasQueryErr := errQuery != nil
	if hasQueryErr {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			subtaskId        int64
			taskSlug         string
			taskCode         string
			title            string
			assignedRole     sql.NullString
			status           string
			totalStepsBudget int
			currentStep      int
			createdAt        string
		)
		if errScan := rows.Scan(&subtaskId, &taskSlug, &taskCode, &title, &assignedRole, &status, &totalStepsBudget, &currentStep, &createdAt); errScan != nil {
			continue
		}

		isHeavy := detectHeavyDutyTask(totalStepsBudget, title, taskCode)
		role := assignedRole.String
		hasNoRole := len(role) == 0
		if hasNoRole {
			role = "unassigned"
		}

		item := TaskItemSummary{
			TaskID:       subtaskId,
			TaskSlug:     taskSlug,
			TaskCode:     taskCode,
			Title:        title,
			AssignedRole: role,
			Status:       status,
			IsHeavyDuty:  isHeavy,
			CurrentStep:  currentStep,
			TotalSteps:   totalStepsBudget,
			CreatedAt:    createdAt,
		}

		isActive := status == "CLAIMED" || status == "ACTIVE" || status == "IN_PROGRESS"
		if isActive {
			report.ActiveTasks = append(report.ActiveTasks, item)
		} else {
			report.PendingTasks = append(report.PendingTasks, item)
		}
	}
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

func printFleetTasksJSON(payload FleetTasksPayload) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "printFleetTasksJSON")
	}
	fmt.Println(string(data))
	return nil
}

func renderFleetTasksTable(out io.Writer, payload FleetTasksPayload) error {
	border := strings.Repeat("═", 116)
	div := resStyleDim.Render(strings.Repeat("─", 118))
	hdr := resStyleHeader.Render(fmt.Sprintf("%s %s %s %s %s %s %s",
		padCell("NODE (ALIAS)", 18),
		padCell("HOST", 18),
		padCell("STATUS", 14),
		padCell("ACTIVE", 10),
		padCell("PENDING", 10),
		padCell("HEAVY", 10),
		padCell("QUEUE OVERVIEW", 30)))

	fmt.Fprintf(out, "\n  %s╔%s╗%s\n", constants.ColorDim, border, constants.ColorReset)
	fmt.Fprintf(out, "  %s║%s%s%s║%s\n", constants.ColorDim, constants.ColorReset, resStyleTitle.Render(padCell(" GITMAP FLEET NODES TASK QUEUE & DELEGATION", 116)), constants.ColorDim, constants.ColorReset)
	fmt.Fprintf(out, "  %s╚%s╝%s\n", constants.ColorDim, border, constants.ColorReset)
	fmt.Fprintf(out, "  Fleet Tasks: %d active, %d pending, %d heavy-duty across %d node(s)\n\n  %s\n  %s\n",
		payload.TotalActive, payload.TotalPending, payload.TotalHeavyDuty, payload.TotalNodes,
		hdr, div)

	for _, r := range payload.NodeReports {
		statusBadge := resStyleIdle.Render("● " + r.Status)
		isUnreachable := r.Status == "UNREACHABLE"
		if isUnreachable {
			statusBadge = resStyleHigh.Render("○ " + r.Status)
		}

		activeStr := strconv.Itoa(r.ActiveTaskCount)
		pendingStr := strconv.Itoa(r.PendingTaskCount)
		heavyStr := strconv.Itoa(r.HeavyTaskCount)

		overview := fmt.Sprintf("%d active / %d queued", r.ActiveTaskCount, r.PendingTaskCount)
		hasHeavy := r.HeavyTaskCount > 0
		if hasHeavy {
			overview += fmt.Sprintf(" (%d heavy)", r.HeavyTaskCount)
		}
		if isUnreachable && r.ErrorMessage != "" {
			overview = resStyleHigh.Render("error: " + r.ErrorMessage)
		}

		fmt.Fprintf(out, "  %s %s %s %s %s %s %s\n",
			padCell(r.NodeAlias, 18),
			resStyleHeader.Render(padCell(r.NodeHost, 18)),
			padCell(statusBadge, 14),
			padCell(activeStr, 10),
			padCell(pendingStr, 10),
			padCell(heavyStr, 10),
			padCell(overview, 30))

		hasSubtasks := len(r.ActiveTasks) > 0 || len(r.PendingTasks) > 0
		if hasSubtasks {
			for _, t := range r.ActiveTasks {
				heavyBadge := ""
				if t.IsHeavyDuty {
					heavyBadge = " " + resStyleExceed.Render("[HEAVY]")
				}
				fmt.Fprintf(out, "    %s %s - %s (%s, step %d/%d)%s\n",
					resStyleIdle.Render("▶ [ACTIVE]"),
					resStyleHeader.Render(t.TaskCode),
					t.Title,
					resStyleDim.Render(t.AssignedRole),
					t.CurrentStep,
					t.TotalSteps,
					heavyBadge)
			}
			for _, t := range r.PendingTasks {
				heavyBadge := ""
				if t.IsHeavyDuty {
					heavyBadge = " " + resStyleExceed.Render("[HEAVY]")
				}
				fmt.Fprintf(out, "    %s %s - %s (%s)%s\n",
					resStyleDim.Render("• [QUEUED]"),
					resStyleDim.Render(t.TaskCode),
					t.Title,
					resStyleDim.Render(t.AssignedRole),
					heavyBadge)
			}
		}
	}

	fmt.Fprintf(out, "\n  %s\n\n", div)
	return nil
}

func printNodesTasksHelp() error {
	fmt.Print("\nUsage: gitmap nodes tasks [flags] [target]\n\nInspect active and pending AI task queues across fleet nodes.\n\nAliases: gitmap nodes tasks, task, queue, task-queue, jobs, delegation\n\nFlags:\n  -l, --limit <N>   Maximum tasks to display per node (default: 10)\n  -j, --json        Output structured JSON telemetry\n  -h, --help        Show this help message\n\n")
	return nil
}
