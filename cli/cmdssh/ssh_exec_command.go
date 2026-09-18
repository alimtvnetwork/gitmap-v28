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
	switch first {
	case "gitmap", "status", "pipeline", "pipe", "pl", "clone", "pull", "sync", "push":
		return true
	case "clean", "log", "branch", "diff", "storage", "macro", "install", "update", "setup":
		return true
	case "chrome", "vscode", "vsc", "vhost", "zip", "service", "os", "schedule", "schedules":
		return true
	case "agy", "ag", "antigravity", "aef", "fix-pipeline", "prompts-template", "prompt-template", "pt":
		return true
	case "ssh", "se", "sj", "cluster", "sc", "mkdir", "cat":
		return true
	default:
		return false
	}
}

func resolveGitmapCommandString(args []string) string {
	if len(args) == 0 {
		return "gitmap"
	}
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
