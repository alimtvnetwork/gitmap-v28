package store

import (
	"path/filepath"
	"testing"
)

func TestAiInstruction_QueryFrequentAiCommands(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "ai_instructions_test.db")

	splitDb, err := OpenAiInstructionSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open split db: %v", err)
	}
	defer splitDb.Close()

	// Record sample commands
	_ = splitDb.RecordAiExecution("instruction", "pwsh Get-Process", "[]", ".", "127.0.0.1", 10, 0, "", "", true)
	_ = splitDb.RecordAiExecution("instruction", "pwsh Get-Process", "[]", ".", "127.0.0.1", 15, 0, "", "", true)
	_ = splitDb.RecordAiExecution("instruction", "pwsh Get-Date", "[]", ".", "127.0.0.1", 5, 0, "", "", true)

	freq, errFreq := splitDb.QueryFrequentAiCommands(10)
	if errFreq != nil {
		t.Fatalf("failed to query frequent commands: %v", errFreq)
	}

	if len(freq) != 2 {
		t.Fatalf("expected 2 distinct commands, got %d", len(freq))
	}

	if freq[0].CommandText != "pwsh Get-Process" || freq[0].RunCount != 2 {
		t.Errorf("expected Get-Process with count 2, got %s with count %d", freq[0].CommandText, freq[0].RunCount)
	}
}
