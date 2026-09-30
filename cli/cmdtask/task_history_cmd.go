// Package cmdtask manages task history inspection and querying.
package cmdtask

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// TaskHistoryWindow holds pagination limits and offsets.
type TaskHistoryWindow struct {
	Limit            int
	Offset           int
	Section          string
	IsNegativeOffset bool
}

func parseSingleHistoryToken(token string, win *TaskHistoryWindow) {
	if strings.HasPrefix(token, "--section=") {
		win.Section = strings.TrimPrefix(token, "--section=")
		return
	}
	hasNeg := strings.HasPrefix(token, "-")
	clean := strings.TrimPrefix(token, "-")
	val, err := strconv.Atoi(clean)
	isNum := err == nil && val > 0
	if isNum && hasNeg {
		win.Offset = val
		win.IsNegativeOffset = true
		return
	}
	if isNum {
		win.Limit = val
	}
}

func parseOffsetToken(token string, win *TaskHistoryWindow) {
	if strings.HasPrefix(token, "--section=") {
		win.Section = strings.TrimPrefix(token, "--section=")
		return
	}
	hasNeg := strings.HasPrefix(token, "-")
	clean := strings.TrimPrefix(token, "-")
	val, err := strconv.Atoi(clean)
	isNum := err == nil && val >= 0
	if isNum && hasNeg {
		win.Offset = val
		win.IsNegativeOffset = true
		return
	}
	if isNum {
		win.Offset = val
	}
}

// ParseHistoryArgs extracts limit, offset, and section from CLI arguments.
func ParseHistoryArgs(args []string) TaskHistoryWindow {
	win := TaskHistoryWindow{Limit: 50, Offset: 0, Section: "all", IsNegativeOffset: false}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--section" && i+1 < len(args) {
			win.Section = args[i+1]
			i++
			continue
		}
		if i == 0 {
			parseSingleHistoryToken(a, &win)
			continue
		}
		parseOffsetToken(a, &win)
	}

	return win
}

func calculateSliceBounds(total int, win TaskHistoryWindow) (int, int) {
	if win.IsNegativeOffset {
		return calculateNegativeBounds(total, win)
	}

	return calculatePositiveBounds(total, win)
}

func calculateNegativeBounds(total int, win TaskHistoryWindow) (int, int) {
	start := total - win.Offset
	if start < 0 {
		start = 0
	}
	end := start + win.Limit
	if end > total {
		end = total
	}

	return start, end
}

func calculatePositiveBounds(total int, win TaskHistoryWindow) (int, int) {
	start := win.Offset
	if start > total {
		start = total
	}
	end := start + win.Limit
	if end > total {
		end = total
	}

	return start, end
}

func formatSectionBadge(sec string) string {
	switch strings.ToLower(sec) {
	case "macro":
		return fmt.Sprintf("%s⚡ MACRO    %s", constants.ColorCyan, constants.ColorReset)
	case "ssh":
		return fmt.Sprintf("%s🔑 SSH      %s", constants.ColorMagenta, constants.ColorReset)
	case "installer", "install":
		return fmt.Sprintf("%s📦 INSTALL  %s", constants.ColorBlue, constants.ColorReset)
	default:
		return fmt.Sprintf("%s📋 %-9s%s", constants.ColorYellow, strings.ToUpper(sec), constants.ColorReset)
	}
}

func formatStatusBadge(status string) string {
	low := strings.ToLower(status)
	switch {
	case strings.Contains(low, "fail"):
		return fmt.Sprintf("%s[✖ FAILED]   %s", constants.ColorRed, constants.ColorReset)
	case strings.Contains(low, "pend") || strings.Contains(low, "run"):
		return fmt.Sprintf("%s[⏳ RUNNING]  %s", constants.ColorYellow, constants.ColorReset)
	default:
		return fmt.Sprintf("%s[✔ COMPLETED]%s", constants.ColorGreen, constants.ColorReset)
	}
}

func printAuditHistoryHeader(count int, secFilter string) {
	fmt.Println()
	filterTag := ""
	if secFilter != "" && secFilter != "all" {
		filterTag = fmt.Sprintf(" [Filter: %s]", secFilter)
	}
	fmt.Printf("  %s╔══ TASK & AUDIT EXECUTION HISTORY (showing %d)%s══════════════════════════════════╗%s\n",
		constants.ColorCyan, count, filterTag, constants.ColorReset)
	fmt.Printf("  %s║ %-5s %-12s %-10s %-14s %-24s %-20s ║%s\n",
		constants.ColorCyan, "ID", "SECTION", "ACTION", "STATUS", "TARGET", "EXECUTED AT", constants.ColorReset)
	fmt.Printf("  %s╠═════════════════════════════════════════════════════════════════════════════════════╣%s\n",
		constants.ColorCyan, constants.ColorReset)
}

func printAuditHistoryRow(r model.TaskHistoryRecord) {
	target := r.Target
	if len(target) > 23 {
		target = target[:20] + "..."
	}
	secBadge := formatSectionBadge(r.Section)
	statBadge := formatStatusBadge(r.Status)
	timeStr := r.ExecutedAt
	if len(timeStr) > 19 {
		timeStr = timeStr[:19]
	}

	fmt.Printf("  │ %-5d %s %-10s %s %-24s %-20s │\n",
		r.TaskHistoryId, secBadge, r.Action, statBadge, target, timeStr)
}

func printAuditHistoryFooter() {
	fmt.Printf("  %s╚═════════════════════════════════════════════════════════════════════════════════════╝%s\n",
		constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Printf("  %sTip: Filter by section: 'gitmap task history --section macro|ssh|installer'%s\n",
		constants.ColorDim, constants.ColorReset)
	fmt.Printf("  %s     Inspect task detail: 'gitmap task view <id>'%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

// RunTaskHistory displays completed task history with limit, offset, and section filter support.
func RunTaskHistory(args []string) error {
	tasksDB, errOpen := store.OpenTasksRootSplitDB()
	if errOpen != nil {
		return apperror.WrapSimple(errOpen, constants.WarnPendingDBOpen)
	}
	defer tasksDB.Close()

	win := ParseHistoryArgs(args)
	records, _ := tasksDB.ListTaskHistory(win.Section, win.Limit, win.Offset)
	if len(records) > 0 {
		printAuditHistoryHeader(len(records), win.Section)
		for _, r := range records {
			printAuditHistoryRow(r)
		}
		printAuditHistoryFooter()

		return nil
	}

	return renderLegacyCompletedTasks(tasksDB, win)
}

func renderLegacyCompletedTasks(tasksDB *store.TasksSplitDB, win TaskHistoryWindow) error {
	completed, errList := tasksDB.ListCompletedTasks()
	if errList != nil {
		return apperror.WrapSimple(errList, constants.ErrPendingTaskQuery)
	}

	start, end := calculateSliceBounds(len(completed), win)
	slice := completed[start:end]

	fmt.Println()
	fmt.Printf("  %s📋 Task Execution History (showing %d of %d total)%s\n",
		constants.ColorCyan, len(slice), len(completed), constants.ColorReset)
	fmt.Printf("    %-6s %-12s %-26s %-20s %s\n", "ID", "TYPE", "COMPLETED_AT", "SOURCE_CMD", "ARGS")
	fmt.Printf("    %s\n", constants.TermTableRule)
	for _, t := range slice {
		fmt.Printf("    %-6d %-12s %-26s %-20s %s\n",
			t.ID, t.TaskTypeName, t.CompletedAt, t.SourceCommand, t.CommandArgs)
	}
	fmt.Println()

	return nil
}
