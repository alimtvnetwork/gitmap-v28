package cmdpurge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func checkWorkTreeAndBranch() error {
	out, err := runPurgeCmd("git", "status", "--porcelain")
	if err != nil || strings.TrimSpace(out) != "" {
		return apperror.NewSimple("EXECUTION", "Working tree not clean.")
	}

	br, err := runPurgeCmd("git", "branch", "--show-current")
	if err != nil || strings.TrimSpace(br) == "" {
		return apperror.NewSimple("EXECUTION", "Could not determine branch.")
	}

	return nil
}

func normalizePurgePattern(pattern string) string {
	return strings.ReplaceAll(pattern, "\\", "/")
}

func validatePurgeState(pattern string) ([]string, error) {
	if err := checkWorkTreeAndBranch(); err != nil {
		return nil, err
	}

	out, err := runPurgeCmd("git", "ls-files", normalizePurgePattern(pattern))

	return strings.Fields(out), err
}

func backupFilesToTemp(tempDir string, files []string) ([]string, error) {
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, apperror.Wrap(err, "temp dir creation failed", nil)
	}

	var backed []string
	for _, f := range files {
		if err := copyPurgeFile(f, filepath.Join(tempDir, f)); err != nil {
			return nil, err
		} else if err := sendToRecycleBin(f); err != nil {
			return nil, err
		}

		backed = append(backed, f)
	}

	return backed, nil
}

func createPurgeBackup(br, tmp string, files []string) ([]string, error) {
	if _, err := runPurgeCmd("git", "branch", br); err != nil {
		return nil, err
	}

	return backupFilesToTemp(tmp, files)
}

func executeFilterRepo(pattern string) error {
	rem, _ := runPurgeCmd("git", "remote", "get-url", "origin")
	if _, err := runPurgeCmd("git", "filter-repo", "--path-glob", normalizePurgePattern(pattern), "--invert-paths", "--force"); err != nil {
		return apperror.Wrap(err, "filter-repo failed", nil)
	}

	if r := strings.TrimSpace(rem); r != "" {
		_, err := runPurgeCmd("git", "remote", "add", "origin", r)

		return err
	}

	return nil
}

func appendGitignore(pattern string) error {
	norm := normalizePurgePattern(pattern)
	f, err := os.OpenFile(".gitignore", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	if _, err := f.WriteString("\n" + norm + "\n"); err != nil {
		return err
	}

	if _, err := runPurgeCmd("git", "add", ".gitignore"); err != nil {
		return err
	}

	_, err = runPurgeCmd("git", "commit", "-m", "chore: add "+norm+" to .gitignore")

	return err
}

func applyPurgeChanges(pattern string) error {
	if err := executeFilterRepo(pattern); err != nil {
		return err
	}

	return appendGitignore(pattern)
}

func recordPurgeLog(db *store.DB, repo, pat, br, tmp string, ts int64, files []string) error {
	raw, err := json.Marshal(files)
	if err != nil {
		return apperror.Wrap(err, "marshal failed", nil)
	}

	return db.InsertPurgeHistoryLog(&store.PurgeHistoryLog{
		RepoPath: repo, Pattern: normalizePurgePattern(pat),
		BackupBranch: br, TempDir: tmp, Files: string(raw), Timestamp: ts,
	})
}

func doPurge(db *store.DB, repoPath, pattern string, isAutoConfirm bool) error {
	files, err := validatePurgeState(pattern)
	if err != nil {
		return err
	}

	ts := time.Now().Unix()
	br, tmp := fmt.Sprintf("backup-purge-%d", ts), filepath.Join(os.TempDir(), fmt.Sprintf("gitmap_purge_%d", ts))
	backed, err := createPurgeBackup(br, tmp, files)
	if err != nil {
		return err
	}

	if err := applyPurgeChanges(pattern); err != nil {
		return err
	}

	return recordPurgeLog(db, repoPath, pattern, br, tmp, ts, backed)
}
