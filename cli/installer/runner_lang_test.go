package installer

import (
	"context"
	"testing"
)

func TestRunnerLang(t *testing.T) {
	ctx := context.Background()

	if errEmpty := RunLanguageScript(ctx, "", "bash"); errEmpty == nil {
		t.Fatal("expected error on empty script")
	}

	if errRun := RunLanguageScript(ctx, "echo test_lang", ""); errRun != nil {
		t.Fatalf("RunLanguageScript default failed: %v", errRun)
	}
}
