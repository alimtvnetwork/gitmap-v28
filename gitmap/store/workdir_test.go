package store

import (
	"testing"
)

func TestWorkDirStoreSuite(t *testing.T) {
	db := setupInstallerCreateTestDB(t)
	defer db.Close()

	if _, err := db.conn.Exec(SQLCreateWorkDirsTable); err != nil {
		t.Fatalf("failed creating work dirs table: %v", err)
	}

	wd, errEnsure := db.EnsureWorkDir("/home/user/git-work", "main-workspace", true)
	if errEnsure != nil {
		t.Fatalf("EnsureWorkDir failed: %v", errEnsure)
	}
	if !wd.IsDefault {
		t.Fatal("expected isDefault to be true")
	}

	list, errList := db.ListWorkDirs()
	if errList != nil || len(list) == 0 {
		t.Fatalf("ListWorkDirs failed: %v (count: %d)", errList, len(list))
	}

	if errDef := db.SetDefaultWorkDir("/home/user/git-work"); errDef != nil {
		t.Fatalf("SetDefaultWorkDir failed: %v", errDef)
	}

	if errDel := db.DeleteWorkDir("/home/user/git-work"); errDel != nil {
		t.Fatalf("DeleteWorkDir failed: %v", errDel)
	}
}
