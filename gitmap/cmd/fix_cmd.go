package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func resolveFixOption(args []string, aliasOverride string) (string, error) {
	if aliasOverride != "" {
		return aliasOverride, nil
	}

	if len(args) == 0 {
		return "", apperror.New("fix", "E_USAGE", map[string]any{"msg": "Usage: gitmap fix <1|2|3|stash|wip|discard>"})
	}

	return args[0], nil
}

func runFix(args []string, aliasOverride string) error {
	option, err := resolveFixOption(args, aliasOverride)
	if err != nil {
		return err
	}

	stateFile := getRemediationStateFile()
	b, err := os.ReadFile(stateFile)
	if err != nil {
		return apperror.New("fix", "E_NO_STATE", map[string]any{"msg": "No previous failed pull remediation state found."})
	}

	var state RemediationState
	if err := json.Unmarshal(b, &state); err != nil {
		return apperror.New("fix", "E_INVALID_STATE", map[string]any{"msg": "Invalid remediation state file."})
	}

	idx := -1
	switch option {
	case "1", "stash":
		idx = 0
	case "2", "wip":
		idx = 1
	case "3", "discard":
		idx = 2
	default:
		// maybe they passed an index
		if i, err := strconv.Atoi(option); err == nil && i > 0 && i <= len(state.Recipes) {
			idx = i - 1
		}
	}

	if idx == -1 || idx >= len(state.Recipes) {
		return apperror.New("fix", "E_INVALID_OPTION", map[string]any{"msg": fmt.Sprintf("Invalid fix option: %s", option)})
	}

	recipe := state.Recipes[idx]
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, state.RepoName)
	fmt.Printf("  Running: %s\n\n", recipe.Command)

	// recipe.Command contains chained git commands like `git -C ... stash && git -C ... pull`
	// We need to execute it via sh -c or cmd.exe /c
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", recipe.Command)
	} else {
		cmd = exec.Command("sh", "-c", recipe.Command)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\n%s Fix failed: %v\n", constants.ColorRed+"✗"+constants.ColorReset, err)
		return nil // don't crash, let user see git output
	}

	fmt.Printf("\n%s Fix applied !successfully\n", constants.ColorGreen+"✓"+constants.ColorReset)

	// Clean up state so we don't accidentally run it again blindly
	os.Remove(stateFile)

	return nil
}
