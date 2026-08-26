package cmd

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerHistoryTreeTestDB(testingT *testing.T) *store.DB {
	tempDirectory := testingT.TempDir()
	dbFilePath := filepath.Join(tempDirectory, "test_history_tree.db")
	dbInstance, errOpen := store.Open(dbFilePath)
	if errOpen != nil {
		testingT.Fatalf("failed to open test database: %v", errOpen)
	}
	testingT.Cleanup(func() { dbInstance.Close() })

	if errMigrate := dbInstance.MigrateInstallers(); errMigrate != nil {
		testingT.Fatalf("failed to migrate installers: %v", errMigrate)
	}
	return dbInstance
}

func TestPrintInstallerHistoryTree_EmptyAndNil(testingT *testing.T) {
	testingT.Parallel()

	// Should not panic on nil DB
	printInstallerHistoryTree(nil)

	// Should handle empty DB gracefully
	dbInstance := setupInstallerHistoryTreeTestDB(testingT)
	printInstallerHistoryTree(dbInstance)
}

func TestPrintInstallerHistoryTree_ProfileAndNonProfile(testingT *testing.T) {
	testingT.Parallel()

	dbInstance := setupInstallerHistoryTreeTestDB(testingT)

	scriptDev := &model.InstallerScript{
		Name:        "Ubuntu Developer Profile",
		Slug:        "ubuntu+dev",
		Description: "Full workstation suite",
		TargetOS:    "ubuntu",
		Version:     "v1.0.0",
	}
	scriptCustom := &model.InstallerScript{
		Name:        "Custom CLI Tool",
		Slug:        "custom-cli",
		Description: "A custom utility binary",
		TargetOS:    "linux",
		Version:     "v2.1.0",
	}

	if errCreate := dbInstance.CreateInstaller(scriptDev); errCreate != nil {
		testingT.Fatalf("failed to insert scriptDev: %v", errCreate)
	}
	if errCreate := dbInstance.CreateInstaller(scriptCustom); errCreate != nil {
		testingT.Fatalf("failed to insert scriptCustom: %v", errCreate)
	}

	// Should render without error or panic
	printInstallerHistoryTree(dbInstance)
}

func TestGroupLatestInstallers(testingT *testing.T) {
	testingT.Parallel()

	scripts := []model.InstallerScript{
		{ID: 1, Slug: "tool-a", UpdatedAt: "2026-08-01 10:00:00", Version: "v1.0.0"},
		{ID: 2, Slug: "tool-a", UpdatedAt: "2026-08-02 10:00:00", Version: "v1.1.0"},
		{ID: 3, Slug: "tool-b", UpdatedAt: "2026-08-03 10:00:00", Version: "v2.0.0"},
	}

	grouped := groupLatestInstallers(scripts)
	if len(grouped) != 2 {
		testingT.Fatalf("expected 2 unique slugs, got %d", len(grouped))
	}
	if grouped[0].Slug != "tool-b" {
		testingT.Errorf("expected tool-b first (newest), got %s", grouped[0].Slug)
	}
	if grouped[1].Slug != "tool-a" || grouped[1].Version != "v1.1.0" {
		testingT.Errorf("expected tool-a v1.1.0, got %s %s", grouped[1].Slug, grouped[1].Version)
	}
}

func TestRenderSingleHistoryTree(testingT *testing.T) {
	testingT.Parallel()

	scriptRecord := model.InstallerScript{
		Slug:        "sample-script",
		Description: "Sample description",
		TargetOS:    "ubuntu",
		Version:     "v1.0.0",
		CreatedAt:   "2026-08-26 12:00:00",
	}

	renderSingleHistoryTree(scriptRecord)
	renderHistoryEntry(scriptRecord)
}
