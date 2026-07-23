package release

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

const gitmapModuleMarker = "module github.com/alimtvnetwork/gitmap-v27/gitmap"

// resolveSourceRepo finds the gitmap source repo root.
// It tries executable path, DB fallback, current directory, then user prompt.
func resolveSourceRepo() (string, error) {
	for _, root := range []string{
		resolveFromExecutable(),
		resolveFromDB(),
		resolveFromCWD(),
	} {
		if root == "" {
			continue
		}

		saveSourceRepoDB(root)

		return root, nil
	}

	return promptForSourceRepo()
}

func resolveFromExecutable() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return ""
	}

	return normalizeSourceRepoRoot(filepath.Dir(exe))
}

func resolveFromDB() string {
	return normalizeSourceRepoRoot(loadSourceRepoDB())
}

func resolveFromCWD() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return normalizeSourceRepoRoot(dir)
}

func promptForSourceRepo() (string, error) {
	if !canPromptForPath() {
		return "", fmt.Errorf("%s", constants.ErrSelfReleaseNoRepo)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stdout, constants.MsgSelfReleasePromptPath)
		input, err := reader.ReadString('\n')
		if err != nil && len(input) == 0 {
			return "", fmt.Errorf("%s", constants.ErrSelfReleaseNoRepo)
		}

		path := strings.TrimSpace(strings.Trim(input, "\"'"))
		if path == "" {
			return "", fmt.Errorf("%s", constants.ErrSelfReleaseNoRepo)
		}

		root := normalizeSourceRepoRoot(path)
		if root != "" {
			saveSourceRepoDB(root)
			fmt.Printf(constants.MsgSelfReleaseSavedPath, root)

			return root, nil
		}

		fmt.Fprintf(os.Stderr, constants.MsgSelfReleaseInvalidPath, path)
	}
}

func canPromptForPath() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

func normalizeSourceRepoRoot(path string) string {
	if path == "" {
		return ""
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return ""
	}

	root := findGitRoot(absPath)
	if root == "" || !isGitmapSourceRepo(root) {
		return ""
	}

	return root
}

func isGitmapSourceRepo(root string) bool {
	if fileExists(filepath.Join(root, "gitmap", "constants", "constants.go")) {
		return true
	}

	if !fileExists(filepath.Join(root, "go.mod")) || !fileExists(filepath.Join(root, "constants", "constants.go")) {
		return false
	}

	goMod, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return false
	}

	return strings.Contains(string(goMod), gitmapModuleMarker)
}

// saveSourceRepoDB persists the source repo path in the Settings table.
func saveSourceRepoDB(path string) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not open DB to save source repo path: %v\n", err)

		return
	}
	defer db.Close()

	if err := db.SetSetting(constants.SettingSourceRepoPath, path); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save source repo path to DB: %v\n", err)
	}
}

// loadSourceRepoDB reads the source repo path from the Settings table.
func loadSourceRepoDB() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	return db.GetSetting(constants.SettingSourceRepoPath)
}

// findGitRoot walks up from dir looking for a .git directory.
func findGitRoot(dir string) string {
	for {
		gitEntry := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitEntry); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}

		dir = parent
	}
}
