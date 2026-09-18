package cmdssh

import (
	"strings"
)

func determineSSHCommand(osType string, args []string) (string, string, bool) {
	if len(args) == 0 {
		return "", "", false
	}

	first := args[0]
	if isGitmapCommand(first) {
		return "", resolveGitmapCommandString(args), true
	}

	if isExplicitShell(first) {
		return first, extractShellCommandArgs(args), false
	}

	shell := "bash"
	if isWindowsOS(osType) {
		shell = "ps"
	}

	return shell, strings.Join(args, " "), false
}

func isGitmapCommand(first string) bool {
	return first == "gitmap" || first == "mkdir" || first == "cat" || first == "ssh" ||
		first == "agy" || first == "ag" || first == "antigravity" ||
		first == "schedule" || first == "schedules" ||
		first == "prompts-template" || first == "prompt-template" || first == "pt"
}

func resolveGitmapCommandString(args []string) string {
	if args[0] == "gitmap" {
		return strings.Join(args, " ")
	}

	return "gitmap " + strings.Join(args, " ")
}

func isExplicitShell(s string) bool {
	return s == "ps" || s == "cmd" || s == "bash" || s == "sh"
}

func extractShellCommandArgs(args []string) string {
	if len(args) > 1 {
		return strings.Join(args[1:], " ")
	}

	return ""
}
