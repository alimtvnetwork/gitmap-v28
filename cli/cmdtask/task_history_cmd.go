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
	IsNegativeOffset bool
}

func parseSingleHistoryToken(token string, win *TaskHistoryWindow) {
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

// ParseHistoryArgs extracts limit and offset from CLI arguments.
func ParseHistoryArgs(args []string) TaskHistoryWindow {
	win := TaskHistoryWindow{Limit: 100, Offset: 0, IsNegativeOffset: false}
	hasArgs := len(args) > 0
	if !hasArgs {
		return win
	}
	parseSingleHistoryToken(args[0], &win)
	hasSecond := len(args) > 1
	if hasSecond {
		parseOffsetToken(args[1], &win)
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

func printTaskHistoryHeader(count, total int, win TaskHistoryWindow) {
	fmt.Println()
	fmt.Printf("  %s📋 Task Execution History (showing %d of %d total)%s\n",
		constants.ColorCyan, count, total, constants.ColorReset)
	fmt.Printf("    %-6s %-12s %-26s %-20s %s\n", "ID", "TYPE", "COMPLETED_AT", "SOURCE_CMD", "ARGS")
	fmt.Printf("    %s\n", constants.TermTableRule)
}

func printTaskHistoryRow(t model.CompletedTaskRecord) {
	fmt.Printf("    %-6d %-12s %-26s %-20s %s\n",
		t.ID, t.TaskTypeName, t.CompletedAt, t.SourceCommand, t.CommandArgs)
}

func printTaskHistoryFooter(win TaskHistoryWindow) {
	fmt.Println()
	fmt.Printf("  %sTip: Use 'gitmap task history [limit] [offset]' (e.g. -10 for negative offset)%s\n",
		constants.ColorDim, constants.ColorReset)
	fmt.Printf("  %s     Use 'gitmap task undo [id]' to revert a task.%s\n\n",
		constants.ColorDim, constants.ColorReset)
}

// RunTaskHistory displays completed task history with limit and offset support.
func RunTaskHistory(args []string) error {
	tasksDB, errOpen := store.OpenTasksRootSplitDB()
	if errOpen != nil {
		return apperror.WrapSimple(errOpen, constants.WarnPendingDBOpen)
	}
	defer tasksDB.Close()

	completed, errList := tasksDB.ListCompletedTasks()
	if errList != nil {
		return apperror.WrapSimple(errList, constants.ErrPendingTaskQuery)
	}

	win := ParseHistoryArgs(args)
	start, end := calculateSliceBounds(len(completed), win)
	slice := completed[start:end]

	printTaskHistoryHeader(len(slice), len(completed), win)
	for _, t := range slice {
		printTaskHistoryRow(t)
	}
	printTaskHistoryFooter(win)

	return nil
}
