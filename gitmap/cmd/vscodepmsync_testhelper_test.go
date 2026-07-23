// Package cmd — vscodepmsync_testhelper_test.go: test fixtures and
// helpers shared by vscodepmsync_test.go. Kept in a separate file so
// the primary test file stays under the 200-line code-style cap.
package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/vscodepm"
)

// vscodePMUserDataRel returns the OS-specific subpath under the home/temp
// dir that ProjectsJSONPath() will look at, so test fixtures place
// projects.json where the resolver actually reads it on every platform.
func vscodePMUserDataRel() string {
	switch runtime.GOOS {
	case "darwin":
		return filepath.FromSlash(constants.VSCodeUserDataMacRel)
	case "windows":
		return constants.VSCodeUserDataRootDirName
	default:
		return filepath.Join(".config", constants.VSCodeUserDataRootDirName)
	}
}

// setupVSCodePMSyncFixture creates a temp HOME, a real repo dir, and
// a single-entry projects.json file pointing at the repo. Returns the
// repo path and a restore function that resets HOME / XDG_CONFIG_HOME.
func setupVSCodePMSyncFixture(t *testing.T) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	repoDir := filepath.Join(tmp, "demo-repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	jsonPath := vscodepmSyncFixturePath(t, tmp)
	writeVSCodePMSyncSeed(t, jsonPath, repoDir)

	restore := swapHomeEnv(tmp)

	return repoDir, restore
}

// vscodepmSyncFixturePath returns the projects.json path inside the
// faux user-data root and ensures the parent directory exists.
func vscodepmSyncFixturePath(t *testing.T, tmp string) string {
	t.Helper()
	dir := filepath.Join(tmp, vscodePMUserDataRel(),
		constants.VSCodePMUserDir, constants.VSCodePMGlobalStorageDir,
		constants.VSCodePMExtensionDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir extdir: %v", err)
	}

	return filepath.Join(dir, constants.VSCodePMProjectsFile)
}

// writeVSCodePMSyncSeed writes a single-entry projects.json with a
// pre-existing "user" tag so we can prove UNION preservation.
func writeVSCodePMSyncSeed(t *testing.T, jsonPath, repoDir string) {
	t.Helper()
	seed := []vscodepm.Entry{{
		Name: "demo-repo", RootPath: repoDir,
		Paths: []string{}, Tags: []string{"user"}, Enabled: true,
	}}
	data, err := json.MarshalIndent(seed, "", "\t")
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		t.Fatalf("write seed: %v", err)
	}
}

// loadVSCodePMSyncProjectsJSON reads the on-disk projects.json after
// the runner finished and unmarshals it for assertion.
func loadVSCodePMSyncProjectsJSON(t *testing.T) []vscodepm.Entry {
	t.Helper()
	got, err := vscodepm.ListEntries()
	if err != nil {
		t.Fatalf("list entries: %v", err)
	}

	return got
}

// containsTag reports whether tag appears in tags.
func containsTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}

	return false
}

// setupVSCodePMSyncFixtureWithTags is a generalisation of
// setupVSCodePMSyncFixture that lets the caller pre-load the entry
// with an arbitrary tag set. Used by the dedupe tests to prove the
// merge layer collapses duplicates regardless of source.
func setupVSCodePMSyncFixtureWithTags(t *testing.T, seedTags []string) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	repoDir := filepath.Join(tmp, "demo-repo")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo: %v", err)
	}

	jsonPath := vscodepmSyncFixturePath(t, tmp)
	seed := []vscodepm.Entry{{
		Name: "demo-repo", RootPath: repoDir,
		Paths: []string{}, Tags: seedTags, Enabled: true,
	}}
	data, err := json.MarshalIndent(seed, "", "\t")
	if err != nil {
		t.Fatalf("marshal seed: %v", err)
	}
	if err := os.WriteFile(jsonPath, data, 0o644); err != nil {
		t.Fatalf("write seed: %v", err)
	}

	return repoDir, swapHomeEnv(tmp)
}

// setupVSCodePMSyncEmptyHome wires HOME / XDG_CONFIG_HOME to a fresh
// temp dir but writes NO projects.json. The extension dir does not
// exist either — exercises the "missing file" code path end-to-end.
func setupVSCodePMSyncEmptyHome(t *testing.T) func() {
	t.Helper()
	tmp := t.TempDir()
	// Pre-create the extension dir so ProjectsJSONPath() returns a
	// valid path, but leave projects.json itself absent.
	_ = vscodepmSyncFixturePath(t, tmp)

	return swapHomeEnv(tmp)
}

// setupVSCodePMSyncMalformedFile writes a deliberately broken
// projects.json (truncated JSON object) and returns its path so the
// caller can snapshot the bytes before invoking the runner. Asserts
// the soft-fail "leave the file untouched" contract.
func setupVSCodePMSyncMalformedFile(t *testing.T) (string, func()) {
	t.Helper()
	tmp := t.TempDir()
	jsonPath := vscodepmSyncFixturePath(t, tmp)
	const malformed = `[{"name": "demo-repo", "rootPath":` // truncated
	if err := os.WriteFile(jsonPath, []byte(malformed), 0o644); err != nil {
		t.Fatalf("write malformed: %v", err)
	}

	return jsonPath, swapHomeEnv(tmp)
}

// swapHomeEnv points HOME / XDG_CONFIG_HOME / APPDATA at tmp and returns
// a restore func that puts the original values back. Centralized so every
// fixture uses the same swap shape across darwin / linux / windows.
func swapHomeEnv(tmp string) func() {
	prevHome := os.Getenv(constants.VSCodeEnvHome)
	prevXDG := os.Getenv(constants.VSCodeEnvXDGConfigHome)
	prevAppData := os.Getenv(constants.VSCodeEnvAppData)
	os.Setenv(constants.VSCodeEnvHome, tmp)
	os.Setenv(constants.VSCodeEnvXDGConfigHome, filepath.Join(tmp, ".config"))
	os.Setenv(constants.VSCodeEnvAppData, tmp)

	return func() {
		os.Setenv(constants.VSCodeEnvHome, prevHome)
		os.Setenv(constants.VSCodeEnvXDGConfigHome, prevXDG)
		os.Setenv(constants.VSCodeEnvAppData, prevAppData)
	}
}
