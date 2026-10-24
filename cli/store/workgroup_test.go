package store

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func initTestTables(db *DB) {
	_, _ = db.conn.Exec(constants.SQLCreateRepo)
	_, _ = db.conn.Exec(constants.SQLCreateSshKey)
}

func openTestDBAt(t *testing.T, name string) *DB {
	dbPath := filepath.Join(t.TempDir(), name)
	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	return db
}

func setupWorkgroupTestDB(t *testing.T) *DB {
	t.Helper()
	db := openTestDBAt(t, "test_workgroup.db")
	initTestTables(db)
	if errMig := db.MigrateWorkgroup(); errMig != nil {
		db.Close()
		t.Fatalf("failed to migrate workgroup: %v", errMig)
	}

	return db
}

func dropLegacyObjects(db *DB) {
	_, _ = db.conn.Exec("DROP TABLE IF EXISTS work_directories")
	_, _ = db.conn.Exec("DROP VIEW IF EXISTS work_directories")
	_, _ = db.conn.Exec("DROP TABLE IF EXISTS WorkDirectory")
}

const sqlCreateLegacyWorkDirs = `CREATE TABLE work_directories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	absolute_path TEXT UNIQUE NOT NULL,
	label TEXT,
	is_default INTEGER DEFAULT 0,
	created_at TEXT DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);`

func seedLegacyWorkDirs(t *testing.T, db *DB) {
	t.Helper()
	dropLegacyObjects(db)
	_, _ = db.conn.Exec(sqlCreateLegacyWorkDirs)
	const sqlInsert = `INSERT INTO work_directories (absolute_path, label, is_default)
		VALUES ('/legacy/workspace', 'Legacy WS', 1);`
	_, _ = db.conn.Exec(sqlInsert)
}

func verifyMigratedRow(t *testing.T, wd *model.WorkDirectory) {
	t.Helper()
	if !wd.IsDefault || wd.Label != "Legacy WS" {
		t.Fatalf("unexpected migrated attributes: %+v", wd)
	}
}

func verifyMigratedData(t *testing.T, db *DB) {
	t.Helper()
	wd, err := db.GetWorkDirectoryByPath("/legacy/workspace")
	if err != nil || wd == nil {
		t.Fatalf("expected migrated workdir, err: %v", err)
	}

	verifyMigratedRow(t, wd)
	var count int
	_ = db.conn.QueryRow("SELECT COUNT(*) FROM work_directories").Scan(&count)
	if count != 1 {
		t.Fatalf("expected 1 row via compat view, count: %d", count)
	}
}

func runSchemaMigrationTest(t *testing.T, db *DB) {
	t.Helper()
	initTestTables(db)
	_, _ = db.conn.Exec(constants.SQLCreateSetting)
	seedLegacyWorkDirs(t, db)

	if errMig := Migration_AddFolderWorkDirAndSshKeyAutomation(db.conn); errMig != nil {
		t.Fatalf("migration failed: %v", errMig)
	}

	verifyMigratedData(t, db)
}

func testWorkgroupSchemaAndMigration(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "test_migration.db")
	db, errOpen := OpenAt(dbPath)
	if errOpen != nil {
		t.Fatalf("failed to open test db: %v", errOpen)
	}

	defer db.Close()

	runSchemaMigrationTest(t, db)
}

func verifyWorkDirCreation(t *testing.T, db *DB) int64 {
	wd, err := db.EnsureWorkDirectory("/home/user/workspace", "Primary", true)
	if err != nil || !wd.IsDefault || wd.Label != "Primary" {
		t.Fatalf("EnsureWorkDirectory failed: %v, attributes: %+v", err, wd)
	}

	byPath, errPath := db.GetWorkDirectoryByPath("/home/user/workspace")
	if errPath != nil || byPath.WorkDirectoryId != wd.WorkDirectoryId {
		t.Fatalf("GetWorkDirectoryByPath mismatch: %v", errPath)
	}

	return wd.WorkDirectoryId
}

func verifyWorkDirListAndDelete(t *testing.T, db *DB) {
	list, errList := db.ListWorkDirectories()
	if errList != nil || len(list) != 1 {
		t.Fatalf("ListWorkDirectories failed: %v, count: %d", errList, len(list))
	}

	if errDel := db.DeleteWorkDirectory("/home/user/workspace"); errDel != nil {
		t.Fatalf("DeleteWorkDirectory failed: %v", errDel)
	}
}

