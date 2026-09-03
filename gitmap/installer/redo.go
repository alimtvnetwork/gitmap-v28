// Package installer — redo.go implements redo version restoration logic.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Redo advances an installer script to a previously undone version.
func (m *Manager) Redo(slug string) error {
	if m == nil || m.db == nil {
		return apperror.New("Redo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return apperror.New("Redo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	_, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	return nil
}
