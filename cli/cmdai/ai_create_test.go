package cmdai

import (
	"strings"
	"testing"
)

func TestGenerateScriptContent_Linter(t *testing.T) {
	opts := CreateScriptOptions{
		Name:        "test-linter",
		Type:        TemplateLinter,
		Description: "A unit test linter",
		IsParallel:  true,
	}

	res := GenerateScriptContent(opts)
	if res.IsError() {
		t.Fatalf("unexpected error: %v", res.Error())
	}

	code := res.Data()
	assertContains(t, code, "check_file")
	assertContains(t, code, "--workers")
	assertContains(t, code, "test-linter")
}

func TestGenerateScriptContent_Fixer(t *testing.T) {
	opts := CreateScriptOptions{
		Name:       "test-fixer",
		Type:       TemplateFixer,
		HasFixMode: true,
	}

	res := GenerateScriptContent(opts)
	if res.IsError() {
		t.Fatalf("unexpected error: %v", res.Error())
	}

	code := res.Data()
	assertContains(t, code, "fix_file")
	assertContains(t, code, "--fix")
}

func TestGenerateScriptContent_Auditor(t *testing.T) {
	opts := CreateScriptOptions{
		Name: "test-auditor",
		Type: TemplateAuditor,
	}

	res := GenerateScriptContent(opts)
	if res.IsError() {
		t.Fatalf("unexpected error: %v", res.Error())
	}

	code := res.Data()
	assertContains(t, code, "audit_component")
	assertContains(t, code, "--strict")
}

func TestSanitizeScriptSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"my_new_script.py", "my-new-script"},
		{"  45-Custom-Linter  ", "45-custom-linter"},
		{"UPPER_CASE", "upper-case"},
	}

	for _, tc := range tests {
		got := sanitizeScriptSlug(tc.input)
		if got != tc.want {
			t.Errorf("sanitizeScriptSlug(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestValidateCreateOptions(t *testing.T) {
	emptyOpts := CreateScriptOptions{}
	err := validateCreateOptions(emptyOpts)
	if err == nil {
		t.Error("expected error for empty name, got nil")
	}

	validOpts := CreateScriptOptions{Name: "my-script"}
	err = validateCreateOptions(validOpts)
	if err != nil {
		t.Errorf("unexpected error for valid opts: %v", err)
	}
}

func TestRunAiCreate_DryRun(t *testing.T) {
	opts := CreateScriptOptions{
		Name:     "dry-run-script",
		Type:     TemplateLinter,
		IsDryRun: true,
	}

	appErr := RunAiCreate(opts)
	if appErr != nil {
		t.Fatalf("unexpected dry-run error: %v", appErr)
	}
}

func assertContains(t *testing.T, s, substr string) {
	t.Helper()
	hasSub := strings.Contains(s, substr)
	if !hasSub {
		t.Errorf("expected string to contain %q, got: %s", substr, s)
	}
}
