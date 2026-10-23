package heatmap

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

func parseSubcmdAndArgs(args []string) (string, map[string]string, []string) {
	subcmd := strings.ToLower(args[0])
	flags, pos := parseHeatmapArgs(args[1:])

	return subcmd, flags, pos
}

func deserializePayload(payloadRaw string) (HeatmapLogPayload, error) {
	var payload HeatmapLogPayload
	if err := json.Unmarshal([]byte(payloadRaw), &payload); err != nil {
		return payload, apperror.WrapSimple(err, "heatmap.parseJSONPayload")
	}

	return payload, nil
}

func mergeActionPayloadIntoFlags(flags map[string]string, p HeatmapActionPayload) {
	setFlagIfMissing(flags, "task", p.TaskSlug)
	setFlagIfMissing(flags, "repo", p.RepoName)
	setFlagIfMissing(flags, "action", p.Action)
	setFlagIfMissing(flags, "file", p.File)
	setFlagIfMissing(flags, "state", p.State)
	setFlagIfMissing(flags, "expected", p.Expected)
	setFlagIfMissing(flags, "actual", p.Actual)
	setFlagIfMissing(flags, "ai", p.AiName)
	setFlagIfMissing(flags, "model", p.ModelName)
	setFlagIfMissing(flags, "tool", p.ToolName)
	setFlagIfMissing(flags, "cmd", p.ToolParameters)
	setFlagIfMissing(flags, "met", resolveExpectationFlag(p.IsExpectationMet))
}

func applyJSONPayloadFlags(flags map[string]string) error {
	payloadRaw := flags["json_payload"]
	if payloadRaw == "" {
		return nil
	}

	payloads, err := ParseActionPayloads(payloadRaw)
	if err == nil && len(payloads) > 0 {
		mergeActionPayloadIntoFlags(flags, payloads[0])

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
