// Package cmdagy — agy_pin_projects_json.go handles JSON output and path verification for pinned projects.
package cmdagy

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func checkPathMissing(path string) bool {
	if path == "" {
		return true
	}

	_, err := os.Stat(path)

	return err != nil
}

func printPinnedSummary(projects []PinnedProject) {
	missing := 0

	for _, p := range projects {
		if checkPathMissing(p.Path) {
			missing++
		}
	}

	active := len(projects) - missing
	fmt.Println()
	fmt.Printf("  %s%s%s\n", constants.ColorDim, constants.TermTableRule, constants.ColorReset)
	fmt.Printf("  %d pinned · %s%d active%s · %s%d missing%s\n\n",
		len(projects),
		constants.ColorGreen, active, constants.ColorReset,
		constants.ColorRed, missing, constants.ColorReset)
}

func outputPinnedProjectsJSON(projects []PinnedProject) error {
	data, err := json.MarshalIndent(projects, "", "  ")

	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
