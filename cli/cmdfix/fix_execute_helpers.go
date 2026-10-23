// Package cmdfix — fix_execute_helpers.go provides recipe synthesis, telemetry persistence,
// and step error recovery helpers for remediation execution.
package cmdfix

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func synthesizeRecipeSteps(recipe gitutil.RemediationRecipe, repoPath string) []gitutil.RemediationStep {
	steps := synthesizeSpecialStateSteps(repoPath)
	if len(steps) > 0 {
		return steps
	}

	return synthesizeActionSteps(recipe, repoPath)
}

func synthesizeSpecialStateSteps(repoPath string) []gitutil.RemediationStep {
	steps := resolveNoInitialCommitSteps(repoPath)
	if len(steps) > 0 {
		return steps
	}

	return resolveLocalOnlySteps(repoPath)
}

func resolveNoInitialCommitSteps(repoPath string) []gitutil.RemediationStep {
	if gitutil.HasInitialCommit(repoPath) {
		return nil
	}

	recipes := gitutil.GenerateNoInitialCommitRecipes(repoPath)
	if len(recipes) == 0 {
		return nil
	}

	return recipes[0].Steps
}

func resolveLocalOnlySteps(repoPath string) []gitutil.RemediationStep {
	if gitutil.HasRemote(repoPath) {
		return nil
	}

	recipes := gitutil.GenerateLocalOnlyRecipes(repoPath)
	if len(recipes) == 0 {
		return nil
	}

	return recipes[0].Steps
}

func synthesizeActionSteps(recipe gitutil.RemediationRecipe, repoPath string) []gitutil.RemediationStep {
	titleLower := strings.ToLower(recipe.Title)
	cmdLower := strings.ToLower(recipe.Command)
	if strings.Contains(titleLower, "case collision") || strings.Contains(titleLower, "normalize case") {
		return synthesizeCaseCollisionSteps(repoPath)
	}
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

func synthesizeCaseCollisionSteps(repoPath string) []gitutil.RemediationStep {
	cols, err := gitutil.DetectCaseCollisions(repoPath)
	if err == nil && len(cols) > 0 {
		return gitutil.GenerateCaseNormalizationRecipe(repoPath, cols[0]).Steps
	}

	return nil
}

func hasCaseCollisionRecipe(recipe gitutil.RemediationRecipe) bool {
	return strings.Contains(strings.ToLower(recipe.Title), "case collision")
}

type lastRemediationAction struct {
	Type    string   `json:"type"`
	Actions []string `json:"actions"`
}

type lastRemediationData struct {
	RepoPath      string                  `json:"repoPath"`
	RepoName      string                  `json:"repoName"`
	SummaryReason string                  `json:"summaryReason"`
	Recipes       []lastRemediationAction `json:"recipes"`
}

func persistCaseRemediationTelemetry(item *cmdremediation.RemediationItem, recipe gitutil.RemediationRecipe) {
	if item == nil || !hasCaseCollisionRecipe(recipe) {
		return
	}
	var actions []string
	for _, s := range recipe.Steps {
		actions = append(actions, s.Name+" "+strings.Join(s.Args, " "))
	}
	entry := lastRemediationAction{Type: "case_normalization", Actions: actions}
	data := lastRemediationData{
		RepoPath:      item.RepoPath,
		RepoName:      item.RepoName,
		SummaryReason: "Git index case collision normalized",
		Recipes:       []lastRemediationAction{entry},
	}
	_ = writeLastRemediationJSON(data)
}

func writeLastRemediationJSON(data lastRemediationData) error {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	home, _ := os.UserHomeDir()
	outDir := filepath.Join(home, ".gitmap", "output")
	_ = os.MkdirAll(outDir, 0755)

	return os.WriteFile(filepath.Join(outDir, "last_remediation.json"), b, 0644)
}

func formatStepCommand(step gitutil.RemediationStep) string {
	if len(step.Args) >= 2 && step.Args[0] == "-C" {
		return step.Name + " " + strings.Join(step.Args[2:], " ")
	}

	return step.Name + " " + strings.Join(step.Args, " ")
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

func summarizeRecipeAction(recipe gitutil.RemediationRecipe) string {
	titleLower := strings.ToLower(recipe.Title)
	if strings.Contains(titleLower, "stash") {
		return "stashed -> pulled -> re-applied"
	}
	if strings.Contains(titleLower, "wip") || strings.Contains(titleLower, "commit") {
		return "committed WIP changes -> pulled"
	}
	if strings.Contains(titleLower, "discard") || strings.Contains(titleLower, "clean") {
		return "discarded changes"
	}
	if strings.Contains(titleLower, "case") {
		return "normalized casing"
	}
	if recipe.Title != "" {
		return recipe.Title
	}

	return "remediated"
}
