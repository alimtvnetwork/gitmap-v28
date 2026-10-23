package heatmap

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func writeIndentedJSON(out io.Writer, payload any) error {
	b, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}

	_, writeErr := fmt.Fprintln(out, string(b))

	return writeErr
}

func buildAgentLikeLogMap(rec store.HeatmapActionLogRecord, dbPath string) map[string]any {
	resp := map[string]any{
		"actionLogId": rec.HeatmapActionLogId,
		"status":      "LOGGED",
	}

	populateLogRecordFields(resp, rec)
	if dbPath != "" {
		resp["databasePath"] = computeRelativeDbPath(dbPath)
	}

	resp["files"] = buildHeatmapFilesSection(rec)

	return resp
}

func buildHeatmapFilesSection(rec store.HeatmapActionLogRecord) fsutil.ClassifiedFiles {
	cwd, _ := os.Getwd()

	return fsutil.ClassifyFileTargets(rec.FilePath, rec.ActionType, cwd)
}

func populateLogRecordIdentity(resp map[string]any, rec store.HeatmapActionLogRecord) {
	if rec.TaskSlug != "" {
		resp["taskSlug"] = rec.TaskSlug
	}

	if rec.ActionType != "" {
		resp["actionType"] = rec.ActionType
	}

	if rec.FilePath != "" {
		resp["filePath"] = rec.FilePath
	}
}

func populateLogRecordReasoning(resp map[string]any, rec store.HeatmapActionLogRecord) {
	if rec.ReasoningState != "" {
		resp["reasoningState"] = rec.ReasoningState
	}

	if rec.ReasoningNotes != "" {
		resp["reasoningNotes"] = rec.ReasoningNotes
	}

	if rec.ExpectedOutcome != "" {
		resp["expectedOutcome"] = rec.ExpectedOutcome
	}
}

func populateLogRecordAttribution(resp map[string]any, rec store.HeatmapActionLogRecord) {
	if rec.ModelName != "" {
		resp["modelName"] = rec.ModelName
	}

	if rec.ToolName != "" {
		resp["toolName"] = rec.ToolName
	}
}

func populateLogRecordFields(resp map[string]any, rec store.HeatmapActionLogRecord) {
	populateLogRecordIdentity(resp, rec)
	populateLogRecordReasoning(resp, rec)
	populateLogRecordAttribution(resp, rec)
}

func renderLogTextIfRequested(out io.Writer, rec store.HeatmapActionLogRecord, flags map[string]string) bool {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Logged heatmap action #%d [task=%s, action=%s, file=%s]\n",
			rec.HeatmapActionLogId, rec.TaskSlug, rec.ActionType, rec.FilePath)
	}

	return hasText
}

func renderHeatmapLogResponse(out io.Writer, rec store.HeatmapActionLogRecord, dbPath string, flags map[string]string) error {
	if renderLogTextIfRequested(out, rec, flags) {
		return nil
	}

	if flags["envelope"] == "true" {
		return writeIndentedJSON(out, BuildHeatmapLogEnvelope(rec, dbPath))
	}

	resp := buildAgentLikeLogMap(rec, dbPath)

	return writeIndentedJSON(out, resp)
}

func renderRcLogTextIfRequested(out io.Writer, rec store.HeatmapActionLogRecord, aiName, repoSlug string, flags map[string]string) bool {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Logged RC heatmap action #%d [ai=%s, repo=%s, task=%s, action=%s, file=%s]\n",
			rec.HeatmapActionLogId, aiName, repoSlug, rec.TaskSlug, rec.ActionType, rec.FilePath)
	}

	return hasText
}

func renderRcHeatmapLogResponse(out io.Writer, rec store.HeatmapActionLogRecord, dbPath, aiName, repoSlug string, flags map[string]string) error {
	if renderRcLogTextIfRequested(out, rec, aiName, repoSlug, flags) {
		return nil
	}

	if flags["envelope"] == "true" {
		return writeIndentedJSON(out, BuildRcLogEnvelope(rec, dbPath, aiName, repoSlug))
	}

	resp := buildAgentLikeLogMap(rec, dbPath)
	resp["aiName"] = aiName
	resp["repoSlug"] = repoSlug

	return writeIndentedJSON(out, resp)
}

func initUpdateResponseMap(updated store.HeatmapActionLogRecord, statusStr string) map[string]any {
	return map[string]any{
		"actionLogId":       updated.HeatmapActionLogId,
		"status":            "UPDATED",
		"actualOutcome":     updated.ActualOutcome,
		"expectedOutcome":   updated.ExpectedOutcome,
		"isExpectationMet":  updated.IsExpectationMet,
		"expectationStatus": statusStr,
	}
}

func buildAgentLikeUpdateMap(updated store.HeatmapActionLogRecord, dbPath, statusStr string) map[string]any {
	resp := initUpdateResponseMap(updated, statusStr)
	populateLogRecordAttribution(resp, updated)
	if dbPath != "" {
		resp["databasePath"] = computeRelativeDbPath(dbPath)
	}

	resp["files"] = buildHeatmapFilesSection(updated)

	return resp
}

func renderUpdateTextIfRequested(out io.Writer, updated store.HeatmapActionLogRecord, statusStr string, flags map[string]string) bool {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Updated heatmap action #%d [outcome=%q, met=%v, status=%s]\n",
			updated.HeatmapActionLogId, updated.ActualOutcome, updated.IsExpectationMet, statusStr)
	}

	return hasText
}

func renderHeatmapUpdateResponse(out io.Writer, updated store.HeatmapActionLogRecord, dbPath, statusStr string, flags map[string]string) error {
	if renderUpdateTextIfRequested(out, updated, statusStr, flags) {
		return nil
	}

	if flags["envelope"] == "true" {
		return writeIndentedJSON(out, BuildHeatmapLogEnvelope(updated, dbPath))
	}

	resp := buildAgentLikeUpdateMap(updated, dbPath, statusStr)

	return writeIndentedJSON(out, resp)
}

func renderRcUpdateResponse(out io.Writer, updated store.HeatmapActionLogRecord, dbPath, aiName, repoSlug, statusStr string, flags map[string]string) error {
	if renderUpdateTextIfRequested(out, updated, statusStr, flags) {
		return nil
	}

	if flags["envelope"] == "true" {
		return writeIndentedJSON(out, BuildRcLogEnvelope(updated, dbPath, aiName, repoSlug))
	}

	resp := buildAgentLikeUpdateMap(updated, dbPath, statusStr)
	resp["aiName"] = aiName
	resp["repoSlug"] = repoSlug

	return writeIndentedJSON(out, resp)
}
