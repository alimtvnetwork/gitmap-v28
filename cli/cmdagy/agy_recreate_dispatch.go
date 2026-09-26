// Package cmdagy — agy_recreate_dispatch.go handles creating new conversation sessions and dispatching prompts.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const defaultRecreatePrompt = "Execute enhanced Read Memory protocol and analyze project structure. Thoroughly inspect all files, documentation, .ai-memory, specifications, and architecture to deeply understand the project."

func dispatchRecreateConversation(p AgyProject, opts AgyRecreateOptions) (string, error) {
	promptText := resolveRecreatePrompt(opts.CustomPrompt)
	title := fmt.Sprintf("Read & Understand: %s", p.Name)
	newRes := AgentAPINewConversationWithOptions(title, opts.Model, opts.Profile, promptText)
	if newRes.IsSuccess() {
		focusRunningIDE()
		return newRes.Value, nil
	}
	return fallbackPromptDispatch(p, title, promptText)
}

func resolveRecreatePrompt(customPrompt string) string {
	if customPrompt != "" {
		return customPrompt
	}
	return defaultRecreatePrompt
}

func focusRunningIDE() {
	ideProc := DetectRunningAntigravityIDE()
	if ideProc.IsSuccess() && ideProc.Value.PID > 0 {
		FocusAntigravityWindow(ideProc.Value.PID)
	}
}

func fallbackPromptDispatch(p AgyProject, title, promptText string) (string, error) {
	projectPath := p.GetPath()
	promptPath := filepath.Join(os.TempDir(), fmt.Sprintf("agy-recreate-prompt-%s.txt", p.ID))
	if writeErr := os.WriteFile(promptPath, []byte(promptText), 0644); writeErr != nil {
		return "", apperror.WrapSimple(writeErr, "write temp recreate prompt")
	}
	ideProc := DetectRunningAntigravityIDE()
	pid := resolveActiveOrZeroPID(ideProc)
	res := DispatchPromptToAntigravity(projectPath, promptPath, title, promptText, pid)
	if res.IsSuccess {
		fmt.Printf("  %s %s\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)
		return "staged", nil
	}
	return stageOfflinePrompt(projectPath, promptText)
}

func stageOfflinePrompt(projectPath, promptText string) (string, error) {
	if projectPath == "" {
		return "offline", nil
	}
	targetFile := filepath.Join(projectPath, activeAgyPromptRelativePath)
	writePromptFile(targetFile, promptText)
	fmt.Printf("  %s Antigravity offline; prompt staged to %s\n",
		constants.ColorYellow+"!"+constants.ColorReset, activeAgyPromptRelativePath)
	return "staged-offline", nil
}
