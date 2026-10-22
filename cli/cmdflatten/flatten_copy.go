package cmdflatten

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveSuccessorPath(parentDir, repoName string) (string, string) {
	suffix := 1
	for {
		name := fmt.Sprintf("%s-flat-v%d", repoName, suffix)
		path := filepath.Join(parentDir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return path, name
		}
		suffix++
	}
}

func countCommits(repoPath string) int {
	cmd := exec.Command("git", "-C", repoPath, "rev-list", "--count", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}

	count, parseErr := strconv.Atoi(strings.TrimSpace(string(out)))
	if parseErr != nil {
		return 0
	}

	return count
}

func computeWorkingTreeSize(root string) (int64, error) {
	var totalBytes int64

	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		if !info.IsDir() {
			totalBytes += info.Size()
			return nil
		}
		if isExcludedDir(info.Name()) && path != root {
			return filepath.SkipDir
		}
		return nil
	})

	return totalBytes, err
}

func isExcludedDir(name string) bool {
	switch strings.ToLower(name) {
	case ".git", "node_modules", "dist", "temp", ".pnpm-store":
		return true
	default:
		return false
	}
}

func copyWorkingTree(srcRoot, dstRoot string) error {
	return filepath.Walk(srcRoot, func(srcPath string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if !info.IsDir() {
			return copySingleFileEntry(srcRoot, dstRoot, srcPath, info)
		}
		if isExcludedDir(info.Name()) && srcPath != srcRoot {
			return filepath.SkipDir
		}
		return nil
	})
}

func copySingleFileEntry(srcRoot, dstRoot, srcPath string, info os.FileInfo) error {
	relPath, err := filepath.Rel(srcRoot, srcPath)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(dstRoot, relPath)
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return err
	}

	return copyFile(srcPath, dstPath, info.Mode())
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func initGitRepo(repoPath string) error {
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = repoPath
	if cmdInit.Run() == nil {
		return nil
	}

	fallbackCmd := exec.Command("git", "init")
	fallbackCmd.Dir = repoPath
	return fallbackCmd.Run()
}

func commitFlatSnapshot(repoPath string) error {
	cmdAdd := exec.Command("git", "add", "-A")
	cmdAdd.Dir = repoPath
	if err := cmdAdd.Run(); err != nil {
		return apperror.WrapSimple(err, "git add failed")
	}

	cmdCommit := exec.Command("git", "-c", "user.name=GitMap", "-c", "user.email=gitmap@local",
		"commit", "-m", "flat: initial squashed repository snapshot")
	cmdCommit.Dir = repoPath
	if err := cmdCommit.Run(); err != nil {
		return apperror.WrapSimple(err, "git commit failed")
	}

	return nil
}
