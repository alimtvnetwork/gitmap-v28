package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsBinaryExtension(t *testing.T) {
	bins := []string{".exe", ".dll", ".zip", ".tar.gz", ".png", ".jpg", ".wasm", ".db"}
	for _, b := range bins {
		if !IsBinaryExtension(b) && !IsBinaryExtension(filepath.Ext(b)) {
			t.Errorf("expected %s to be recognized as binary extension", b)
		}
	}
	if IsBinaryExtension(".go") || IsBinaryExtension(".ts") || IsBinaryExtension(".md") {
		t.Error("code extensions should not be flagged as binary")
	}
}

func TestIsAllowedLargeWaiver(t *testing.T) {
	if !IsAllowedLargeWaiver("src/data/specTree.json") {
		t.Error("expected specTree.json to have waiver")
	}
	if !IsAllowedLargeWaiver(".ai-memory/test-inventory.json") {
		t.Error("expected test-inventory.json to have waiver")
	}
	if IsAllowedLargeWaiver("random/big.json") {
		t.Error("random file should not have waiver")
	}
}

func TestIsPathExcluded(t *testing.T) {
	excs := []string{"tmp/cache", "huge.bin"}
	if !IsPathExcluded("tmp/cache/item.txt", excs) {
		t.Error("expected path in tmp/cache to be excluded")
	}
	if !IsPathExcluded("huge.bin", excs) {
		t.Error("expected huge.bin to be excluded")
	}
	if IsPathExcluded("src/index.ts", excs) {
		t.Error("src/index.ts should not be excluded")
	}
}

func TestRunSequenceAuditor_GapsAndFixes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "gitmap-seq-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sub := filepath.Join(tmpDir, "docs")
	_ = os.MkdirAll(sub, 0755)

	f1 := filepath.Join(sub, "01-intro.md")
	f2 := filepath.Join(sub, "03-missing-two.md")
	_ = os.WriteFile(f1, []byte("# 01 Introduction\n\nSome text.\n"), 0644)
	_ = os.WriteFile(f2, []byte("# 02 Wrong Title\n\nSome text.\n"), 0644)

	res, appErr := RunSequenceAuditor(SequenceAuditorOptions{Dir: tmpDir, IsFixMode: false})
	if appErr != nil {
		t.Fatalf("RunSequenceAuditor error: %v", appErr)
	}

	if res.TotalGaps == 0 {
		t.Error("expected gap between 01 and 03 to be detected")
	}
	if res.TotalMismatches == 0 {
		t.Error("expected title mismatch (03 prefix vs # 02) to be detected")
	}

	// Now run in fix mode
	resFix, fixErr := RunSequenceAuditor(SequenceAuditorOptions{Dir: tmpDir, IsFixMode: true})
	if fixErr != nil {
		t.Fatalf("RunSequenceAuditor fix error: %v", fixErr)
	}
	if resFix.TotalFixed == 0 {
		t.Error("expected mismatched title to be fixed")
	}

	fixedData, _ := os.ReadFile(f2)
	if !containsStr(string(fixedData), "# 03 Wrong Title") {
		t.Errorf("expected H1 to be updated to # 03, got: %s", string(fixedData))
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || filepath.ToSlash(s) != "" && len(s) > 0 && len(substr) > 0 && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || (len(s) > len(substr) && stringContains(s, substr))))
}

func stringContains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
