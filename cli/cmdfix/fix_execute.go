// Package cmd — fix_execute.go runs remediation recipe steps with native argument isolation.
package cmdfix

import (
	"bytes"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func executeFixRecipe(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
	fmt.Printf("%s Applying Fix: %s on %s\n", constants.ColorCyan+"ℹ"+constants.ColorReset, recipe.Title, item.RepoName)
	if recipe.Description != "" {
		fmt.Printf("  Plan:    %s\n\n", recipe.Description)
	}

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

func executeStructuredRecipe(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) error {
	total := len(recipe.Steps)
	for i, step := range recipe.Steps {
		err := executeSingleStep(item.RepoName, i+1, total, step)
		if err != nil {
			return err
		}
	}

	fmt.Printf("\n%s Fix applied successfully on %s\n", constants.ColorGreen+"✓"+constants.ColorReset, item.RepoName)
	cmdremediation.RemoveRemediationItem(item.RepoName)

	return nil
}

func formatStepCommand(step gitutil.RemediationStep) string {
	if len(step.Args) >= 2 && step.Args[0] == "-C" {
		return step.Name + " " + strings.Join(step.Args[2:], " ")
	}

	return step.Name + " " + strings.Join(step.Args, " ")
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

	fmt.Printf("%s failed\n", constants.ColorRed+"✖"+constants.ColorReset)
	printBluntRemediationFailure(repoName, step, outStr, err)

	return err
}

func isBenignStashClean(step gitutil.RemediationStep, output string) bool {
	hasPop := false
	for _, arg := range step.Args {
		if arg == "pop" {
			hasPop = true
			break
		}
	}
	if !hasPop {
		return false
	}

	lower := strings.ToLower(output)

	return strings.Contains(lower, "no stash entries found") || strings.Contains(lower, "no stash found")
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

func recoverStashCollision(step gitutil.RemediationStep, output string) bool {
	hasPop := checkHasPop(step.Args)
	if !hasPop {
		return false
	}

	lowerOut := strings.ToLower(output)
	hasUntracked := strings.Contains(lowerOut, "could not restore untracked files from stash")
	hasCollisionNotice := strings.Contains(lowerOut, "already exists, no checkout")

	if hasUntracked && hasCollisionNotice {
		return dropStashCollision(step.Args)
	}

	return false
}

func checkHasPop(args []string) bool {
	hasPop := false
	for _, arg := range args {
		if arg == "pop" {
			hasPop = true
		}
	}

	return hasPop
}

func dropStashCollision(args []string) bool {
	hasArgs := len(args) >= 2
	if !hasArgs {
		return false
	}

	isDirFlag := args[0] == "-C"
	if !isDirFlag {
		return false
	}

	err := exec.Command("git", "-C", args[1], "stash", "drop").Run()
	isSuccess := err == nil

	return isSuccess
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
