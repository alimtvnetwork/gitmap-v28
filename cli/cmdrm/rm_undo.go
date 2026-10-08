package cmdrm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ExecuteUndo performs atomic restoration of staged files from a manifest.
func ExecuteUndo(opts RmUndoOptions) error {
	sessionDir, manifest, err := findUndoSession(opts.TaskId)
	if err != nil {
		return err
	}

	repoRoot := resolveUndoRepoRoot(opts.RepoRoot, manifest.RepoRoot)
	if err := restoreManifestFiles(sessionDir, repoRoot, manifest.Files); err != nil {
		return err
	}

	_ = os.RemoveAll(sessionDir)
	cleanParentTaskDir(sessionDir)
	return nil
}

func findUndoSession(targetTaskID string) (string, *RmManifest, error) {
	manifests, err := ListRmBackups()
	if err != nil {
		return "", nil, err
	}

	if targetTaskID != "" {
		return findSessionByTaskID(manifests, targetTaskID)
	}

	if len(manifests) == 0 {
		return "", nil, apperror.NewValidationError("no recoverable file removal sessions found in temporary storage")
	}

	latest := manifests[0]
	sessionDir, err := locateManifestDirectory(latest.TaskId)
	if err != nil {
		return "", nil, err
	}

	return sessionDir, &latest, nil
}

func findSessionByTaskID(manifests []RmManifest, taskID string) (string, *RmManifest, error) {
	for _, m := range manifests {
		if !strings.EqualFold(m.TaskId, taskID) {
			continue
		}
		dir, err := locateManifestDirectory(m.TaskId)
		if err != nil {
			break
		}
		return dir, &m, nil
	}

	errMsg := fmt.Sprintf("Backup directory for task '%s' not found in OS temporary storage. The files were permanently purged by OS temp cleaner and cannot be reverted.", taskID)
	return "", nil, apperror.NewValidationError(errMsg)
}

func locateManifestDirectory(taskID string) (string, error) {
	baseDirs := []string{GetBackupBaseDir(), GetLegacyBackupBaseDir()}
	for _, base := range baseDirs {
		foundDir := searchTaskDirInBase(base, taskID)
		if foundDir != "" {
			return foundDir, nil
		}
	}

	errMsg := fmt.Sprintf("Backup directory for task '%s' not found in OS temporary storage. The files were permanently purged by OS temp cleaner and cannot be reverted.", taskID)
	return "", apperror.NewValidationError(errMsg)
}

func searchTaskDirInBase(base, taskID string) string {
	taskRoot := filepath.Join(base, taskID)
	info, err := os.Stat(taskRoot)
	if err != nil || !info.IsDir() {
		return ""
	}

	var candidate string
	_ = filepath.Walk(taskRoot, func(current string, info os.FileInfo, err error) error {
		if err != nil || info == nil || info.IsDir() {
			return nil
		}
		if info.Name() == "manifest.json" {
			candidate = filepath.Dir(current)
			return filepath.SkipDir
		}
		return nil
	})
	return candidate
}

func resolveUndoRepoRoot(optRoot, manifestRoot string) string {
	if optRoot != "" {
		return optRoot
	}
	if isValidDirectory(manifestRoot) {
		return manifestRoot
	}
	return "."
}

func isValidDirectory(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func restoreManifestFiles(sessionDir, repoRoot string, files []RmFileRecord) error {
	for _, record := range files {
		if err := restoreSingleFile(sessionDir, repoRoot, record); err != nil {
			return err
		}
	}
	return nil
}

func restoreSingleFile(sessionDir, repoRoot string, record RmFileRecord) error {
	srcPath := resolveStagedSource(sessionDir, record)
	if _, err := os.Stat(srcPath); err != nil {
		return apperror.WrapSimple(err, "rm: staged backup file missing")
	}

	destPath := filepath.Join(repoRoot, filepath.FromSlash(record.RelativePath))
	if err := CopyFilePreserve(srcPath, destPath); err != nil {
		return err
	}

	if record.FileMode != 0 {
		_ = os.Chmod(destPath, record.FileMode)
	}

	return verifyRestoredChecksum(destPath, record.Sha256)
}

func resolveStagedSource(sessionDir string, record RmFileRecord) string {
	if record.StagedPath != "" && pathExists(record.StagedPath) {
		return record.StagedPath
	}
	return filepath.Join(sessionDir, "files", filepath.FromSlash(record.RelativePath))
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func verifyRestoredChecksum(destPath, expectedSha string) error {
	actualSha, err := CalculateSHA256(destPath)
	if err != nil {
		return err
	}

	if !strings.EqualFold(actualSha, expectedSha) {
		return apperror.NewValidationError("checksum mismatch after file restoration")
	}

	return nil
}

func cleanParentTaskDir(sessionDir string) {
	parent := filepath.Dir(sessionDir)
	entries, err := os.ReadDir(parent)
	if err == nil && len(entries) == 0 {
		_ = os.Remove(parent)
	}
}
