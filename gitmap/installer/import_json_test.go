// Package installer — import_json_test.go tests raw JSON import logic.
package installer

import (
	"testing"
)

func TestImportFromJson(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	jsonStr := `{"name": "JSON App", "slug": "json-app", "target_os": "all", "version": "v1.0.0"}`
	if errImport := mgr.ImportFromJson(jsonStr); errImport != nil {
		t.Fatalf("ImportFromJson failed: %v", errImport)
	}

	imported, errGet := db.GetInstallerBySlug("json-app")
	if errGet != nil || imported == nil {
		t.Fatalf("failed to retrieve imported script: %v", errGet)
	}
}
