package cmd

// Test helpers shared by the clone-side projects.json dedup integration
// tests. Split out of clonepmsync_dedup_integration_test.go to keep
// every file under the project's <200-line cap (mem://style/code-constraints).

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// sandboxVSCodePMRoot redirects vscodepm.ProjectsJSONPath into a temp
// dir by overriding the OS-specific user-data env var, then creates
// the extension storage dir so Sync's "extension missing" soft-fail
// does not short-circuit the test.
func sandboxVSCodePMRoot(t *testing.T) string {
	t.Helper()

	root := t.TempDir()

	switch runtime.GOOS {
	case "windows":
		t.Setenv(constants.VSCodeEnvAppData, root)
	case "darwin":
		t.Setenv(constants.VSCodeEnvHome, root)
	default:
		t.Setenv(constants.VSCodeEnvXDGConfigHome, root)
	}

	extDir := vscodePMExtensionDir(root)
	if err := os.MkdirAll(extDir, 0o755); err != nil {
		t.Fatalf("mkdir ext dir: %v", err)
	}

	return filepath.Join(extDir, constants.VSCodePMProjectsFile)
}

// vscodePMExtensionDir mirrors vscodepm.ProjectsJSONPath's join logic
// for the OS-specific user-data root that sandboxVSCodePMRoot just set.
func vscodePMExtensionDir(root string) string {
	var userDataRoot string

	switch runtime.GOOS {
	case "windows":
		userDataRoot = filepath.Join(root, constants.VSCodeUserDataRootDirName)
	case "darwin":
		userDataRoot = filepath.Join(root,
			filepath.FromSlash(constants.VSCodeUserDataMacRel))
	default:
		userDataRoot = filepath.Join(root, constants.VSCodeUserDataRootDirName)
	}

	return filepath.Join(userDataRoot,
		constants.VSCodePMUserDir,
		constants.VSCodePMGlobalStorageDir,
		constants.VSCodePMExtensionDir)
}

// readProjectsJSONEntries loads projects.json off disk so the assertions
// see exactly what a real VS Code restart would see. Missing file -> nil.
func readProjectsJSONEntries(t *testing.T, path string) []map[string]any {
	t.Helper()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		t.Fatalf("read %s: %v", path, err)
	}

	var out []map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("parse %s: %v", path, err)
	}

	return out
}

// runDoubleClone invokes the real clone-side sync helper twice with two
// spellings of the same physical path, then returns the resulting
// projects.json entries.
func runDoubleClone(t *testing.T, projectsPath, spellingA, spellingB,
	repoName string) []map[string]any {
	t.Helper()

	syncSingleClonedRepoToVSCodePM(spellingA, repoName, false)
	syncSingleClonedRepoToVSCodePM(spellingB, repoName, false)

	return readProjectsJSONEntries(t, projectsPath)
}
