package cmdagy

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/workspacesync"
)

func syncMissingAgyProjects(repos []localRepo, dirPath string) {
	existing, _ := loadAllAgyProjects(dirPath)
	pathMap := make(map[string]bool)
	for _, p := range existing {
		pathMap[strings.ToLower(filepath.Clean(p.GetPath()))] = true
	}

	for _, r := range repos {
		clean := strings.ToLower(filepath.Clean(r.Path))
		if !pathMap[clean] {
			if workspacesync.SyncAntigravity(r.Path, r.Name) {
				fmt.Printf("  %s Registered into Antigravity: %s (%s)\n",
					constants.ColorGreen+"✓"+constants.ColorReset, r.Name, r.Path)
			}
		}
	}
}

func broadcastReadMemoryToAll(dirPath string) error {
	projects, loadErr := loadAllAgyProjects(dirPath)
	if loadErr != nil {
		return loadErr
	}
	tokens := parseAgyExceptTokens(agyRprpExcept)
	targets, excluded := partitionPromptProjects(projects, tokens)
	if len(targets) == 0 {
		fmt.Println("No eligible Antigravity projects found.")

		return nil
	}
	agyAprmpPrompt = agyRprpPrompt
	agyAprmpDryRun = agyRprpDryRun
	agyAprmpYes = agyRprpYes

	return executePromptBroadcast(targets, excluded)
}
