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
	if m == nil || m.db == nil {
		return apperror.New("ExportToGitFolder", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(folderPath) == "" {
		return apperror.New("ExportToGitFolder", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "folderPath is required"})
	}

	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}

	var scripts []model.InstallerScript
	if slug != "" {
		s, errGet := m.db.GetInstallerBySlug(slug)
		if errGet != nil {
			return errGet
		}
		scripts = append(scripts, *s)
	} else {
		all, errList := m.db.ListInstallers()
		if errList != nil {
			return errList
		}
		scripts = all
	}

	targetName := filename
	if targetName == "" {
		if slug != "" {
			targetName = fmt.Sprintf("%s.json", slug)
		} else {
			targetName = "gitmap-installers.json"
		}
	}

	targetFile := filepath.Join(folderPath, targetName)
	data, errMarshal := json.MarshalIndent(scripts, "", "  ")
	if errMarshal != nil {
		return errMarshal
	}

	if errWrite := os.WriteFile(targetFile, data, 0644); errWrite != nil {
		return errWrite
	}

	// Git add and commit if it is a git repository
	msg := commitMsg
	if msg == "" {
		msg = fmt.Sprintf("chore: update gitmap installer definitions (%s)", targetName)
	}

	cmdAdd := exec.Command("git", "-C", folderPath, "add", targetName)
	_ = cmdAdd.Run()

	cmdCommit := exec.Command("git", "-C", folderPath, "commit", "-m", msg)
	_ = cmdCommit.Run()

	return nil
}
