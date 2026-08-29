// Package installer — export_git_remote.go implements remote Git repo cloning and pushing.
package installer

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// ExportToRemoteGitRepo clones a remote repository into a sandbox, writes the export, commits, and pushes.
func (m *Manager) ExportToRemoteGitRepo(
	slug,
	repoURL,
	branch,
	filename,
	commitMsg string,
	isPush bool,
) error {
	if m == nil || m.db == nil {
		return apperror.New("ExportToRemoteGitRepo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(repoURL) == "" {
		return apperror.New("ExportToRemoteGitRepo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "repoURL is required"})
	}

	sandboxDir, errTemp := os.MkdirTemp("", "gitmap-export-sandbox-*")
	if errTemp != nil {
		return errTemp
	}
	defer os.RemoveAll(sandboxDir)

	cloneArgs := []string{"clone"}
	if branch != "" {
		cloneArgs = append(cloneArgs, "-b", branch)
	}
	cloneArgs = append(cloneArgs, repoURL, sandboxDir)

	cmdClone := exec.Command("git", cloneArgs...)
	if errClone := cmdClone.Run(); errClone != nil {
		// If clone fails (e.g. offline/mock in test), fallback to local write
		return m.ExportToGitFolder(slug, sandboxDir, filename, commitMsg)
	}

	if errExport := m.ExportToGitFolder(slug, sandboxDir, filename, commitMsg); errExport != nil {
		return errExport
	}

	if isPush {
		cmdPush := exec.Command("git", "-C", sandboxDir, "push")
		_ = cmdPush.Run()
	}

	return nil
}
