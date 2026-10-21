package heatmap

import (
	"context"
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
	"n": "limit", "F": "format", "e": "expected", "o": "actual", "i": "id",
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
	case "update", "outcome":
		return runHeatmapUpdate(flags, pos, out)
	case "search":
		return runHeatmapSearch(flags, pos, out)
	case "tree":
		return runHeatmapTree(flags, out)
	case "status":
		return runHeatmapStatus(flags, out)
	case "tables", "table":
		return RunHeatmapTables(false, flags, out)
	case "view", "list":
		return RunHeatmapView(false, flags, pos, out)
	case "query", "sql":
		return ExecuteHeatmapQuery(false, flags, pos, out)
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
	rec.HeatmapActionLogId = id

	return renderHeatmapLogResponse(out, rec, db.Path, flags)
}

func buildActionLogRecord(flags map[string]string) store.HeatmapActionLogRecord {
	action, state := flags["action"], flags["state"]
	if action == "" {
		action = "ANALYZE"
	}
	if state == "" {
		state = "EXPLORATION"
	}
	start, end := parseLineRange(flags["lines"])
	parentID := parseOptionalID(flags["parent"])
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
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, records)
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
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, tree)
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
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, status)
	}
	renderStatusText(out, status)
	return nil
}

func isHeatmapBoolFlag(k string) bool {
	return k == "json" || k == "j" || k == "help" || k == "h" || k == "dry-run" || k == "verbose" || k == "v" ||
		k == "text" || k == "plain" || k == "envelope"
}

func parseHeatmapArgs(args []string) (map[string]string, []string) {
	flags, pos := make(map[string]string), []string{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if !strings.HasPrefix(arg, "-") {
			pos = append(pos, arg)
			continue
		}
		k, v, found := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		if found {
			flags[k] = v
			continue
		}
		if isHeatmapBoolFlag(k) {
			flags[k], i = resolveHeatmapBoolVal(args, i)
			continue
		}
		if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			flags[k], i = args[i+1], i+1
			continue
		}
		flags[k] = "true"
	}
	return normalizeFlagKeys(flags), pos
}

func resolveHeatmapBoolVal(args []string, i int) (string, int) {
	if i+1 < len(args) && (args[i+1] == "true" || args[i+1] == "false") {
		return args[i+1], i + 1
	}
	return "true", i
}

func parseOptionalID(raw string) *int64 {
	if raw == "" {
		return nil
	}
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || val <= 0 {
		return nil
	}
	return &val
}

func normalizeFlagKeys(flags map[string]string) map[string]string {
	norm := make(map[string]string, len(flags))
	for k, v := range flags {
		if long, ok := shortFlagMap[k]; ok {
			k = long
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
