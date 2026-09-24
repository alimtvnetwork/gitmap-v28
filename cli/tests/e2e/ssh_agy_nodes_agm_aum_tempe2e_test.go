//go:build tempe2e

package e2e

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmd"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/searcher"
)

func requireTempE2EEnvSpec150(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("Skipping temporary E2E test: set RUN_TEMP_E2E=1 with -tags=tempe2e to execute on-demand")
	}
}

func TestTempE2E_SSHAgyNodesExportImportOnelinerAndDeployNodeConfig(t *testing.T) {
	requireTempE2EEnvSpec150(t)

	tmpDir := t.TempDir()
	origDir, _ := os.Getwd()
	_ = os.Chdir(tmpDir)
	defer func() { _ = os.Chdir(origDir) }()

	sampleConns := []db.SSHConnection{
		{Alias: "alpha-win", IPAddress: "10.20.0.11", Username: "admin", OS: "windows", KeyPath: "~/.ssh/id_ed25519"},
		{Alias: "beta-linux", IPAddress: "10.20.0.12", Username: "ubuntu", OS: "linux", KeyPath: "~/.ssh/id_ed25519"},
		{Alias: "gamma-mac", IPAddress: "10.20.0.13", Username: "devops", OS: "darwin", EncryptedPassword: "enc"},
	}

	// 1. Test FilterSSHConnectionsByExcept by numeric ID, worker-ID, IP, and alias
	filtered := cmdssh.FilterSSHConnectionsByExcept(sampleConns, "1,10.20.0.13")
	if len(filtered) != 1 || filtered[0].Alias != "beta-linux" {
		t.Fatalf("expected only beta-linux after --except 1,10.20.0.13, got %+v", filtered)
	}

	// 2. Test default gitmap-ssh-nodes.json export and import
	env := cmdssh.SSHNodesExportEnvelope{
		SchemaVersion: "1.0",
		ExportedAt:    time.Now().UTC().Format(time.RFC3339),
		TotalNodes:    len(sampleConns),
		Connections:   sampleConns,
	}
	rawBytes, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal envelope: %v", err)
	}
	defaultFile := filepath.Join(tmpDir, cmdssh.DefaultSSHNodesJSONFile)
	if err := os.WriteFile(defaultFile, rawBytes, 0644); err != nil {
		t.Fatalf("failed to write %s: %v", defaultFile, err)
	}
	if err := cmdssh.RunSSHNodesImportJSON(nil); err != nil {
		t.Fatalf("RunSSHNodesImportJSON default file failed: %v", err)
	}
	if err := cmdssh.RunSSHNodesExportJSON([]string{"custom-nodes.json"}); err != nil {
		t.Fatalf("RunSSHNodesExportJSON custom-nodes.json failed: %v", err)
	}
	if _, statErr := os.Stat(filepath.Join(tmpDir, "custom-nodes.json")); statErr != nil {
		t.Fatalf("expected custom-nodes.json to exist: %v", statErr)
	}

	// 3. Test export-oneliner and --base64 import
	compactBytes, _ := json.Marshal(env)
	b64 := base64.StdEncoding.EncodeToString(compactBytes)
	if err := cmdssh.RunSSHNodesImportJSON([]string{"--base64", b64}); err != nil {
		t.Fatalf("RunSSHNodesImportJSON --base64 failed: %v", err)
	}
	if err := cmdssh.RunSSHExportOnelinerCLI([]string{"--raw"}); err != nil {
		t.Fatalf("RunSSHExportOnelinerCLI failed: %v", err)
	}

	// 4. Test ssh deploy node-config (nc) --except and ssh agy / agy ssh dry-run
	if err := cmdssh.RunSSHDeployNodeConfigCLI([]string{"nc", "--except", "1,gamma-mac", "--dry-run", "--json"}); err != nil {
		t.Fatalf("RunSSHDeployNodeConfigCLI failed: %v", err)
	}
	if err := cmdssh.RunSSHAgyCLI([]string{"--except", "worker-1,10.20.0.13", "status", "--dry-run"}); err != nil {
		t.Fatalf("RunSSHAgyCLI failed: %v", err)
	}
}

func TestTempE2E_AUMSearchHistoryDH2DAndAIMemoryMultiPortServer(t *testing.T) {
	requireTempE2EEnvSpec150(t)

	query := "SSHConnection"
	dh2d := searcher.ComputeDH2D(query, "keyword")
	if !strings.HasPrefix(dh2d, "DH2D-") || len(dh2d) != 13 {
		t.Fatalf("expected DH2D-XXXXXXXX format, got %q", dh2d)
	}

	fakeMatches := []searcher.SearchResult{
		{MatchedText: "SSHConnection", FilePath: "cli/cmdssh/ssh.go", RelativePath: "cli/cmdssh/ssh.go", StartPosition: 10, EndPosition: 23},
	}

	// Hit 1: Warm tier
	id1, err := searcher.RecordAUMSearchExecution(query, "keyword", true, 1, fakeMatches)
	if err != nil || id1 != dh2d {
		t.Fatalf("RecordAUMSearchExecution hit 1 failed: id=%s err=%v", id1, err)
	}

	// Hit 2: Auto-promoted to HOT_MEMORY_CACHE (<0.04ms)
	_, _ = searcher.RecordAUMSearchExecution(query, "keyword", true, 0, fakeMatches)
	cached, gotID, isHot := searcher.LookupHotCachedSearch(query, "keyword")
	if !isHot || gotID != dh2d || len(cached) != 1 {
		t.Fatalf("expected query %q (%s) to be hot-cached after 2 hits, got isHot=%v len=%d", query, gotID, isHot, len(cached))
	}

	// Start In-Memory AI Server on the 4 unique candidate ports (47831..47834)
	ln, boundPort, srv, srvErr := cmd.StartInMemoryAIServerOnCandidatePorts(cmd.UniqueAICandidatePorts)
	if srvErr != nil {
		t.Fatalf("StartInMemoryAIServerOnCandidatePorts failed: %v", srvErr)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_ = srv.Shutdown(ctx)
		_ = ln.Close()
	}()

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%d/api/v1/ai/search-cache?q=%s", boundPort, query))
	if err != nil {
		t.Fatalf("GET /api/v1/ai/search-cache failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), dh2d) {
		t.Fatalf("expected AI memory server response to contain DH2D ID %s, got %s", dh2d, string(body))
	}
}
