package downloaderconfig

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	dirPermissionMode  = 0755
	filePermissionMode = 0644
)

// CreateInFlightStageDir creates an isolated temporary directory in /tmp.
func CreateInFlightStageDir() (string, error) {
	stageDir, err := os.MkdirTemp("", constants.DirInFlightStagePrefix)
	if err != nil {
		return "", apperror.WrapSimple(err, "downloaderconfig.CreateInFlightStageDir")
	}

	return stageDir, nil
}

// ResolvePersistentKeepDir returns the root keep directory path (~/.gitmap-installation).
func ResolvePersistentKeepDir() (string, error) {
	sudoUser := os.Getenv("SUDO_USER")
	if sudoUser != "" {
		return filepath.Join("/home", sudoUser, constants.DirPersistentKeep), nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", apperror.WrapSimple(err, "downloaderconfig.ResolvePersistentKeepDir")
	}

	return filepath.Join(homeDir, constants.DirPersistentKeep), nil
}

// EnsureKeepDirectoryStructure creates downloads, scripts, and logs subdirectories.
func EnsureKeepDirectoryStructure(baseDir string) error {
	subdirs := []string{
		constants.DirKeepDownloads,
		constants.DirKeepScripts,
		constants.DirKeepLogs,
	}

	for _, sub := range subdirs {
		targetPath := filepath.Join(baseDir, sub)
		if err := os.MkdirAll(targetPath, dirPermissionMode); err != nil {
			return apperror.WrapSimple(err, "downloaderconfig.EnsureKeepDirectoryStructure")
		}
	}

	return nil
}

func copyFileContents(srcFile, dstFile *os.File) error {
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.copyFileContents.copy")
	}

	if err := dstFile.Sync(); err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.copyFileContents.sync")
	}

	return nil
}

func copyAndRemoveFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.copyAndRemoveFile.openSrc")
	}

	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePermissionMode)
	if err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.copyAndRemoveFile.openDst")
	}

	defer dstFile.Close()

	if err := copyFileContents(srcFile, dstFile); err != nil {
		return err
	}

	return os.Remove(src)
}

func isCrossDeviceError(err error) bool {
	var linkErr *os.LinkError
	if errors.As(err, &linkErr) {
		return errors.Is(linkErr.Err, syscall.EXDEV)
	}

	return errors.Is(err, syscall.EXDEV)
}

// PromoteStagedFile moves a staged file from /tmp to persistent keep directory.
func PromoteStagedFile(srcTmpPath, dstKeepPath string) error {
	if _, err := os.Stat(srcTmpPath); err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.PromoteStagedFile.srcNotFound")
	}

	parentDir := filepath.Dir(dstKeepPath)
	if err := os.MkdirAll(parentDir, dirPermissionMode); err != nil {
		return apperror.WrapSimple(err, "downloaderconfig.PromoteStagedFile.mkdirDst")
	}

	err := os.Rename(srcTmpPath, dstKeepPath)
	if err == nil {
		return nil
	}

	if isCrossDeviceError(err) {
		return copyAndRemoveFile(srcTmpPath, dstKeepPath)
	}

	return apperror.WrapSimple(err, "downloaderconfig.PromoteStagedFile.rename")
}
