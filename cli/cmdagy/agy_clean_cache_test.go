// Package cmdagy — agy_clean_cache_test.go provides unit tests for agy clean-cache.
package cmdagy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{512, "512 B"},
		{1024, "1.0 KB"},
		{1536, "1.5 KB"},
		{1048576, "1.0 MB"},
		{1073741824, "1.0 GB"},
	}

	for _, tt := range tests {
		got := FormatBytes(tt.bytes)
		if got != tt.expected {
			t.Errorf("FormatBytes(%d) = %q; want %q", tt.bytes, got, tt.expected)
		}
	}
}

func TestCalculateDirStatsAndClean(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agy-clean-cache-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	subDir := filepath.Join(tempDir, "subdir")
	_ = os.MkdirAll(subDir, 0755)

	file1 := filepath.Join(tempDir, "file1.txt")
	file2 := filepath.Join(subDir, "file2.txt")
	_ = os.WriteFile(file1, []byte("hello world"), 0644)
	_ = os.WriteFile(file2, []byte("test content 12345"), 0644)

	size, count, statErr := CalculateDirStats(tempDir)
	if statErr != nil {
		t.Fatalf("CalculateDirStats error: %v", statErr)
	}

	expectedFiles := 2
	expectedSize := int64(len("hello world") + len("test content 12345"))
	if count != expectedFiles || size != expectedSize {
		t.Errorf("got (%d bytes, %d files), want (%d bytes, %d files)", size, count, expectedSize, expectedFiles)
	}

	freed, deleted, warnings := CleanDirectoryContents(tempDir)
	if len(warnings) > 0 {
		t.Errorf("unexpected warnings during clean: %v", warnings)
	}

	if freed != expectedSize || deleted != expectedFiles {
		t.Errorf("got cleaned (%d bytes, %d files), want (%d bytes, %d files)", freed, deleted, expectedSize, expectedFiles)
	}

	entries, _ := os.ReadDir(tempDir)
	if len(entries) != 0 {
		t.Errorf("expected tempDir to be empty after clean, got %d entries", len(entries))
	}
}

func TestParseTasklistCSV(t *testing.T) {
	mockCSV := `"Image Name","PID","Session Name","Session#","Mem Usage"
"Antigravity.exe","1001","Console","1","132,924 K"
"msedgewebview2.exe","1002","Console","1","20,328 K"
"chrome.exe","1003","Console","1","92,976 K"
"electron.exe","1004","Console","1","42,028 K"
"msedge.exe","1005","Console","1","20,492 K"
"powershell.exe","1006","Console","1","109,828 K"
"Antigravity.exe","9999","Console","1","80,216 K"
`
	procs := parseTasklistCSV(mockCSV, 9999) // ignore PID 9999 (simulating current process)

	if len(procs) != 4 {
		t.Fatalf("expected 4 matched procs, got %d: %+v", len(procs), procs)
	}

	expectedPids := []int{1001, 1002, 1004, 1005}
	for i, expectedPid := range expectedPids {
		if procs[i].PID != expectedPid {
			t.Errorf("procs[%d].PID = %d; want %d", i, procs[i].PID, expectedPid)
		}
	}
}

func TestParsePsOutput(t *testing.T) {
	mockPs := `  PID COMMAND
 2001 antigravity
 2002 msedgewebview2
 2003 bash
 2004 electron
 2005 msedge
 2006 python3
 8888 antigravity
`
	procs := parsePsOutput(mockPs, 8888) // ignore PID 8888 (current process)

	if len(procs) != 4 {
		t.Fatalf("expected 4 matched procs, got %d: %+v", len(procs), procs)
	}

	expectedPids := []int{2001, 2002, 2004, 2005}
	for i, expectedPid := range expectedPids {
		if procs[i].PID != expectedPid {
			t.Errorf("procs[%d].PID = %d; want %d", i, procs[i].PID, expectedPid)
		}
	}
}

func TestNormalizeAgySubcommand_CleanCache(t *testing.T) {
	aliases := []string{"clean-cache", "CLEAN-CACHE", "cleancache", "CleanCache", "clean_cache", "cc", "CC"}
	for _, alias := range aliases {
		normalized := normalizeAgySubcommand(alias)
		if normalized != "clean-cache" {
			t.Errorf("normalizeAgySubcommand(%q) = %q; want 'clean-cache'", alias, normalized)
		}
	}
}

func TestExecuteCleanCache_DryRunJSON(t *testing.T) {
	opts := CleanCacheOptions{
		DryRun:      true,
		JSON:        true,
		NoKill:      true,
		IncludeTemp: false,
	}

	err := ExecuteCleanCache(opts)
	if err != nil {
		t.Fatalf("ExecuteCleanCache dry-run failed: %v", err)
	}
}
