// Package installer — create.go implements installer script creation business logic.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// slugify converts a human-readable name into a CLI-safe slug.
func slugify(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(name)) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		} else if r == ' ' || r == '_' || r == '.' {
			b.WriteRune('-')
		}
	}

	return strings.Trim(b.String(), "-")
}

// Create registers a new installer script in the database with defaults.
func (m *Manager) Create(name, desc string) (*model.InstallerScript, error) {
	if m == nil || m.db == nil {
		return nil, apperror.New("Create", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(name) == "" {
		return nil, apperror.New("Create", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "name cannot be empty"})
	}

	slug := slugify(name)
	script := &model.InstallerScript{
		Name:        strings.TrimSpace(name),
		Slug:        slug,
		Description: strings.TrimSpace(desc),
		TargetOS:    "all",
		Version:     "v1.0.0",
	}

	if err := m.db.CreateInstaller(script); err != nil {
		return nil, err
	}

	return script, nil
}
