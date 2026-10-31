package heatmap

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func dispatchHeatmapAction(subcmd string, flags map[string]string, pos []string, out io.Writer) (bool, error) {
	switch subcmd {
	case "touch", "touch-file", "analyze", "inspect":
		return true, runHeatmapTouch(flags, pos, out)
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
	case "live", "watch", "in-flight", "inflight":
		return true, runHeatmapLive(flags, out)
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

// DispatchHeatmapSubsystem routes heatmap subcommands or direct file reads.
func DispatchHeatmapSubsystem(ctx context.Context, args []string, out io.Writer) error {
	if isHeatmapOverviewRequested(args) {
		return DispatchHeatmapOverview(ctx, out, false)
	}

	if isHeatmapFileReadTarget(args[0]) {
		return cmdmacro.RunCat(args)
	}

	subcmd, flags, pos := parseSubcmdAndArgs(args)
	isHandled, err := dispatchHeatmapAction(subcmd, flags, pos, out)
	if isHandled {
		return err
	}

	return dispatchHeatmapAdmin(ctx, subcmd, args, flags, pos, out)
}

func isHeatmapFileReadTarget(target string) bool {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" || strings.HasPrefix(trimmed, "-") {
		return false
	}

	if isKnownHeatmapSubcommand(trimmed) {
		return false
	}

	if info, err := os.Stat(trimmed); err == nil && !info.IsDir() {
		return true
	}

	hasFileExt := strings.Contains(filepath.Base(trimmed), ".")
	hasPathSep := strings.ContainsAny(trimmed, "/\\")

	return hasFileExt || hasPathSep
}

func isKnownHeatmapSubcommand(sub string) bool {
	switch strings.ToLower(sub) {
	case "overview", "tables", "table", "view", "list", "query", "sql",
		"failures", "broken", "failed", "unmet", "hotfiles", "hot", "hotspots",
		"recent", "latest", "tasks", "task-list", "rc", "rc-heatmap",
		"status", "tree", "live", "watch", "in-flight", "inflight",
		"search", "log", "update", "outcome", "touch", "touch-file",
		"analyze", "inspect", "help", "sync":
		return true
	default:
		return false
	}
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

func runHeatmapTouch(flags map[string]string, pos []string, out io.Writer) error {
	filePath := ""
	if len(pos) > 0 {
		filePath = pos[0]
	}
	if filePath == "" {
		filePath = flags["file"]
	}
	if filePath == "" {
		return apperror.NewValidation("heatmap.touch", "E1000", "usage: gitmap heatmap touch <filepath> [--reason <why>]")
	}

	flags["action"] = "analyze"
	flags["file"] = filePath
	if flags["reason"] == "" && len(pos) > 1 {
		flags["reason"] = strings.Join(pos[1:], " ")
	}

	if err := runHeatmapLog(flags, out); err != nil {
		return err
	}

	if _, err := os.Stat(filePath); err == nil {
		return cmdmacro.RunCat([]string{filePath})
	}

	return nil
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

func runHeatmapLive(flags map[string]string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	limit := resolveLiveLimit(flags["limit"])
	records, err := db.SearchActionLogs(flags["task"], "", flags["file"], limit)
	if err != nil {
		return err
	}

	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, records)
	}

	title := fmt.Sprintf("Root Heatmap: %s", flags["task"])
	renderLiveTouches(out, records, title)

	return nil
}
