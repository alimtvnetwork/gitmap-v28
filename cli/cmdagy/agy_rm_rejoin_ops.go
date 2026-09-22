// Package cmdagy — agy_rm_rejoin_ops.go handles rm-rejoin-read and rm-rejoin-pin-read lifecycle.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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

	for _, p := range targets {
		if err := processSingleRmRejoin(p, isPin); err != nil {
			fmt.Printf("  %s✗ Failed rm-rejoin for %s: %v%s\n",
				constants.ColorRed, p.Name, err, constants.ColorReset)
			continue
		}
	}

	printAgyUndoGuidance()
	return nil
}

func processSingleRmRejoin(p AgyProject, isPin bool) error {
	inv, _ := json.Marshal(p)
	_ = recordAgyTask("rm-rejoin-read", p.Name, p.ID, string(inv))

	purged, _ := purgeProjectConversations(p)
	_ = deleteProjectFile(p.ID)

	if err := rejoinAgyProject(p); err != nil {
		return err
	}
	fmt.Printf("  %s✓ Rejoined project '%s' (purged %d convs)%s\n",
		constants.ColorGreen, p.Name, purged, constants.ColorReset)

	if isPin {
		_, _ = addPinnedProjectTarget(p.ID)
		fmt.Printf("  %s✓ Pinned project '%s'%s\n", constants.ColorGreen, p.Name, constants.ColorReset)
	}

	return dispatchProjectReadPrompt(p)
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

func dispatchProjectReadPrompt(p AgyProject) error {
	ideProc := DetectRunningAntigravityIDE()
	pid := 0
	if ideProc.IsSuccess() {
		pid = ideProc.Value.PID
	}
	promptPath := filepath.Join(os.TempDir(), fmt.Sprintf("agy-read-prompt-%s.txt", p.ID))
	_ = os.WriteFile(promptPath, []byte(defaultReadMemoryPrompt), 0644)

	res := DispatchPromptToAntigravity(p.GetPath(), promptPath, p.Name, defaultReadMemoryPrompt, pid)
	if res.IsSuccess {
		fmt.Printf("  %s %s\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)
	}
	_ = renameProjectInitialConversation(p, p.Name)
	return nil
}
