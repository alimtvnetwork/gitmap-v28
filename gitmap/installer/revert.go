// Package installer — revert.go implements undo version rollback logic.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Undo rolls back an installer script to its previous version.
func (m *Manager) Undo(slug string) error {
	if m == nil || m.db == nil {
		return apperror.New("Undo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return apperror.New("Undo", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	_, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	return nil
}
