package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func createTestSnapshotJSON(t *testing.T, dir, filename, name, displayName, email string) string {
	t.Helper()
	content := fmt.Sprintf(`{"schema_version":1,"name":%q,"display_name":%q,"email":%q}`, name, displayName, email)
	target := filepath.Join(dir, filename)
	if err := os.WriteFile(target, []byte(content), 0644); err != nil {
		t.Fatalf("write test snapshot: %v", err)
	}

	return target
}

func TestProfileImportRouting(t *testing.T) {
	tempUserData := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempUserData)

	workDir := t.TempDir()
	createTestSnapshotJSON(t, workDir, "Default.json", "Default", "Personal", "test@test.com")

	if err := runProfile([]string{"import", workDir}); err != nil {
		t.Errorf("gitmap profile import failed: %v", err)
	}

	if err := runProfile([]string{"import-all", workDir}); err != nil {
		t.Errorf("gitmap profile import-all failed: %v", err)
	}

	if err := runProfile([]string{"inspect", workDir}); err != nil {
		t.Errorf("gitmap profile inspect failed: %v", err)
	}

	if err := runProfile([]string{"preview", workDir}); err != nil {
		t.Errorf("gitmap profile preview failed: %v", err)
	}

	if err := runProfile([]string{"check", workDir}); err != nil {
		t.Errorf("gitmap profile check failed: %v", err)
	}
}

func TestProfileInvalidSubcommandReturnsAppError(t *testing.T) {
	err := runProfile([]string{"nonexistent-subcommand-123"})
	if err == nil {
		t.Fatalf("expected error for invalid subcommand, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *apperror.AppError, got %T: %v", err, err)
	}

	if appErr.Code != "E9000" {
		t.Errorf("expected error code E9000, got %s", appErr.Code)
	}
}
