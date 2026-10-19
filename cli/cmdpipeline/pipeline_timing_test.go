package cmdpipeline

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestFormatExecutionDuration(t *testing.T) {
	testCases := []struct {
		d    time.Duration
		want string
	}{
		{0, "1ms"},
		{500 * time.Microsecond, "1ms"},
		{45 * time.Millisecond, "45ms"},
		{500 * time.Millisecond, "500ms"},
		{1 * time.Second, "1.00s"},
		{2150 * time.Millisecond, "2.15s"},
		{10500 * time.Millisecond, "10.50s"},
	}

	for _, tc := range testCases {
		got := formatExecutionDuration(tc.d)
		if got != tc.want {
			t.Errorf("formatExecutionDuration(%v) = %q; want %q", tc.d, got, tc.want)
		}
	}
}

func TestFormatTimingBanner_CacheHit(t *testing.T) {
	SetActiveCommandTagForTest("[gitmap te]")
	defer SetActiveCommandTagForTest("")

	banner := formatTimingBanner("[gitmap te]", 45*time.Millisecond, true)
	expected := "⚡ [gitmap te] Request completed in 45ms (from cache)"
	if banner != expected {
		t.Errorf("got %q; want %q", banner, expected)
	}
}

func TestFormatTimingBanner_LiveQuery(t *testing.T) {
	SetActiveCommandTagForTest("[gitmap te]")
	defer SetActiveCommandTagForTest("")

	banner := formatTimingBanner("[gitmap te]", 2150*time.Millisecond, false)
	expected := "⏱️ [gitmap te] Request completed in 2.15s"
	if banner != expected {
		t.Errorf("got %q; want %q", banner, expected)
	}
}

func TestFormatTimingBanner_PeCommand(t *testing.T) {
	SetActiveCommandTagForTest("[gitmap pe]")
	defer SetActiveCommandTagForTest("")

	banner := formatTimingBanner("[gitmap pe]", 45*time.Millisecond, true)
	expected := "⚡ [gitmap pe] Request completed in 45ms (from cache)"
	if banner != expected {
		t.Errorf("got %q; want %q", banner, expected)
	}
}

func TestFormatTimingBanner_FromPayload(t *testing.T) {
	SetActiveCommandTagForTest("[gitmap te]")
	defer SetActiveCommandTagForTest("")

	cachePayload := PipelineErrorLogsPayload{
		IsFromCache:     true,
		ElapsedSeconds:  0.045,
		ElapsedDuration: "45ms",
	}
	gotCache := FormatTimingBanner(cachePayload)
	wantCache := "⚡ [gitmap te] Request completed in 45ms (from cache)"
	if gotCache != wantCache {
		t.Errorf("FormatTimingBanner(cache) = %q; want %q", gotCache, wantCache)
	}

	freshPayload := PipelineErrorLogsPayload{
		IsFromCache:     false,
		ElapsedSeconds:  2.15,
		ElapsedDuration: "2.15s",
	}
	gotFresh := FormatTimingBanner(freshPayload)
	wantFresh := "⏱️ [gitmap te] Request completed in 2.15s"
	if gotFresh != wantFresh {
		t.Errorf("FormatTimingBanner(fresh) = %q; want %q", gotFresh, wantFresh)
	}
}

func TestAttachExecutionDuration(t *testing.T) {
	past := time.Now().Add(-100 * time.Millisecond)
	SetPipelineExecutionStartTime(past)

	var payload PipelineErrorLogsPayload
	attachExecutionDuration(&payload)

	if payload.ElapsedSeconds <= 0 {
		t.Errorf("expected ElapsedSeconds > 0, got %f", payload.ElapsedSeconds)
	}
	if len(payload.ElapsedDuration) == 0 {
		t.Errorf("expected non-empty ElapsedDuration")
	}
}

func TestJSONOutputIncludesDurationFields(t *testing.T) {
	payload := PipelineErrorLogsPayload{
		Repo:            "test-org/test-repo",
		ElapsedDuration: "45ms",
		ElapsedSeconds:  0.045,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("json.Marshal failed: %v", err)
	}

	jsonStr := string(data)
	if !strings.Contains(jsonStr, `"elapsedDuration":"45ms"`) {
		t.Errorf("expected json to contain elapsedDuration, got: %s", jsonStr)
	}
	if !strings.Contains(jsonStr, `"elapsedSeconds":0.045`) {
		t.Errorf("expected json to contain elapsedSeconds, got: %s", jsonStr)
	}
}

func TestMatchArgCommandTag(t *testing.T) {
	tagCases := []struct {
		arg  string
		want string
	}{
		{"te", "[gitmap te]"},
		{"TE", "[gitmap te]"},
		{"pe", "[gitmap pe]"},
		{"PE", "[gitmap pe]"},
		{"pipeline-errors", "[gitmap pe]"},
		{"pipeline_errors", "[gitmap pe]"},
		{"ee", "[gitmap pe]"},
		{"status", ""},
		{"--json", ""},
	}

	for _, tc := range tagCases {
		got := matchArgCommandTag(tc.arg)
		if got != tc.want {
			t.Errorf("matchArgCommandTag(%q) = %q; want %q", tc.arg, got, tc.want)
		}
	}
}
