// Package installer — execute_order.go orchestrates multi-OS execution ordering.
package installer

import (
	"context"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ExecuteOrdered runs the installer script respecting the configured order mode.

func (m *Manager) ExecuteOrdered(ctx context.Context, slug, osTarget string) error {
	if m == nil || m.db == nil {
		return apperror.New("ExecuteOrdered", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "manager or db is nil"})
	}

	if strings.TrimSpace(slug) == "" {
		return apperror.New("ExecuteOrdered", "E_INSTALLER_INVALID_INPUT", map[string]any{"error": "slug cannot be empty"})
	}

	script, errGet := m.db.GetInstallerBySlug(slug)
	if errGet != nil {
		return errGet
	}

	order := script.OrderMode
	if order == "" {
		order = constants.OrderFallback
	}

	return m.dispatchOrder(ctx, script, osTarget, order)
}

func (m *Manager) dispatchOrder(
	ctx context.Context,
	s *model.InstallerScript,
	osTarget,
	order string,
) error {
	switch order {
	case constants.OrderUnixFirst:
		_ = RunLanguageScript(ctx, "echo [gitmap] running unix pre-verify", "sh")

		return RunInstallerCommand("echo [gitmap] installing for " + osTarget)
	case constants.OrderOSFirst:
		_ = RunInstallerCommand("echo [gitmap] installing for " + osTarget)

		return RunLanguageScript(ctx, "echo [gitmap] running unix post-verify", "sh")
	case constants.OrderOSOnly:
		return RunInstallerCommand("echo [gitmap] installing for " + osTarget)
	default: // fallback

		return RunInstallerCommand("echo [gitmap] executing fallback for " + osTarget)
	}
}
