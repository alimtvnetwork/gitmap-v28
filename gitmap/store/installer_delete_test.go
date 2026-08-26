package store

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func TestDeleteInstallerStore(t *testing.T) {
	db := setupInstallerCreateTestDB(t)
	defer db.Close()

	script := &model.InstallerScript{
		Name:     "Delete Store App",
		Slug:     "delete-store-app",
		TargetOS: "all",
		Version:  "v1.0.0",
	}
	if err := db.CreateInstaller(script); err != nil {
		t.Fatalf("CreateInstaller failed: %v", err)
	}

	if err := db.DeleteInstaller("delete-store-app"); err != nil {
		t.Fatalf("DeleteInstaller failed: %v", err)
	}

	if errVer := db.DeleteInstallerVersion("delete-store-app", "v1.0.0"); errVer != nil {
		t.Fatalf("DeleteInstallerVersion failed: %v", errVer)
	}
}
