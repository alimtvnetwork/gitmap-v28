package cmdnodes

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

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
			renderNodeActiveTasks(out, r.ActiveTasks)
			renderNodePendingTasks(out, r.PendingTasks)
		}
	}

	fmt.Fprintf(out, "\n  %s\n\n", div)
	return nil
}

func renderNodeActiveTasks(out io.Writer, tasks []TaskItemSummary) {
	for _, t := range tasks {
		heavyBadge := renderHeavyBadge(t.IsHeavyDuty)
		fmt.Fprintf(out, "    %s %s - %s (%s, step %d/%d)%s\n",
			resStyleIdle.Render("▶ [ACTIVE]"),
			resStyleHeader.Render(t.TaskCode),
			t.Title,
			resStyleDim.Render(t.AssignedRole),
			t.CurrentStep,
			t.TotalSteps,
			heavyBadge)
	}
}

func renderNodePendingTasks(out io.Writer, tasks []TaskItemSummary) {
	for _, t := range tasks {
		heavyBadge := renderHeavyBadge(t.IsHeavyDuty)
		fmt.Fprintf(out, "    %s %s - %s (%s)%s\n",
			resStyleDim.Render("• [QUEUED]"),
			resStyleDim.Render(t.TaskCode),
			t.Title,
			resStyleDim.Render(t.AssignedRole),
			heavyBadge)
	}
}

func renderHeavyBadge(isHeavy bool) string {
	if isHeavy {
		return " " + resStyleExceed.Render("[HEAVY]")
	}
	return ""
}

func printNodesTasksHelp() error {
	fmt.Print("\nUsage: gitmap nodes tasks [flags] [target]\n\nInspect active and pending AI task queues across fleet nodes.\n\nAliases: gitmap nodes tasks, task, queue, task-queue, jobs, delegation\n\nFlags:\n  -l, --limit <N>   Maximum tasks to display per node (default: 10)\n  -j, --json        Output structured JSON telemetry\n  -h, --help        Show this help message\n\n")
	return nil
}
