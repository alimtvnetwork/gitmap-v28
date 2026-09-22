// Package cmdagy — agy_prompt_dispatch.go dispatches read memory prompt to Antigravity sessions.
package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func dispatchProjectReadPrompt(p AgyProject) error {
	ideProc := DetectRunningAntigravityIDE()
	pid := 0
	if ideProc.IsSuccess() {
		pid = ideProc.Value.PID
	}

	promptPath := filepath.Join(os.TempDir(), fmt.Sprintf("agy-read-prompt-%s.txt", p.ID))
	if writeErr := os.WriteFile(promptPath, []byte(defaultReadMemoryPrompt), 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write temp read prompt")
	}

	res := DispatchPromptToAntigravity(p.GetPath(), promptPath, p.Name, defaultReadMemoryPrompt, pid)
	if res.IsSuccess {
		fmt.Printf("  %s %s\n", constants.ColorGreen+"✓"+constants.ColorReset, res.Message)
	}
	if renErr := renameProjectInitialConversation(p, p.Name); renErr != nil {
		fmt.Printf("  %s! Notice: rename conversation: %v%s\n", constants.ColorYellow, renErr, constants.ColorReset)
	}

	return nil
}
