// Package cmdfix — fix_execute.go runs remediation recipe steps with native argument isolation.
package cmdfix

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func executeFixRecipe(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
	if isRepoAlreadyClean(item) && !hasCaseCollisionRecipe(recipe) {
		fmt.Printf("%s Repository %s is already clean (no pending changes)\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
		cmdremediation.RemoveRemediationItem(item.RepoName)

		return nil
	}

	recipe = prepareRecipeForExecution(item, recipe)
	if len(recipe.Steps) > 0 {
		return executeStructuredRecipe(item, recipe)
	}

	return executeShellFallback(item, recipe)
}

func prepareRecipeForExecution(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) gitutil.RemediationRecipe {
	recipe = sanitizeRecipeForRepo(item.RepoPath, recipe)
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, item.RepoName)
	if recipe.Description != "" {
		fmt.Printf("  Plan:    %s\n\n", recipe.Description)
	}

	if len(recipe.Steps) == 0 && item.RepoPath != "" {
		recipe.Steps = synthesizeRecipeSteps(recipe, item.RepoPath)
	}

	return recipe
}

func executeStructuredRecipe(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
	total := len(recipe.Steps)
	for i, step := range recipe.Steps {
		err := executeSingleStep(item.RepoName, i+1, total, step)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	persistCaseRemediationTelemetry(item, recipe)
	cmdremediation.RemoveRemediationItem(item.RepoName)

	return nil
}

func executeSingleStep(repoName string, idx, total int, step gitutil.RemediationStep) error {
	stepCmd := formatStepCommand(step)
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

	return handleStepFailure(repoName, step, outBuf.String(), err)
}

func handleStepFailure(repoName string, step gitutil.RemediationStep, outStr string, err error) error {
	if isBenignCommitClean(step, outStr) || isBenignStashClean(step, outStr) {
		fmt.Printf("%s clean (nothing to commit/apply)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	isRecovered := recoverStashCollision(step, outStr)
	if isRecovered {
		fmt.Printf("%s recovered (dropped stashed untracked collision)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	isNoInit, msg := recoverNoInitialCommit(step, outStr)
	if isNoInit {
		fmt.Printf("%s %s\n", constants.ColorYellow+"•"+constants.ColorReset, msg)

		return nil
	}

	if isBenignPullNoTracking(step, outStr) {
		fmt.Printf("%s clean (no remote to pull)\n", constants.ColorYellow+"•"+constants.ColorReset)

		return nil
	}

	fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
	printBluntRemediationFailure(repoName, step, outStr, err)

	return err
}

func executeShellFallback(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
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
	cmdremediation.RemoveRemediationItem(item.RepoName)

	return nil
}

// FixStepResult captures the outcome of a batch recipe execution.
type FixStepResult struct {
	RepoName      string
	ActionSummary string
	Duration      time.Duration
	IsClean       bool
	Err           error
}

func executeFixRecipeBatch(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) FixStepResult {
	start := time.Now()
	if isRepoAlreadyClean(item) && !hasCaseCollisionRecipe(recipe) {
		cmdremediation.RemoveRemediationItem(item.RepoName)

		return FixStepResult{
			RepoName:      item.RepoName,
			ActionSummary: "already clean",
			Duration:      time.Since(start),
			IsClean:       true,
		}
	}
	recipe = sanitizeRecipeForRepo(item.RepoPath, recipe)
	if len(recipe.Steps) == 0 && item.RepoPath != "" {
		recipe.Steps = synthesizeRecipeSteps(recipe, item.RepoPath)
	}
	if len(recipe.Steps) > 0 {
		return runStructuredRecipeBatch(item, recipe, start)
	}

	return runShellFallbackBatch(item, recipe, start)
}

func runStructuredRecipeBatch(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe, start time.Time) FixStepResult {
	total := len(recipe.Steps)
	for i, step := range recipe.Steps {
		err := executeSingleStepSilent(item.RepoName, i+1, total, step)
		if err != nil {
			return FixStepResult{
				RepoName:      item.RepoName,
				ActionSummary: formatStepCommand(step),
				Duration:      time.Since(start),
				Err:           err,
			}
		}
	}
	persistCaseRemediationTelemetry(item, recipe)
	cmdremediation.RemoveRemediationItem(item.RepoName)

	return FixStepResult{
		RepoName:      item.RepoName,
		ActionSummary: summarizeRecipeAction(recipe),
		Duration:      time.Since(start),
	}
}

func runShellFallbackBatch(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe, start time.Time) FixStepResult {
	err := executeShellFallbackSilent(item, recipe)
	if err != nil {
		return FixStepResult{
			RepoName:      item.RepoName,
			ActionSummary: recipe.Command,
			Duration:      time.Since(start),
			Err:           err,
		}
	}
	cmdremediation.RemoveRemediationItem(item.RepoName)

	return FixStepResult{
		RepoName:      item.RepoName,
		ActionSummary: summarizeRecipeAction(recipe),
		Duration:      time.Since(start),
	}
}

func executeSingleStepSilent(repoName string, idx, total int, step gitutil.RemediationStep) error {
	cmd := exec.Command(step.Name, step.Args...)
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	err := cmd.Run()
	if err == nil {
		return nil
	}

	return handleStepFailureSilent(repoName, step, outBuf.String(), err)
}

func handleStepFailureSilent(repoName string, step gitutil.RemediationStep, outStr string, err error) error {
	if isBenignCommitClean(step, outStr) || isBenignStashClean(step, outStr) {
		return nil
	}
	if recoverStashCollision(step, outStr) {
		return nil
	}
	if isNoInit, _ := recoverNoInitialCommit(step, outStr); isNoInit {
		return nil
	}
	if isBenignPullNoTracking(step, outStr) {
		return nil
	}

	return err
}

func executeShellFallbackSilent(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", recipe.Command)
	} else {
		cmd = exec.Command("sh", "-c", recipe.Command)
	}
	var outBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &outBuf

	return cmd.Run()
}
