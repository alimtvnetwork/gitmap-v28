package cmd

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
	"github.com/spf13/cobra"
)

func setupTestDB(t *testing.T) *store.DB {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_cmd_installer_create.db")

	db, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.MigrateInstallers(); err != nil {
		db.Close()
		t.Fatalf("failed to migrate installer tables: %v", err)
	}

	return db
}

func TestInstallerCreateCmd(t *testing.T) {
	t.Parallel()

	if installerCreateCmd == nil {
		t.Fatal("expected installerCreateCmd to be initialized")
	}

	if installerCreateCmd.Use != "create <name>" {
		t.Fatalf("expected Use to be 'create <name>', got %q", installerCreateCmd.Use)
	}

	if installerCreateCmd.Short == "" {
		t.Fatal("expected Short description to be non-empty")
	}

	// Test failure boundary with nil command
	errNil := runInstallerCreate(nil, []string{"my-pkg"})
	if errNil == nil {
		t.Fatal("expected error when passing nil command")
	}

	var appErr *apperror.AppError
	if !errors.As(errNil, &appErr) {
		t.Fatalf("expected *apperror.AppError, got %T", errNil)
	}
	if appErr.Code != "E_INSTALLER_NIL_COMMAND" {
		t.Fatalf("expected code E_INSTALLER_NIL_COMMAND, got %q", appErr.Code)
	}

	// Test failure boundary with empty args
	cmd := &cobra.Command{}
	errEmpty := runInstallerCreate(cmd, []string{})
	if errEmpty == nil {
		t.Fatal("expected error when running with empty args")
	}
	if !errors.As(errEmpty, &appErr) {
		t.Fatalf("expected *apperror.AppError, got %T", errEmpty)
	}
	if appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Fatalf("expected code E_INSTALLER_INVALID_INPUT, got %q", appErr.Code)
	}
}

func TestParseCreateFlagsSuccess(t *testing.T) {
	t.Parallel()

	// 1. Positional argument
	flags, err := parseCreateFlags([]string{"my-awesome-tool"})
	if err != nil {
		t.Fatalf("parseCreateFlags failed: %v", err)
	}
	if flags.Name != "my-awesome-tool" {
		t.Errorf("expected Name 'my-awesome-tool', got %q", flags.Name)
	}
	if flags.Slug != "my-awesome-tool" {
		t.Errorf("expected Slug 'my-awesome-tool', got %q", flags.Slug)
	}
	if flags.TargetOS != "win" {
		t.Errorf("expected default TargetOS 'win', got %q", flags.TargetOS)
	}
	if flags.Version != "v1.0.0" {
		t.Errorf("expected default Version 'v1.0.0', got %q", flags.Version)
	}

	// 2. Full flags
	fullArgs := []string{
		"-n", "Golang CLI Tool",
		"--desc", "Installs go binaries",
		"--os", "ubuntu",
		"-v", "v2.1.0",
		"-i", `{"step":"go install"}`,
		"-s", "go-cli-custom",
	}
	flagsFull, errFull := parseCreateFlags(fullArgs)
	if errFull != nil {
		t.Fatalf("parseCreateFlags with full args failed: %v", errFull)
	}
	if flagsFull.Name != "Golang CLI Tool" {
		t.Errorf("expected Name 'Golang CLI Tool', got %q", flagsFull.Name)
	}
	if flagsFull.Slug != "go-cli-custom" {
		t.Errorf("expected Slug 'go-cli-custom', got %q", flagsFull.Slug)
	}
	if flagsFull.Description != "Installs go binaries" {
		t.Errorf("expected Description 'Installs go binaries', got %q", flagsFull.Description)
	}
	if flagsFull.TargetOS != "ubuntu" {
		t.Errorf("expected TargetOS 'ubuntu', got %q", flagsFull.TargetOS)
	}
	if flagsFull.Version != "v2.1.0" {
		t.Errorf("expected Version 'v2.1.0', got %q", flagsFull.Version)
	}
	if flagsFull.Instructions != `{"step":"go install"}` {
		t.Errorf("expected Instructions %q, got %q", `{"step":"go install"}`, flagsFull.Instructions)
	}

	// 3. Name with space auto-slugifies
	flagsSpace, errSpace := parseCreateFlags([]string{"Docker Desktop Suite", "-d", "Docker stuff"})
	if errSpace != nil {
		t.Fatalf("parseCreateFlags failed: %v", errSpace)
	}
	if flagsSpace.Slug != "docker-desktop-suite" {
		t.Errorf("expected Slug 'docker-desktop-suite', got %q", flagsSpace.Slug)
	}
}

