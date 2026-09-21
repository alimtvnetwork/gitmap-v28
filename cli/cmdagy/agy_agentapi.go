package cmdagy

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type rawAgentAPIResponse struct {
	Response rawAgentAPIBody `json:"response"`
	Error    string          `json:"error"`
}

type rawAgentAPIBody struct {
	SendMessage     *rawSendMessageResult     `json:"sendMessage,omitempty"`
	NewConversation *rawNewConversationResult `json:"newConversation,omitempty"`
}

type rawSendMessageResult struct {
	RecipientID string `json:"recipientId"`
	Content     string `json:"content"`
}

type rawNewConversationResult struct {
	ConversationID string `json:"conversationId"`
	Prompt         string `json:"prompt"`
}

// ResolveAgentAPI finds the agentapi binary or language_server fallback on the system.
func ResolveAgentAPI() (string, []string, bool) {
	homePath, homeArgs, hasHome := findAgentAPIFromHome()
	if hasHome {
		return homePath, homeArgs, true
	}
	installPath, installArgs, hasInstall := findAgentAPIFromInstallDir()
	if hasInstall {
		return installPath, installArgs, true
	}

	return findAgentAPIFromPath()
}

func findAgentAPIFromHome() (string, []string, bool) {
	home, err := os.UserHomeDir()
	hasHome := err == nil
	if hasHome == false {
		return "", nil, false
	}
	winBat := filepath.Join(home, ".gemini", "antigravity", "bin", "agentapi.bat")
	hasWinBat := checkFileExists(winBat)
	if hasWinBat {
		return winBat, []string{}, true
	}
	unixBin := filepath.Join(home, ".gemini", "antigravity", "bin", "agentapi")
	hasUnixBin := checkFileExists(unixBin)
	if hasUnixBin {
		return unixBin, []string{}, true
	}

	return "", nil, false
}

func findAgentAPIFromInstallDir() (string, []string, bool) {
	localApp := os.Getenv("LOCALAPPDATA")
	hasLocal := len(localApp) > 0
	if hasLocal == false {
		return "", nil, false
	}
	langServer := filepath.Join(localApp, "Programs", "Antigravity", "resources", "bin", "language_server.exe")
	hasLangServer := checkFileExists(langServer)
	if hasLangServer {
		return langServer, []string{"agentapi"}, true
	}

	return "", nil, false
}

func findAgentAPIFromPath() (string, []string, bool) {
	batPath, err := exec.LookPath("agentapi.bat")
	hasBat := err == nil
	if hasBat {
		return batPath, []string{}, true
	}
	exePath, err := exec.LookPath("agentapi")
	hasExe := err == nil
	if hasExe {
		return exePath, []string{}, true
	}

	return "", nil, false
}

func executeAgentAPICmd(subArgs []string) result.Result[[]byte] {
	binPath, baseArgs, isResolved := ResolveAgentAPI()
	if isResolved == false {
		return result.Fail[[]byte](apperror.NewSimple("agentapi binary not found", "E9030"))
	}
	fullArgs := append(baseArgs, subArgs...)
	cmd := exec.Command(binPath, fullArgs...)
	out, err := cmd.CombinedOutput()
	hasErr := err != nil
	if hasErr {
		errMsg := strings.TrimSpace(string(out))
		return result.Fail[[]byte](apperror.NewSimple("agentapi execution failed: "+errMsg, "E9031"))
	}

	return result.Ok(out)
}
