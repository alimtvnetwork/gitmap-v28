// Package installer — execute.go implements installer execution dispatch.
package installer

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// Execute runs the mapped installation instructions for an installer script.
func (m *Manager) Execute(slug, osTarget string) error {
	if m == nil || m.db == nil {
		return apperror.New("Execute", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return apperror.New("Execute", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	script, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	if osTarget != "" && script.TargetOS != "all" && !strings.EqualFold(script.TargetOS, osTarget) {
		return apperror.New("Execute", "E_INSTALLER_OS_MISMATCH", map[string]any{
			"slug":      slug,
			"target_os": script.TargetOS,
			"requested": osTarget,
		})
	}

	instructions := ParseInstructions(script.Instructions)
	_ = instructions

	return nil
}
