package cmdautomation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRelPathAudit_DetectsAndFixes(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "doc.md")

	absTarget := filepath.ToSlash(tempDir) + "/sub/file.go"
	initialContent := "# Test\nSee file:///" + absTarget + " for info.\n"
	if err := os.WriteFile(testFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	opts := RelPathOptions{
		Dir:       tempDir,
		IsFixMode: false,
	}
	monad := RunRelPathAudit(opts)
	if monad.IsFailure() {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	if len(res.Violations) == 0 {
		t.Errorf("expected violations for file URI, got none")
	}
}

func TestNamingAudit_DetectsViolations(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "sample.go")
	badContent := "package sample\n\nfunc check(trigger bool) bool {\n\tif isValid == true {\n\t\treturn true\n\t}\n\treturn false\n}\n"
	if err := os.WriteFile(testFile, []byte(badContent), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	opts := NamingOptions{
		Dir: tempDir,
	}
	monad := RunNamingAudit(opts)
	if monad.IsFailure() {
		t.Fatalf("unexpected error: %v", monad.Err)
	}
	res := monad.Value
	if len(res.Violations) < 2 {
		t.Errorf("expected at least 2 violations (EXPLICIT_TRUE and NON_AFFIRMATIVE_BOOLEAN), got: %d", len(res.Violations))
	}
}

func TestRuleAudit_ResultWrapperAndEnums(t *testing.T) {
	tempDir := t.TempDir()
	testGo := filepath.Join(tempDir, "bad.go")
	content := "package bad\n\nfunc getMap() (map[string]int, error) {\n\t_ = " + "rune" + "(10)\n\treturn nil, nil\n}\n"
	if err := os.WriteFile(testGo, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	resWrapMonad := RunRuleAudit(RuleAuditOptions{Dir: tempDir, RuleName: "result-wrapper"})
	if resWrapMonad.IsFailure() {
		t.Fatalf("unexpected failure: %v", resWrapMonad.Err)
	}
	if len(resWrapMonad.Value.Violations) == 0 {
		t.Errorf("expected RESULT_MAP_TUPLE violation")
	}

	enumMonad := RunRuleAudit(RuleAuditOptions{Dir: tempDir, RuleName: "enums"})
	if enumMonad.IsFailure() {
		t.Fatalf("unexpected failure: %v", enumMonad.Err)
	}
	if len(enumMonad.Value.Violations) == 0 {
		t.Errorf("expected RAW_RUNE_CAST violation")
	}
}
