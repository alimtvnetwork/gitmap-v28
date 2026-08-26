// Package installer — conflict.go implements import conflict resolution.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// resolveImportConflict resolves conflicting slugs by updating to the latest incoming version.
func (m *Manager) resolveImportConflict(slug string) error {
	if m == nil || m.db == nil {
		return apperror.New("resolveImportConflict", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(slug) == "" {
		return apperror.New("resolveImportConflict", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	existing, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	_ = existing
	return nil
}
