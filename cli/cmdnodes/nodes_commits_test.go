// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type mockSSHExecutor struct {
	responses map[string]string
	errors    map[string]error
	calls     []string
}

func (m *mockSSHExecutor) Execute(conn db.SSHConnection, command string, shell string) (string, error) {
	m.calls = append(m.calls, conn.Alias+":"+command)
	if err, exists := m.errors[conn.Alias]; exists && err != nil {
		return "", err
	}

	if resp, exists := m.responses[conn.Alias]; exists {
		return resp, nil
	}

	return "{}", nil
}

func TestExtractJSONPayload_Resilience(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "object with motd prefix and suffix",
			input:    "Welcome to Ubuntu 22.04 LTS\n{\"status\":\"ok\",\"dirty\":2}\nLogout",
			expected: "{\"status\":\"ok\",\"dirty\":2}",
		},
		{
			name:     "array with shell noise",
			input:    "Last login: Sun Oct 4 10:00:00\n[{\"repo_name\":\"gitmap\"}]\nConnection closed.",
			expected: "[{\"repo_name\":\"gitmap\"}]",
		},
		{
			name:     "nested array inside object",
			input:    "Banner text\n{\"repos\":[\"a\",\"b\"],\"count\":2}\nDone",
			expected: "{\"repos\":[\"a\",\"b\"],\"count\":2}",
		},
		{
			name:     "nested object inside array",
			input:    "Noise\n[{\"id\":1,\"tags\":{\"env\":\"prod\"}}]\nBye",
			expected: "[{\"id\":1,\"tags\":{\"env\":\"prod\"}}]",
		},
		{
			name:     "clean json object without noise",
			input:    "{\"key\":\"val\"}",
			expected: "{\"key\":\"val\"}",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "plain text without json braces",
			input:    "Plain text output without braces",
			expected: "Plain text output without braces",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractJSONPayload(tc.input)
			if actual != tc.expected {
				t.Errorf("extractJSONPayload(%q) = %q; want %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestQuoteBashArg(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "simple", expected: "'simple'"},
		{input: "hello world", expected: "'hello world'"},
		{input: "it's working", expected: `'it'\''s working'`},
		{input: "", expected: "''"},
		{input: "Feature: add nodes commits", expected: "'Feature: add nodes commits'"},
	}

	for _, tc := range tests {
		actual := quoteBashArg(tc.input)
		if actual != tc.expected {
			t.Errorf("quoteBashArg(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestQuotePowerShellArg(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "simple", expected: `\"simple\"`},
		{input: "hello world", expected: `\"hello world\"`},
		{input: `with "quotes"`, expected: `\"with \"quotes\"\"`},
		{input: "", expected: `\"\"`},
	}

	for _, tc := range tests {
		actual := quotePowerShellArg(tc.input)
		if actual != tc.expected {
			t.Errorf("quotePowerShellArg(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestParseNodesCommitsArgs(t *testing.T) {
	args := []string{"gitmap", "add new feature", "-t", "u1", "--dry-run", "--no-push"}
	opts, err := parseNodesCommitsOptions("cpf", args)
	if err != nil {
		t.Fatalf("unexpected error parsing commits args: %v", err)
	}

	if opts.Action != "cpf" {
		t.Errorf("expected action cpf, got %s", opts.Action)
	}

	if opts.TargetScope != "gitmap" {
		t.Errorf("expected targetScope gitmap, got %s", opts.TargetScope)
	}

	if opts.CommitMessage != "add new feature" {
		t.Errorf("expected commit message 'add new feature', got %s", opts.CommitMessage)
	}

	if opts.FilterOpts.Target != "u1" {
		t.Errorf("expected filter target u1, got %s", opts.FilterOpts.Target)
	}

	if !opts.IsDryRun {
		t.Errorf("expected IsDryRun true, got false")
	}

	if opts.IsPushEnabled {
		t.Errorf("expected IsPushEnabled false with --no-push, got true")
	}
}

func TestParseNodesCommitsArgs_MissingMessage(t *testing.T) {
	_, err := parseNodesCommitsOptions("cpf", []string{"gitmap"})
	if err == nil {
		t.Fatal("expected error for missing commit message, got nil")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Code != "E9080" {
		t.Errorf("expected error code E9080, got %s", appErr.Code)
	}
}

func TestParseNodesCommitsArgs_MissingTarget(t *testing.T) {
	_, err := parseNodesCommitsOptions("cpf", []string{})
	if err == nil {
		t.Fatal("expected error for missing target, got nil")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}

	if appErr.Code != "E9081" {
		t.Errorf("expected error code E9081, got %s", appErr.Code)
	}
}

func TestResolveRemoteCommitsCommand(t *testing.T) {
	opts := NodesCommitsOptions{
		Action:        "cpf",
		TargetScope:   "gitmap",
		CommitMessage: "add suite",
		IsDryRun:      true,
		IsPushEnabled: false,
	}

	linuxConn := db.SSHConnection{OS: "linux"}
	cmdLinux, shellLinux := resolveRemoteCommitsCommand(linuxConn, opts)
	if shellLinux != "bash" {
		t.Errorf("expected bash shell for linux, got %s", shellLinux)
	}

	if !strings.Contains(cmdLinux, "gitmap sends cpf 'gitmap' 'add suite' --json --dry-run --no-push") {
		t.Errorf("unexpected linux remote command: %s", cmdLinux)
	}

	winConn := db.SSHConnection{OS: "windows"}
	cmdWin, shellWin := resolveRemoteCommitsCommand(winConn, opts)
	if shellWin != "ps" {
		t.Errorf("expected ps shell for windows, got %s", shellWin)
	}

	if !strings.Contains(cmdWin, "powershell.exe -NoProfile -Command") {
		t.Errorf("expected powershell invocation, got %s", cmdWin)
	}

	if !strings.Contains(cmdWin, `gitmap sends cpf \"gitmap\" \"add suite\" --json --dry-run --no-push`) {
		t.Errorf("unexpected windows remote command: %s", cmdWin)
	}
}

func TestRenderNodesCommitsTable(t *testing.T) {
	results := []NodeCommitResult{
		{
			NodeAlias:    "u1",
			Host:         "192.168.1.101",
			OS:           "linux",
			Action:       "cpf",
			IsNodeOnline: true,
			IsSuccess:    true,
			IsDryRun:     false,
			Latency:      150 * time.Millisecond,
			Repositories: []RemoteRepoCommitOutcome{
				{
					RepoName:     "gitmap",
					Branch:       "main",
					CommitHash:   "a1b2c3d4e5f6",
					CommitMsg:    "Feature: add cluster node commit delegation suite",
					IsSuccess:    true,
					IsPushed:     true,
					HasChanges:   true,
					FilesChanged: 4,
				},
			},
			CommittedCount: 1,
		},
		{
			NodeAlias:    "w1",
			Host:         "192.168.1.102",
			OS:           "windows",
			Action:       "cpf",
			IsNodeOnline: true,
			IsSuccess:    true,
			IsDryRun:     true,
			Latency:      200 * time.Millisecond,
			Repositories: []RemoteRepoCommitOutcome{
				{
					RepoName:     "service-a",
					Branch:       "feature",
					CommitHash:   "",
					CommitMsg:    "dry run",
					IsSuccess:    true,
					IsPushed:     false,
					HasChanges:   true,
					FilesChanged: 2,
				},
			},
			CommittedCount: 1,
		},
		{
			NodeAlias:    "off1",
			Host:         "192.168.1.103",
			OS:           "linux",
			Action:       "cpf",
			IsNodeOnline: false,
			IsSuccess:    false,
		},
	}

	renderNodesCommitsTable(results, "cpf")
}

func TestRenderNodesPendingCommitsTable(t *testing.T) {
	results := []NodePendingCommitsResult{
		{
			NodeAlias:    "u1",
			Host:         "192.168.1.101",
			OS:           "linux",
			IsNodeOnline: true,
			IsSuccess:    true,
			Latency:      80 * time.Millisecond,
			PendingRepos: []NodePendingRepoItem{
				{
					RepoName:           "gitmap",
					RelativePath:       "work/gitmap",
					Branch:             "main",
					IsClean:            false,
					HasChanges:         true,
					DirtyFileCount:     3,
					HasUnpushedCommits: true,
					UnpushedCount:      1,
					PriorityScore:      104,
					ChangedFiles:       []string{"cli/cmdnodes/nodes_commits.go"},
				},
			},
			TotalDirtyFiles: 3,
			TotalUnpushed:   1,
		},
		{
			NodeAlias:    "u2",
			Host:         "192.168.1.102",
			OS:           "linux",
			IsNodeOnline: true,
			IsSuccess:    true,
			Latency:      50 * time.Millisecond,
			PendingRepos: nil,
		},
		{
			NodeAlias:    "off1",
			Host:         "192.168.1.103",
			OS:           "linux",
			IsNodeOnline: false,
			IsSuccess:    false,
		},
	}

	renderNodesPendingCommitsTable(results)
}

func TestParseNodesPendingCommitsOptions(t *testing.T) {
	args := []string{"-t", "u1", "--sort=name", "--detail", "--dirty-only", "--json"}
	opts := parseNodesPendingCommitsOptions(args)

	if opts.FilterOpts.Target != "u1" {
		t.Errorf("expected filter target u1, got %s", opts.FilterOpts.Target)
	}

	if opts.SortMode != "name" {
		t.Errorf("expected sort mode name, got %s", opts.SortMode)
	}

	if !opts.IsDetail {
		t.Errorf("expected IsDetail true, got false")
	}

	if !opts.IsDirtyOnly {
		t.Errorf("expected IsDirtyOnly true, got false")
	}

	if !opts.IsJSON {
		t.Errorf("expected IsJSON true, got false")
	}
}

func TestParseRemotePendingRepos(t *testing.T) {
	payload := `{
		"pending_repos": [
			{
				"repo_name": "repo-b",
				"branch": "main",
				"dirty_file_count": 1,
				"unpushed_count": 0,
				"has_changes": true
			},
			{
				"repo_name": "repo-a",
				"branch": "dev",
				"dirty_file_count": 5,
				"unpushed_count": 2,
				"has_changes": true,
				"has_unpushed_commits": true
			}
		]
	}`

	reposPriority := parseRemotePendingRepos(payload, "priority", false)
	if len(reposPriority) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(reposPriority))
	}

	if reposPriority[0].RepoName != "repo-a" {
		t.Errorf("expected repo-a highest priority, got %s", reposPriority[0].RepoName)
	}

	reposName := parseRemotePendingRepos(payload, "name", false)
	if reposName[0].RepoName != "repo-a" || reposName[1].RepoName != "repo-b" {
		t.Errorf("expected alphabetical order repo-a, repo-b")
	}

	reposCount := parseRemotePendingRepos(payload, "count", false)
	if reposCount[0].RepoName != "repo-a" {
		t.Errorf("expected repo-a highest count, got %s", reposCount[0].RepoName)
	}
}

func TestMockSSHExecution_PendingCommits(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	mock := &mockSSHExecutor{
		responses: map[string]string{
			"u1": `{"pending_repos":[{"repo_name":"gitmap","branch":"main","dirty_file_count":2,"unpushed_count":1}]}`,
		},
	}
	currentSSHExecutor = mock

	conn := db.SSHConnection{Alias: "u1", IPAddress: "192.168.1.101", OS: "linux"}
	opts := NodesPendingCommitsOptions{SortMode: "priority"}
	res := executeSingleNodePendingCommits(conn, opts)

	if !res.IsNodeOnline {
		t.Errorf("expected node online true, got false")
	}

	if res.IsSuccess {
	} else {
		t.Errorf("expected node success true, got false")
	}

	if len(res.PendingRepos) != 1 {
		t.Fatalf("expected 1 pending repo, got %d", len(res.PendingRepos))
	}

	if res.TotalDirtyFiles != 2 {
		t.Errorf("expected 2 dirty files, got %d", res.TotalDirtyFiles)
	}

	if res.TotalUnpushed != 1 {
		t.Errorf("expected 1 unpushed commit, got %d", res.TotalUnpushed)
	}
}

func TestMockSSHExecution_Commits(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	mock := &mockSSHExecutor{
		responses: map[string]string{
			"u1": `{"repositories":[{"repo_name":"gitmap","branch":"main","commit_hash":"1234567890ab","is_success":true,"files_changed":3}]}`,
		},
	}
	currentSSHExecutor = mock

	conn := db.SSHConnection{Alias: "u1", IPAddress: "192.168.1.101", OS: "linux"}
	opts := NodesCommitsOptions{
		Action:        "cpf",
		TargetScope:   "gitmap",
		CommitMessage: "Feature: test",
		IsPushEnabled: true,
	}
	res := executeSingleNodeCommits(conn, opts)

	if !res.IsNodeOnline {
		t.Errorf("expected node online true, got false")
	}

	if res.IsSuccess {
	} else {
		t.Errorf("expected node success true, got false")
	}

	if len(res.Repositories) != 1 {
		t.Fatalf("expected 1 repository result, got %d", len(res.Repositories))
	}

	if res.CommittedCount != 1 {
		t.Errorf("expected committed count 1, got %d", res.CommittedCount)
	}
}

func TestMockSSHExecution_Offline(t *testing.T) {
	oldExec := currentSSHExecutor
	defer func() { currentSSHExecutor = oldExec }()

	mock := &mockSSHExecutor{
		errors: map[string]error{
			"u1": errors.New("dial tcp 192.168.1.101:22: connect: connection refused"),
		},
	}
	currentSSHExecutor = mock

	conn := db.SSHConnection{Alias: "u1", IPAddress: "192.168.1.101", OS: "linux"}
	opts := NodesCommitsOptions{
		Action:        "cpf",
		TargetScope:   "all",
		CommitMessage: "Feature: test",
	}
	res := executeSingleNodeCommits(conn, opts)

	if res.IsNodeOnline {
		t.Errorf("expected node online false for connection refused, got true")
	}

	if res.IsSuccess {
		t.Errorf("expected node success false, got true")
	}
}