func testWorkgroupWorkDirectoryCRUD(t *testing.T) {
	db := setupWorkgroupTestDB(t)
	defer db.Close()

	wdId := verifyWorkDirCreation(t, db)
	byId, errId := db.GetWorkDirectoryById(wdId)
	if errId != nil || byId.AbsolutePath != "/home/user/workspace" {
		t.Fatalf("GetWorkDirectoryById failed: %v", errId)
	}

	verifyWorkDirListAndDelete(t, db)
}

func seedFolderWorkDirHierarchy(t *testing.T, db *DB) int64 {
	wd, err := db.EnsureWorkDirectory("/work/roots/corp", "Corp Root", true)
	if err != nil {
		t.Fatalf("EnsureWorkDirectory failed: %v", err)
	}

	_, errBind := db.BindFolderWorkDir("/work/roots/corp/clients", wd.WorkDirectoryId, true, 100)
	if errBind != nil {
		t.Fatalf("BindFolderWorkDir failed: %v", errBind)
	}

	return wd.WorkDirectoryId
}

func verifyChildResolution(t *testing.T, db *DB, wdId int64) {
	childUnix, errUnix := db.ResolveWorkDirectoryForPath("/work/roots/corp/clients/acme/project")
	if errUnix != nil || childUnix == nil || childUnix.WorkDirectoryId != wdId {
		t.Fatalf("ResolveWorkDirectoryForPath childUnix failed: %v", errUnix)
	}

	childWin, errWin := db.ResolveWorkDirectoryForPath("/work/roots/corp/clients\\acme\\project")
	if errWin != nil || childWin == nil || childWin.WorkDirectoryId != wdId {
		t.Fatalf("ResolveWorkDirectoryForPath childWin failed: %v", errWin)
	}
}

func verifyFolderResolution(t *testing.T, db *DB, wdId int64) {
	exact, errExact := db.ResolveWorkDirectoryForPath("/work/roots/corp/clients")
	if errExact != nil || exact == nil || exact.WorkDirectoryId != wdId {
		t.Fatalf("ResolveWorkDirectoryForPath exact failed: %v", errExact)
	}

	verifyChildResolution(t, db, wdId)
}

func testWorkgroupFolderWorkDirResolution(t *testing.T) {
	db := setupWorkgroupTestDB(t)
	defer db.Close()

	wdId := seedFolderWorkDirHierarchy(t, db)
	verifyFolderResolution(t, db, wdId)
	if errUnbind := db.UnbindFolderWorkDir("/work/roots/corp/clients"); errUnbind != nil {
		t.Fatalf("UnbindFolderWorkDir failed: %v", errUnbind)
	}
}

func seedTestGroups(t *testing.T, db *DB) {
	_, err1 := db.AddFolderToGroup("team-core", "/repos/api", "API Service")
	if err1 != nil {
		t.Fatalf("AddFolderToGroup 1 failed: %v", err1)
	}

	_, err2 := db.AddFolderToGroup("team-core", "/repos/web", "Web Frontend")
	if err2 != nil {
		t.Fatalf("AddFolderToGroup 2 failed: %v", err2)
	}
}

func verifyGroupLists(t *testing.T, db *DB) {
	items, errList := db.ListFoldersInGroup("team-core")
	if errList != nil || len(items) != 2 {
		t.Fatalf("ListFoldersInGroup failed: %v, count: %d", errList, len(items))
	}

	summaries, errSum := db.ListActiveGroups()
	if errSum != nil || len(summaries) != 1 || summaries[0].FolderCount != 2 {
		t.Fatalf("ListActiveGroups failed: %v, summaries: %+v", errSum, summaries)
	}
}

func testWorkgroupFolderGroupCRUD(t *testing.T) {
	db := setupWorkgroupTestDB(t)
	defer db.Close()

	seedTestGroups(t, db)
	verifyGroupLists(t, db)
	if errRm := db.RemoveFolderFromGroup("team-core", "/repos/api"); errRm != nil {
		t.Fatalf("RemoveFolderFromGroup failed: %v", errRm)
	}
}

func insertMockSshKey(t *testing.T, db *DB, name, keyPath string, isDefault int) int64 {
	t.Helper()
	const sqlInsert = `INSERT INTO SshKey (Name, PrivatePath, PublicKey, Fingerprint, IsDefault)
		VALUES (?, ?, 'ssh-ed25519 AAAA...', 'SHA256:mock', ?);`
	res, err := db.conn.Exec(sqlInsert, name, keyPath, isDefault)
	if err != nil {
		t.Fatalf("insertMockSshKey failed: %v", err)
	}

	id, _ := res.LastInsertId()

	return id
}

