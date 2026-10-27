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

func dispatchHeatmapAction(subcmd string, flags map[string]string, pos []string, out io.Writer) (bool, error) {
	switch subcmd {
	case "log":
		return true, runHeatmapLog(flags, out)
	case "update", "outcome":
		return true, runHeatmapUpdate(flags, pos, out)
	case "search":
		return true, runHeatmapSearch(flags, pos, out)
	case "tree":
		return true, runHeatmapTree(flags, out)
	case "status":
		return true, runHeatmapStatus(flags, out)
	}

	return false, nil
}

func dispatchHeatmapAdmin(ctx context.Context, subcmd string, args []string, flags map[string]string, pos []string, out io.Writer) error {
	switch subcmd {
	case "overview":
		return DispatchHeatmapOverview(ctx, out, false)
	case "tables", "table":
		return RunHeatmapTables(false, flags, out)
	case "view", "list":
		return RunHeatmapView(false, flags, pos, out)
	case "query", "sql":
		return ExecuteHeatmapQuery(false, flags, pos, out)
	case "failures", "broken", "failed", "unmet":
		return RunHeatmapFailures(false, flags, out)
	case "hotfiles", "hot", "hotspots":
		return RunHeatmapHotfiles(false, flags, out)
	case "recent", "latest":
		return RunHeatmapRecent(false, flags, out)
	case "tasks", "task-list":
		return RunHeatmapTasks(false, flags, out)
	case "rc", "rc-heatmap":
		return DispatchRcHeatmap(ctx, args[1:], out)
	default:
		return apperror.NewValidation("heatmap", "E1000", fmt.Sprintf("unknown heatmap command: %s", subcmd))
	}
}

// DispatchHeatmapSubsystem routes heatmap subcommands (log, search, tree, status, help).
func DispatchHeatmapSubsystem(ctx context.Context, args []string, out io.Writer) error {
	if isHeatmapOverviewRequested(args) {
		return DispatchHeatmapOverview(ctx, out, false)
	}

	subcmd, flags, pos := parseSubcmdAndArgs(args)
	isHandled, err := dispatchHeatmapAction(subcmd, flags, pos, out)
	if isHandled {
		return err
	}

	return dispatchHeatmapAdmin(ctx, subcmd, args, flags, pos, out)
}

func normalizeJSONInputFlag(flags map[string]string) {
	val := strings.TrimSpace(flags["json"])
	if strings.HasPrefix(val, "{") || strings.HasPrefix(val, "[") {
		flags["json_payload"] = val
		flags["json"] = "true"
	}
}

func openRootAndApplyPayload(flags map[string]string) (*store.HeatmapDB, error) {
	normalizeJSONInputFlag(flags)
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return nil, err
	}

	if err := applyJSONPayloadFlags(flags); err != nil {
		_ = db.Close()

		return nil, err
	}

	return db, nil
}

func insertActionAndSetID(db *store.HeatmapDB, rec *store.HeatmapActionLogRecord) error {
	id, err := db.InsertActionLog(*rec)
	if err != nil {
		return err
	}

	rec.HeatmapActionLogId = id

	return nil
}

func isArrayPayload(raw string) bool {
	return strings.HasPrefix(strings.TrimSpace(raw), "[")
}

func insertMultipleActionLogs(db *store.HeatmapDB, raw string, out io.Writer, flags map[string]string) error {
	payloads, err := ParseActionPayloads(raw)
	if err != nil {
		return err
	}

	var lastRec store.HeatmapActionLogRecord
	for _, p := range payloads {
		rec := p.ToRecord()
		if err := insertActionAndSetID(db, &rec); err != nil {
			return err
		}
		lastRec = rec
	}

	return renderHeatmapLogResponse(out, lastRec, db.Path, flags)
}

func runHeatmapLog(flags map[string]string, out io.Writer) error {
	db, err := openRootAndApplyPayload(flags)
	if err != nil {
		return err
	}

	defer db.Close()

	if isArrayPayload(flags["json_payload"]) {
		return insertMultipleActionLogs(db, flags["json_payload"], out, flags)
	}

	rec := buildActionLogRecord(flags)
	if err := insertActionAndSetID(db, &rec); err != nil {
		return err
	}

	return renderHeatmapLogResponse(out, rec, db.Path, flags)
}

func resolveModelAlias(flags map[string]string) string {
	if m := flags["model"]; m != "" {
		return m
	}

	return flags["modelName"]
}

func resolveModelFlag(flags map[string]string) string {
	if m := resolveModelAlias(flags); m != "" {
		return m
	}

	if m := flags["ai"]; m != "" {
		return m
	}

	return flags["aiName"]
}

