// Package cmdzsh provides Oh-My-Zsh repository cloning and plugin installation.
package cmdzsh

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const (
	ohMyZshGitUrl         = "https://github.com/ohmyzsh/ohmyzsh.git"
	autosuggestionsGitUrl = "https://github.com/zsh-users/zsh-autosuggestions.git"
)

// InstallOhMyZsh installs Oh-My-Zsh into the target home directory.
func InstallOhMyZsh(targetHome string, isDryRun bool) *apperror.AppError {
	home := ResolveTargetHome(targetHome)
	dstDir := filepath.Join(home, ".oh-my-zsh")
	if HasDirectory(dstDir) {
		return nil
	}

	if isDryRun {
		fmt.Printf("[dry-run] Would clone Oh-My-Zsh into %s%s", dstDir, constants.NewLineUnix)

		return nil
	}

	return cloneAndInitOhMyZsh(home, dstDir)
}

func cloneAndInitOhMyZsh(home, dstDir string) *apperror.AppError {
	cmd := exec.Command("git", "clone", "--depth=1", ohMyZshGitUrl, dstDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "cmdzsh.cloneAndInitOhMyZsh", map[string]any{"output": string(out)})
	}

	initZshrcTemplate(home, dstDir)
	chmodOhMyZshScript(dstDir)

	return nil
}

func initZshrcTemplate(home, dstDir string) {
	zshrcPath := filepath.Join(home, ".zshrc")
	if HasFile(zshrcPath) {
		return
	}

	templatePath := filepath.Join(dstDir, "templates", "zshrc.zsh-template")
	if HasFile(templatePath) {
		_ = CopyFile(templatePath, zshrcPath)
	}
}

func chmodOhMyZshScript(dstDir string) {
	scriptPath := filepath.Join(dstDir, "oh-my-zsh.sh")
	if HasFile(scriptPath) {
		_ = os.Chmod(scriptPath, 0755)
	}
}

// InstallAutosuggestions clones zsh-autosuggestions plugin into custom plugins directory.
func InstallAutosuggestions(targetHome string, isDryRun bool) *apperror.AppError {
	home := ResolveTargetHome(targetHome)
	pluginDir := filepath.Join(home, ".oh-my-zsh", "custom", "plugins", "zsh-autosuggestions")
	if HasDirectory(pluginDir) {
		return nil
	}

	if isDryRun {
		fmt.Printf("[dry-run] Would clone zsh-autosuggestions into %s%s", pluginDir, constants.NewLineUnix)

		return nil
	}

	return cloneAutosuggestionsRepo(pluginDir)
}

func cloneAutosuggestionsRepo(pluginDir string) *apperror.AppError {
	_ = EnsureDirectory(filepath.Dir(pluginDir), 0755)
	cmd := exec.Command("git", "clone", "--depth=1", autosuggestionsGitUrl, pluginDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.Wrap(err, "cmdzsh.cloneAutosuggestionsRepo", map[string]any{"output": string(out)})
	}

	return nil
}
