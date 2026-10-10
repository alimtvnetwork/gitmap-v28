package cmdnodes

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

type mockSummarySSHExecutor struct {
	mu        sync.Mutex
	responses map[string]string
	executed  []string
}

func (m *mockSummarySSHExecutor) Execute(conn db.SSHConnection, command string, shell string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executed = append(m.executed, command)
	if resp, ok := m.responses[command]; ok {
		return resp, nil
	}
	return "{}", nil
}

func TestParseNodesSummaryOptions(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		wantWithPE   bool
		wantJSON     bool
		wantForceAll bool
		wantReleases int
	}{
		{
			name:         "default fs",
			args:         []string{"fs"},
			wantWithPE:   false,
			wantJSON:     false,
			wantReleases: 3,
		},
		{
			name:         "fs+pe preserves PE",
			args:         []string{"fs+pe"},
			wantWithPE:   true,
			wantJSON:     false,
			wantReleases: 3,
		},
		{
			name:         "fspe alias preserves PE",
			args:         []string{"fspe"},
			wantWithPE:   true,
			wantJSON:     false,
			wantReleases: 3,
		},
		{
			name:         "full summary+pe preserves PE and json",
			args:         []string{"full", "summary+pe", "--json"},
			wantWithPE:   true,
			wantJSON:     true,
			wantReleases: 3,
		},
		{
			name:         "full status+pe preserves PE",
			args:         []string{"full", "status+pe"},
			wantWithPE:   true,
			wantJSON:     false,
			wantReleases: 3,
		},
		{
			name:         "custom release limit and force all",
			args:         []string{"fs+pe", "8", "--force-all", "--json"},
			wantWithPE:   true,
			wantJSON:     true,
			wantForceAll: true,
			wantReleases: 8,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			opts := parseNodesSummaryOptions(tc.args)
			if opts.withPE != tc.wantWithPE {
				t.Errorf("withPE = %v; want %v", opts.withPE, tc.wantWithPE)
			}
			if opts.isJSON != tc.wantJSON {
				t.Errorf("isJSON = %v; want %v", opts.isJSON, tc.wantJSON)
			}
			if opts.forceAll != tc.wantForceAll {
				t.Errorf("forceAll = %v; want %v", opts.forceAll, tc.wantForceAll)
			}
			if opts.releasesCount != tc.wantReleases {
				t.Errorf("releasesCount = %d; want %d", opts.releasesCount, tc.wantReleases)
			}
		})
	}
}

func TestParseNodesPEAllOptions(t *testing.T) {
	opts := parseNodesPEAllOptions([]string{"--json", "--force-all", "--verbose"})
	if !opts.isJSON {
		t.Errorf("expected isJSON true")
	}
	if !opts.forceAll {
		t.Errorf("expected forceAll true")
	}
	if !opts.verbose {
		t.Errorf("expected verbose true")
	}
}

func TestIsNodesFullSummary_Preservation(t *testing.T) {
	cases := [][]string{
		{"fs"},
		{"fs+pe"},
		{"fspe"},
		{"full-status"},
		{"full-summary"},
		{"full-status+pe"},
		{"full-summary+pe"},
		{"full", "status"},
		{"full", "summary"},
		{"full", "status+pe"},
		{"full", "summary+pe"},
	}

	for _, c := range cases {
		if !isNodesFullSummary(c) {
			t.Errorf("expected isNodesFullSummary(%v) = true", c)
		}
	}
}

func TestIsNodesPipelineAll(t *testing.T) {
	cases := [][]string{
		{"pe", "all"},
		{"pipe", "all"},
		{"pipeline", "all"},
		{"pipe-error-all"},
		{"pipeline-error-all"},
		{"pipe-errors-all"},
		{"pe-all"},
		{"pe_all"},
		{"pipe", "error", "all"},
		{"pipeline", "errors", "all"},
	}

	for _, c := range cases {
		if !isNodesPipelineAll(c) {
			t.Errorf("expected isNodesPipelineAll(%v) = true", c)
		}
	}
}

