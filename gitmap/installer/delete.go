// Package installer — delete.go implements installer removal and OS target stripping.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Delete removes an installer script and all its version history.
func (m *Manager) Delete(slug string) error {
	if m == nil || m.db == nil {
		return apperror.New("Delete", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(slug) == "" {
		return apperror.New("Delete", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	return m.db.DeleteInstaller(slug)
}

// DeleteVersion removes a specific historical version tag for an installer.
func (m *Manager) DeleteVersion(slug, version string) error {
	if m == nil || m.db == nil {
		return apperror.New("DeleteVersion", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(slug) == "" || strings.TrimSpace(version) == "" {
		return apperror.New("DeleteVersion", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug and version are required"})
	}

	return m.db.DeleteInstallerVersion(slug, version)
}
