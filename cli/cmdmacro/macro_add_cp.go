// Package cmdmacro — macro_add_cp.go: in-builder cp/copy helper for interactive macro creation and edit.
package cmdmacro

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
)

func runCpFileCmd(args []string) error {
	srcRaw, dstRaw, ok := extractCpPositionalArgs(args)
	if !ok {
		fmt.Printf("  %s▲ Usage: cp [-r] <src> <dst>%s\n\n", constants.ColorYellow, constants.ColorReset)
		return apperror.NewValidationError("missing src or dst for cp")
	}
	src := macro.ExpandPathAndEnv(srcRaw)
	dst := macro.ExpandPathAndEnv(dstRaw)
	_, _, _, err := executeCopyOperation(src, dst)
	return err
}

func executeInteractiveCopy(line string, args []string) error {
	srcRaw, dstRaw, ok := extractCpPositionalArgs(args)
	if !ok {
		fmt.Printf("  %s▲ Usage: cp [-r] <src> <dst>%s\n\n", constants.ColorYellow, constants.ColorReset)
		return apperror.NewValidationError("missing src or dst for cp")
	}
	src := macro.ExpandPathAndEnv(srcRaw)
	dst := macro.ExpandPathAndEnv(dstRaw)
	isDir, count, bytes, err := executeCopyOperation(src, dst)
	if err != nil {
		fmt.Printf("  %s▲ cp %s %s: %v%s\n\n", constants.ColorRed, srcRaw, dstRaw, err, constants.ColorReset)
		return err
	}
	printInteractiveCopySuccess(isDir, src, dst, line, count, bytes)
	return nil
}

func executeCopyOperation(src, dst string) (bool, int, int64, error) {
	srcStat, err := os.Stat(src)
	if err != nil {
		return false, 0, 0, apperror.WrapSimple(err, fmt.Sprintf("cannot stat %q", src))
	}
	targetDst := resolveTargetDestination(src, dst)
	if srcStat.IsDir() {
		count, totalBytes, copyErr := copyDirectoryRecursive(src, targetDst)
		return true, count, totalBytes, copyErr
	}
	copied, copyErr := copySingleFile(src, targetDst, srcStat.Mode())
	return false, 1, copied, copyErr
}

func resolveTargetDestination(src, dst string) string {
	if isDirTarget(dst) {
		return filepath.Join(dst, filepath.Base(src))
	}
	return dst
}

func isDirTarget(dst string) bool {
	if dst == "." || dst == ".." || strings.HasSuffix(dst, "/") || strings.HasSuffix(dst, "\\") {
		return true
	}
	stat, err := os.Stat(dst)
	return err == nil && stat.IsDir()
}

func extractCpPositionalArgs(args []string) (string, string, bool) {
	var positional []string
	for _, a := range args {
		trimmed := strings.Trim(a, "\"'")
		if !strings.HasPrefix(trimmed, "-") {
			positional = append(positional, trimmed)
		}
	}
	if len(positional) < 2 {
		return "", "", false
	}
	return positional[0], positional[1], true
}

func copySingleFile(src, dst string, mode os.FileMode) (int64, error) {
	ensureParentDir(dst)
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}

func copyDirectoryRecursive(src, dst string) (int, int64, error) {
	var count int
	var totalBytes int64
	absDst, _ := filepath.Abs(dst)
	err := filepath.Walk(src, makeCopyWalkFunc(src, dst, absDst, &count, &totalBytes))
	return count, totalBytes, err
}

func makeCopyWalkFunc(src, dst, absDst string, count *int, totalBytes *int64) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil || isInsideDestDir(path, absDst) {
			return err
		}
		return copyWalkItem(src, dst, path, info, count, totalBytes)
	}
}

func isInsideDestDir(path, absDst string) bool {
	absPath, _ := filepath.Abs(path)
	return absPath == absDst || (len(absPath) > len(absDst) && strings.HasPrefix(absPath, absDst+string(filepath.Separator)))
}

func copyWalkItem(src, dst, path string, info os.FileInfo, count *int, totalBytes *int64) error {
	rel, relErr := filepath.Rel(src, path)
	if relErr != nil {
		return relErr
	}
	targetPath := filepath.Join(dst, rel)
	if info.IsDir() {
		return os.MkdirAll(targetPath, 0755)
	}
	copied, copyErr := copySingleFile(path, targetPath, info.Mode())
	if copyErr == nil {
		*count++
		*totalBytes += copied
	}
	return copyErr
}

func printInteractiveCopySuccess(isDir bool, src, dst, line string, count int, totalBytes int64) {
	targetDst := resolveTargetDestination(src, dst)
	if isDir {
		fmt.Printf("  %s✓ Copied directory: %s → %s (%d files, %s)%s\n",
			constants.ColorGreen, src, targetDst, count, formatByteSize(totalBytes), constants.ColorReset)
		fmt.Printf("  %s(copied directory live. Enter command for Step, or '+add' to record '%s')%s\n\n",
			constants.ColorDim, line, constants.ColorReset)
		return
	}
	fmt.Printf("  %s✓ Copied file: %s → %s (%s)%s\n",
		constants.ColorGreen, src, targetDst, formatByteSize(totalBytes), constants.ColorReset)
	fmt.Printf("  %s(copied file live. Enter command for Step, or '+add' to record '%s')%s\n\n",
		constants.ColorDim, line, constants.ColorReset)
}