func TestParseRemoteManifest_FlexibleFormats(t *testing.T) {
	standardJSON := `[
		{"repoName": "auth-service", "slug": "auth-service", "remoteUrl": "https://github.com/org/auth-service.git", "path": "/home/dev/auth"},
		{"repoName": "billing-service", "slug": "billing-service", "remoteUrl": "https://github.com/org/billing-service.git", "path": "/home/dev/billing"}
	]`
	items := parseRemoteManifest(standardJSON)
	if len(items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(items))
	}
	if items[0].RepoName != "auth-service" || items[0].Slug != "auth-service" {
		t.Errorf("unexpected item 0: %+v", items[0])
	}

	snakeJSON := `[
		{"repo_name": "worker-queue", "slug": "worker-queue", "remote_url": "git@github.com:org/worker-queue.git"}
	]`
	itemsSnake := parseRemoteManifest(snakeJSON)
	if len(itemsSnake) != 1 || itemsSnake[0].RemoteURL != "git@github.com:org/worker-queue.git" {
		t.Errorf("unexpected snake_case parse: %+v", itemsSnake)
	}

	envelopeJSON := `{"data": [
		{"name": "env-service", "url": "https://github.com/org/env-service.git"}
	]}`
	itemsEnv := parseRemoteManifest(envelopeJSON)
	if len(itemsEnv) != 1 || itemsEnv[0].RepoName != "env-service" {
		t.Errorf("unexpected envelope parse: %+v", itemsEnv)
	}
}

func TestNodesSummary_DeduplicationAndPrecedence(t *testing.T) {
	oldExec := currentSSHExecutor
	oldFetch := fetchAllSSHConnectionsHook
	oldResolve := resolveWorkspaceReposHook
	defer func() {
		currentSSHExecutor = oldExec
		fetchAllSSHConnectionsHook = oldFetch
		resolveWorkspaceReposHook = oldResolve
	}()

	localRepos := make([]model.ScanRecord, 65)
	for i := 0; i < 65; i++ {
		name := fmt.Sprintf("shared-repo-%02d", i+1)
		localRepos[i] = model.ScanRecord{
			RepoName:     name,
			Slug:         name,
			HTTPSUrl:     fmt.Sprintf("https://github.com/org/%s.git", name),
			AbsolutePath: "/local/" + name,
		}
	}
	resolveWorkspaceReposHook = func(cwd string) []model.ScanRecord {
		return localRepos
	}

	testConn := db.SSHConnection{
		Alias:     "worker-1",
		IPAddress: "192.168.1.50",
		OS:        "linux",
	}
	fetchAllSSHConnectionsHook = func() ([]db.SSHConnection, error) {
		return []db.SSHConnection{testConn}, nil
	}

	var remoteManifest []RemoteRepoManifestItem
	for i := 0; i < 65; i++ {
		name := fmt.Sprintf("shared-repo-%02d", i+1)
		remoteManifest = append(remoteManifest, RemoteRepoManifestItem{
			RepoName:  name,
			Slug:      name,
			RemoteURL: fmt.Sprintf("https://github.com/org/%s.git", name),
			Path:      "/remote/" + name,
		})
	}
	remoteManifest = append(remoteManifest, RemoteRepoManifestItem{
		RepoName:  "unique-worker-repo-1",
		Slug:      "unique-worker-repo-1",
		RemoteURL: "https://github.com/org/unique-worker-repo-1.git",
		Path:      "/remote/unique-worker-repo-1",
	})
	remoteManifest = append(remoteManifest, RemoteRepoManifestItem{
		RepoName:  "unique-worker-repo-2",
		Slug:      "unique-worker-repo-2",
		RemoteURL: "https://github.com/org/unique-worker-repo-2.git",
		Path:      "/remote/unique-worker-repo-2",
	})

	manifestBytes, _ := json.Marshal(remoteManifest)

	mockExec := &mockSummarySSHExecutor{
		responses: map[string]string{
			"gitmap list --json": string(manifestBytes),
			"gitmap fs 3":        `{"status": "ok", "summary": "remote summary evaluated"}`,
		},
	}
	currentSSHExecutor = mockExec

	err := RunNodesFullSummary([]string{"fs", "3"})
	if err != nil {
		t.Fatalf("unexpected error running nodes full summary: %v", err)
	}

	mockExec.mu.Lock()
	executedCommands := append([]string{}, mockExec.executed...)
	mockExec.mu.Unlock()

	hasDiscovery := false
	hasDelegation := false
	for _, cmd := range executedCommands {
		if cmd == "gitmap list --json" {
			hasDiscovery = true
		}
		if cmd == "gitmap fs 3" {
			hasDelegation = true
		}
	}

	if !hasDiscovery {
		t.Errorf("expected remote discovery command 'gitmap list --json' to be executed")
	}
	if !hasDelegation {
		t.Errorf("expected targeted remote command 'gitmap fs 3' to be executed for unique repos")
	}
}

