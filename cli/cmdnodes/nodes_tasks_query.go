package cmdnodes

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		isPending := strings.EqualFold(status, "PENDING")
		if isPending {
			report.PendingTasks = append(report.PendingTasks, item)
		} else {
			report.ActiveTasks = append(report.ActiveTasks, item)
		}
	}
}
