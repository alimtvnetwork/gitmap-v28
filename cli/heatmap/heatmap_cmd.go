package heatmap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var shortFlagMap = map[string]string{
	"t": "task", "a": "action", "f": "file", "l": "lines",
	"c": "code", "s": "state", "r": "reason", "p": "parent", "j": "json",
}

// DispatchHeatmapSubsystem routes heatmap subcommands (log, search, tree, status, help).
func DispatchHeatmapSubsystem(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHeatmapHelp(out)
		return nil
	}

	subcmd := strings.ToLower(args[0])
	flags, pos := parseHeatmapArgs(args[1:])

	switch subcmd {
	case "log":
		return runHeatmapLog(flags, out)
	case "search":
		return runHeatmapSearch(flags, pos, out)
	case "tree":
		return runHeatmapTree(flags, out)
	case "status":
		return runHeatmapStatus(flags, out)
	case "rc", "rc-heatmap":
		return DispatchRcHeatmap(ctx, args[1:], out)
	default:
		return apperror.NewValidation("heatmap", "E1000", fmt.Sprintf("unknown heatmap command: %s", subcmd))
	}
}

func runHeatmapLog(flags map[string]string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	rec := buildActionLogRecord(flags)
	id, err := db.LogAction(rec)
	if err != nil {
		return err
	}

	fmt.Fprintf(out, "Logged heatmap action #%d [task=%s, action=%s, file=%s]\n",
		id, rec.TaskSlug, rec.ActionType, rec.FilePath)
	return nil
}

func buildActionLogRecord(flags map[string]string) store.HeatmapActionLogRecord {
	action := flags["action"]
	if action == "" {
		action = "ANALYZE"
	}
	state := flags["state"]
	if state == "" {
		state = "EXPLORATION"
	}
	start, end := parseLineRange(flags["lines"])
	var parentID *int64
	if pStr := flags["parent"]; pStr != "" {
		if pVal, err := strconv.ParseInt(pStr, 10, 64); err == nil && pVal > 0 {
			parentID = &pVal
		}
	}

	return store.HeatmapActionLogRecord{
		ParentLogId:      parentID,
		TaskSlug:         flags["task"],
		RepoName:         flags["repo"],
		ActionType:       strings.ToUpper(action),
		FilePath:         flags["file"],
		StartLine:        start,
		EndLine:          end,
		CodeSnippet:      flags["code"],
		CommandText:      flags["cmd"],
		ReasoningState:   strings.ToUpper(state),
		ReasoningNotes:   flags["reason"],
		ExpectedOutcome:  flags["expected"],
		ActualOutcome:    flags["actual"],
		IsExpectationMet: flags["met"] != "false",
	}
}

func runHeatmapSearch(flags map[string]string, pos []string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	limit := 50
	if l, err := strconv.Atoi(flags["limit"]); err == nil && l > 0 {
		limit = l
	}

	records, err := db.SearchActionLogs(flags["task"], strings.Join(pos, " "), flags["file"], limit)
	if err != nil {
		return err
	}

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(records)
	}
	renderSearchText(out, records)
	return nil
}

func runHeatmapTree(flags map[string]string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	tree, err := db.GetExecutionTree(flags["task"])
	if err != nil {
		return err
	}

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(tree)
	}
	renderTreeText(out, tree, "")
	return nil
}

func runHeatmapStatus(flags map[string]string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	status, err := db.GetHeatmapStatus(flags["task"])
	if err != nil {
		return err
	}

	if flags["json"] == "true" {
		return json.NewEncoder(out).Encode(status)
	}
	renderStatusText(out, status)
	return nil
}

func parseHeatmapArgs(args []string) (map[string]string, []string) {
	flags := make(map[string]string)
	var pos []string

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			pos = append(pos, arg)
			continue
		}
		trimmed := strings.TrimLeft(arg, "-")
		k, v, found := strings.Cut(trimmed, "=")
		if found {
			flags[k] = v
			continue
		}
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			flags[k] = args[i+1]
			i++
			continue
		}
		flags[k] = "true"
	}

	return normalizeFlagKeys(flags), pos
}

func normalizeFlagKeys(flags map[string]string) map[string]string {
	norm := make(map[string]string, len(flags))
	for k, v := range flags {
		if long, ok := shortFlagMap[k]; ok {
			norm[long] = v
			continue
		}
		norm[k] = v
	}
	return norm
}

