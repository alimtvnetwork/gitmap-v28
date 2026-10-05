package cmdpull

import (
	"os"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// HealRepoPathSeparators normalizes path separators for the host OS.
// On Unix/Linux, it converts Windows backslashes to forward slashes.
func HealRepoPathSeparators(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}

	if !strings.Contains(path, "\\") {
		return path
	}

	normalized := strings.ReplaceAll(path, "\\", "/")
	if len(normalized) >= 3 && normalized[1] == ':' && normalized[2] == '/' {
		return resolveWindowsDriveOnUnix(normalized)
	}

	return normalized
}

func resolveWindowsDriveOnUnix(normalized string) string {
	drive := strings.ToLower(string(normalized[0]))
	wslPath := "/mnt/" + drive + normalized[2:]
	if _, err := os.Stat(wslPath); err == nil {
		return wslPath
	}

	tailPath := normalized[2:]
	if _, err := os.Stat(tailPath); err == nil {
		return tailPath
	}

	return tailPath
}

// AutoHealRepoPathInDB updates the repository record in gitmap.db with the healed path.
func AutoHealRepoPathInDB(oldPath, newPath string) error {
	if oldPath == newPath || newPath == "" {
		return nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		return err
	}

	defer db.Close()

	if db.Conn() != nil {
		_, err = db.Conn().Exec("UPDATE Repo SET AbsolutePath = ? WHERE AbsolutePath = ?", newPath, oldPath)
	}

	return err
}

// TryHealMissingRecordPath attempts to repair a missing record's path if it has backslashes on Unix.
// If the healed path exists on disk, it updates rec.AbsolutePath and persists the healed path in DB.
func TryHealMissingRecordPath(rec *model.ScanRecord) bool {
	if rec == nil {
		return false
	}

	healed := HealRepoPathSeparators(rec.AbsolutePath)
	if healed == rec.AbsolutePath {
		return false
	}

	isFound := !cloner.IsMissingRepo(healed)
	if !isFound {
		return false
	}

	oldPath := rec.AbsolutePath
	rec.AbsolutePath = healed
	_ = AutoHealRepoPathInDB(oldPath, healed)

	return true
}
