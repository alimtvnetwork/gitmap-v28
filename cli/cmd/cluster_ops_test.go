package cmd

import (
	"context"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func setupClusterOpsTestDB(t *testing.T) (string, *store.DB) {
	t.Helper()
	dbPath := filepath.Join(t.TempDir(), "test_cluster_ops.db")
	dbConn, err := store.OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	_ = dbConn.Migrate()
	_ = store.EnsureSSHTables(dbConn.Conn())
	return dbPath, dbConn
}

func withMockClusterOpsStore(t *testing.T, fn func(testDB *store.DB)) {
	t.Helper()
	dbPath, testDB := setupClusterOpsTestDB(t)
	defer testDB.Close()
	prevOpener := openClusterStore
	openClusterStore = func() (*store.DB, error) { return store.OpenAt(dbPath) }
	defer func() { openClusterStore = prevOpener }()
	fn(testDB)
}

func seedTestHost(t *testing.T, dbConn *store.DB, host store.SSHHost) {
	t.Helper()
	ctx := context.Background()
	err := store.UpsertSSHHost(ctx, host, dbConn.Conn())
	if err != nil {
		t.Fatalf("failed to seed host: %v", err)
	}
}

func TestRunClusterNodes_Help(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	flags := []string{"--help", "-h", "help"}
	for _, flag := range flags {
		err := runClusterNodes([]string{flag})
		if err != nil {
			t.Fatalf("expected nil error on help flag %s, got: %v", flag, err)
		}
	}
}

func TestRunClusterNodes_EmptyHosts(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		out := captureStdoutForTest(t, func() {
			err := runClusterNodes([]string{})
			if err != nil {
				t.Fatalf("expected nil error on empty nodes, got: %v", err)
			}
		})
		hasNotice := strings.Contains(out, "No nodes currently registered in cluster")
		if !hasNotice {
			t.Errorf("expected empty nodes notice in output, got: %s", out)
		}
	})
}

func TestRunClusterNodes_WithHosts(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		h1 := store.SSHHost{ID: "h-1", Alias: "alpha", IP: "10.0.0.1", Username: "root", Port: 22, ClusterRole: "control", CreatedAt: time.Now().UTC()}
		h2 := store.SSHHost{ID: "h-2", Alias: "beta", IP: "10.0.0.2", Username: "worker", Port: 2222, ClusterRole: "worker", CreatedAt: time.Now().UTC()}
		seedTestHost(t, testDB, h1)
		seedTestHost(t, testDB, h2)
		out := captureStdoutForTest(t, func() {
			_ = runClusterNodes([]string{})
		})
		assertNodesTableOutput(t, out)
	})
}

func assertNodesTableOutput(t *testing.T, out string) {
	t.Helper()
	hasHeader := strings.Contains(out, "ALIAS") && strings.Contains(out, "IP")
	hasAlpha := strings.Contains(out, "alpha") && strings.Contains(out, "10.0.0.1")
	hasBeta := strings.Contains(out, "beta") && strings.Contains(out, "10.0.0.2")
	if !hasHeader || !hasAlpha || !hasBeta {
		t.Errorf("missing expected table output, got: %s", out)
	}
}

func seedSecretTestHost(t *testing.T, dbConn *store.DB) {
	t.Helper()
	h := store.SSHHost{
		ID: "h-sec", Alias: "node-sec", IP: "10.0.0.99", Username: "admin",
		Port: 22, EncryptedPassword: "supersecretpassword", ClusterRole: "worker", CreatedAt: time.Now().UTC(),
	}
	seedTestHost(t, dbConn, h)
}

func TestRunClusterNodes_JSON(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		seedSecretTestHost(t, testDB)
		assertNodesJSONOutput(t)
	})
}

func assertNodesJSONOutput(t *testing.T) {
	t.Helper()
	out := captureStdoutForTest(t, func() {
		_ = runClusterNodes([]string{"--json"})
	})
	hasAlias := strings.Contains(out, "node-sec")
	hasSecret := strings.Contains(out, "supersecretpassword")
	if !hasAlias {
		t.Errorf("expected node-sec in JSON output, got: %s", out)
	}
	if hasSecret {
		t.Errorf("expected encrypted password to be redacted in JSON output")
	}
}

func TestRunClusterRemove_Help(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	flags := []string{"--help", "-h", "help"}
	for _, flag := range flags {
		err := runClusterRemove([]string{flag})
		if err != nil {
			t.Fatalf("expected nil error on help flag %s, got: %v", flag, err)
		}
	}
}

func assertPositionalRemove(t *testing.T, target string) {
	t.Helper()
	out := captureStdoutForTest(t, func() {
		err := runClusterRemove([]string{target})
		if err != nil {
			t.Fatalf("expected nil error on positional remove, got: %v", err)
		}
	})
	expected := "✓ Node '" + target + "' removed from cluster registry."
	hasConfirmation := strings.Contains(out, expected)
	if !hasConfirmation {
		t.Errorf("expected confirmation message, got: %s", out)
	}
}

func TestRunClusterRemove_PositionalByAlias(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		h := store.SSHHost{ID: "h-rm1", Alias: "remove-me", IP: "10.0.0.20", Username: "root", Port: 22, ClusterRole: "worker", CreatedAt: time.Now().UTC()}
		seedTestHost(t, testDB, h)
		assertPositionalRemove(t, "remove-me")
	})
}

func TestRunClusterRemove_PositionalByIP(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		h := store.SSHHost{ID: "h-rm2", Alias: "node-ip", IP: "10.0.0.30", Username: "root", Port: 22, ClusterRole: "worker", CreatedAt: time.Now().UTC()}
		seedTestHost(t, testDB, h)
		assertPositionalRemove(t, "10.0.0.30")
	})
}

func TestRunClusterRemove_MissingArgs(t *testing.T) {
	prev := cliexit.SetExitFunc(func(code int) {})
	defer cliexit.SetExitFunc(prev)

	err := runClusterRemove([]string{})
	hasErr := err != nil
	if !hasErr {
		t.Error("expected error on missing args for runClusterRemove, got nil")
	}
}

func TestRunClusterRemove_LegacyID(t *testing.T) {
	withMockClusterOpsStore(t, func(testDB *store.DB) {
		out := captureStdoutForTest(t, func() {
			err := runClusterRemove([]string{"--id", "test-node", "--confirm"})
			if err != nil {
				t.Fatalf("expected nil on legacy remove, got: %v", err)
			}
		})
		hasDeleted := strings.Contains(out, "Node deleted successfully.")
		if !hasDeleted {
			t.Errorf("expected deletion message, got: %s", out)
		}
	})
}

func TestHasPositionalNodeTarget(t *testing.T) {
	hasPos1 := hasPositionalNodeTarget([]string{"my-node"})
	hasPos2 := hasPositionalNodeTarget([]string{"--id", "foo"})
	hasPos3 := hasPositionalNodeTarget([]string{})
	if !hasPos1 || hasPos2 || hasPos3 {
		t.Errorf("unexpected hasPositionalNodeTarget results: %v, %v, %v", hasPos1, hasPos2, hasPos3)
	}
}

func TestHasClusterJSONFlag(t *testing.T) {
	hasJ1 := hasClusterJSONFlag([]string{"--json"})
	hasJ2 := hasClusterJSONFlag([]string{"--other"})
	hasJ3 := hasClusterJSONFlag([]string{})
	if !hasJ1 || hasJ2 || hasJ3 {
		t.Errorf("unexpected hasClusterJSONFlag results: %v, %v, %v", hasJ1, hasJ2, hasJ3)
	}
}
