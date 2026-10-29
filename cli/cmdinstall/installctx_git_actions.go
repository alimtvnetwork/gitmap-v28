package cmdinstall

import (
	"os/exec"
	"runtime"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	gitActionsMu     sync.RWMutex
	gitActionsHidden bool
)

// IsGitActionsHidden reports whether git direct actions are hidden from context menu.
func IsGitActionsHidden() bool {
	if runtime.GOOS != "windows" {
		return isGitActionsHiddenMem()
	}

	return isWindowsGitActionsHidden()
}

func isGitActionsHidden() bool {
	return IsGitActionsHidden()
}

func isGitActionsHiddenMem() bool {
	gitActionsMu.RLock()
	defer gitActionsMu.RUnlock()

	return gitActionsHidden
}

// SetGitActionsHidden configures git direct action visibility.
func SetGitActionsHidden(hidden bool) error {
	gitActionsMu.Lock()
	gitActionsHidden = hidden
	gitActionsMu.Unlock()
	if runtime.GOOS == "windows" {
		return setWindowsGitActionsHidden(hidden)
	}

	return nil
}

func isWindowsGitActionsHidden() bool {
	out, err := exec.Command("reg", "query", constants.CtxRootKeyDirectory, "/v", "HideGitActions").Output()
	if err != nil {
		return isGitActionsHiddenMem()
	}

	return strings.Contains(string(out), "0x1")
}

func setWindowsGitActionsHidden(hidden bool) error {
	val := "0"
	if hidden {
		val = "1"
	}
	cmd := exec.Command("reg", "add", constants.CtxRootKeyDirectory, "/v", "HideGitActions", "/t", "REG_DWORD", "/d", val, "/f")

	return cmd.Run()
}
