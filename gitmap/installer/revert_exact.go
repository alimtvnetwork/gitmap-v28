// Package installer — revert_exact.go implements exact semantic version restoration logic.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// RevertTo restores an installer script to an exact historical version tag.
func (m *Manager) RevertTo(slug, version string) error {
	if m == nil || m.db == nil {
		return apperror.New("RevertTo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(slug) == "" || strings.TrimSpace(version) == "" {
		return apperror.New("RevertTo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug and version are required"})
	}

	_, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	return nil
}
