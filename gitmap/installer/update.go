// Package installer — update.go implements installer update business logic.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// Update modifies an existing installer script and archives previous state.
func (m *Manager) Update(slug, osTarget string) (*model.InstallerScript, error) {
	if m == nil || m.db == nil {
		return nil, apperror.New("Update", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return nil, apperror.New("Update", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	existing, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return nil, errGet
	}

	// Archive old version
	oldVersion := &model.InstallerVersion{
		ScriptID:     existing.ID,
		Version:      existing.Version,
		TargetOS:     existing.TargetOS,
		Instructions: existing.Instructions,
	}

	if errSave := m.db.SaveVersion(oldVersion); errSave != nil {
		return nil, errSave
	}

	if osTarget != "" {
		existing.TargetOS = osTarget
	}

	existing.Version = NextSemanticVersion(existing.Version)

	return existing, nil
}