func parseLineRange(raw string) (int, int) {
	if raw == "" {
		return 0, 0
	}
	sStr, eStr, hasEnd := strings.Cut(raw, "-")
	start, _ := strconv.Atoi(strings.TrimSpace(sStr))
	if hasEnd {
		end, _ := strconv.Atoi(strings.TrimSpace(eStr))
		return start, end
	}
	return start, start
}

func renderSearchText(out io.Writer, records []store.HeatmapActionLogRecord) {
	if len(records) == 0 {
		fmt.Fprintln(out, "No heatmap action logs found matching query.")
		return
	}
	fmt.Fprintf(out, "Found %d heatmap action log(s):\n", len(records))
	fmt.Fprintf(out, "%-6s %-14s %-8s %-30s %s\n", "ID", "STATE", "ACTION", "FILE", "REASON")
	fmt.Fprintf(out, "%-6s %-14s %-8s %-30s %s\n", "------", "--------------", "--------", "------------------------------", "--------------------")
	for _, r := range records {
		fStr := r.FilePath
		if r.StartLine > 0 {
			fStr += fmt.Sprintf(":%d", r.StartLine)
		}
		if len(fStr) > 30 {
			fStr = "..." + fStr[len(fStr)-27:]
		}
		fmt.Fprintf(out, "%-6d %-14s %-8s %-30s %s\n", r.HeatmapActionLogId, r.ReasoningState, r.ActionType, fStr, r.ReasoningNotes)
	}
}

func renderTreeText(out io.Writer, nodes []store.HeatmapExecutionTreeNode, indent string) {
	for i, node := range nodes {
		branch, nextIndent := "├── ", indent+"│   "
		if i == len(nodes)-1 {
			branch, nextIndent = "└── ", indent+"    "
		}
		if indent == "" {
			branch, nextIndent = "• ", "  "
		}
		lineInfo := ""
		if node.FilePath != "" {
			lineInfo = fmt.Sprintf(" [%s:%d-%d]", node.FilePath, node.StartLine, node.EndLine)
		}
		fmt.Fprintf(out, "%s%s#%d (%s) %s%s: %s\n", indent, branch, node.HeatmapActionLogId,
			node.ReasoningState, node.ActionType, lineInfo, node.ReasoningNotes)
		if len(node.Children) > 0 {
			renderTreeText(out, node.Children, nextIndent)
		}
	}
}

func renderStatusText(out io.Writer, s store.HeatmapStatusSummary) {
	fmt.Fprintf(out, "Heatmap Status Summary (Task: %s)\n--------------------------------------------------\n", s.TaskSlug)
	fmt.Fprintf(out, "Total Actions:    %d\nExpectations Met: %d (%.1f%%)\n\nReasoning States:\n",
		s.TotalActions, s.ExpectationsMetCount, s.ExpectationsMetPct)
	for state, count := range s.ReasoningStateCounts {
		fmt.Fprintf(out, "  %-18s: %d\n", state, count)
	}
	if len(s.TopTouchedFiles) > 0 {
		fmt.Fprintln(out, "\nTop Touched Files:")
		for _, f := range s.TopTouchedFiles {
			fmt.Fprintf(out, "  %-40s (%d touches)\n", f.FilePath, f.Count)
		}
	}
}

func printHeatmapHelp(out io.Writer) {
	fmt.Fprintln(out, "Usage: gitmap heatmap <command> [flags]\n\nCommands:")
	fmt.Fprintln(out, "  log     Log an AI cognitive decision or file action")
	fmt.Fprintln(out, "  search  Search logged actions by query, task, or file")
	fmt.Fprintln(out, "  tree    Display hierarchical causal execution tree")
	fmt.Fprintln(out, "  status  Summarize action volume, states, and file touches")
	fmt.Fprintln(out, "  rc      Manage repo-cache isolated multi-AI heatmap telemetry")
	fmt.Fprintln(out, "  help    Show this help message\n\nExamples:")
	fmt.Fprintln(out, "  gitmap heatmap log --task t1 --action READ --file cli/store/path.go --state EXPLORATION --reason \"Checking base paths\"")
	fmt.Fprintln(out, "  gitmap heatmap log --task t1 --action MODIFY --file cli/store/heatmap.go --lines 10-50 --parent 1 --state CODE_EDITING")
	fmt.Fprintln(out, "  gitmap heatmap tree --task t1")
	fmt.Fprintln(out, "  gitmap heatmap status --task t1")
	fmt.Fprintln(out, "  gitmap heatmap search \"base paths\" --json")
}
