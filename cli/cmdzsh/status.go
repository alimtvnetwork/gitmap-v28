// Package cmdzsh provides inspection and reporting of local ZSH installation state.
package cmdzsh

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// InspectStatus inspects local ZSH, Oh-My-Zsh, theme, and default shell status.
func InspectStatus(targetHome string) ZshStatusResult {
	home := ResolveTargetHome(targetHome)
	var s ZshStatus
	populateBinaryStatus(&s)
	populateOhMyZshStatus(&s, home)
	populateZshrcStatus(&s, home)

	return result.Ok(s)
}

func populateBinaryStatus(s *ZshStatus) {
	res := ResolveZshBinary()
	s.IsZshInstalled = res.IsSuccess()
	s.ZshPath = res.Value
	s.ZshVersion = queryZshVersion(s.IsZshInstalled)
	s.DefaultShell = os.Getenv("SHELL")
	s.IsDefaultShell = isCurrentShellZsh(s.DefaultShell)
}

func populateOhMyZshStatus(s *ZshStatus, home string) {
	s.OhMyZshPath = filepath.Join(home, ".oh-my-zsh")
	s.IsOhMyZshInstalled = HasFile(filepath.Join(s.OhMyZshPath, "oh-my-zsh.sh"))
}

func populateZshrcStatus(s *ZshStatus, home string) {
	zshrcPath := filepath.Join(home, ".zshrc")
	s.CurrentTheme = inspectThemeInZshrc(zshrcPath)
	s.ActivePlugins = inspectPluginsInZshrc(zshrcPath)
}

func queryZshVersion(isInstalled bool) string {
	if !isInstalled {
		return ""
	}

	out, err := exec.Command("zsh", "--version").Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func isCurrentShellZsh(defaultShell string) bool {
	return strings.HasSuffix(defaultShell, "/zsh") || defaultShell == "zsh"
}
