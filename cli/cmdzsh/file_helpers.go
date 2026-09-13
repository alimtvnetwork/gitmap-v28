// Package cmdzsh provides file and directory manipulation helpers.
package cmdzsh

import (
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// ResolveTargetHome resolves the user target home directory.
func ResolveTargetHome(customHome string) string {
	if customHome != "" {
		return customHome
	}

	homeDir, err := os.UserHomeDir()
	if err == nil && homeDir != "" {
		return homeDir
	}

	return os.Getenv("HOME")
}

// HasDirectory reports whether the path exists and is a directory.
func HasDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

// HasFile reports whether the path exists and is a regular file.
func HasFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// EnsureDirectory creates directory path if missing with specified permissions.
func EnsureDirectory(path string, perm os.FileMode) *apperror.AppError {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.EnsureDirectory")
	}

	return nil
}

// CopyFile copies file content from source to destination.
func CopyFile(src, dst string) *apperror.AppError {
	srcFile, err := os.Open(src)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.CopyFile.open")
	}

	defer srcFile.Close()

	return writeCopiedFile(srcFile, dst)
}

func writeCopiedFile(r io.Reader, dst string) *apperror.AppError {
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.writeCopiedFile.create")
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, r)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.writeCopiedFile.copy")
	}

	return nil
}
