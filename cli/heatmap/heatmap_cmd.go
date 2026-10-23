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
	"n": "limit", "F": "format", "e": "expected", "o": "actual", "i": "id",
	"m": "model",
}

func isHeatmapHelpRequested(args []string) bool {
	return isHeatmapOverviewRequested(args)
}

func isHeatmapOverviewRequested(args []string) bool {
	if len(args) == 0 {
		return true
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	if first == "help" || first == "--help" || first == "-h" || first == "overview" {
		return true
	}

	return strings.HasPrefix(first, "-")
}

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
	case "rc", "rc-heatmap":
		return DispatchRcHeatmap(ctx, args[1:], out)
	default:
		return apperror.NewValidation("heatmap", "E1000", fmt.Sprintf("unknown heatmap command: %s", subcmd))
	}
}

func parseSubcmdAndArgs(args []string) (string, map[string]string, []string) {
	subcmd := strings.ToLower(args[0])
	flags, pos := parseHeatmapArgs(args[1:])

	return subcmd, flags, pos
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

func deserializePayload(payloadRaw string) (HeatmapLogPayload, error) {
	var payload HeatmapLogPayload
	if err := json.Unmarshal([]byte(payloadRaw), &payload); err != nil {
		return payload, apperror.WrapSimple(err, "heatmap.parseJSONPayload")
	}

	return payload, nil
}

func applyJSONPayloadFlags(flags map[string]string) error {
	payloadRaw := flags["json_payload"]
	if payloadRaw == "" {
		return nil
	}

	payload, err := deserializePayload(payloadRaw)
	if err != nil {
		return err
	}

	mergePayloadIntoFlags(flags, payload)

	return nil
}

func setFlagIfMissing(flags map[string]string, key, val string) {
	hasVal := val != ""
	isUnset := flags[key] == ""
	if hasVal && isUnset {
		flags[key] = val
	}
}

func mergePayloadIdentity(flags map[string]string, p HeatmapLogPayload) {
	setFlagIfMissing(flags, "task", p.ResolveTask())
	setFlagIfMissing(flags, "repo", p.ResolveRepo())
	setFlagIfMissing(flags, "action", p.ResolveAction())
	setFlagIfMissing(flags, "file", p.ResolveFile())
	setFlagIfMissing(flags, "lines", p.Lines)
	setFlagIfMissing(flags, "code", p.ResolveCode())
	setFlagIfMissing(flags, "cmd", p.ResolveCmd())
}

func resolveMetString(p HeatmapLogPayload) string {
	if p.HasMetExpectation() {
		return "true"
	}

	return "false"
}

func mergePayloadMetStatus(flags map[string]string, p HeatmapLogPayload) {
	hasMetFlag := flags["met"] != ""
	hasExplicitMet := p.IsMet != nil || p.IsExpectationMet != nil
	if hasMetFlag || !hasExplicitMet {
		return
	}

	flags["met"] = resolveMetString(p)
}

func mergePayloadReasoning(flags map[string]string, p HeatmapLogPayload) {
	setFlagIfMissing(flags, "state", p.ResolveState())
	setFlagIfMissing(flags, "reason", p.ResolveReason())
	setFlagIfMissing(flags, "expected", p.ResolveExpected())
	setFlagIfMissing(flags, "actual", p.ResolveActual())
	setFlagIfMissing(flags, "model", p.ResolveModel())
	setFlagIfMissing(flags, "tool", p.ResolveTool())
	mergePayloadMetStatus(flags, p)
}

func resolveParentIDStr(parentLogID *int64) string {
	if parentLogID != nil {
		return strconv.FormatInt(*parentLogID, 10)
	}

	return ""
}

func mergePayloadParent(flags map[string]string, p HeatmapLogPayload) {
	if flags["parent"] != "" {
		return
	}

	if idStr := resolveParentIDStr(p.ParentLogId); idStr != "" {
		flags["parent"] = idStr

		return
	}

	if p.Parent != "" {
		flags["parent"] = p.Parent
	}
}

func mergePayloadIntoFlags(flags map[string]string, p HeatmapLogPayload) {
	mergePayloadIdentity(flags, p)
	mergePayloadReasoning(flags, p)
	mergePayloadParent(flags, p)
}

func openRootAndApplyPayload(flags map[string]string) (*store.HeatmapDB, error) {
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

func runHeatmapLog(flags map[string]string, out io.Writer) error {
	db, err := openRootAndApplyPayload(flags)
	if err != nil {
		return err
	}

	defer db.Close()

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

func renderStatusResults(status store.HeatmapStatusSummary, flags map[string]string, out io.Writer) error {
	if flags["json"] == "true" || isJsonRequested(flags) {
		return writeIndentedJSON(out, status)
	}

	renderStatusText(out, status)

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

	return renderStatusResults(status, flags, out)
}

func isHeatmapBoolFlag(k string) bool {
	return k == "json" || k == "j" || k == "help" || k == "h" || k == "dry-run" || k == "verbose" || k == "v" ||
		k == "text" || k == "plain" || k == "envelope"
}

func isJSONFlag(k string) bool {
	return k == "json" || k == "j"
}

func isJSONPayloadToken(arg string) bool {
	trimmed := strings.TrimSpace(arg)
	hasObj := strings.HasPrefix(trimmed, "{")
	hasArr := strings.HasPrefix(trimmed, "[")

	return hasObj || hasArr
}

func tryConsumeJSONPayload(flags map[string]string, args []string, i int, k, v string, hasCut bool) (bool, int) {
	if !isJSONFlag(k) {
		return false, i
	}

	if hasCut && isJSONPayloadToken(v) {
		flags["json_payload"] = v
		flags["json"] = "true"

		return true, i
	}

	hasNext := i+1 < len(args) && isJSONPayloadToken(args[i+1])
	if hasNext {
		flags["json_payload"] = args[i+1]
		flags["json"] = "true"

		return true, i + 1
	}

	return false, i
}

func resolveHeatmapVal(args []string, i int) (string, int) {
	hasNextVal := i+1 < len(args) && !strings.HasPrefix(args[i+1], "-")
	if hasNextVal {
		return args[i+1], i + 1
	}

	return "true", i
}

func processFlagWithValue(flags map[string]string, args []string, i int, k, v string, hasCut bool) int {
	if hasCut {
		flags[k] = v

		return i
	}

	if isHeatmapBoolFlag(k) {
		flags[k], i = resolveHeatmapBoolVal(args, i)

		return i
	}

	flags[k], i = resolveHeatmapVal(args, i)

	return i
}

func processHeatmapFlag(flags map[string]string, args []string, i int) int {
	k, v, hasCut := strings.Cut(strings.TrimLeft(args[i], "-"), "=")
	isJSON, nextIdx := tryConsumeJSONPayload(flags, args, i, k, v, hasCut)
	if isJSON {
		return nextIdx
	}

	return processFlagWithValue(flags, args, i, k, v, hasCut)
}

func parseHeatmapArgs(args []string) (map[string]string, []string) {
	flags, pos := make(map[string]string), []string{}
	for i := 0; i < len(args); i++ {
		if !strings.HasPrefix(args[i], "-") {
			pos = append(pos, args[i])
			continue
		}

		i = processHeatmapFlag(flags, args, i)
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
		key := k
		if long, ok := shortFlagMap[k]; ok {
			key = long
		}

		norm[key] = v
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
