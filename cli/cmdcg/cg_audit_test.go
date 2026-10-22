package cmdcg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCalculateAuditScoreEmpty(t *testing.T) {
	score := calculateAuditScore(nil)
	if score != 0 {
		t.Fatalf("expected 0 for empty checks, got %d", score)
	}
}

func TestCalculateAuditScoreVariations(t *testing.T) {
	checks := []CGAuditCheck{{IsPass: true}, {IsPass: true}, {IsPass: false}}
	score := calculateAuditScore(checks)
	if score != 66 {
		t.Fatalf("expected score 66, got %d", score)
	}
}

func TestAuditSingleRepoEmptyDir(t *testing.T) {
	tempDir := t.TempDir()
	summary := auditSingleRepo(tempDir)
	if summary.Score != 0 {
		t.Fatalf("expected 0 score for empty repo, got %d", summary.Score)
	}

	if len(summary.Checks) != 5 {
		t.Fatalf("expected 5 checks, got %d", len(summary.Checks))
	}
}

func setupCompliantRepo(dir string) {
	vData := `{"coding-guidelines":{"version":"v2.24.0","status":"active"}}`
	_ = os.WriteFile(filepath.Join(dir, "version.json"), []byte(vData), 0644)
	_ = os.MkdirAll(filepath.Join(dir, "02-spec", "02-coding-guidelines"), 0755)
	_ = os.MkdirAll(filepath.Join(dir, ".ai-memory"), 0755)
	_ = os.WriteFile(filepath.Join(dir, ".ai-memory", "strictly-avoid.md"), []byte("# Rules"), 0644)
	_ = os.MkdirAll(filepath.Join(dir, ".agents", "rules"), 0755)
}

func TestAuditSingleRepoCompliant(t *testing.T) {
	tempDir := t.TempDir()
	setupCompliantRepo(tempDir)

	summary := auditSingleRepo(tempDir)
	if summary.Score != 100 {
		t.Fatalf("expected score 100 for compliant repo, got %d", summary.Score)
	}

	assertAllChecksPassed(t, summary.Checks)
}

func assertAllChecksPassed(t *testing.T, checks []CGAuditCheck) {
	for _, check := range checks {
		if !check.IsPass {
			t.Fatalf("expected check %s to pass, got detail: %s", check.Name, check.Detail)
		}
	}
}

func TestAuditCheckVersionJSONMissing(t *testing.T) {
	tempDir := t.TempDir()
	cMissing := auditCheckVersionJSON(tempDir)
	if cMissing.IsPass {
		t.Fatalf("expected version.json check to fail when missing")
	}
}

func TestAuditCheckVersionJSONValid(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "version.json"), []byte(`{"name":"test"}`), 0644)
	cValid := auditCheckVersionJSON(tempDir)
	if !cValid.IsPass {
		t.Fatalf("expected version.json check to pass when valid")
	}
}

func TestAuditCheckCGMetadata(t *testing.T) {
	tempDir := t.TempDir()
	_ = os.WriteFile(filepath.Join(tempDir, "version.json"), []byte(`{}`), 0644)
	cMissing := auditCheckCGMetadata(tempDir)
	if cMissing.IsPass {
		t.Fatalf("expected cg metadata check to fail when missing")
	}

	vData := `{"coding-guidelines":{"version":"v2.24.0"}}`
	_ = os.WriteFile(filepath.Join(tempDir, "version.json"), []byte(vData), 0644)
	cValid := auditCheckCGMetadata(tempDir)
	if !cValid.IsPass {
		t.Fatalf("expected cg metadata check to pass when present")
	}
}

func TestAuditCheckSpecDir(t *testing.T) {
	tempDir := t.TempDir()
	cMissing := auditCheckSpecDir(tempDir)
	if cMissing.IsPass {
		t.Fatalf("expected spec dir check to fail when missing")
	}

	_ = os.MkdirAll(filepath.Join(tempDir, "02-spec", "02-coding-guidelines"), 0755)
	cValid := auditCheckSpecDir(tempDir)
	if !cValid.IsPass {
		t.Fatalf("expected spec dir check to pass when present")
	}
}

func TestAuditCheckMemoryFiles(t *testing.T) {
	tempDir := t.TempDir()
	cMissing := auditCheckMemoryFiles(tempDir)
	if cMissing.IsPass {
		t.Fatalf("expected memory files check to fail when missing")
	}

	_ = os.MkdirAll(filepath.Join(tempDir, ".ai-memory"), 0755)
	_ = os.WriteFile(filepath.Join(tempDir, ".ai-memory", "strictly-avoid.md"), []byte("# Avoid"), 0644)
	cValid := auditCheckMemoryFiles(tempDir)
	if !cValid.IsPass {
		t.Fatalf("expected memory files check to pass when present")
	}
}

func TestAuditCheckRulesDir(t *testing.T) {
	tempDir := t.TempDir()
	cMissing := auditCheckRulesDir(tempDir)
	if cMissing.IsPass {
		t.Fatalf("expected rules dir check to fail when missing")
	}

	_ = os.MkdirAll(filepath.Join(tempDir, ".agents", "rules"), 0755)
	cValid := auditCheckRulesDir(tempDir)
	if !cValid.IsPass {
		t.Fatalf("expected rules dir check to pass when present")
	}
}
