package cmdagy

import (
	"fmt"
)

const maxDirectPromptChars = 24000

func resolvePayloadContentForAgent(promptPath, content string) string {
	hasExceeded := len(content) > maxDirectPromptChars
	if hasExceeded == false {
		return content
	}

	return formatLargePayloadInstruction(promptPath, content)
}

func formatLargePayloadInstruction(promptPath, content string) string {
	limit := 1500
	hasMore := len(content) > limit
	excerpt := content
	if hasMore {
		excerpt = content[:limit] + "\n...[truncated, see staged file for full logs]..."
	}

	return fmt.Sprintf("Execute 4-Part RCA on CI/CD pipeline errors. Full error payload staged at:\n%s\n\nPreview:\n%s\n\nPlease inspect the staged payload file and begin fixing the pipeline.", promptPath, excerpt)
}

func tryInjectExistingConversation(convID, title, content string, pid int, repoRoot, promptPath string) (AgyInjectionResult, bool) {
	sendRes := AgentAPISendMessage(convID, title, content)
	if sendRes.IsFailure() {
		return AgyInjectionResult{}, false
	}
	FocusAntigravityWindow(pid)

	return makeAgentAPISuccessResult(pid, repoRoot, promptPath, convID, false), true
}

func tryInjectNewConversation(title, content string, pid int, repoRoot, promptPath string) (AgyInjectionResult, bool) {
	newRes := AgentAPINewConversation(title, content)
	if newRes.IsFailure() {
		return AgyInjectionResult{}, false
	}
	FocusAntigravityWindow(pid)

	return makeAgentAPISuccessResult(pid, repoRoot, promptPath, newRes.Value, true), true
}

func tryDispatchExisting(repoRoot, promptPath, title, sendContent string, pid int) (AgyInjectionResult, bool) {
	conv, err := SelectMatchingConversation(repoRoot)
	hasConv := err == nil && len(conv.ID) > 0
	if !hasConv {
		return AgyInjectionResult{}, false
	}

	return tryInjectExistingConversation(conv.ID, title, sendContent, pid, repoRoot, promptPath)
}

// DispatchPromptToAntigravity sends the prompt to active or new Antigravity session via agentapi.
func DispatchPromptToAntigravity(repoRoot, promptPath, title, content string, pid int) AgyInjectionResult {
	sendContent := resolvePayloadContentForAgent(promptPath, content)
	existRes, hasExistSuccess := tryDispatchExisting(repoRoot, promptPath, title, sendContent, pid)
	if hasExistSuccess {
		return existRes
	}

	res, isSuccess := tryInjectNewConversation(title, sendContent, pid, repoRoot, promptPath)
	if isSuccess {
		return res
	}

	return makeOfflineFallbackResult(pid, repoRoot, promptPath)
}

func makeAgentAPISuccessResult(pid int, repoDir, promptPath, convID string, isNew bool) AgyInjectionResult {
	msg := formatAgentAPISuccessMsg(convID, pid, isNew)

	return AgyInjectionResult{
		IsSuccess:  true,
		Mode:       AgyInjectionModeIDE,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func formatAgentAPISuccessMsg(convID string, pid int, isNew bool) string {
	if isNew {
		return fmt.Sprintf("Created new Antigravity session (%s) and injected prompt via agentapi!", convID)
	}

	return fmt.Sprintf("Injected prompt into active Antigravity session (%s) via agentapi!", convID)
}

func makeOfflineFallbackResult(pid int, repoDir, promptPath string) AgyInjectionResult {
	copyClipboardIfNotSkipped(promptPath, false)
	msg := formatOfflineFallbackMsg(pid, promptPath)

	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeNone,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
	}
}

func formatOfflineFallbackMsg(pid int, promptPath string) string {
	hasPID := pid > 0
	if hasPID {
		return fmt.Sprintf("Antigravity IDE detected (PID: %d), but agentapi is unreachable; staged prompt in %s and copied to clipboard", pid, promptPath)
	}

	return fmt.Sprintf("Antigravity IDE offline; staged prompt in %s and copied to clipboard", promptPath)
}
