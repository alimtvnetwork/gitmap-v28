// Package cmdagy — agy_recreate_test.go provides unit tests for agy recreate-project.
package cmdagy

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func TestNormalizeAgySubcommand_Recreate(t *testing.T) {
	aliases := []string{"recreate-project", "recreate", "rcp", "rec", "RECREATE-PROJECT", "Recreate"}
	for _, alias := range aliases {
		normalized := normalizeAgySubcommand(alias)
		if normalized != "recreate-project" {
			t.Errorf("normalizeAgySubcommand(%q) = %q; want 'recreate-project'", alias, normalized)
		}
	}

	rpNormalized := normalizeAgySubcommand("rp")
	if rpNormalized != "running-projects" {
		t.Errorf("normalizeAgySubcommand('rp') = %q; want 'running-projects'", rpNormalized)
	}
}

func TestResolveRecreatePrompt(t *testing.T) {
	def := resolveRecreatePrompt("")
	if !strings.Contains(def, "Read Memory protocol") {
		t.Errorf("expected default prompt to mention Read Memory protocol, got: %s", def)
	}
	if !strings.Contains(def, "architecture") {
		t.Errorf("expected default prompt to mention architecture, got: %s", def)
	}

	custom := resolveRecreatePrompt("Custom instructions for testing")
	if custom != "Custom instructions for testing" {
		t.Errorf("expected custom prompt to be preserved, got: %s", custom)
	}
}

func TestResolveAgyRecreateTargets_Cwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}

	// When no args provided, should target cwd / repo root
	targets, err := ResolveAgyRecreateTargets(nil, nil)
	if err != nil {
		t.Fatalf("ResolveAgyRecreateTargets(nil) failed: %v", err)
	}
	if len(targets) != 1 {
		t.Fatalf("expected 1 target, got %d", len(targets))
	}
	if targets[0].GetPath() == "" {
		t.Errorf("expected non-empty path for cwd target")
	}
	_ = cwd
}

func TestResolveAgyRecreateTargets_MultiArgs(t *testing.T) {
	tempDir1, err := os.MkdirTemp("", "agy-rec-test1-*")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(tempDir1)

	tempDir2, err := os.MkdirTemp("", "agy-rec-test2-*")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(tempDir2)

	projects := []AgyProject{
		{
			ID:   "proj-uuid-1",
			Name: "alpha-project",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI:     buildFolderURI(tempDir1),
							DefaultBranch: "main",
						},
					},
				},
			},
		},
		{
			ID:   "proj-uuid-2",
			Name: "beta-project",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI:     buildFolderURI(tempDir2),
							DefaultBranch: "main",
						},
					},
				},
			},
		},
	}

	// Test targeting by sequence
	seqTargets, err := ResolveAgyRecreateTargets([]string{"1"}, projects)
	if err != nil || len(seqTargets) != 1 || seqTargets[0].Name != "alpha-project" {
		t.Errorf("expected sequence 1 to resolve to alpha-project, got: %v, err: %v", seqTargets, err)
	}

	// Test targeting by comma-separated sequences and slug
	multiTargets, err := ResolveAgyRecreateTargets([]string{"1, 2", "beta-project"}, projects)
	if err != nil {
		t.Fatalf("multi target resolution failed: %v", err)
	}
	if len(multiTargets) != 2 {
		t.Errorf("expected 2 unique targets, got %d", len(multiTargets))
	}

	// Test targeting by direct directory path
	pathTargets, err := ResolveAgyRecreateTargets([]string{tempDir1}, projects)
	if err != nil || len(pathTargets) != 1 {
		t.Errorf("expected directory path to resolve, got: %v, err: %v", pathTargets, err)
	}
}

func TestExecuteAgyRecreate_DryRun(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agy-rec-dryrun-*")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	_ = os.WriteFile(testFile, []byte("preserve me"), 0644)

	p := AgyProject{
		ID:   "dry-run-id",
		Name: "dry-run-test",
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{
					GitFolder: &AgyGitFolder{
						FolderURI:     buildFolderURI(tempDir),
						DefaultBranch: "main",
					},
				},
			},
		},
	}

	opts := AgyRecreateOptions{
		IsDryRun: true,
	}

	err = ExecuteAgyRecreate([]AgyProject{p}, opts)
	if err != nil {
		t.Errorf("expected dry-run to succeed, got: %v", err)
	}

	// Verify test file was preserved
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Errorf("dry-run deleted test file")
	}
}

func TestExecuteAgyRecreate_RealLifecycle(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agy-rec-real-*")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configDir, err := getProjectsDirPath()
	if err == nil {
		_ = os.MkdirAll(configDir, 0755)
	}

	p := buildAdHocProject(tempDir)
	opts := AgyRecreateOptions{
		IsDryRun: false,
	}

	err = ExecuteAgyRecreate([]AgyProject{p}, opts)
	if err != nil {
		t.Fatalf("real recreate failed: %v", err)
	}

	// Verify project was registered in Antigravity
	configDir, err = getProjectsDirPath()
	if err != nil {
		return
	}
	fileURI := buildFolderURI(tempDir)
	id := workspacesync.FindExistingProjectID(configDir, fileURI)
	if id == "" {
		t.Errorf("expected project to be registered in %s", configDir)
		return
	}
	_ = deleteProjectFile(id)
}

func TestIsRestrictedSystemOrHomeDir(t *testing.T) {
	if !IsRestrictedSystemOrHomeDir("") || !IsRestrictedSystemOrHomeDir("/") {
		t.Errorf("expected empty path and / to be restricted")
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return
	}
	if !IsRestrictedSystemOrHomeDir(home) || !IsRestrictedSystemOrHomeDir(filepath.Dir(home)) {
		t.Errorf("expected home %q and parent to be restricted", home)
	}
	vol := filepath.VolumeName(home)
	if vol != "" && !IsRestrictedSystemOrHomeDir(vol+"\\") {
		t.Errorf("expected drive root %s\\ to be restricted", vol)
	}
}

func TestResolveAgyRecreateTargets_Restricted(t *testing.T) {
	home, err := os.UserHomeDir()
	if err == nil && home != "" {
		assertRecreateFails(t, home, "home directory")
	}
	vol := filepath.VolumeName(home)
	if vol != "" {
		assertRecreateFails(t, vol+"\\", "drive root")
	}
}

func assertRecreateFails(t *testing.T, target, desc string) {
	if _, recErr := ResolveAgyRecreateTargets([]string{target}, nil); recErr == nil {
		t.Errorf("expected recreation of %s %q to fail", desc, target)
	}
}

func TestResolveAgyRecreateTargets_NonGitDir(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "agy-rec-nongit-*")
	if err != nil {
		t.Fatalf("temp dir error: %v", err)
	}
	defer os.RemoveAll(tempDir)

	if _, recErr := ResolveAgyRecreateTargets([]string{tempDir}, nil); recErr == nil {
		t.Errorf("expected non-git ad-hoc folder recreation to fail")
	}
}
