package cmdnodes

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	_ "modernc.org/sqlite"
)

func discoverAgentTaskDBs() []string {
	var paths []string
	root := filepath.Join(".ai-memory", "temp-agents")
	entries, err := os.ReadDir(root)
	hasErr := err != nil
	if hasErr {
		return paths
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		dbFile := filepath.Join(root, e.Name(), "agent-task.db")
		if _, errStat := os.Stat(dbFile); errStat == nil {
			paths = append(paths, dbFile)
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
		hasRole := len(role) > 0
		if hasRole {
			role = fmt.Sprintf("[%s] ?", role)
		} else {
			role = "?"
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

		isPending := strings.EqualFold(status, "PENDING")
		if isPending {
			report.PendingTasks = append(report.PendingTasks, item)
		} else {
			report.ActiveTasks = append(report.ActiveTasks, item)
		}
	}
}

func resolveNodeSSHConnection(node UnifiedFleetNode, connMap map[string]db.SSHConnection) db.SSHConnection {
	if conn, ok := connMap[strings.ToLower(node.Host)]; ok && len(conn.IPAddress) > 0 {
		return conn
	}
	if conn, ok := connMap[strings.ToLower(node.Alias)]; ok && len(conn.IPAddress) > 0 {
		return conn
	}
	return db.SSHConnection{
		Alias:     node.Alias,
		IPAddress: node.Host,
		Username:  node.User,
	}
}

func resolveFleetNodeEndpoint(node UnifiedFleetNode, connMap map[string]db.SSHConnection) (string, int, db.SSHConnection) {
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
	conn := resolveNodeSSHConnection(node, connMap)
	conn.Alias = alias
	return alias, port, conn
}

func queryRemoteNodeTasksAndHardware(ctx context.Context, node UnifiedFleetNode, conn db.SSHConnection, limit int) NodeTaskReport {
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
		NodeAlias:       alias,
		NodeHost:        node.Host,
		NodePort:        port,
		IsLocal:         false,
		Status:          "ONLINE",
		RolePlaceholder: "?",
		ActiveTasks:     make([]TaskItemSummary, 0),
		PendingTasks:    make([]TaskItemSummary, 0),
	}

	isWin := isWindowsNode(conn) || strings.Contains(strings.ToLower(node.Role), "win")
	remoteCmd := fmt.Sprintf("gitmap tasks --json --limit %d", limit)
	shell := "sh"
	if isWin {
		shell = "powershell"
	}

	out, errExec := currentSSHExecutor.Execute(conn, remoteCmd, shell)
	hasExecErr := errExec != nil
	if hasExecErr {
		return handleRemoteFallbackProbe(alias, node.Host, port, isWin, conn, errExec.Error())
	}

	trimmedOut := strings.TrimSpace(out)
	parsedReport, hasReport := parseRemoteNodeTaskReport(trimmedOut, alias, node.Host, port)
	if hasReport {
		ensureRemoteHardwareMetrics(&parsedReport, isWin, conn)
		return parsedReport
	}

	parsedPayload, hasPayload := parseRemoteFleetPayload(trimmedOut, alias, node.Host, port)
	if hasPayload {
		ensureRemoteHardwareMetrics(&parsedPayload, isWin, conn)
		return parsedPayload
	}

	hwRec := probeNodeExec(alias, fmt.Sprintf("%s:%d", node.Host, port), isWin, false, 80, &conn)
	if hwRec.IsOnline {
		report.Status = "ONLINE"
		report.CPUPercent = hwRec.CPUPercent
		report.RAMPercent = hwRec.RAMPercent
		report.TotalRAMMB = hwRec.TotalRAMMB
		report.UsedRAMMB = hwRec.UsedRAMMB
	} else {
		report.Status = "ACTIVE"
	}

	return report
}

func handleRemoteFallbackProbe(alias, host string, port int, isWin bool, conn db.SSHConnection, errMsg string) NodeTaskReport {
	rep := NodeTaskReport{
		NodeAlias:       alias,
		NodeHost:        host,
		NodePort:        port,
		IsLocal:         false,
		RolePlaceholder: "?",
		ActiveTasks:     make([]TaskItemSummary, 0),
		PendingTasks:    make([]TaskItemSummary, 0),
	}
	hwRec := probeNodeExec(alias, fmt.Sprintf("%s:%d", host, port), isWin, false, 80, &conn)
	if hwRec.IsOnline {
		rep.Status = "ONLINE"
		rep.CPUPercent = hwRec.CPUPercent
		rep.RAMPercent = hwRec.RAMPercent
		rep.TotalRAMMB = hwRec.TotalRAMMB
		rep.UsedRAMMB = hwRec.UsedRAMMB
		rep.ErrorMessage = "gitmap tasks unavailable"
		return rep
	}
	rep.Status = "UNREACHABLE"
	rep.ErrorMessage = errMsg
	return rep
}

func parseRemoteNodeTaskReport(raw, alias, host string, port int) (NodeTaskReport, bool) {
	var rep NodeTaskReport
	errJSON := json.Unmarshal([]byte(raw), &rep)
	hasValidJSON := errJSON == nil && rep.NodeHost != ""
	if !hasValidJSON {
		return NodeTaskReport{}, false
	}
	rep.NodeAlias = alias
	rep.NodeHost = host
	rep.NodePort = port
	rep.IsLocal = false
	hasPlaceholder := rep.RolePlaceholder != ""
	if !hasPlaceholder {
		rep.RolePlaceholder = "?"
	}
	return rep, true
}

func parseRemoteFleetPayload(raw, alias, host string, port int) (NodeTaskReport, bool) {
	var payload FleetTasksPayload
	errPayload := json.Unmarshal([]byte(raw), &payload)
	hasPayload := errPayload == nil && len(payload.NodeReports) > 0
	if !hasPayload {
		return NodeTaskReport{}, false
	}
	first := payload.NodeReports[0]
	first.NodeAlias = alias
	first.NodeHost = host
	first.NodePort = port
	first.IsLocal = false
	hasPlaceholder := first.RolePlaceholder != ""
	if !hasPlaceholder {
		first.RolePlaceholder = "?"
	}
	return first, true
}

func ensureRemoteHardwareMetrics(report *NodeTaskReport, isWin bool, conn db.SSHConnection) {
	needsHardware := report.CPUPercent == 0 && report.RAMPercent == 0
	if !needsHardware {
		return
	}
	hwRec := probeNodeExec(report.NodeAlias, fmt.Sprintf("%s:%d", report.NodeHost, report.NodePort), isWin, false, 80, &conn)
	if hwRec.IsOnline {
		report.CPUPercent = hwRec.CPUPercent
		report.RAMPercent = hwRec.RAMPercent
		report.TotalRAMMB = hwRec.TotalRAMMB
		report.UsedRAMMB = hwRec.UsedRAMMB
	}
}