func TestParseCreateFlagsFailure(t *testing.T) {
	t.Parallel()

	// Missing name
	_, errEmpty := parseCreateFlags([]string{})
	if errEmpty == nil {
		t.Fatal("expected error on empty args")
	}
	var appErr *apperror.AppError
	if !errors.As(errEmpty, &appErr) || appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Fatalf("expected E_INSTALLER_INVALID_INPUT, got %v", errEmpty)
	}

	// Unknown flag
	_, errFlag := parseCreateFlags([]string{"--invalid-flag-123"})
	if errFlag == nil {
		t.Fatal("expected error on invalid flag")
	}
	if !errors.As(errFlag, &appErr) || appErr.Code != "E_INSTALLER_INVALID_FLAGS" {
		t.Fatalf("expected E_INSTALLER_INVALID_FLAGS, got %v", errFlag)
	}

	// Whitespace only name
	_, errBlank := parseCreateFlags([]string{"-n", "   "})
	if errBlank == nil {
		t.Fatal("expected error on blank name")
	}
	if !errors.As(errBlank, &appErr) || appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Fatalf("expected E_INSTALLER_INVALID_INPUT, got %v", errBlank)
	}
}

func TestExecuteCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// 1. Success execution
	flags := &CreateInstallerFlags{
		Name:         "NodeJS LTS",
		Slug:         "nodejs-lts",
		Description:  "Node.js runtime",
		TargetOS:     "win",
		Version:      "v18.0.0",
		Instructions: `{"pkg":"nodejs"}`,
	}

	err := executeCreate(ctx, db, flags)
	if err != nil {
		t.Fatalf("executeCreate failed: %v", err)
	}

	// Verify row in DB
	saved, errGet := db.GetInstallerBySlug("nodejs-lts")
	if errGet != nil {
		t.Fatalf("GetInstallerBySlug failed: %v", errGet)
	}
	if saved.Name != "NodeJS LTS" || saved.Version != "v18.0.0" {
		t.Errorf("saved script mismatch: %+v", saved)
	}

	// 2. Duplicate slug failure
	errDup := executeCreate(ctx, db, flags)
	if errDup == nil {
		t.Fatal("expected error on duplicate slug")
	}
	var appErr *apperror.AppError
	if !errors.As(errDup, &appErr) || appErr.Code != "E_INSTALLER_CREATE_FAILED" {
		t.Fatalf("expected E_INSTALLER_CREATE_FAILED, got %v", errDup)
	}

	// 3. Nil DB failure
	errNilDB := executeCreate(ctx, nil, flags)
	if errNilDB == nil {
		t.Fatal("expected error on nil db")
	}
	if !errors.As(errNilDB, &appErr) || appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Fatalf("expected E_INSTALLER_INVALID_INPUT, got %v", errNilDB)
	}

	// 4. Nil flags failure
	errNilFlags := executeCreate(ctx, db, nil)
	if errNilFlags == nil {
		t.Fatal("expected error on nil flags")
	}
	if !errors.As(errNilFlags, &appErr) || appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		t.Fatalf("expected E_INSTALLER_INVALID_INPUT, got %v", errNilFlags)
	}
}

func TestSlugify(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input    string
		expected string
	}{
		{"NodeJS LTS", "nodejs-lts"},
		{"Docker Desktop v2", "docker-desktop-v2"},
		{"my_custom_tool", "my-custom-tool"},
		{"   spaced   ", "spaced"},
		{"Tool.With.Dots", "tool-with-dots"},
		{"Special!@#Characters", "specialcharacters"},
		{"---dashes---", "dashes"},
	}

	for _, tc := range cases {
		actual := slugify(tc.input)
		if actual != tc.expected {
			t.Errorf("slugify(%q) = %q, expected %q", tc.input, actual, tc.expected)
		}
	}
}
