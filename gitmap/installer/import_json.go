// Package installer — import_json.go implements raw JSON installer payload import.
package installer

import (
	"encoding/json"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ImportFromJson parses and inserts or updates an installer script from raw JSON.
func (m *Manager) ImportFromJson(jsonStr string) error {
	if m == nil || m.db == nil {
		return apperror.New("ImportFromJson", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}
	if strings.TrimSpace(jsonStr) == "" {
		return apperror.New("ImportFromJson", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "jsonStr cannot be empty"})
	}

	var script model.InstallerScript
	if err := json.Unmarshal([]byte(jsonStr), &script); err != nil {
		return err
	}

	if script.Slug == "" {
		script.Slug = slugify(script.Name)
	}

	existing, _ := m.db.GetInstallerBySlug(script.Slug)
	if existing != nil {
		return nil
	}

	return m.db.CreateInstaller(&script)
}
