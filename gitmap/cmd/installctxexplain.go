package cmd

import (
	"fmt"
	"strings"
)

// ctxExplainPrefixCmdExe returns a `cmd.exe`-style echo prefix for the
// Windows pwsh -Command body: when explain is on, prints
// `> gitmap <args>` before the real command runs.
//
// Empty string (no prefix) when explain is off, so callers can
// concatenate unconditionally without polluting the registry value
// with a no-op.
func ctxExplainPrefixPwsh(target string, args []string) string {
	if !ctxExplainEnabled {
		return ""
	}
	resolved := strings.Join(args, " ")

	return fmt.Sprintf(`Write-Host '> %s %s'; `, target, resolved)
}

// ctxExplainPrefixSh returns a POSIX `echo` prefix for macOS/Linux
// terminal-mode entries. Same on/off semantics as the pwsh variant.
func ctxExplainPrefixSh(target string, args []string) string {
	if !ctxExplainEnabled {
		return ""
	}
	resolved := strings.Join(args, " ")

	return fmt.Sprintf(`echo '> %s %s'; `, target, resolved)
}

// ctxExplainAnnounce returns "[explain mode] " for Silent-mode entries
// so the OS notification carries the resolved invocation in its
// payload (since there is no terminal to print to).
func ctxExplainAnnounce(target string, args []string) string {
	if !ctxExplainEnabled {
		return ""
	}

	return fmt.Sprintf("> %s %s\n", target, strings.Join(args, " "))
}