func TestNodesSummary_AllReposLocal_ZeroRemoteExecution(t *testing.T) {
	oldExec := currentSSHExecutor
	oldFetch := fetchAllSSHConnectionsHook
	oldResolve := resolveWorkspaceReposHook
	defer func() {
		currentSSHExecutor = oldExec
		fetchAllSSHConnectionsHook = oldFetch
		resolveWorkspaceReposHook = oldResolve
	}()

	localRepos := []model.ScanRecord{
		{
			RepoName: "repo-alpha",
			Slug:     "repo-alpha",
			HTTPSUrl: "https://github.com/org/repo-alpha.git",
		},
		{
			RepoName: "repo-beta",
			Slug:     "repo-beta",
			HTTPSUrl: "https://github.com/org/repo-beta.git",
		},
	}
	resolveWorkspaceReposHook = func(cwd string) []model.ScanRecord {
		return localRepos
	}

	testConn := db.SSHConnection{
		Alias:     "worker-2",
		IPAddress: "192.168.1.60",
		OS:        "linux",
	}
	fetchAllSSHConnectionsHook = func() ([]db.SSHConnection, error) {
		return []db.SSHConnection{testConn}, nil
	}

	remoteManifest := []RemoteRepoManifestItem{
		{
			RepoName:  "repo-alpha",
			Slug:      "repo-alpha",
			RemoteURL: "https://github.com/org/repo-alpha.git",
		},
		{
			RepoName:  "repo-beta",
			Slug:      "repo-beta",
			RemoteURL: "https://github.com/org/repo-beta.git",
		},
	}
	manifestBytes, _ := json.Marshal(remoteManifest)

	mockExec := &mockSummarySSHExecutor{
		responses: map[string]string{
			"gitmap list --json": string(manifestBytes),
		},
	}
	currentSSHExecutor = mockExec

	err := RunNodesFullSummary([]string{"fs", "3"})
	if err != nil {
		t.Fatalf("unexpected error running nodes full summary: %v", err)
	}

	mockExec.mu.Lock()
	executedCommands := append([]string{}, mockExec.executed...)
	mockExec.mu.Unlock()

	for _, cmd := range executedCommands {
		if strings.HasPrefix(cmd, "gitmap fs") {
			t.Errorf("expected zero remote execution calls when all repos exist locally, got: %s", cmd)
		}
	}
}

