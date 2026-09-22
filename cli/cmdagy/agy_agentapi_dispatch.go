package cmdagy

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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

func tryInjectExistingConversation(convID, title, content string, pid int, repoRoot, promptPath string) (AgyInjectionResult, *apperror.AppError, bool) {
	sendRes := AgentAPISendMessage(convID, title, content)
	if sendRes.IsFailure() {
		return AgyInjectionResult{}, sendRes.Err, false
	}
	FocusAntigravityWindow(pid)

	return makeAgentAPISuccessResult(pid, repoRoot, promptPath, convID, false), nil, true
}

func tryInjectNewConversation(title, content string, pid int, repoRoot, promptPath string) (AgyInjectionResult, *apperror.AppError, bool) {
	newRes := AgentAPINewConversation(title, content)
	if newRes.IsFailure() {
		return AgyInjectionResult{}, newRes.Err, false
	}
	FocusAntigravityWindow(pid)

	return makeAgentAPISuccessResult(pid, repoRoot, promptPath, newRes.Value, true), nil, true
}

func tryDispatchExisting(repoRoot, promptPath, title, sendContent string, pid int) (AgyInjectionResult, *apperror.AppError, bool) {
	conv, err := SelectMatchingConversation(repoRoot)
	hasConv := err == nil && len(conv.ID) > 0
	if !hasConv {
		return AgyInjectionResult{}, nil, false
	}

	return tryInjectExistingConversation(conv.ID, title, sendContent, pid, repoRoot, promptPath)
}

// DispatchPromptToAntigravity sends the prompt to active or new Antigravity session via agentapi.
func DispatchPromptToAntigravity(repoRoot, promptPath, title, content string, pid int) AgyInjectionResult {
	sendContent := resolvePayloadContentForAgent(promptPath, content)
	existRes, existErr, hasExistSuccess := tryDispatchExisting(repoRoot, promptPath, title, sendContent, pid)
	if hasExistSuccess {
		return existRes
	}

	res, newErr, isSuccess := tryInjectNewConversation(title, sendContent, pid, repoRoot, promptPath)
	if isSuccess {
		return res
	}

	var appErr *apperror.AppError
	if newErr != nil {
		appErr = newErr
	} else if existErr != nil {
		appErr = existErr
	} else if pid <= 0 {
		appErr = apperror.NewNotFound("detect_ide", "E9002", "Antigravity IDE process is not running")
	} else {
		appErr = apperror.NewExecutionError("Antigravity agentapi is unreachable")
	}

	return makeOfflineFallbackResult(pid, repoRoot, promptPath, appErr)
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

func makeOfflineFallbackResult(pid int, repoDir, promptPath string, appErr *apperror.AppError) AgyInjectionResult {
	copyClipboardIfNotSkipped(promptPath, false)
	msg := formatOfflineFallbackMsg(pid, promptPath, appErr)

	return AgyInjectionResult{
		IsSuccess:  false,
		Mode:       AgyInjectionModeNone,
		PID:        pid,
		Message:    msg,
		PromptPath: promptPath,
		RepoDir:    repoDir,
		Err:        appErr,
	}
}

func formatOfflineFallbackMsg(pid int, promptPath string, appErr *apperror.AppError) string {
	var sb strings.Builder
	hasPID := pid > 0
	if hasPID {
		sb.WriteString(fmt.Sprintf("Antigravity IDE detected (PID: %d), but agentapi injection failed.", pid))
	} else {
		sb.WriteString("Antigravity IDE offline (process not running).")
	}
	sb.WriteString(fmt.Sprintf("\n  Staged prompt in: %s (copied to clipboard)", promptPath))
	if appErr != nil {
		sb.WriteString(fmt.Sprintf("\n  Error: %s", appErr.Error()))
		if len(appErr.Stack) > 0 {
			sb.WriteString(fmt.Sprintf("\n  Stack Trace:%s", appErr.Stack))
		}
	}
	sb.WriteString("\n  " + strings.ReplaceAll(DiagnoseAntigravityIDEAndCLI(), "\n", "\n  "))

	return sb.String()
}
