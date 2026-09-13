// Package cmdzsh provides text formatting for ZSH suite status output.
package cmdzsh

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// FormatStatusReport formats the ZshStatus struct into a readable summary string.
func FormatStatusReport(s ZshStatus) string {
	var sb strings.Builder
	sb.WriteString("=== ZSH Suite Status ===" + constants.NewLineUnix)
	sb.WriteString(fmt.Sprintf("  ZSH Installed:       %v%s", s.IsZshInstalled, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  ZSH Version:         %s%s", s.ZshVersion, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  ZSH Binary Path:     %s%s", s.ZshPath, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  Oh-My-Zsh Installed:  %v%s", s.IsOhMyZshInstalled, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  Oh-My-Zsh Path:      %s%s", s.OhMyZshPath, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  Current Theme:       %s%s", s.CurrentTheme, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  Active Plugins:      %v%s", s.ActivePlugins, constants.NewLineUnix))
	sb.WriteString(fmt.Sprintf("  Default Shell:       %s (active: %v)%s", s.DefaultShell, s.IsDefaultShell, constants.NewLineUnix))

	return sb.String()
}
