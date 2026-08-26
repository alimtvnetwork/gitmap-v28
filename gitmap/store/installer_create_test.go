package store

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func setupInstallerCreateTestDB(testingT *testing.T) *DB {
	testingT.Helper()
	tempDir := testingT.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer_create.db")

	dbInstance, errOpen := OpenAt(dbPath)
	if errOpen != nil {
		testingT.Fatalf("failed to open test db: %v", errOpen)
	}

	if errMigrate := dbInstance.MigrateInstallers(); errMigrate != nil {
		dbInstance.Close()
		testingT.Fatalf("failed to migrate installers: %v", errMigrate)
	}

	return dbInstance
}

func TestCreateInstallerSuccess(testingT *testing.T) {
	dbInstance := setupInstallerCreateTestDB(testingT)
	defer dbInstance.Close()

	script := &model.InstallerScript{
		Name:         "Golang Suite",
		Slug:         "go-suite",
		Description:  "Go tooling installer",
		TargetOS:     "win",
		Version:      "v1.0.0",
		Instructions: `{"steps":[{"action":"install","pkg":"go"}]}`,
	}

	errCreate := dbInstance.CreateInstaller(script)
	if errCreate != nil {
		testingT.Fatalf("expected CreateInstaller to succeed, got: %v", errCreate)
	}

	if script.ID <= 0 {
		testingT.Fatalf("expected script.ID to be populated > 0, got %d", script.ID)
	}

	row := dbInstance.conn.QueryRow("SELECT id, name, slug, description, target_os, version, instructions, created_at, updated_at FROM installer_scripts WHERE slug = ?", "go-suite")
	var (
		fetchedID           int64
		fetchedName         string
		fetchedSlug         string
		fetchedDescription  string
		fetchedTargetOS     string
		fetchedVersion      string
		fetchedInstructions string
		fetchedCreatedAt    string
		fetchedUpdatedAt    string
	)

	errScan := row.Scan(
		&fetchedID,
		&fetchedName,
		&fetchedSlug,
		&fetchedDescription,
		&fetchedTargetOS,
		&fetchedVersion,
		&fetchedInstructions,
		&fetchedCreatedAt,
		&fetchedUpdatedAt,
	)
	if errScan != nil {
		testingT.Fatalf("failed to scan inserted script: %v", errScan)
	}

	if fetchedID != script.ID {
		testingT.Errorf("expected ID %d, got %d", script.ID, fetchedID)
	}
	if fetchedName != "Golang Suite" {
		testingT.Errorf("expected name 'Golang Suite', got %q", fetchedName)
	}
	if fetchedSlug != "go-suite" {
		testingT.Errorf("expected slug 'go-suite', got %q", fetchedSlug)
	}
	if fetchedDescription != "Go tooling installer" {
		testingT.Errorf("expected description 'Go tooling installer', got %q", fetchedDescription)
	}
	if fetchedTargetOS != "win" {
		testingT.Errorf("expected targetOS 'win', got %q", fetchedTargetOS)
	}
	if fetchedVersion != "v1.0.0" {
		testingT.Errorf("expected version 'v1.0.0', got %q", fetchedVersion)
	}
	if fetchedInstructions != script.Instructions {
		testingT.Errorf("expected instructions %q, got %q", script.Instructions, fetchedInstructions)
	}
	if fetchedCreatedAt == "" || fetchedUpdatedAt == "" {
		testingT.Errorf("expected timestamps to be populated")
	}
}

func TestCreateInstallerNilScript(testingT *testing.T) {
	dbInstance := setupInstallerCreateTestDB(testingT)
	defer dbInstance.Close()

	errCreate := dbInstance.CreateInstaller(nil)
	if errCreate == nil {
		testingT.Fatalf("expected error on nil script, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errCreate, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errCreate)
	}

	if appErr.Code != "E_INSTALLER_INVALID_INPUT" {
		testingT.Errorf("expected error code E_INSTALLER_INVALID_INPUT, got %s", appErr.Code)
	}
}

func TestCreateInstallerDuplicateSlug(testingT *testing.T) {
	dbInstance := setupInstallerCreateTestDB(testingT)
	defer dbInstance.Close()

	script1 := &model.InstallerScript{
		Name: "Docker CE",
		Slug: "docker",
	}
	script2 := &model.InstallerScript{
		Name: "Docker Desktop",
		Slug: "docker",
	}

	if err1 := dbInstance.CreateInstaller(script1); err1 != nil {
		testingT.Fatalf("first insert failed: %v", err1)
	}

	err2 := dbInstance.CreateInstaller(script2)
	if err2 == nil {
		testingT.Fatalf("expected duplicate slug insert to fail, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(err2, &appErr) {
		testingT.Fatalf("expected AppError, got %T", err2)
	}

	if appErr.Code != "E_INSTALLER_CREATE_FAILED" {
		testingT.Errorf("expected error code E_INSTALLER_CREATE_FAILED, got %s", appErr.Code)
	}
}

func TestCreateInstallerClosedDB(testingT *testing.T) {
	dbInstance := setupInstallerCreateTestDB(testingT)
	dbInstance.Close()

	script := &model.InstallerScript{
		Name: "NodeJS",
		Slug: "nodejs",
	}

	errCreate := dbInstance.CreateInstaller(script)
	if errCreate == nil {
		testingT.Fatalf("expected error on closed db, got nil")
	}

	var appErr *apperror.AppError
	if !errors.As(errCreate, &appErr) {
		testingT.Fatalf("expected AppError, got %T", errCreate)
	}

	if appErr.Code != "E_INSTALLER_CREATE_FAILED" {
		testingT.Errorf("expected error code E_INSTALLER_CREATE_FAILED, got %s", appErr.Code)
	}
}
