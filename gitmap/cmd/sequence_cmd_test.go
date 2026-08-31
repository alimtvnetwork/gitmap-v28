package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScanDirectorySequence(t *testing.T) {
	tempDir := t.TempDir()

	_ = os.WriteFile(filepath.Join(tempDir, "01-first.txt"), []byte("one"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "02-second.txt"), []byte("two"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "unsequenced.txt"), []byte("three"), 0644)

	payload, err := scanDirectorySequence(tempDir)
	if err != nil {
		t.Fatalf("scanDirectorySequence failed: %v", err)
	}

	if payload.TotalFiles != 3 {
		t.Errorf("expected 3 total files, got %d", payload.TotalFiles)
	}
	if payload.SequencedFiles != 2 {
		t.Errorf("expected 2 sequenced files, got %d", payload.SequencedFiles)
	}
}

func TestApplySequenceOrderingWithPin(t *testing.T) {
	tempDir := t.TempDir()

	f1 := filepath.Join(tempDir, "index.md")
	f2 := filepath.Join(tempDir, "shared-engine.py")
	f3 := filepath.Join(tempDir, "file-manipulator.py")

	_ = os.WriteFile(f1, []byte("index"), 0644)
	_ = os.WriteFile(f2, []byte("shared"), 0644)
	_ = os.WriteFile(f3, []byte("manipulator"), 0644)

	entries, _ := os.ReadDir(tempDir)
	parsedFiles := parseSeqFiles(entries, tempDir)

	flags := SequenceFlags{
		StartNum: 1,
		PinMap: map[string]int{
			"index":         1,
			"shared-engine": 2,
		},
	}

	applySequenceOrdering(parsedFiles, flags)

	report := executeSequenceRenames(parsedFiles, tempDir, true)
	if len(report.Operations) != 3 {
		t.Errorf("expected 3 rename operations, got %d", len(report.Operations))
	}

	for _, op := range report.Operations {
		if op.From == "index.md" && op.To != "01-index.md" {
			t.Errorf("expected index.md -> 01-index.md, got %s", op.To)
		}
		if op.From == "shared-engine.py" && op.To != "02-shared-engine.py" {
			t.Errorf("expected shared-engine.py -> 02-shared-engine.py, got %s", op.To)
		}
		if op.From == "file-manipulator.py" && op.To != "03-file-manipulator.py" {
			t.Errorf("expected file-manipulator.py -> 03-file-manipulator.py, got %s", op.To)
		}
	}
}
