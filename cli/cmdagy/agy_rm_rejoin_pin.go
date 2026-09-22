// Package cmdagy — agy_rm_rejoin_pin.go handles pinning and rejoining projects.
package cmdagy

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func handleRejoinPin(p AgyProject, isPin bool) error {
	if !isPin {
		return nil
	}
	if _, pinErr := addPinnedProjectTarget(p.ID); pinErr != nil {
		return apperror.WrapSimple(pinErr, "pin project on rejoin")
	}
	fmt.Printf("  %s✓ Pinned project '%s'%s\n", constants.ColorGreen, p.Name, constants.ColorReset)
	return nil
}

func rejoinAgyProject(p AgyProject) error {
	if err := createProjectFile(p.ID, p.Name); err != nil {
		return apperror.WrapSimple(err, "re-create project file")
	}
	if p.GetPath() != "" {
		workspacesync.SyncAntigravity(p.GetPath(), p.Name)
	}
	return nil
}
