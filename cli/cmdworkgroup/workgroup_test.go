package cmdworkgroup

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func setupTestDB(t *testing.T) (*store.DB, func()) {
	t.Helper()
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")
	db, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	origOpen := openDBFunc
	openDBFunc = func() (*store.DB, error) {
		return store.OpenAt(dbPath)
	}

	return db, func() {
		openDBFunc = origOpen
		_ = db.Close()
	}
}

func TestWorkgroupFlagParsing(t *testing.T) {
	args := []string{
		"add", "client-alpha", "/tmp/client-alpha",
		"--ssh-key", "/home/user/.ssh/id_client",
		"--desc", "Client Alpha Workspace",
		"--user-name", "John Doe",
		"--user-email", "john@example.com",
		"--default",
		"-j",
	}

	opts := ParseWorkgroupOptions(args)
	if opts.Subcommand != "add" {
		t.Fatalf("expected subcommand add, got %s", opts.Subcommand)
	}
	if opts.Name != "client-alpha" {
		t.Fatalf("expected name client-alpha, got %s", opts.Name)
	}
	if opts.Path != "/tmp/client-alpha" {
		t.Fatalf("expected path /tmp/client-alpha, got %s", opts.Path)
	}
	if opts.SSHKeyPath != "/home/user/.ssh/id_client" {
		t.Fatalf("expected ssh key, got %s", opts.SSHKeyPath)
	}
	if opts.Description != "Client Alpha Workspace" {
		t.Fatalf("expected desc, got %s", opts.Description)
	}
	if opts.UserName != "John Doe" || opts.UserEmail != "john@example.com" {
		t.Fatalf("expected git identity mismatch")
	}
	if !opts.IsDefault || !opts.IsJSON {
		t.Fatalf("expected isDefault and isJSON to be true")
	}
}

func TestWorkgroupCRUD(t *testing.T) {
	_, cleanup := setupTestDB(t)
	defer cleanup()

	rec := WorkgroupRecord{
		Name:          "team-core",
		RootDirectory: filepath.Clean("/workspace/team-core"),
		SSHKeyPath:    "/root/.ssh/id_rsa",
		Description:   "Core services",
		GitUserName:   "Core Bot",
		GitUserEmail:  "bot@core.org",
		IsDefault:     true,
	}

	created, err := AddWorkgroup(rec)
	if err != nil {
		t.Fatalf("failed to add workgroup: %v", err)
	}
	if created.ID <= 0 {
		t.Fatalf("expected positive ID, got %d", created.ID)
	}

	list, errList := ListWorkgroups()
	if errList != nil || len(list) != 1 {
		t.Fatalf("expected 1 workgroup, got %d, err: %v", len(list), errList)
	}

	fetched, errGet := GetWorkgroup("team-core")
	if errGet != nil || fetched.Name != "team-core" {
		t.Fatalf("failed to get workgroup: %v", errGet)
	}
	if !fetched.IsDefault {
		t.Fatalf("expected isDefault=true")
	}

	if errDel := DeleteWorkgroup("team-core"); errDel != nil {
		t.Fatalf("failed to delete workgroup: %v", errDel)
	}

	listAfter, _ := ListWorkgroups()
	if len(listAfter) != 0 {
		t.Fatalf("expected 0 workgroups after delete, got %d", len(listAfter))
	}
}

func TestFolderBindingCRUD(t *testing.T) {
	_, cleanup := setupTestDB(t)
	defer cleanup()

	rec := FolderBindingRecord{
		FolderPath:   filepath.Clean("/workspace/repos/sub-repo"),
		SSHKeyPath:   "/keys/id_ed25519",
		GitUserName:  "Sub Developer",
		GitUserEmail: "sub@dev.org",
	}

	bound, err := BindFolder(rec)
	if err != nil {
		t.Fatalf("failed to bind folder: %v", err)
	}
	if bound.ID <= 0 {
		t.Fatalf("expected positive ID, got %d", bound.ID)
	}

	bindings, errList := ListFolderBindings()
	if errList != nil || len(bindings) != 1 {
		t.Fatalf("expected 1 binding, got %d, err: %v", len(bindings), errList)
	}

	if errUnbind := UnbindFolder(rec.FolderPath); errUnbind != nil {
		t.Fatalf("failed to unbind folder: %v", errUnbind)
	}

	bindingsAfter, _ := ListFolderBindings()
	if len(bindingsAfter) != 0 {
		t.Fatalf("expected 0 bindings after unbind, got %d", len(bindingsAfter))
	}
}