func resolveToolFlag(flags map[string]string) string {
	if t := flags["tool"]; t != "" {
		return t
	}

	return flags["toolName"]
}

func resolveActionFlag(flags map[string]string) string {
	if action := flags["action"]; action != "" {
		return strings.ToUpper(action)
	}

	return "ANALYZE"
}

func resolveStateFlag(flags map[string]string) string {
	if state := flags["state"]; state != "" {
		return strings.ToUpper(state)
	}

	return "EXPLORATION"
}

func populateRecordDetails(rec *store.HeatmapActionLogRecord, flags map[string]string) {
	rec.CodeSnippet = flags["code"]
	rec.CommandText = flags["cmd"]
	rec.ReasoningNotes = flags["reason"]
	rec.ExpectedOutcome = flags["expected"]
	rec.ActualOutcome = flags["actual"]
	rec.IsExpectationMet = flags["met"] != "false"
	rec.ModelName = resolveModelFlag(flags)
	rec.ToolName = resolveToolFlag(flags)
}

func initActionLogRecord(flags map[string]string) store.HeatmapActionLogRecord {
	start, end := parseLineRange(flags["lines"])

	return store.HeatmapActionLogRecord{
		ParentLogId: parseOptionalID(flags["parent"]),
		TaskSlug:    flags["task"], RepoName: flags["repo"],
		ActionType: resolveActionFlag(flags), FilePath: flags["file"],
		StartLine: start, EndLine: end, ReasoningState: resolveStateFlag(flags),
	}
}

func buildActionLogRecord(flags map[string]string) store.HeatmapActionLogRecord {
	rec := initActionLogRecord(flags)
	populateRecordDetails(&rec, flags)

	return rec
}

func resolveSearchLimit(flags map[string]string) int {
	if l, err := strconv.Atoi(flags["limit"]); err == nil && l > 0 {
		return l
	}

	return 50
}

func renderSearchResults(records []store.HeatmapActionLogRecord, flags map[string]string, out io.Writer) error {
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, records)
	}

	renderSearchText(out, records)

	return nil
}

func executeActionSearch(db *store.HeatmapDB, flags map[string]string, pos []string) ([]store.HeatmapActionLogRecord, error) {
	limit := resolveSearchLimit(flags)

	return db.SearchActionLogs(flags["task"], strings.Join(pos, " "), flags["file"], limit)
}

func runHeatmapSearch(flags map[string]string, pos []string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}

	defer db.Close()

	records, err := executeActionSearch(db, flags, pos)
	if err != nil {
		return err
	}

	return renderSearchResults(records, flags, out)
}

func renderTreeResults(tree []store.HeatmapExecutionTreeNode, flags map[string]string, out io.Writer) error {
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, tree)
	}

	renderTreeText(out, tree, "")

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

	return renderTreeResults(tree, flags, out)
}

type statusContext struct {
	inv        TaskInventory
	activeSlug string
	alerts     []HeatmapAiConfusion
	cmds       SuggestedCommandsPayload
}

func resolveStatusLimit(flags map[string]string) int {
	if l, err := strconv.Atoi(flags["limit"]); err == nil && l > 0 {
		return l
	}

	return 5
}

func loadEnrichedTaskInventory(db *store.HeatmapDB, flags map[string]string) TaskInventory {
	limit := resolveStatusLimit(flags)
	inv, _ := DiscoverHeatmapTasks(".", limit)
	inv.Completed = enrichTasksWithDbMetrics(db, inv.Completed)
	inv.Pending = enrichTasksWithDbMetrics(db, inv.Pending)

	return inv
}

func buildStatusContext(status store.HeatmapStatusSummary, db *store.HeatmapDB, flags map[string]string) statusContext {
	inv := loadEnrichedTaskInventory(db, flags)
	activeSlug := ResolveActiveTaskSlug(status.TaskSlug, inv)

	return statusContext{
		inv:        inv,
		activeSlug: activeSlug,
		alerts:     CollectConfusionAlerts(db, activeSlug),
		cmds:       BuildSuggestedCommands(activeSlug),
	}
}

func renderStatusResults(status store.HeatmapStatusSummary, db *store.HeatmapDB, flags map[string]string, out io.Writer) error {
	ctx := buildStatusContext(status, db, flags)
	if flags["json"] == "true" || isJsonRequested(flags) {
		wsSlug := ResolveWorkspaceSlug(".")
		payload := BuildHeatmapStatusPayload(status, wsSlug, ctx.activeSlug, ctx.inv, ctx.alerts, ctx.cmds)

		return writeIndentedJSON(out, payload)
	}

	renderStatusDashboard(out, status, ctx.inv, ctx.alerts, ctx.cmds)

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

	return renderStatusResults(status, db, flags, out)
}
