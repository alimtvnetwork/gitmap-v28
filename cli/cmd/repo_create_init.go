package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func initLocalRepo(p createRepoParams) error {
	absDir, absErr := filepath.Abs(p.LocalDir)
	if absErr != nil {
		return apperror.WrapSimple(absErr, "resolve absolute dir:")
	}

	isExisting := checkDirExists(absDir)
	errMk := maybeMkdirAll(isExisting, absDir)
	if errMk != nil {
		return errMk
	}

	gitDir := filepath.Join(absDir, ".git")
	errInit := maybeGitInit(gitDir, absDir)
	if errInit != nil {
		return errInit
	}

	errInitFiles := maybeWriteInitialFiles(isExisting, absDir, p)
	if errInitFiles != nil {
		return errInitFiles
	}

	errCommon := maybeApplyCommonFiles(p.IsCommon, absDir)
	if errCommon != nil {
		return errCommon
	}

	errCG := maybeApplyCGFiles(p.IsCG, absDir)
	if errCG != nil {
		return errCG
	}

	return nil
}

func maybeMkdirAll(isExisting bool, absDir string) error {
	if isExisting {
		return nil
	}
	err := os.MkdirAll(absDir, 0755)
	if err != nil {
		return apperror.WrapSimple(err, "create directory:")
	}
	return nil
}

func maybeGitInit(gitDir, absDir string) error {
	_, err := os.Stat(gitDir)
	if !os.IsNotExist(err) {
		return nil
	}
	return runGitInit(absDir)
}

func maybeWriteInitialFiles(isExisting bool, absDir string, p createRepoParams) error {
	if isExisting {
		return stageAndCommitExistingFiles(absDir, p)
	}
	writeInitialFiles(absDir, p)

	return commitInitialFiles(absDir)
}

func stageAndCommitExistingFiles(absDir string, p createRepoParams) error {
	readmePath := filepath.Join(absDir, "README.md")
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		readmeContent := fmt.Sprintf("# %s\n\n%s\n", p.Name, p.Description)
		_ = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	}

	gitignorePath := filepath.Join(absDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		gitignoreContent := ".DS_Store\nThumbs.db\nnode_modules/\nbin/\n*.log\n"
		_ = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	}

	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdStatus := exec.Command("git", "status", "--porcelain")
	cmdStatus.Dir = absDir
	out, err := cmdStatus.Output()
	if err == nil && len(strings.TrimSpace(string(out))) > 0 {
		cmdCommit := exec.Command("git", "commit", "-m", "feat: initial commit")
		cmdCommit.Dir = absDir
		_ = cmdCommit.Run()
	}

	return nil
}

func checkDirExists(absDir string) bool {
	_, statErr := os.Stat(absDir)
	return statErr == nil
}

func runGitInit(absDir string) error {
	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = absDir
	initErr := cmdInit.Run()
	if initErr != nil {
		return apperror.WrapSimple(initErr, "git init:")
	}
	return nil
}

func maybeApplyCommonFiles(isCommon bool, absDir string) error {
	if !isCommon {
		return nil
	}
	return applyCommonFiles(absDir)
}

func maybeApplyCGFiles(isCG bool, absDir string) error {
	if !isCG {
		return nil
	}
	return applyCGFiles(absDir)
}

func writeInitialFiles(absDir string, p createRepoParams) {
	readmePath := filepath.Join(absDir, "README.md")
	readmeLower := filepath.Join(absDir, "readme.md")
	if !fileExists(readmePath) && !fileExists(readmeLower) {
		readmeContent := fmt.Sprintf("# %s\n\n%s\n", p.Name, p.Description)
		_ = os.WriteFile(readmePath, []byte(readmeContent), 0644)
	}

	gitignorePath := filepath.Join(absDir, ".gitignore")
	if !fileExists(gitignorePath) {
		gitignoreContent := ".DS_Store\nThumbs.db\nnode_modules/\nbin/\n*.log\n"
		_ = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	}
}

func fileExists(p string) bool {
	info, err := os.Stat(p)

	return err == nil && !info.IsDir()
}

func commitInitialFiles(absDir string) error {
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdStatus := exec.Command("git", "status", "--porcelain")
	cmdStatus.Dir = absDir
	out, _ := cmdStatus.Output()
	if len(strings.TrimSpace(string(out))) > 0 {
		cmdCommit := exec.Command("git", "commit", "-m", "feat: initial commit")
		cmdCommit.Dir = absDir
		_ = cmdCommit.Run()
	}

	return nil
}

func applyCommonFiles(absDir string) error {
	exe, err := os.Executable()
	if err != nil {
		exe = "gitmap"
	}

	cmdCommon := exec.Command(exe, "common", "all")
	cmdCommon.Dir = absDir
	_ = cmdCommon.Run()

	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdStatus := exec.Command("git", "status", "--porcelain")
	cmdStatus.Dir = absDir
	out, _ := cmdStatus.Output()
	if len(strings.TrimSpace(string(out))) > 0 {
		cmdCommit := exec.Command("git", "commit", "-m", "chore(common): apply common baselines and git-lfs configuration")
		cmdCommit.Dir = absDir
		_ = cmdCommit.Run()
	}

	return nil
}

func applyCGFiles(absDir string) error {
	backupBranch := fmt.Sprintf("backup/cg-sync-%s", time.Now().Format("20060102-150405"))
	cmdBranch := exec.Command("git", "branch", backupBranch)
	cmdBranch.Dir = absDir
	_ = cmdBranch.Run()

	baseDir := resolveCGBaseDir()
	cgSrc := filepath.Join(baseDir, "02-spec", "02-coding-guidelines")
	cgDst := filepath.Join(absDir, "02-spec", "02-coding-guidelines")

	if dirExists(cgSrc) {
		_ = copyDir(cgSrc, cgDst)
	} else {
		ensureCGOverviewFallback(cgDst)
	}

	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdStatus := exec.Command("git", "status", "--porcelain")
	cmdStatus.Dir = absDir
	out, _ := cmdStatus.Output()
	if len(strings.TrimSpace(string(out))) > 0 {
		cmdCommit := exec.Command("git", "commit", "-m", "docs(spec): synchronize latest coding guidelines to 02-spec/02-coding-guidelines/")
		cmdCommit.Dir = absDir
		_ = cmdCommit.Run()
	}

	return nil
}

func dirExists(p string) bool {
	info, err := os.Stat(p)

	return err == nil && info.IsDir()
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	return err
}

func resolveCGBaseDir() string {
	if constants.RepoPath != "" {
		return constants.RepoPath
	}
	exe, err := os.Executable()
	if err == nil {
		return filepath.Dir(filepath.Dir(exe))
	}
	return "."
}

func ensureCGOverviewFallback(cgDst string) {
	_ = os.MkdirAll(cgDst, 0755)
	readmePath := filepath.Join(cgDst, "00-overview.md")
	if fileExists(readmePath) {
		return
	}
	_ = os.WriteFile(readmePath, []byte("# Coding Guidelines\n\nCanonical coding guidelines standard.\n"), 0644)
}
