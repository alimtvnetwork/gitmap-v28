package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
		return nil
	}
	writeInitialFiles(absDir, p)
	return commitInitialFiles(absDir)
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
	readmeContent := fmt.Sprintf("# %s\n\n%s\n", p.Name, p.Description)
	_ = os.WriteFile(readmePath, []byte(readmeContent), 0644)

	gitignorePath := filepath.Join(absDir, ".gitignore")
	gitignoreContent := ".DS_Store\nThumbs.db\nnode_modules/\nbin/\n*.log\n"
	_ = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}

func commitInitialFiles(absDir string) error {
	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "commit", "-m", "feat: initial commit")
	cmdCommit.Dir = absDir
	_ = cmdCommit.Run()

	return nil
}

func applyCommonFiles(absDir string) error {
	exe, err := os.Executable()
	if err != nil {
		exe = "gitmap"
	}

	cmdCommon := exec.Command(exe, "common", "all")
	cmdCommon.Dir = absDir
	if err := cmdCommon.Run(); err != nil {
		return apperror.WrapSimple(err, "gitmap common all:")
	}

	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "commit", "-m", "chore: apply common files")
	cmdCommit.Dir = absDir
	_ = cmdCommit.Run()

	return nil
}

func applyCGFiles(absDir string) error {
	cmdBranch := exec.Command("git", "checkout", "-b", "backup-cg-sync")
	cmdBranch.Dir = absDir
	if err := cmdBranch.Run(); err != nil {
		return apperror.WrapSimple(err, "git checkout -b backup-cg-sync:")
	}

	cgSrc := filepath.Join("02-spec", "02-coding-guidelines")
	cgDst := filepath.Join(absDir, "02-spec", "02-coding-guidelines")
	if err := copyDir(cgSrc, cgDst); err != nil {
		return apperror.WrapSimple(err, "copy cg files:")
	}

	cmdAdd := exec.Command("git", "add", ".")
	cmdAdd.Dir = absDir
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "commit", "-m", "chore: sync coding guidelines")
	cmdCommit.Dir = absDir
	_ = cmdCommit.Run()

	return nil
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
