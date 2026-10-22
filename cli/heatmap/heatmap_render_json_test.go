package heatmap

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func sampleTestRecord() store.HeatmapActionLogRecord {
	return store.HeatmapActionLogRecord{
		HeatmapActionLogId: 44,
		TaskSlug:           "task-spec-creation",
		RepoName:           "gitmap",
		ActionType:         "READ",
		FilePath:           "02-spec/21-app/spec.md",
		ReasoningState:     "EXPLORATION",
		ReasoningNotes:     "Reading spec",
		ExpectedOutcome:    "Inspect spec",
		ActualOutcome:      "Spec inspected",
		IsExpectationMet:   true,
	}
}

func TestRenderHeatmapLogResponse_Default(t *testing.T) {
	rec := sampleTestRecord()
	var buf bytes.Buffer
	flags := map[string]string{}
	err := renderHeatmapLogResponse(&buf, rec, "data/heatmap.db", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v, raw: %s", unmarshalErr, buf.String())
	}

	if parsed["status"] != "LOGGED" || parsed["actionLogId"] != float64(44) {
		t.Fatalf("unexpected json fields: %+v", parsed)
	}

	if parsed["taskSlug"] != "task-spec-creation" || parsed["actionType"] != "READ" {
		t.Fatalf("missing action details in json: %+v", parsed)
	}
}

func TestRenderHeatmapLogResponse_Text(t *testing.T) {
	rec := sampleTestRecord()
	var buf bytes.Buffer
	flags := map[string]string{"text": "true"}
	err := renderHeatmapLogResponse(&buf, rec, "data/heatmap.db", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	outStr := buf.String()
	if !strings.HasPrefix(outStr, "Logged heatmap action #44") {
		t.Fatalf("expected text output prefix, got: %s", outStr)
	}
}

func TestRenderHeatmapLogResponse_Envelope(t *testing.T) {
	rec := sampleTestRecord()
	var buf bytes.Buffer
	flags := map[string]string{"envelope": "true"}
	err := renderHeatmapLogResponse(&buf, rec, "data/heatmap.db", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v", unmarshalErr)
	}

	if parsed["status"] != "success" || parsed["where_saved"] == nil {
		t.Fatalf("expected envelope structure, got: %+v", parsed)
	}
}

func TestRenderRcHeatmapLogResponse_Default(t *testing.T) {
	rec := sampleTestRecord()
	var buf bytes.Buffer
	flags := map[string]string{}
	err := renderRcHeatmapLogResponse(&buf, rec, "data/rc.db", "antigravity", "gitmap", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v", unmarshalErr)
	}

	if parsed["status"] != "LOGGED" || parsed["aiName"] != "antigravity" {
		t.Fatalf("unexpected fields: %+v", parsed)
	}
}

func TestRenderHeatmapUpdateResponse_Default(t *testing.T) {
	rec := sampleTestRecord()
	var buf bytes.Buffer
	flags := map[string]string{}
	err := renderHeatmapUpdateResponse(&buf, rec, "data/heatmap.db", "MET", flags)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var parsed map[string]any
	if unmarshalErr := json.Unmarshal(buf.Bytes(), &parsed); unmarshalErr != nil {
		t.Fatalf("output is not valid JSON: %v", unmarshalErr)
	}

	if parsed["status"] != "UPDATED" || parsed["expectationStatus"] != "MET" {
		t.Fatalf("unexpected update fields: %+v", parsed)
	}

	if parsed["databasePath"] == nil || parsed["databasePath"] == "" {
		t.Fatalf("expected non-empty databasePath, got %+v", parsed["databasePath"])
	}

	filesRaw, hasFiles := parsed["files"].(map[string]any)
	if !hasFiles {
		t.Fatalf("expected files map in response, got %+v", parsed["files"])
	}

	readList, isList := filesRaw["read"].([]any)
	if !isList || len(readList) != 1 || readList[0] != "02-spec/21-app/spec.md" {
		t.Fatalf("expected read file in files.read, got %+v", filesRaw["read"])
	}
}
