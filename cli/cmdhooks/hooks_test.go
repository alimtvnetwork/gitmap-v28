package cmdhooks

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdhooks/enums"
)

func TestDispatchHooks(t *testing.T) {
	testDispatchEmpty(t)
	testDispatchUnknown(t)
	testDispatchHelp(t)
	testDispatchInstall(t)
}

func testDispatchEmpty(t *testing.T) {
	err := DispatchHooks(nil)
	if err == nil {
		t.Fatal("expected error on empty args")
	}
	if !strings.Contains(err.Error(), "missing subcommand") {
		t.Errorf("expected missing subcommand error, got: %v", err)
	}
}

func testDispatchUnknown(t *testing.T) {
	err := DispatchHooks([]string{"invalid-subcmd"})
	if err == nil {
		t.Fatal("expected error on invalid subcommand")
	}
	if !strings.Contains(err.Error(), "unknown subcommand") {
		t.Errorf("expected unknown subcommand error, got: %v", err)
	}
}

func testDispatchHelp(t *testing.T) {
	helpArgs := [][]string{
		{"help"},
		{"--help"},
		{"-h"},
	}
	for _, args := range helpArgs {
		if err := DispatchHooks(args); err != nil {
			t.Errorf("expected nil for help arg %v, got: %v", args, err)
		}
	}
}

func testDispatchInstall(t *testing.T) {
	err := DispatchHooks([]string{"install", "--help"})
	if err != nil {
		t.Errorf("expected nil for install --help, got: %v", err)
	}

	err = DispatchHooks([]string{"install"})
	if err == nil {
		t.Error("expected error for install without --all")
	}
}

func TestRunHooksInstallValidation(t *testing.T) {
	testInstallFlagsMissingAll(t)
	testInstallFlagsInvalid(t)
	testInstallFlagsHelp(t)
}

func testInstallFlagsMissingAll(t *testing.T) {
	err := RunHooksInstall([]string{"--dry-run"})
	if err == nil {
		t.Fatal("expected error when missing --all")
	}
	if !strings.Contains(err.Error(), "missing required flag: --all") {
		t.Errorf("expected missing --all error, got: %v", err)
	}
}

func testInstallFlagsInvalid(t *testing.T) {
	err := RunHooksInstall([]string{"--invalid-flag-xyz"})
	if err == nil {
		t.Fatal("expected error on invalid flag")
	}
}

func testInstallFlagsHelp(t *testing.T) {
	if err := RunHooksInstall([]string{"--help"}); err != nil {
		t.Errorf("expected nil for --help, got: %v", err)
	}
	if err := RunHooksInstall([]string{"-h"}); err != nil {
		t.Errorf("expected nil for -h, got: %v", err)
	}
}

func setupMockRepos(t *testing.T) string {
	tempDir := t.TempDir()
	repo1Git := filepath.Join(tempDir, "repo1", ".git", "hooks")
	repo2Git := filepath.Join(tempDir, "repo2", ".git", "hooks")
	repo3Git := filepath.Join(tempDir, "repo3", ".git", "hooks")
	notARepo := filepath.Join(tempDir, "not_repo")

	if err := os.MkdirAll(repo1Git, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo2Git, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo3Git, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(notARepo, 0755); err != nil {
		t.Fatal(err)
	}

	existingHook := filepath.Join(repo3Git, "commit-msg")
	if err := os.WriteFile(existingHook, []byte("#!/bin/sh\n# gitmap commit-msg guard\nexit 0\n"), 0755); err != nil {
		t.Fatal(err)
	}

	return tempDir
}

func TestRunHooksInstallExecution(t *testing.T) {
	tempDir := setupMockRepos(t)
	testListHookRepos(t, tempDir)
	testPlanHookInstall(t, tempDir)
	testDryRunExecution(t, tempDir)
	testActualInstallExecution(t, tempDir)
	testAlreadyGuardedExecution(t, tempDir)
}

func testListHookRepos(t *testing.T, tempDir string) {
	repos := listHookRepos(tempDir)
	if len(repos) != 3 {
		t.Fatalf("expected 3 repos, got %d: %v", len(repos), repos)
	}
	if repos[0] != "repo1" || repos[1] != "repo2" || repos[2] != "repo3" {
		t.Errorf("unexpected repo list: %v", repos)
	}
}

func testPlanHookInstall(t *testing.T, tempDir string) {
	plan1 := planHookInstall(tempDir, "repo1")
	if !plan1.action.IsInstall() {
		t.Errorf("expected install action for repo1, got: %v", plan1.action)
	}

	plan3 := planHookInstall(tempDir, "repo3")
	if !plan3.action.IsAlready() {
		t.Errorf("expected already action for repo3, got: %v", plan3.action)
	}
}

func testDryRunExecution(t *testing.T, tempDir string) {
	args := []string{"--all", "--workdir", tempDir, "--dry-run"}
	if err := RunHooksInstall(args); err != nil {
		t.Fatalf("expected nil from dry run, got: %v", err)
	}

	hookPath := filepath.Join(tempDir, "repo1", ".git", "hooks", "commit-msg")
	if _, err := os.Stat(hookPath); err == nil {
		t.Error("expected hook file to not exist after dry run")
	}
}

func testActualInstallExecution(t *testing.T, tempDir string) {
	args := []string{"--all", "--workdir", tempDir, "--yes"}
	if err := RunHooksInstall(args); err != nil {
		t.Fatalf("expected nil from real install, got: %v", err)
	}

	hookPath := filepath.Join(tempDir, "repo1", ".git", "hooks", "commit-msg")
	data, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("expected hook file to exist: %v", err)
	}
	if !strings.Contains(string(data), constants.GuardMarker) {
		t.Error("expected hook to contain guard marker")
	}
}

func testAlreadyGuardedExecution(t *testing.T, tempDir string) {
	args := []string{"--all", "--workdir", tempDir, "--yes"}
	if err := RunHooksInstall(args); err != nil {
		t.Fatalf("expected nil from repeat install, got: %v", err)
	}
}

func TestResolveHooksWorkdir(t *testing.T) {
	testWorkdirExplicit(t)
	testWorkdirEnv(t)
}

func testWorkdirExplicit(t *testing.T) {
	custom := "/tmp/custom/repos"
	if res := resolveHooksWorkdir(custom); res != custom {
		t.Errorf("expected %s, got %s", custom, res)
	}
}

func testWorkdirEnv(t *testing.T) {
	envKey := constants.EnvGitmapWorkdir
	orig := os.Getenv(envKey)
	defer os.Setenv(envKey, orig)

	testDir := "/env/test/workdir"
	os.Setenv(envKey, testDir)
	if res := resolveHooksWorkdir(""); res != testDir {
		t.Errorf("expected %s from env, got %s", testDir, res)
	}
}

func TestCountPlanStats(t *testing.T) {
	plans := []hookPlan{
		{repo: "r1", action: enums.HookPlanActionInstall},
		{repo: "r2", action: enums.HookPlanActionAlready},
		{repo: "r3", action: enums.HookPlanActionSkip},
	}
	stats := countPlanStats(plans)
	if stats.installed != 1 || stats.already != 1 || stats.skipped != 1 {
		t.Errorf("unexpected stats: %+v", stats)
	}
}