func TestLineageResolver_AncestorAndDefault(t *testing.T) {
	_, cleanup := setupTestDB(t)
	defer cleanup()

	parentDir := filepath.Clean("/workspace/client-corp")
	subRepoDir := filepath.Join(parentDir, "backend", "auth-service")

	rec := WorkgroupRecord{
		Name:          "client-corp",
		RootDirectory: parentDir,
		SSHKeyPath:    "/keys/client_key",
		GitUserName:   "Corp Dev",
		GitUserEmail:  "corp@corp.com",
		IsDefault:     false,
	}
	_, err := AddWorkgroup(rec)
	if err != nil {
		t.Fatalf("failed to add workgroup: %v", err)
	}

	eff, errRes := ResolveWorkgroupForPath(subRepoDir)
	if errRes != nil {
		t.Fatalf("failed to resolve ancestor: %v", errRes)
	}
	if !eff.HasActiveGroup || !eff.IsInherited {
		t.Fatalf("expected active group and isInherited=true")
	}
	if eff.SSHKeyPath != "/keys/client_key" {
		t.Fatalf("expected inherited key /keys/client_key, got %s", eff.SSHKeyPath)
	}
	if eff.GitUserName != "Corp Dev" {
		t.Fatalf("expected inherited git author Corp Dev, got %s", eff.GitUserName)
	}

	unrelatedDir := filepath.Clean("/some/unrelated/path")
	effUnrelated, errUnrelated := ResolveWorkgroupForPath(unrelatedDir)
	if errUnrelated != nil {
		t.Fatalf("failed unrelated resolve: %v", errUnrelated)
	}
	if effUnrelated.HasActiveGroup {
		t.Fatalf("expected no active group for unrelated path")
	}
}

func TestRunWorkgroupDispatch(t *testing.T) {
	_, cleanup := setupTestDB(t)
	defer cleanup()

	if err := RunWorkgroup([]string{"help"}); err != nil {
		t.Fatalf("RunWorkgroup help failed: %v", err)
	}

	if err := RunWorkgroup([]string{"ls"}); err != nil {
		t.Fatalf("RunWorkgroup ls failed: %v", err)
	}

	errAdd := RunWorkgroup([]string{
		"add", "testgroup", "/tmp/testdir",
		"--ssh-key", "/tmp/testkey",
		"--desc", "Testing group",
	})
	if errAdd != nil {
		t.Fatalf("RunWorkgroup add failed: %v", errAdd)
	}

	errShow := RunWorkgroup([]string{"show", "testgroup"})
	if errShow != nil {
		t.Fatalf("RunWorkgroup show failed: %v", errShow)
	}

	errRm := RunWorkgroup([]string{"rm", "testgroup"})
	if errRm != nil {
		t.Fatalf("RunWorkgroup rm failed: %v", errRm)
	}
}

func TestWorkgroupPositionalBind(t *testing.T) {
	opts := ParseWorkgroupOptions([]string{"bind", "/tmp/client-a", "/tmp/id_rsa"})
	if opts.Subcommand != "bind" {
		t.Fatalf("expected subcommand bind, got %s", opts.Subcommand)
	}
	if opts.Name != "/tmp/client-a" {
		t.Fatalf("expected name /tmp/client-a, got %s", opts.Name)
	}
	if opts.SSHKeyPath != "/tmp/id_rsa" {
		t.Fatalf("expected SSHKeyPath /tmp/id_rsa, got %s", opts.SSHKeyPath)
	}
}
