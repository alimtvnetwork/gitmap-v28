// Package cmd — fix_execute.go runs remediation recipe steps with native argument isolation.
package cmd

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

func executeFixRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, item.RepoName)
	fmt.Printf("  Plan:    %s\n\n", recipe.Command)

	if len(recipe.Steps) == 0 && item.RepoPath != "" {
		recipe.Steps = synthesizeRecipeSteps(recipe, item.RepoPath)
	}

	if len(recipe.Steps) > 0 {
		return executeStructuredRecipe(item, recipe)
	}

	return executeShellFallback(item, recipe)
}

func synthesizeRecipeSteps(recipe gitutil.RemediationRecipe, repoPath string) []gitutil.RemediationStep {
	titleLower := strings.ToLower(recipe.Title)
	cmdLower := strings.ToLower(recipe.Command)
	if strings.Contains(titleLower, "wip") || strings.Contains(titleLower, "commit") || strings.Contains(cmdLower, "commit") {
		return gitutil.GenerateCommitRecipe(repoPath).Steps
	}

	if strings.Contains(titleLower, "discard") || strings.Contains(titleLower, "clean") || strings.Contains(cmdLower, "reset --hard") {
		return gitutil.GenerateDiscardRecipe(repoPath).Steps
	}

	if strings.Contains(titleLower, "stash") || strings.Contains(cmdLower, "stash") {
		return gitutil.GenerateStashRecipe(repoPath).Steps
	}

	return nil
}

func executeStructuredRecipe(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	total := len(recipe.Steps)
	for i, step := range recipe.Steps {
		err := executeSingleStep(item.RepoName, i+1, total, step)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	RemoveRemediationItem(item.RepoName)

	return nil
}

func executeSingleStep(repoName string, idx, total int, step gitutil.RemediationStep) error {
	stepCmd := step.Name + " " + strings.Join(step.Args, " ")
	fmt.Printf("  [%d/%d] ➜ %s ... ", idx, total, stepCmd)

	cmd := exec.Command(step.Name, step.Args...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err == nil {
		fmt.Printf("%s ok\n", constants.ColorGreen+"✔"+constants.ColorReset)

		return nil
	}

	outStr := outBuf.String()
	if isBenignCommitClean(step, outStr) {
		fmt.Printf("%s clean (nothing to commit)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
	printBluntRemediationFailure(repoName, step, outStr, err)

	return err
}

func isBenignCommitClean(step gitutil.RemediationStep, output string) bool {
	hasCommit := false
	for _, arg := range step.Args {
		if arg == "commit" {
			hasCommit = true
			break
		}
	}

	if !hasCommit {
		return false
	}

	lower := strings.ToLower(output)

	return strings.Contains(lower, "nothing to commit") || strings.Contains(lower, "working tree clean")
}

func executeShellFallback(item *RemediationItem, recipe gitutil.RemediationRecipe) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", recipe.Command)
	} else {
		cmd = exec.Command("sh", "-c", recipe.Command)
	}

	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err != nil {
		step := gitutil.RemediationStep{Name: "shell", Args: []string{recipe.Command}}
		printBluntRemediationFailure(item.RepoName, step, outBuf.String(), err)

		return err
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	RemoveRemediationItem(item.RepoName)

	return nil
}
