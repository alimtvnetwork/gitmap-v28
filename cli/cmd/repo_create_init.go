package cmd

import (
	"fmt"
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

	if mkErr := os.MkdirAll(absDir, 0755); mkErr != nil {
		return apperror.WrapSimple(mkErr, "create directory:")
	}

	cmdInit := exec.Command("git", "init", "-b", "main")
	cmdInit.Dir = absDir
	if initErr := cmdInit.Run(); initErr != nil {
		return apperror.WrapSimple(initErr, "git init:")
	}

	writeInitialFiles(absDir, p)

	return commitInitialFiles(absDir)
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
