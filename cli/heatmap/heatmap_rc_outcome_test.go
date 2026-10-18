package heatmap

import (
	"encoding/json"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestBuildRcLogEnvelope(t *testing.T) {
	rec := store.HeatmapActionLogRecord{
		HeatmapActionLogId: 10,
		TaskSlug:           "unit-test-task",
		RepoName:           "gitmap-test",
		ActionType:         "EDIT",
		FilePath:           "test_file.go",
		ExpectedOutcome:    "File edited cleanly",
		ActualOutcome:      "All assertions passed",
		IsExpectationMet:   true,
	}

	env := BuildRcLogEnvelope(rec, "D:/work/repo-cache/.heatmap/01-antigravity/gitmap-test/heatmap.db", "antigravity", "gitmap-test")

	if env.Status != "success" {
		t.Fatalf("expected status 'success', got: %s", env.Status)
	}

	if env.ActionId != 10 {
		t.Fatalf("expected action_id 10, got: %d", env.ActionId)
	}

	if env.WhereSaved.DbPath == "" {
		t.Fatal("expected non-empty where_saved.db_path")
	}

	if env.DbInfo.Driver != "sqlite" || env.DbInfo.Table != "HeatmapActionLog" {
		t.Fatalf("unexpected db_info: %+v", env.DbInfo)
	}

	if env.WhatSaved.ActionType != "EDIT" || env.WhatSaved.FilePath != "test_file.go" {
		t.Fatalf("unexpected what_saved: %+v", env.WhatSaved)
	}

	if env.AiExpectation.ExpectedOutcome != "File edited cleanly" {
		t.Fatalf("unexpected expected_outcome: %s", env.AiExpectation.ExpectedOutcome)
	}

	if env.AiExpectation.Status != "MET" {
		t.Fatalf("expected status 'MET', got: %s", env.AiExpectation.Status)
	}

	data, err := json.Marshal(env)
	if err != nil {
		t.Fatalf("failed to marshal envelope: %v", err)
	}

	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	for _, key := range []string{"status", "action_id", "where_saved", "db_info", "what_saved", "ai_expectation"} {
		if _, ok := raw[key]; !ok {
			t.Fatalf("missing key %q in JSON output", key)
		}
	}
}

func TestComputeExpectationStatus(t *testing.T) {
	if s := computeExpectationStatus("plan A", "", true); s != "PENDING" {
		t.Fatalf("expected PENDING, got: %s", s)
	}

	if s := computeExpectationStatus("plan A", "done", true); s != "MET" {
		t.Fatalf("expected MET, got: %s", s)
	}

	if s := computeExpectationStatus("plan A", "failed with error", false); s != "FAILED" {
		t.Fatalf("expected FAILED, got: %s", s)
	}

	if s := computeExpectationStatus("", "", false); s != "FAILED" {
		t.Fatalf("expected FAILED when isMet is false, got: %s", s)
	}
}

func TestResolveIsExpectationMet(t *testing.T) {
	cases := []struct {
		flags    map[string]string
		expected bool
	}{
		{map[string]string{"failed": "true"}, false},
		{map[string]string{"fail": "true"}, false},
		{map[string]string{"success": "true"}, true},
		{map[string]string{"status": "failed"}, false},
		{map[string]string{"status": "success"}, true},
		{map[string]string{"status": "met"}, true},
		{map[string]string{"met": "false"}, false},
		{map[string]string{"met": "true"}, true},
		{map[string]string{"outcome": "failed"}, false},
		{map[string]string{}, true},
	}

	for _, tc := range cases {
		actual := resolveIsExpectationMet(tc.flags)
		if actual != tc.expected {
			t.Fatalf("for flags %+v expected %v, got %v", tc.flags, tc.expected, actual)
		}
	}
}

func TestParseActionId(t *testing.T) {
	if id := parseActionId(map[string]string{"id": "42"}, nil); id != 42 {
		t.Fatalf("expected 42, got %d", id)
	}

	if id := parseActionId(map[string]string{"action-id": "99"}, nil); id != 99 {
		t.Fatalf("expected 99, got %d", id)
	}

	if id := parseActionId(map[string]string{}, []string{"77"}); id != 77 {
		t.Fatalf("expected 77, got %d", id)
	}
}