func verifyStageResolutions(t *testing.T, db *DB, kFolder, kRepo int64) {
	credFolder, errF := db.ResolveEffectiveCredentialsForPath("/projects/finance/svc", 0)
	if errF != nil || !credFolder.HasCustomKey || credFolder.SshKeyId != kFolder {
		t.Fatalf("expected folder cred %d, got: %+v, err: %v", kFolder, credFolder, errF)
	}

	if credFolder.GitAuthorName != "Finance Dev" {
		t.Fatalf("expected GitAuthorName override 'Finance Dev', got: %s", credFolder.GitAuthorName)
	}

	credRepo, errR := db.ResolveEffectiveCredentialsForPath("/projects/finance/svc", 99)
	if errR != nil || !credRepo.HasCustomKey || credRepo.SshKeyId != kRepo {
		t.Fatalf("expected repo cred %d, got: %+v, err: %v", kRepo, credRepo, errR)
	}
}

func seedSshMockEntities(t *testing.T, db *DB) (int64, int64, int64) {
	kDefault := insertMockSshKey(t, db, "key-default", "/home/.ssh/id_default", 1)
	kFolder := insertMockSshKey(t, db, "key-folder", "/home/.ssh/id_folder", 0)
	kRepo := insertMockSshKey(t, db, "key-repo", "/home/.ssh/id_repo", 0)

	const sqlInsertRepo = `INSERT INTO Repo (RepoId, Slug, RepoName, HttpsUrl, SshUrl, Branch, RelativePath, AbsolutePath, CloneInstruction, IdentifiedTransport)
		VALUES (99, 'svc-fin', 'svc-fin', 'https://github.com/org/svc-fin', 'git@github.com:org/svc-fin.git', 'main', 'projects/finance/svc', '/projects/finance/svc', 'git clone', 'ssh');`
	_, _ = db.conn.Exec(sqlInsertRepo)

	return kDefault, kFolder, kRepo
}

func verifyFallbackCred(t *testing.T, db *DB, kDefault int64) {
	cred, err := db.ResolveEffectiveCredentialsForPath("/other/path", 0)
	if err != nil || !cred.HasCustomKey || cred.SshKeyId != kDefault {
		t.Fatalf("expected fallback default key %d, got: %+v", kDefault, cred)
	}
}

func testWorkgroupSshKeyCascadingResolution(t *testing.T) {
	db := setupWorkgroupTestDB(t)
	defer db.Close()

	kDefault, kFolder, kRepo := seedSshMockEntities(t, db)
	_, _ = db.AssociateSSHKeyWithFolder(kFolder, "/projects/finance", "Finance Dev", "dev@fin.corp")
	_, _ = db.AssociateSSHKeyWithRepo(kRepo, 99)

	verifyStageResolutions(t, db, kFolder, kRepo)
	verifyFallbackCred(t, db, kDefault)
}

func seedDeletionWorkDir(t *testing.T, db *DB) *model.WorkDirectory {
	wd, _ := db.EnsureWorkDirectory("/to/delete", "ToDelete", false)
	_, _ = db.BindFolderWorkDir("/to/delete/sub", wd.WorkDirectoryId, true, 100)
	keyId := insertMockSshKey(t, db, "key-cascade", "/home/.ssh/id_casc", 0)
	_, _ = db.AssociateSSHKeyWithWorkDir(keyId, wd.WorkDirectoryId)

	return wd
}

func verifyCascadeResult(t *testing.T, db *DB, wdId int64) {
	var folderCount int
	_ = db.conn.QueryRow("SELECT COUNT(*) FROM FolderWorkDir WHERE WorkDirectoryId = ?", wdId).Scan(&folderCount)
	if folderCount != 0 {
		t.Fatalf("expected FolderWorkDir cascade delete, count: %d", folderCount)
	}
}

func testWorkgroupCascadeDeletion(t *testing.T) {
	db := setupWorkgroupTestDB(t)
	defer db.Close()

	wd := seedDeletionWorkDir(t, db)
	if errDel := db.DeleteWorkDirectory("/to/delete"); errDel != nil {
		t.Fatalf("DeleteWorkDirectory failed: %v", errDel)
	}

	verifyCascadeResult(t, db, wd.WorkDirectoryId)
}

func TestWorkgroupSuite(t *testing.T) {
	t.Run("SchemaAndMigration", testWorkgroupSchemaAndMigration)
	t.Run("WorkDirectoryCRUD", testWorkgroupWorkDirectoryCRUD)
	t.Run("FolderWorkDirResolution", testWorkgroupFolderWorkDirResolution)
	t.Run("FolderGroupCRUD", testWorkgroupFolderGroupCRUD)
	t.Run("SshKeyCascadingResolution", testWorkgroupSshKeyCascadingResolution)
	t.Run("CascadeDeletion", testWorkgroupCascadeDeletion)
}
