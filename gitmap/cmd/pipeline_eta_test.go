package cmd

import (
	"strings"
	"testing"
	"time"
)

func TestHistoricalSuccessETAExcludesFailedRuns(t *testing.T) {
	now := time.Now().UTC()
	created30sAgo := now.Add(-30 * time.Second).Format(time.RFC3339)

	runs := []ghRunItem{
		// Active running workflow
		{
			DatabaseId: 100,
			Name:       "CI",
			Status:     "in_progress",
			Conclusion: "",
			CreatedAt:  created30sAgo,
			UpdatedAt:  created30sAgo,
		},
		// Past successful run 1: 120s duration
		{
			DatabaseId: 101,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.Add(-200 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-80 * time.Second).Format(time.RFC3339),
		},
		// Past successful run 2: 120s duration
		{
			DatabaseId: 102,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "success",
			CreatedAt:  now.Add(-400 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-280 * time.Second).Format(time.RFC3339),
		},
		// Past failed run: should be IGNORED (failed after 10s)
		{
			DatabaseId: 103,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "failure",
			CreatedAt:  now.Add(-600 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-590 * time.Second).Format(time.RFC3339),
		},
		// Past canceled run: should be IGNORED (canceled after 15s)
		{
			DatabaseId: 104,
			Name:       "CI",
			Status:     "completed",
			Conclusion: "canceled",
			CreatedAt:  now.Add(-800 * time.Second).Format(time.RFC3339),
			UpdatedAt:  now.Add(-785 * time.Second).Format(time.RFC3339),
		},
	}

	eta := calculateETA(runs)
	// Average duration of success runs is (120 + 120) / 2 = 120s.
	// Elapsed is 30s.
	// Remaining ETA should be 120 - 30 = 90s.
	if eta < 85 || eta > 95 {
		t.Fatalf("expected ETA around 90s, got %d (failed runs were not properly excluded!)", eta)
	}
}

func TestExtractCleanErrorLines(t *testing.T) {
	sampleLog := `
2026-09-03T05:35:10Z Step 1: Set up Go
2026-09-03T05:35:11Z go version go1.24.0 windows/amd64
2026-09-03T05:35:20Z Step 2: Run tests
=== RUN   TestSample
--- PASS: TestSample (0.00s)
=== RUN   TestFailureCase
--- FAIL: TestFailureCase (0.01s)
    main_test.go:42: unexpected error code: got 404, want 200
2026-09-03T05:35:25Z ##[error]Process completed with exit code 1.
2026-09-03T05:35:26Z Post job cleanup.
`

	cleaned := extractCleanErrorLines(sampleLog)
	if !strings.Contains(cleaned, "--- FAIL: TestFailureCase") {
		t.Errorf("expected cleaned log to contain test failure, got:\n%s", cleaned)
	}
	if !strings.Contains(cleaned, "Process completed with exit code 1.") {
		t.Errorf("expected cleaned log to contain exit code error, got:\n%s", cleaned)
	}
	if strings.Contains(cleaned, "go version go1.24.0") {
		t.Errorf("cleaned error log should NOT contain setup lines, got:\n%s", cleaned)
	}
}
