package cmdupdate

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func extractBinaryFromArchive(archivePath, extractDir string) (string, error) {
	lower := strings.ToLower(archivePath)
	if isStandaloneBinary(lower) {
		return extractDirectExecutable(archivePath, extractDir)
	}

	return extractCompressedArchive(lower, archivePath, extractDir)
}

func extractCompressedArchive(lower, archivePath, extractDir string) (string, error) {
	if strings.HasSuffix(lower, ".zip") {
		return extractZipBinary(archivePath, extractDir)
	}

	if strings.HasSuffix(lower, ".tar.gz") || strings.HasSuffix(lower, ".tgz") {
		return extractTarGzBinary(archivePath, extractDir)
	}

	return "", apperror.NewWithDetails("cmd.update.unsupportedArchive", "E1146", "unsupported archive format", "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func extractDirectExecutable(archivePath, extractDir string) (string, error) {
	target := filepath.Join(extractDir, resolveTargetBinaryName(archivePath))
	if errCopy := copyFile(archivePath, target); errCopy != nil {
		return "", apperror.WrapWithDetails(errCopy, "cmd.update.copyDirect", "E1152", "failed direct copy", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	_ = os.Chmod(target, 0o755)

	return target, nil
}

func extractZipBinary(archivePath, extractDir string) (string, error) {
	zr, errZip := zip.OpenReader(archivePath)
	if errZip != nil {
		return "", apperror.WrapWithDetails(errZip, "cmd.update.zipOpen", "E1147", "failed to open zip", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}
	defer zr.Close()

	for _, f := range zr.File {
		if path, hasFound, errExt := extractZipFileEntry(f, extractDir); hasFound || errExt != nil {
			return path, errExt
		}
	}

	return "", apperror.NewSimple("gitmap binary not found in zip", "E1148")
}

func extractZipFileEntry(f *zip.File, extractDir string) (string, bool, error) {
	name := filepath.Base(filepath.Clean(f.Name))
	if !isZipBinaryTarget(name) {
		return "", false, nil
	}

	path, err := writeZipEntryFile(f, filepath.Join(extractDir, name+".new"))

	return path, true, err
}

func isZipBinaryTarget(name string) bool {
	return strings.EqualFold(name, "gitmap.exe") || strings.EqualFold(name, "gitmap")
}

func writeZipEntryFile(f *zip.File, target string) (string, error) {
	rc, errOpen := f.Open()
	if errOpen != nil {
		return "", errOpen
	}
	defer rc.Close()

	return writeExtractedBinary(rc, target)
}

func extractTarGzBinary(archivePath, extractDir string) (string, error) {
	f, errOpen := os.Open(archivePath)
	if errOpen != nil {
		return "", errOpen
	}
	defer f.Close()

	gzr, errGz := gzip.NewReader(f)
	if errGz != nil {
		return "", errGz
	}
	defer gzr.Close()

	return iterateTarEntries(tar.NewReader(gzr), extractDir)
}

func iterateTarEntries(tr *tar.Reader, extractDir string) (string, error) {
	for {
		h, err := tr.Next()
		if err != nil {
			break
		}

		name := filepath.Base(filepath.Clean(h.Name))
		if isReleaseBinaryName(name) {
			return writeExtractedBinary(tr, filepath.Join(extractDir, name+".new"))
		}
	}

	return "", apperror.NewSimple("gitmap binary not found in tar", "E1148")
}

func isReleaseBinaryName(name string) bool {
	return name == "gitmap" || name == "gitmap.exe" || strings.HasPrefix(name, "gitmap-")
}

func writeExtractedBinary(r io.Reader, target string) (string, error) {
	out, errCreate := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if errCreate != nil {
		return "", errCreate
	}
	defer out.Close()

	_, errCopy := io.Copy(out, r)

	return target, errCopy
}

func applyBinaryHandoff(newBinaryPath string) error {
	activeExe, err := resolveActiveExecutable()
	if err != nil {
		return err
	}

	if runtime.GOOS == "windows" {
		return applyWindowsHandoff(activeExe, newBinaryPath)
	}

	return applyUnixHandoff(activeExe, newBinaryPath)
}

func resolveActiveExecutable() (string, error) {
	activeExe, err := os.Executable()
	if err != nil {
		return "", apperror.WrapWithDetails(err, "cmd.update.activeExe", "E1149", "failed active exe", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	if resolved, evalErr := filepath.EvalSymlinks(activeExe); evalErr == nil {
		return resolved, nil
	}

	return activeExe, nil
}

func applyUnixHandoff(activeExe, newBinaryPath string) error {
	if errMove := os.Rename(newBinaryPath, activeExe); errMove != nil && copyFile(newBinaryPath, activeExe) != nil {
		return apperror.NewSimple("failed to place new executable", "E1151")
	}

	_ = os.Chmod(activeExe, 0o755)

	return nil
}

func applyWindowsHandoff(activeExe, newBinaryPath string) error {
	backupPath := activeExe + ".old"
	if errBackup := backupActiveBinary(activeExe, backupPath); errBackup != nil {
		return errBackup
	}

	if errPlace := placeNewBinary(newBinaryPath, activeExe, backupPath); errPlace != nil {
		return errPlace
	}

	launchCleanupProcess(activeExe)

	return nil
}

func backupActiveBinary(activeExe, backupPath string) error {
	_ = os.Remove(backupPath)
	if errRename := os.Rename(activeExe, backupPath); errRename != nil {
		return apperror.WrapWithDetails(errRename, "cmd.update.renameBackup", "E1150", "failed backup rename", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, map[string]any{"active": activeExe, "backup": backupPath})
	}

	return nil
}

func placeNewBinary(newBinaryPath, activeExe, backupPath string) error {
	if os.Rename(newBinaryPath, activeExe) != nil && copyFile(newBinaryPath, activeExe) != nil {
		_ = os.Rename(backupPath, activeExe)

		return apperror.NewSimple("failed placing new exe, restored backup", "E1151")
	}

	return nil
}

func launchCleanupProcess(activeExe string) {
	cmd := exec.Command(activeExe, constants.CmdUpdateCleanup)
	cmd.Env = append(os.Environ(), constants.EnvUpdateCleanupDelayMS+"=1500")
	setHiddenProcessAttr(cmd)
	_ = cmd.Start()
}
