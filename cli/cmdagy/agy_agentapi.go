package cmdagy

import (
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
	installPath, installArgs, hasInstall := findAgentAPIFromInstallDir()
	if hasInstall {
		return installPath, installArgs, true
	}
	homePath, homeArgs, hasHome := findAgentAPIFromHome()
	if hasHome {
		return homePath, homeArgs, true
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
	return executeAgentAPICmdInDir("", subArgs)
}

func executeAgentAPICmdInDir(workDir string, subArgs []string) result.Result[[]byte] {
	binPath, baseArgs, isResolved := ResolveAgentAPI()
	if isResolved == false {
		return result.Fail[[]byte](apperror.NewSimple("agentapi binary not found", "E9030"))
	}
	fullArgs := append([]string{}, baseArgs...)
	fullArgs = append(fullArgs, subArgs...)
	cmd := exec.Command(binPath, fullArgs...)
	if workDir != "" && checkDirExists(workDir) {
		cmd.Dir = workDir
	} else {
		cmd.Dir = resolveAgentAPIWorkingDir()
	}
	cmd.Env = os.Environ()
	addr, token, hasEnv := ResolveAntigravityLSEnv()
	if hasEnv {
		cmd.Env = append(cmd.Env,
			"ANTIGRAVITY_LS_ADDRESS="+addr,
			"ANTIGRAVITY_CSRF_TOKEN="+token,
		)
	}
	out, err := cmd.CombinedOutput()
	if err == nil {
		return result.Ok(out)
	}
	if hasEnv {
		InvalidateAntigravityLSEnvCache()
	}
	errMsg := strings.TrimSpace(string(out))

	return result.Fail[[]byte](apperror.NewSimple("agentapi execution failed: "+errMsg, "E9031"))
}