func TestNodesSummary_JSONOutputEnvelope(t *testing.T) {
	oldExec := currentSSHExecutor
	oldFetch := fetchAllSSHConnectionsHook
	oldResolve := resolveWorkspaceReposHook
	defer func() {
		currentSSHExecutor = oldExec
		fetchAllSSHConnectionsHook = oldFetch
		resolveWorkspaceReposHook = oldResolve
	}()

	resolveWorkspaceReposHook = func(cwd string) []model.ScanRecord {
		return []model.ScanRecord{
			{
				RepoName: "core-repo",
				Slug:     "core-repo",
				HTTPSUrl: "https://github.com/org/core-repo.git",
			},
		}
	}

	fetchAllSSHConnectionsHook = func() ([]db.SSHConnection, error) {
		return []db.SSHConnection{
			{
				Alias:     "worker-test",
				IPAddress: "10.0.0.1",
				OS:        "linux",
			},
		}, nil
	}

	remoteManifest := []RemoteRepoManifestItem{
		{
			RepoName:  "remote-unique",
			Slug:      "remote-unique",
			RemoteURL: "https://github.com/org/remote-unique.git",
		},
	}
	manifestBytes, _ := json.Marshal(remoteManifest)

	mockExec := &mockSummarySSHExecutor{
		responses: map[string]string{
			"gitmap list --json": string(manifestBytes),
			"gitmap fs 3 --json": `{"attributes":{"scanned":1},"data":{"status":"ok"}}`,
		},
	}
	currentSSHExecutor = mockExec

	out, err := captureStdout(func() error {
		return RunNodesFullSummary([]string{"fs", "3", "--json"})
	})
	if err != nil {
		t.Fatalf("unexpected error running with --json: %v", err)
	}

	trimmed := strings.TrimSpace(out)
	if !strings.HasPrefix(trimmed, "{") || !strings.HasSuffix(trimmed, "}") {
		t.Fatalf("expected valid JSON object on stdout, got:\n%s", out)
	}

	var envelope FleetNodesEnvelope
	if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil {
		t.Fatalf("failed to unmarshal JSON envelope: %v\nOutput was:\n%s", err, trimmed)
	}

	if envelope.Attributes.TotalNodes != 2 {
		t.Errorf("expected TotalNodes 2, got %d", envelope.Attributes.TotalNodes)
	}
	if envelope.Attributes.LocalReposCount != 1 {
		t.Errorf("expected LocalReposCount 1, got %d", envelope.Attributes.LocalReposCount)
	}
	if envelope.Attributes.DelegatedCount != 1 {
		t.Errorf("expected DelegatedCount 1, got %d", envelope.Attributes.DelegatedCount)
	}
	if !envelope.Data.LocalNode.IsLocalMaster {
		t.Errorf("expected LocalNode.IsLocalMaster to be true")
	}
	if len(envelope.Data.RemoteNodes) != 1 {
		t.Fatalf("expected 1 remote node result, got %d", len(envelope.Data.RemoteNodes))
	}
	if envelope.Data.RemoteNodes[0].NodeAlias != "worker-test" {
		t.Errorf("expected remote node alias 'worker-test', got '%s'", envelope.Data.RemoteNodes[0].NodeAlias)
	}
}

func TestNodesPipelineAll_JSONOutputEnvelope(t *testing.T) {
	oldExec := currentSSHExecutor
	oldFetch := fetchAllSSHConnectionsHook
	oldResolve := resolveWorkspaceReposHook
	defer func() {
		currentSSHExecutor = oldExec
		fetchAllSSHConnectionsHook = oldFetch
		resolveWorkspaceReposHook = oldResolve
	}()

	resolveWorkspaceReposHook = func(cwd string) []model.ScanRecord {
		return []model.ScanRecord{
			{
				RepoName: "pe-local-repo",
				Slug:     "pe-local-repo",
				HTTPSUrl: "https://github.com/org/pe-local-repo.git",
			},
		}
	}

	fetchAllSSHConnectionsHook = func() ([]db.SSHConnection, error) {
		return []db.SSHConnection{
			{
				Alias:     "worker-ci",
				IPAddress: "10.0.0.2",
				OS:        "linux",
			},
		}, nil
	}

	remoteManifest := []RemoteRepoManifestItem{
		{
			RepoName:  "pe-remote-repo",
			Slug:      "pe-remote-repo",
			RemoteURL: "https://github.com/org/pe-remote-repo.git",
		},
	}
	manifestBytes, _ := json.Marshal(remoteManifest)

	mockExec := &mockSummarySSHExecutor{
		responses: map[string]string{
			"gitmap list --json":               string(manifestBytes),
			"gitmap pe all --force-all --json": `{"attributes":{"activeScanned":1},"data":{"cleanCount":1}}`,
		},
	}
	currentSSHExecutor = mockExec

	out, err := captureStdout(func() error {
		return RunNodesPipelineErrorsAll([]string{"pe", "all", "--force-all", "--json"})
	})
	if err != nil {
		t.Fatalf("unexpected error running nodes pe all with --json: %v", err)
	}

	trimmed := strings.TrimSpace(out)
	var envelope FleetNodesEnvelope
	if err := json.Unmarshal([]byte(trimmed), &envelope); err != nil {
		t.Fatalf("failed to unmarshal pe all JSON envelope: %v\nOutput:\n%s", err, trimmed)
	}

	if envelope.Attributes.DelegatedCount != 1 {
		t.Errorf("expected DelegatedCount 1, got %d", envelope.Attributes.DelegatedCount)
	}
}
