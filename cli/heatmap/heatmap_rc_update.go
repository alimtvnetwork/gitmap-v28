package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
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

	if isJsonRequested(flags) {
		env := BuildRcLogEnvelope(*updated, st.RepoDbPath(), aiName, repoSlug)
		return json.NewEncoder(out).Encode(env)
	}

	statusStr := "MET"
	if !updated.IsExpectationMet {
		statusStr = "FAILED"
	}

	fmt.Fprintf(out, "Updated RC heatmap action #%d [outcome=%q, met=%v, status=%s]\n",
		updated.HeatmapActionLogId, updated.ActualOutcome, updated.IsExpectationMet, statusStr)
	return nil
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

	if isJsonRequested(flags) {
		env := BuildHeatmapLogEnvelope(*updated, db.Path)
		return json.NewEncoder(out).Encode(env)
	}

	statusStr := "MET"
	if !updated.IsExpectationMet {
		statusStr = "FAILED"
	}

	fmt.Fprintf(out, "Updated heatmap action #%d [outcome=%q, met=%v, status=%s]\n",
		updated.HeatmapActionLogId, updated.ActualOutcome, updated.IsExpectationMet, statusStr)
	return nil
}

func isJsonRequested(flags map[string]string) bool {
	return flags["json"] == "true" || strings.ToLower(flags["format"]) == "json"
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

	if s := strings.ToLower(flags["status"]); s == "fail" || s == "failed" {
		return false
	}

	if s := strings.ToLower(flags["status"]); s == "success" || s == "met" || s == "passed" {
		return true
	}

	if flags["met"] == "false" {
		return false
	}

	if flags["met"] == "true" {
		return true
	}

	if o := strings.ToLower(flags["outcome"]); o == "fail" || o == "failed" {
		return false
	}

	return true
}

func parseActionId(flags map[string]string, pos []string) int64 {
	if idStr := flags["id"]; idStr != "" {
		if val, err := strconv.ParseInt(idStr, 10, 64); err == nil && val > 0 {
			return val
		}
	}

	if idStr := flags["action-id"]; idStr != "" {
		if val, err := strconv.ParseInt(idStr, 10, 64); err == nil && val > 0 {
			return val
		}
	}

	if len(pos) > 0 {
		if val, err := strconv.ParseInt(pos[0], 10, 64); err == nil && val > 0 {
			return val
		}
	}

	return 0
}
