// Package cmdagy — agy_rm_rejoin_ops.go handles rm-rejoin-read and rm-rejoin-pin-read lifecycle.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func runAgyRmRejoin(args []string, isPin bool) error {
	dirPath, err := getProjectsDirPath()
	if err != nil {
		return apperror.WrapSimple(err, "get projects dir path")
	}

	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return apperror.WrapSimple(loadErr, "load agy projects")
	}
	sortAgyProjects(projects, "name")

	targets, resolveErr := ResolveAgyProjectTargets(args, projects)
	if resolveErr != nil {
		return resolveErr
	}

	return executeRmRejoinTargets(targets, isPin)
}

func executeRmRejoinTargets(targets []AgyProject, isPin bool) error {
	_, _ = snapshotAgyProjects("pre-rrr")
	var lastErr error
	rejoinedCount := 0

	for _, p := range targets {
		if err := processSingleRmRejoin(p, isPin); err != nil {
			lastErr = err
			fmt.Printf("  %s✗ Failed rm-rejoin for %s: %v%s\n",
				constants.ColorRed, p.Name, err, constants.ColorReset)
			continue
		}
		rejoinedCount++
	}

	printAgyUndoGuidance()
	if rejoinedCount == 0 && lastErr != nil {
		return apperror.WrapSimple(lastErr, "rm-rejoin failed")
	}
	return nil
}

func processSingleRmRejoin(p AgyProject, isPin bool) error {
	inv, _ := json.Marshal(p)
	_ = recordAgyTask("rm-rejoin-read", p.Name, p.ID, string(inv))

	purged, purgeErr := purgeProjectConversations(p)
	if purgeErr != nil {
		fmt.Printf("  %s! Notice: conversation purge: %v%s\n",
			constants.ColorYellow, purgeErr, constants.ColorReset)
	}
	if err := deleteProjectFile(p.ID); err != nil && !os.IsNotExist(err) {
		return apperror.WrapSimple(err, "delete project file before rejoin")
	}

	if err := rejoinAgyProject(p); err != nil {
		return err
	}
	fmt.Printf("  %s✓ Rejoined project '%s' (purged %d convs)%s\n",
		constants.ColorGreen, p.Name, purged, constants.ColorReset)

	if err := handleRejoinPin(p, isPin); err != nil {
		return err
	}

	return dispatchProjectReadPrompt(p)
}

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
