// Package cmdagy — agy_prompt_dispatch.go dispatches prompts to Antigravity sessions.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ExecuteSendPrompt dispatches a prompt to a target project by slug or path.
func ExecuteSendPrompt(projectTarget, promptMsg, title string) error {
	cleanMsg := strings.TrimSpace(promptMsg)
	if cleanMsg == "" {
		return apperror.NewSimple("prompt message is required (--prompt / -m or positional argument)", "E9040")
	}

	repoRoot, projName, err := resolveTargetProjectPath(projectTarget)
	if err != nil {
		return err
	}

	return prepareAndDispatchPrompt(repoRoot, projName, cleanMsg, title)
}

func prepareAndDispatchPrompt(repoRoot, projName, promptMsg, title string) error {
	content, stageFile, resolvedTitle, err := resolvePromptPayloadAndStage(promptMsg, projName, title)
	if err != nil {
		return err
	}

	pid := resolveActiveAntigravityPID()

	return dispatchAndRenderResult(repoRoot, stageFile, resolvedTitle, content, pid)
}

func resolveTargetProjectPath(target string) (string, string, error) {
	cleanTarget := strings.TrimSpace(target)
	if cleanTarget == "" {
		cleanTarget = "."
	}

	if fi, err := os.Stat(cleanTarget); err == nil && fi.IsDir() {
		absPath, _ := filepath.Abs(cleanTarget)

		return absPath, filepath.Base(absPath), nil
	}

	return lookupTargetProjectInRegistry(cleanTarget)
}

func lookupTargetProjectInRegistry(target string) (string, string, error) {
	projects, err := loadActiveSortedProjects()
	if err != nil || len(projects) == 0 {
		return "", "", apperror.NewNotFound("project", "E9041", "no registered projects found to match: "+target)
	}

	matchedProj, matchErr := findProjectByFlexibleTarget(projects, target)
	if matchErr != nil {
		return "", "", apperror.NewNotFound("project", "E9042", "project target not found: "+target)
	}

	return matchedProj.GetPath(), matchedProj.Name, nil
}

func resolvePromptPayloadAndStage(promptMsg, projName, title string) (string, string, string, error) {
	if fi, err := os.Stat(promptMsg); err == nil && !fi.IsDir() {
		return loadPromptFromFile(promptMsg, title)
	}

	return stageInlinePrompt(promptMsg, projName, title)
}

func loadPromptFromFile(filePath, title string) (string, string, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", "", apperror.WrapSimple(err, "read prompt file")
	}

	resolvedTitle := title
	if resolvedTitle == "" {
		resolvedTitle = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return string(data), filePath, resolvedTitle, nil
}

func stageInlinePrompt(promptText, projName, title string) (string, string, string, error) {
	resolvedTitle := title
	if resolvedTitle == "" {
		resolvedTitle = fmt.Sprintf("Prompt for %s", projName)
	}

	stageFile := filepath.Join(os.TempDir(), fmt.Sprintf("agy-stage-prompt-%d.txt", time.Now().UnixNano()))
	if err := os.WriteFile(stageFile, []byte(promptText), 0644); err != nil {
		return "", "", "", apperror.WrapSimple(err, "write staged prompt")
	}

	return promptText, stageFile, resolvedTitle, nil
}

func resolveActiveAntigravityPID() int {
	ideProcRes := DetectRunningAntigravityIDE()

	return resolveActiveOrZeroPID(ideProcRes)
}

func dispatchAndRenderResult(repoRoot, stageFile, title, content string, pid int) error {
	res := DispatchPromptToAntigravity(repoRoot, stageFile, title, content, pid)
	if !res.IsSuccess {
		return apperror.NewExecutionError(res.Message)
	}

	printDispatchSuccess(res, repoRoot, title)

	return nil
}

func printDispatchSuccess(res AgyInjectionResult, repoRoot, title string) {
	fmt.Printf("\n  %s %s\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)
	fmt.Printf("    • Project: %s\n", repoRoot)
	fmt.Printf("    • Title:   %s\n", title)
	if res.PromptPath != "" {
		fmt.Printf("    • Staged:  %s\n", res.PromptPath)
	}
	fmt.Println()
}

func dispatchProjectReadPrompt(p AgyProject) error {
	pid := resolveActiveAntigravityPID()
	promptPath := filepath.Join(os.TempDir(), fmt.Sprintf("agy-read-prompt-%s.txt", p.ID))
	if writeErr := os.WriteFile(promptPath, []byte(defaultReadMemoryPrompt), 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write temp read prompt")
	}

	res := DispatchPromptToAntigravity(p.GetPath(), promptPath, p.Name, defaultReadMemoryPrompt, pid)
	if res.IsSuccess {
		fmt.Printf("  %s %s\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)
	}

	return handleRenameNotice(p)
}

func handleRenameNotice(p AgyProject) error {
	if renErr := renameProjectInitialConversation(p, p.Name); renErr != nil {
		fmt.Printf("  %s! Notice: rename conversation: %v%s\n", constants.ColorYellow, renErr, constants.ColorReset)
	}

	return nil
}
