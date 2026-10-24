package cmdupdate

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// resolveRepoPathFromFlag validates the --repo-path flag if present.
func resolveRepoPathFromFlag() string {
	return normalizeRepoPath(getFlagValue(constants.FlagRepoPath))
}

// resolveRepoPathFromEmbedded validates the embedded repo path.
func resolveRepoPathFromEmbedded() string {
	return normalizeRepoPath(constants.RepoPath)
}

// resolveRepoPathFromDB validates the saved repo path from SQLite.
func resolveRepoPathFromDB() string {
	return normalizeRepoPath(loadRepoPathFromDB())
}

// resolveRepoPathFromCwd checks if the current working directory is inside a gitmap repo.
func resolveRepoPathFromCwd() string {
	cwd, err := os.Getwd()
	if err != nil {
		return ""
	}

	return normalizeRepoPath(cwd)
}

// resolveRepoPathFromRepoTable checks the SQLite Repo table for an indexed gitmap repo.
func resolveRepoPathFromRepoTable() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}

	defer db.Close()

	rows, err := db.SQL().Query(
		"SELECT AbsolutePath FROM Repo WHERE Slug LIKE 'gitmap%' OR RepoName LIKE 'gitmap%' ORDER BY RepoId DESC",
	)
	if err != nil {
		return ""
	}

	defer rows.Close()

	for rows.Next() {
		if root := scanRepoTableRow(rows); root != "" {
			return root
		}
	}

	return ""
}

func scanRepoTableRow(rows *sql.Rows) string {
	var path string
	if err := rows.Scan(&path); err != nil {
		return ""
	}
	return normalizeRepoPath(path)
}

// resolveRepoPathFromWorkDirs checks work_directories in the SQLite database for a gitmap repo.
func resolveRepoPathFromWorkDirs() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}

	defer db.Close()

	rows, err := db.SQL().Query("SELECT path FROM work_directories")
	if err != nil {
		return ""
	}

	defer rows.Close()

	for rows.Next() {
		if root := scanWorkDirRow(rows); root != "" {
			return root
		}
	}

	return ""
}

func scanWorkDirRow(rows *sql.Rows) string {
	var baseDir string
	if err := rows.Scan(&baseDir); err != nil {
		return ""
	}
	return normalizeRepoPath(filepath.Join(baseDir, "gitmap"))
}

// normalizeRepoPath resolves a candidate path to a valid gitmap source root.
func normalizeRepoPath(path string) string {
	if len(path) == 0 {
		return ""
	}

	cleaned := strings.TrimSpace(strings.Trim(path, "\"'"))
	if len(cleaned) == 0 {
		return ""
	}

	cleaned = expandTilde(cleaned)

	absPath, err := filepath.Abs(cleaned)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to resolve absolute path for %s: %v\n", cleaned, err)

		return ""
	}

	return findRepoRoot(absPath)
}

// expandTilde replaces a leading ~ with the user's home directory.
func expandTilde(path string) string {
	if !strings.HasPrefix(path, "~") {
		return path
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}

	if path == "~" {
		return home
	}

	if strings.HasPrefix(path, "~/") || strings.HasPrefix(path, "~\\") {
		return filepath.Join(home, path[2:])
	}

	return path
}

// findRepoRoot walks upward until it finds a valid gitmap source root.
func findRepoRoot(path string) string {
	current := path
	for {
		if isGitmapSourceRepo(current) {
			return current
		}

		parent := filepath.Dir(current)
		if parent == current {
			return ""
		}

		current = parent
	}
}

// isGitmapSourceRepo checks for the update script and source markers.
// Uses fsutil predicates so the strict-file vs strict-dir distinction is
// enforced by a single shared contract (see gitmap/fsutil/exists.go).
func isGitmapSourceRepo(path string) bool {
	hasScript := fsutil.FileExists(filepath.Join(path, "run.ps1")) || fsutil.FileExists(filepath.Join(path, "run.sh"))
	if !fsutil.DirExists(path) || !hasScript {
		return false
	}

	markers := []string{
		filepath.Join(path, "cli", "constants", "constants.go"),
		filepath.Join(path, "cli", constants.GoModFile),
		filepath.Join(path, "gitmap", "constants", "constants.go"),
		filepath.Join(path, "constants", "constants.go"),
	}
	for _, marker := range markers {
		if fsutil.FileExists(marker) {
			return true
		}
	}

	return false
}

// isRepoPathPromptAllowed checks whether stdin is interactive.
func isRepoPathPromptAllowed() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return info.Mode()&os.ModeCharDevice != 0
}

// promptRepoPath asks the user to enter the source repo path interactively.
// If the path does not exist, it clones the gitmap repo into that location.
func promptRepoPath() string {
	if !isRepoPathPromptAllowed() {
		return ""
	}

	saved := loadRepoPathFromDB()
	if len(saved) > 0 {
		fmt.Fprintf(os.Stderr, "\n  ⚠ The saved source repository path (%s) no longer exists on disk.\n", saved)
	} else {
		fmt.Fprint(os.Stderr, "\n  ⚠ No gitmap source repository found on disk.\n")
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Fprint(os.Stderr, constants.MsgUpdatePathPrompt)

		input, err := reader.ReadString('\n')
		if err != nil && len(input) == 0 {
			return ""
		}

		path := strings.TrimSpace(strings.Trim(input, "\"'"))
		if len(path) == 0 {
			return ""
		}

		// Try existing path first.
		if root := normalizeRepoPath(path); len(root) > 0 {
			return root
		}

		// Path doesn't exist — clone into it.
		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrUpdatePathInvalid, path)
			continue
		}

		if !cloneRepoInto(absPath) {
			fmt.Fprintf(os.Stderr, constants.ErrUpdatePathInvalid, path)
			continue
		}
		if root := normalizeRepoPath(absPath); len(root) > 0 {
			return root
		}

		fmt.Fprintf(os.Stderr, constants.ErrUpdatePathInvalid, path)
	}
}

// cloneRepoInto clones the gitmap source repository into the given directory.
func cloneRepoInto(targetPath string) bool {
	fmt.Fprintf(os.Stderr, constants.MsgUpdateCloning, targetPath)

	cmd := exec.Command(constants.GitBin, constants.GitClone, constants.SourceRepoCloneURL, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateCloneFailed, err)

		return false
	}

	fmt.Fprint(os.Stderr, constants.MsgUpdateCloneOK)

	return true
}

// saveRepoPathToDB persists the source repo path in the Settings table.
func saveRepoPathToDB(path string) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not save repo path to database: %v\n", err)
		return
	}
	defer db.Close()

	if err := db.SetSetting(constants.SettingSourceRepoPath, path); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save repo path to DB: %v\n", err)
	}
}

// loadRepoPathFromDB reads the source repo path from the Settings table.
func loadRepoPathFromDB() string {
	db, err := store.OpenDefault()
	if err != nil {
		return ""
	}
	defer db.Close()

	return db.GetSetting(constants.SettingSourceRepoPath)
}
