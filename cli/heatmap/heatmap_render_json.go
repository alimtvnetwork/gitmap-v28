package heatmap

import (
	"encoding/json"
	"fmt"
	"io"

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
	hasDbPath := dbPath != ""
	if hasDbPath {
		resp["databasePath"] = computeRelativeDbPath(dbPath)
	}

	return resp
}

func populateLogRecordFields(resp map[string]any, rec store.HeatmapActionLogRecord) {
	if rec.TaskSlug != "" {
		resp["taskSlug"] = rec.TaskSlug
	}
	if rec.ActionType != "" {
		resp["actionType"] = rec.ActionType
	}
	if rec.FilePath != "" {
		resp["filePath"] = rec.FilePath
	}
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

func renderHeatmapLogResponse(out io.Writer, rec store.HeatmapActionLogRecord, dbPath string, flags map[string]string) error {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Logged heatmap action #%d [task=%s, action=%s, file=%s]\n",
			rec.HeatmapActionLogId, rec.TaskSlug, rec.ActionType, rec.FilePath)

		return nil
	}

	hasEnvelope := flags["envelope"] == "true"
	if hasEnvelope {
		env := BuildHeatmapLogEnvelope(rec, dbPath)

		return writeIndentedJSON(out, env)
	}

	resp := buildAgentLikeLogMap(rec, dbPath)

	return writeIndentedJSON(out, resp)
}

func renderRcHeatmapLogResponse(out io.Writer, rec store.HeatmapActionLogRecord, dbPath, aiName, repoSlug string, flags map[string]string) error {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Logged RC heatmap action #%d [ai=%s, repo=%s, task=%s, action=%s, file=%s]\n",
			rec.HeatmapActionLogId, aiName, repoSlug, rec.TaskSlug, rec.ActionType, rec.FilePath)

		return nil
	}

	hasEnvelope := flags["envelope"] == "true"
	if hasEnvelope {
		env := BuildRcLogEnvelope(rec, dbPath, aiName, repoSlug)

		return writeIndentedJSON(out, env)
	}

	resp := buildAgentLikeLogMap(rec, dbPath)
	resp["aiName"] = aiName
	resp["repoSlug"] = repoSlug

	return writeIndentedJSON(out, resp)
}

func buildAgentLikeUpdateMap(updated store.HeatmapActionLogRecord, statusStr string) map[string]any {
	return map[string]any{
		"actionLogId":       updated.HeatmapActionLogId,
		"status":            "UPDATED",
		"actualOutcome":     updated.ActualOutcome,
		"expectedOutcome":   updated.ExpectedOutcome,
		"isExpectationMet":  updated.IsExpectationMet,
		"expectationStatus": statusStr,
	}
}

func renderHeatmapUpdateResponse(out io.Writer, updated store.HeatmapActionLogRecord, dbPath, statusStr string, flags map[string]string) error {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Updated heatmap action #%d [outcome=%q, met=%v, status=%s]\n",
			updated.HeatmapActionLogId, updated.ActualOutcome, updated.IsExpectationMet, statusStr)

		return nil
	}

	hasEnvelope := flags["envelope"] == "true"
	if hasEnvelope {
		env := BuildHeatmapLogEnvelope(updated, dbPath)

		return writeIndentedJSON(out, env)
	}

	resp := buildAgentLikeUpdateMap(updated, statusStr)

	return writeIndentedJSON(out, resp)
}

func renderRcUpdateResponse(out io.Writer, updated store.HeatmapActionLogRecord, dbPath, aiName, repoSlug, statusStr string, flags map[string]string) error {
	hasText := flags["text"] == "true" || flags["plain"] == "true"
	if hasText {
		fmt.Fprintf(out, "Updated RC heatmap action #%d [outcome=%q, met=%v, status=%s]\n",
			updated.HeatmapActionLogId, updated.ActualOutcome, updated.IsExpectationMet, statusStr)

		return nil
	}

	hasEnvelope := flags["envelope"] == "true"
	if hasEnvelope {
		env := BuildRcLogEnvelope(updated, dbPath, aiName, repoSlug)

		return writeIndentedJSON(out, env)
	}

	resp := buildAgentLikeUpdateMap(updated, statusStr)
	resp["aiName"] = aiName
	resp["repoSlug"] = repoSlug

	return writeIndentedJSON(out, resp)
}
