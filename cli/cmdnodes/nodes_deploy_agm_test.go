package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func TestPackageLocalAGMAccounts_Success(t *testing.T) {
	tmpDir := t.TempDir()

	accountsDir := filepath.Join(tmpDir, "accounts")
	if err := os.MkdirAll(accountsDir, 0755); err != nil {
		t.Fatalf("failed to create accounts dir: %v", err)
	}

	acc1Content := []byte(`{"id": "acc-1", "email": "user1@example.com"}`)
	acc2Content := []byte(`{"id": "acc-2", "email": "user2@example.com"}`)
	_ = os.WriteFile(filepath.Join(accountsDir, "acc1.json"), acc1Content, 0600)
	_ = os.WriteFile(filepath.Join(accountsDir, "acc2.json"), acc2Content, 0600)

	accountsJSONContent := []byte(`{"version": "2.0", "accounts": [{"id": "acc-1"}]}`)
	_ = os.WriteFile(filepath.Join(tmpDir, "accounts.json"), accountsJSONContent, 0600)

	tokensContent := []byte(`SQLite format 3 tokens mock data`)
	_ = os.WriteFile(filepath.Join(tmpDir, "user_tokens.db"), tokensContent, 0600)

	tarData, count, err := PackageLocalAGMAccounts(tmpDir)
	if err != nil {
		t.Fatalf("PackageLocalAGMAccounts returned error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 accounts, got %d", count)
	}

	filesInArchive := extractTarGzFileMap(t, tarData)
	assertFileInArchive(t, filesInArchive, "accounts/acc1.json", acc1Content)
	assertFileInArchive(t, filesInArchive, "accounts/acc2.json", acc2Content)
	assertFileInArchive(t, filesInArchive, "accounts.json", accountsJSONContent)
	assertFileInArchive(t, filesInArchive, "user_tokens.db", tokensContent)
}

func TestPackageLocalAGMAccounts_AccountsJSONOnly(t *testing.T) {
	tmpDir := t.TempDir()

	accountsJSONContent := []byte(`{"accounts": [{"id": "1"}, {"id": "2"}, {"id": "3"}]}`)
	_ = os.WriteFile(filepath.Join(tmpDir, "accounts.json"), accountsJSONContent, 0600)

	tarData, count, err := PackageLocalAGMAccounts(tmpDir)
	if err != nil {
		t.Fatalf("PackageLocalAGMAccounts returned error: %v", err)
	}
	if count != 3 {
		t.Errorf("expected 3 accounts parsed from accounts.json, got %d", count)
	}
	if len(tarData) == 0 {
		t.Errorf("expected non-empty tar.gz bytes")
	}
}

func TestPackageLocalAGMAccounts_NoAccountsError(t *testing.T) {
	tmpDir := t.TempDir()

	_, _, err := PackageLocalAGMAccounts(tmpDir)
	if err == nil {
		t.Fatalf("expected error for empty dir with no accounts, got nil")
	}
}

func TestFilterDeployAGMNodes_DefaultExceptMain(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "local", IPAddress: "127.0.0.1"},
		{Alias: "main", IPAddress: "10.254.1.10"},
		{Alias: "worker-1", IPAddress: "10.254.1.11"},
		{Alias: "worker-2", IPAddress: "10.254.1.12"},
	}

	filtered := FilterDeployAGMNodes(conns, DeployAGMOptions{})
	if len(filtered) != 2 {
		t.Fatalf("expected 2 nodes (excluding local and main), got %d", len(filtered))
	}

	for _, c := range filtered {
		if c.Alias == "main" || c.Alias == "local" {
			t.Errorf("unexpected node in filtered list: %s", c.Alias)
		}
	}
}

func TestFilterDeployAGMNodes_IncludeMain(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "local", IPAddress: "127.0.0.1"},
		{Alias: "main", IPAddress: "10.254.1.10"},
		{Alias: "worker-1", IPAddress: "10.254.1.11"},
	}

	filtered := FilterDeployAGMNodes(conns, DeployAGMOptions{IncludeMain: true})
	if len(filtered) != 2 {
		t.Fatalf("expected 2 nodes (main and worker-1), got %d", len(filtered))
	}

	hasMain := false
	for _, c := range filtered {
		if c.Alias == "main" {
			hasMain = true
		}
	}
	if !hasMain {
		t.Errorf("expected 'main' node to be included when IncludeMain=true")
	}
}

