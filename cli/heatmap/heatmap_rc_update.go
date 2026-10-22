package heatmap

import (
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/strutil"
)

func runRcUpdate(flags map[string]string, pos []string, out io.Writer) error {
	aiName, repoSlug := resolveTargetAi(flags["ai"]), resolveTargetRepo(flags["repo"])
	st, err := store.OpenRcHeatmapStore(aiName, repoSlug)
	if err != nil {
		return err
	}
	defer st.Close()

	actionId := parseActionId(flags, pos)
	if actionId <= 0 {
		return apperror.NewValidation("heatmap_rc", "E1001", "missing or invalid --id for heatmap update")
	}

	expected := flags["expected"]
	actual := resolveActualOutcomeFlag(flags)
	isMet := resolveIsExpectationMet(flags)
	notes := flags["notes"]
	if notes == "" {
		notes = flags["reason"]
	}

	updated, err := st.UpdateActionOutcome(actionId, expected, actual, isMet, notes)
	if err != nil {
		return err
	}

	statusStr := "MET"
	if !updated.IsExpectationMet {
		statusStr = "FAILED"
	}

	return renderRcUpdateResponse(out, *updated, st.RepoDbPath(), aiName, repoSlug, statusStr, flags)
}

func runHeatmapUpdate(flags map[string]string, pos []string, out io.Writer) error {
	db, err := store.OpenHeatmapRootDB("")
	if err != nil {
		return err
	}
	defer db.Close()

	actionId := parseActionId(flags, pos)
	if actionId <= 0 {
		return apperror.NewValidation("heatmap", "E1001", "missing or invalid --id for heatmap update")
	}

	expected := flags["expected"]
	actual := resolveActualOutcomeFlag(flags)
	isMet := resolveIsExpectationMet(flags)
	notes := flags["notes"]
	if notes == "" {
		notes = flags["reason"]
	}

	updated, err := db.UpdateActionOutcome(actionId, expected, actual, isMet, notes)
	if err != nil {
		return err
	}

	statusStr := "MET"
	if !updated.IsExpectationMet {
		statusStr = "FAILED"
	}

	return renderHeatmapUpdateResponse(out, *updated, db.Path, statusStr, flags)
}

func isJsonRequested(flags map[string]string) bool {
	return flags["json"] == "true" || strings.EqualFold(flags["format"], "json")
}

func resolveActualOutcomeFlag(flags map[string]string) string {
	if a := flags["actual"]; a != "" {
		return a
	}

	if o := flags["outcome"]; o != "" && o != "success" && o != "failed" && o != "fail" {
		return o
	}

	if r := flags["result"]; r != "" {
		return r
	}

	if n := flags["notes"]; n != "" {
		return n
	}

	return flags["reason"]
}

func resolveIsExpectationMet(flags map[string]string) bool {
	if flags["failed"] == "true" || flags["fail"] == "true" {
		return false
	}

	if flags["success"] == "true" {
		return true
	}

	if strutil.EqualFoldAny(flags["status"], "fail", "failed") {
		return false
	}

	if strutil.EqualFoldAny(flags["status"], "success", "met", "passed") {
		return true
	}

	if flags["met"] == "false" {
		return false
	}

	if flags["met"] == "true" {
		return true
	}

	if strutil.EqualFoldAny(flags["outcome"], "fail", "failed") {
		return false
	}

	return true
}

func parseActionId(flags map[string]string, pos []string) int64 {
	if val := parsePositiveID(flags["id"]); val > 0 {
		return val
	}

	if val := parsePositiveID(flags["action-id"]); val > 0 {
		return val
	}

	if len(pos) > 0 {
		return parsePositiveID(pos[0])
	}

	return 0
}

func parsePositiveID(raw string) int64 {
	val, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || val <= 0 {
		return 0
	}

	return val
}
