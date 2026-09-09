// Package installer — export_git.go implements Git directory direct export and commit.
package installer

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ExportToGitFolder exports installers into a local Git folder and commits the changes.
func (m *Manager) ExportToGitFolder(slug, folderPath, filename, commitMsg string) error {
	if errValidate := validateExportInputs(m, folderPath); errValidate != nil {
		return errValidate
	}
	if errWrite := m.writeExportScriptFile(slug, folderPath, filename); errWrite != nil {
		return errWrite
	}

	return commitIfGitRepo(folderPath, filename, slug, commitMsg)
}

func validateExportInputs(m *Manager, folderPath string) error {
	if m == nil || m.db == nil {
		return apperror.New("ExportToGitFolder", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(folderPath) == "" {
		return apperror.New("ExportToGitFolder", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "folderPath is required"})
	}

	return os.MkdirAll(folderPath, 0755)
}

func (m *Manager) writeExportScriptFile(slug, folderPath, filename string) error {
	scripts, errScripts := m.resolveExportScripts(slug)
	if errScripts != nil {
		return errScripts
	}
	targetName := resolveExportTargetFilename(filename, slug)
	targetFile := filepath.Join(folderPath, targetName)
	data, errMarshal := json.MarshalIndent(scripts, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	return os.WriteFile(targetFile, data, 0644)
}

func commitIfGitRepo(folderPath, filename, slug, commitMsg string) error {
	gitDir := filepath.Join(folderPath, ".git")
	if _, errStat := os.Stat(gitDir); errStat != nil {
		return nil
	}
	targetName := resolveExportTargetFilename(filename, slug)
	msg := resolveExportCommitMsg(commitMsg, targetName)
	cmdAdd := exec.Command("git", "-C", folderPath, "add", targetName)
	_ = cmdAdd.Run()
	cmdCommit := exec.Command("git", "-C", folderPath, "commit", "-m", msg)
	_ = cmdCommit.Run()

	return nil
}

func resolveExportCommitMsg(commitMsg, targetName string) string {
	if commitMsg != "" {
		return commitMsg
	}

	return fmt.Sprintf("chore: update gitmap installer definitions (%s)", targetName)
}

func (m *Manager) resolveExportScripts(slug string) ([]model.InstallerScript, error) {
	if slug == "" {
		return m.db.ListInstallers()
	}

	s, err := m.db.GetInstallerBySlug(slug)
	if err != nil {
		return nil, err
	}

	return []model.InstallerScript{*s}, nil
}

func resolveExportTargetFilename(filename, slug string) string {
	if filename != "" {
		return filename
	}

	if slug != "" {
		return fmt.Sprintf("%s.json", slug)
	}

	return "gitmap-installers.json"
}