func TestFilterDeployAGMNodes_Target(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "worker-1", IPAddress: "10.254.1.11"},
		{Alias: "worker-2", IPAddress: "10.254.1.12"},
	}

	filtered := FilterDeployAGMNodes(conns, DeployAGMOptions{Target: "worker-2"})
	if len(filtered) != 1 {
		t.Fatalf("expected 1 node, got %d", len(filtered))
	}
	if filtered[0].Alias != "worker-2" {
		t.Errorf("expected worker-2, got %s", filtered[0].Alias)
	}
}

func TestFilterDeployAGMNodes_CustomExcept(t *testing.T) {
	conns := []db.SSHConnection{
		{Alias: "main", IPAddress: "10.254.1.10"},
		{Alias: "worker-1", IPAddress: "10.254.1.11"},
		{Alias: "worker-2", IPAddress: "10.254.1.12"},
	}

	opts := DeployAGMOptions{Except: []string{"worker-1"}}
	filtered := FilterDeployAGMNodes(conns, opts)
	if len(filtered) != 1 {
		t.Fatalf("expected 1 node (worker-2), got %d", len(filtered))
	}
	if filtered[0].Alias != "worker-2" {
		t.Errorf("expected worker-2, got %s", filtered[0].Alias)
	}
}

func TestParseDeployAGMOptions(t *testing.T) {
	args := []string{
		"--target", "worker-1",
		"--except", "worker-2,worker-3",
		"--include-main",
		"--open-only",
		"--dry-run",
		"--json",
	}

	opts, err := ParseDeployAGMOptions(args)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	if opts.Target != "worker-1" {
		t.Errorf("expected target 'worker-1', got %q", opts.Target)
	}
	if len(opts.Except) != 2 || opts.Except[0] != "worker-2" || opts.Except[1] != "worker-3" {
		t.Errorf("unexpected except list: %v", opts.Except)
	}
	if !opts.IncludeMain {
		t.Errorf("expected IncludeMain to be true")
	}
	if !opts.OpenOnly {
		t.Errorf("expected OpenOnly to be true")
	}
	if !opts.DryRun {
		t.Errorf("expected DryRun to be true")
	}
	if !opts.IsJSON {
		t.Errorf("expected IsJSON to be true")
	}
}

func TestFilterOpenOnlyResults(t *testing.T) {
	results := []DeployAGMResult{
		{NodeAlias: "w1", Status: "SUCCESS"},
		{NodeAlias: "w2", Status: "OFFLINE"},
		{NodeAlias: "w3", Status: "DRY-RUN"},
	}

	all := filterOpenOnlyResults(results, false)
	if len(all) != 3 {
		t.Errorf("expected all 3 results when openOnly=false, got %d", len(all))
	}

	openOnly := filterOpenOnlyResults(results, true)
	if len(openOnly) != 2 {
		t.Fatalf("expected 2 open results when openOnly=true, got %d", len(openOnly))
	}
	if openOnly[0].NodeAlias != "w1" || openOnly[1].NodeAlias != "w3" {
		t.Errorf("unexpected filtered results: %v", openOnly)
	}
}

func extractTarGzFileMap(t *testing.T, data []byte) map[string][]byte {
	t.Helper()
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("gzip.NewReader failed: %v", err)
	}
	defer gr.Close()

	tr := tar.NewReader(gr)
	files := make(map[string][]byte)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("tar.Next failed: %v", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		content, readErr := io.ReadAll(tr)
		if readErr != nil {
			t.Fatalf("read tar file %q failed: %v", hdr.Name, readErr)
		}
		files[hdr.Name] = content
	}
	return files
}

func assertFileInArchive(t *testing.T, files map[string][]byte, name string, expected []byte) {
	t.Helper()
	data, ok := files[name]
	if !ok {
		t.Fatalf("missing expected file %q in archive", name)
	}
	if !bytes.Equal(data, expected) {
		t.Errorf("content mismatch for %q: expected %q, got %q", name, string(expected), string(data))
	}
}
