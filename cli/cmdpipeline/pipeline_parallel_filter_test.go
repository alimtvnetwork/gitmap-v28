package cmdpipeline

import (
	"fmt"
	"testing"
)

func TestParallelFilterLines_RemovesOkAndPassLinesOnLargeSlice(t *testing.T) {
	lines := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		switch i % 3 {
		case 0:
			lines[i] = fmt.Sprintf("PASS: step #%d passed successfully", i)
		case 1:
			lines[i] = fmt.Sprintf("✔ ok: task %d finished", i)
		default:
			lines[i] = fmt.Sprintf("Error at line %d: syntax error", i)
		}
	}

	filtered := ParallelFilterLines(lines, isKeepLogLine)
	expectedCount := 333 // 1000/3 rounded down for remainder == 2
	if len(filtered) != expectedCount {
		t.Fatalf("expected %d filtered lines, got %d", expectedCount, len(filtered))
	}

	for _, line := range filtered {
		if isOkLogLine(line) {
			t.Errorf("expected no ok lines in filtered result, found: %s", line)
		}
	}
}

func TestParallelFilterLines_HandlesEmptyAndSmallSlices(t *testing.T) {
	empty := ParallelFilterLines([]string{}, isKeepLogLine)
	if len(empty) != 0 {
		t.Fatalf("expected 0 lines for empty slice, got %d", len(empty))
	}

	small := []string{"PASS", "error line 1", "ok", "error line 2"}
	filteredSmall := ParallelFilterLines(small, isKeepLogLine)
	if len(filteredSmall) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(filteredSmall))
	}
	if filteredSmall[0] != "error line 1" || filteredSmall[1] != "error line 2" {
		t.Fatalf("unexpected contents: %v", filteredSmall)
	}
}

func TestSliceLineMarker_ComputesExactCapacity(t *testing.T) {
	lines := []string{"ok", "error 1", "ok", "error 2", "error 3"}
	marker := NewSliceLineMarker(len(lines), 1)
	marker.MarkWorkerPartition(lines, 0, len(lines), 0, isKeepLogLine)

	if marker.TotalKeptCount() != 3 {
		t.Fatalf("expected kept count 3, got %d", marker.TotalKeptCount())
	}

	materialized := marker.MaterializeKeptLines(lines)
	if len(materialized) != 3 {
		t.Fatalf("expected 3 lines materialized, got %d", len(materialized))
	}
	if materialized[0] != "error 1" || materialized[1] != "error 2" || materialized[2] != "error 3" {
		t.Fatalf("unexpected materialized items: %v", materialized)
	}
}
